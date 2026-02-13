package email

import (
	"bytes"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Ccccraz/cogmoteGO/internal/commonTypes"
	"github.com/Ccccraz/cogmoteGO/internal/keyring"
	"github.com/Ccccraz/cogmoteGO/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/wneessen/go-mail"
)

type emailPayload struct {
	Subject     string            `json:"subject"`
	HTMLBody    string            `json:"html_body"`
	Attachments []emailAttachment `json:"attachments"`
}

type emailAttachment struct {
	Filename string `json:"filename"`
	Content  []byte `json:"content"`
}

type emailConfig struct {
	From       string
	Password   string
	Host       string
	Port       int
	Recipients []string
}

const logKey = "email"

func PostEmailHandler(c *gin.Context) {
	payload, ok := parseEmailPayload(c)
	if !ok {
		return
	}

	cfg, ok := loadEmailConfig(c)
	if !ok {
		return
	}

	message, ok := buildEmailMessage(c, cfg, payload)
	if !ok {
		return
	}

	if !deliverEmail(c, cfg, message) {
		return
	}

	c.Status(http.StatusCreated)
}

func parseEmailPayload(c *gin.Context) (emailPayload, bool) {
	var payload emailPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Logger.Error("invalid email payload",
			slog.Group(logKey,
				slog.String("detail", err.Error()),
			),
		)
		respondError(c, http.StatusBadRequest, "invalid email payload", err.Error())
		return emailPayload{}, false
	}

	payload.Subject = strings.TrimSpace(payload.Subject)
	payload.HTMLBody = strings.TrimSpace(payload.HTMLBody)

	if payload.Subject == "" {
		logger.Logger.Error("email subject cannot be empty")
		respondError(c, http.StatusBadRequest, "email subject cannot be empty", "")
		return emailPayload{}, false
	}

	if payload.HTMLBody == "" {
		logger.Logger.Error("email body cannot be empty")
		respondError(c, http.StatusBadRequest, "email body cannot be empty", "")
		return emailPayload{}, false
	}

	for idx := range payload.Attachments {
		payload.Attachments[idx].Filename = strings.TrimSpace(payload.Attachments[idx].Filename)
		if payload.Attachments[idx].Filename == "" {
			logger.Logger.Error("attachment filename cannot be empty",
				slog.Group(logKey,
					slog.Int("index", idx),
				),
			)
			respondError(c, http.StatusBadRequest, "attachment filename cannot be empty", "")
			return emailPayload{}, false
		}
		if len(payload.Attachments[idx].Content) == 0 {
			logger.Logger.Error("attachment content cannot be empty",
				slog.Group(logKey,
					slog.String("filename", payload.Attachments[idx].Filename),
				),
			)
			respondError(c, http.StatusBadRequest, "attachment content cannot be empty", "")
			return emailPayload{}, false
		}
	}

	return payload, true
}

func loadEmailConfig(c *gin.Context) (emailConfig, bool) {
	// Environment variables take priority over config/keyring.
	sendEmail := strings.TrimSpace(os.Getenv("COGMOTE_EMAIL_ADDRESS"))
	smtpHost := strings.TrimSpace(os.Getenv("COGMOTE_SMTP_HOST"))
	smtpPortStr := strings.TrimSpace(os.Getenv("COGMOTE_SMTP_PORT"))
	password := strings.TrimSpace(os.Getenv("COGMOTE_EMAIL_PASSWORD"))

	var smtpPort int
	if smtpPortStr != "" {
		if v, err := strconv.Atoi(smtpPortStr); err == nil {
			smtpPort = v
		}
	}

	// Fall back to config file values when env vars are not set.
	emailSection := viper.Sub("email")

	if sendEmail == "" && emailSection != nil {
		sendEmail = strings.TrimSpace(emailSection.GetString("address"))
	}
	if smtpHost == "" && emailSection != nil {
		smtpHost = strings.TrimSpace(emailSection.GetString("smtp_host"))
	}
	if smtpPort <= 0 && emailSection != nil {
		smtpPort = emailSection.GetInt("smtp_port")
	}

	// Validate required fields.
	if sendEmail == "" {
		logger.Logger.Error("email address not configured")
		respondError(c, http.StatusInternalServerError, "email address not configured", "")
		return emailConfig{}, false
	}
	if smtpHost == "" {
		logger.Logger.Error("smtp_host not configured")
		respondError(c, http.StatusInternalServerError, "smtp_host not configured", "")
		return emailConfig{}, false
	}
	if smtpPort <= 0 {
		logger.Logger.Error("smtp_port not configured",
			slog.Group(logKey,
				slog.Int("value", smtpPort),
			),
		)
		respondError(c, http.StatusInternalServerError, "smtp_port not configured", "")
		return emailConfig{}, false
	}

	// Recipients come from the config file only.
	var rawRecipients []string
	if emailSection != nil {
		rawRecipients = emailSection.GetStringSlice("recipients")
	}
	recipients := make([]string, 0, len(rawRecipients))
	for _, recipient := range rawRecipients {
		recipient = strings.TrimSpace(recipient)
		if recipient != "" {
			recipients = append(recipients, recipient)
		}
	}
	if len(recipients) == 0 {
		logger.Logger.Error("recipients not configured")
		respondError(c, http.StatusInternalServerError, "recipients not configured", "")
		return emailConfig{}, false
	}

	// Fall back to keyring when the password env var is not set.
	if password == "" {
		var err error
		password, err = keyring.GetPassword(sendEmail)
		if err != nil {
			logger.Logger.Error("email password not found",
				slog.Group(logKey,
					slog.String("detail", err.Error()),
				),
			)
			respondError(c, http.StatusInternalServerError, "email password not found", err.Error())
			return emailConfig{}, false
		}
	}

	return emailConfig{
		From:       sendEmail,
		Password:   password,
		Host:       smtpHost,
		Port:       smtpPort,
		Recipients: recipients,
	}, true
}

func buildEmailMessage(c *gin.Context, cfg emailConfig, payload emailPayload) (*mail.Msg, bool) {
	message := mail.NewMsg()
	if err := message.From(cfg.From); err != nil {
		logger.Logger.Error("failed to prepare email",
			slog.Group(logKey,
				slog.String("detail", err.Error()),
			),
		)
		respondError(c, http.StatusInternalServerError, "failed to prepare email", err.Error())
		return nil, false
	}

	if err := message.To(cfg.Recipients...); err != nil {
		logger.Logger.Error("failed to prepare email",
			slog.Group(logKey,
				slog.String("detail", err.Error()),
			),
		)
		respondError(c, http.StatusInternalServerError, "failed to prepare email", err.Error())
		return nil, false
	}

	message.Subject(payload.Subject)

	message.SetBodyString(mail.TypeTextHTML, payload.HTMLBody)

	for _, attachment := range payload.Attachments {
		if err := message.AttachReader(attachment.Filename, bytes.NewReader(attachment.Content)); err != nil {
			logger.Logger.Error("invalid attachment",
				slog.Group(logKey,
					slog.String("detail", err.Error()),
					slog.String("filename", attachment.Filename),
				),
			)
			respondError(c, http.StatusBadRequest, "invalid attachment", err.Error())
			return nil, false
		}
	}

	return message, true
}

func deliverEmail(c *gin.Context, cfg emailConfig, message *mail.Msg) bool {
	client, err := mail.NewClient(
		cfg.Host,
		mail.WithPort(cfg.Port),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithUsername(cfg.From),
		mail.WithPassword(cfg.Password),
	)
	if err != nil {
		logger.Logger.Error("failed to send email",
			slog.Group(logKey,
				slog.String("detail", err.Error()),
			),
		)
		respondError(c, http.StatusInternalServerError, "failed to send email", err.Error())
		return false
	}

	if err := client.DialAndSend(message); err != nil {
		logger.Logger.Error("failed to send email",
			slog.Group(logKey,
				slog.String("detail", err.Error()),
			),
		)
		respondError(c, http.StatusInternalServerError, "failed to send email", err.Error())
		return false
	}

	return true
}

func respondError(c *gin.Context, status int, userMessage string, detail string) {
	c.JSON(status, commonTypes.APIError{
		Error:  userMessage,
		Detail: detail,
	})
}

func RegisterRoutes(r gin.IRouter) {
	r.POST("/email", PostEmailHandler)
}
