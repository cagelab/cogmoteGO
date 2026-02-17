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
	Embeds      []emailEmbed      `json:"embeds"`
	InReplyTo   string            `json:"in_reply_to"`
}

type emailAttachment struct {
	Filename string `json:"filename"`
	Content  []byte `json:"content"`
}

type emailEmbed struct {
	ContentID string `json:"content_id"`
	Filename  string `json:"filename"`
	Content   []byte `json:"content"`
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

	c.JSON(http.StatusCreated, gin.H{"message_id": message.GetMessageID()})
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

	for idx := range payload.Embeds {
		payload.Embeds[idx].ContentID = strings.TrimSpace(payload.Embeds[idx].ContentID)
		payload.Embeds[idx].Filename = strings.TrimSpace(payload.Embeds[idx].Filename)

		if payload.Embeds[idx].ContentID == "" {
			logger.Logger.Error("embed content_id cannot be empty",
				slog.Group(logKey,
					slog.Int("index", idx),
				),
			)
			respondError(c, http.StatusBadRequest, "embed content_id cannot be empty", "")
			return emailPayload{}, false
		}
		if payload.Embeds[idx].Filename == "" {
			logger.Logger.Error("embed filename cannot be empty",
				slog.Group(logKey,
					slog.Int("index", idx),
					slog.String("content_id", payload.Embeds[idx].ContentID),
				),
			)
			respondError(c, http.StatusBadRequest, "embed filename cannot be empty", "")
			return emailPayload{}, false
		}
		if len(payload.Embeds[idx].Content) == 0 {
			logger.Logger.Error("embed content cannot be empty",
				slog.Group(logKey,
					slog.String("content_id", payload.Embeds[idx].ContentID),
					slog.String("filename", payload.Embeds[idx].Filename),
				),
			)
			respondError(c, http.StatusBadRequest, "embed content cannot be empty", "")
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

	if payload.InReplyTo != "" {
		message.SetGenHeader(mail.HeaderInReplyTo, payload.InReplyTo)
	}

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

	for _, embed := range payload.Embeds {
		if err := message.EmbedReader(embed.Filename, bytes.NewReader(embed.Content), mail.WithFileContentID(embed.ContentID)); err != nil {
			logger.Logger.Error("invalid embed",
				slog.Group(logKey,
					slog.String("detail", err.Error()),
					slog.String("content_id", embed.ContentID),
					slog.String("filename", embed.Filename),
				),
			)
			respondError(c, http.StatusBadRequest, "invalid embed", err.Error())
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

type recipientPayload struct {
	Email string `json:"email"`
}

func GetRecipientsHandler(c *gin.Context) {
	recipients := viper.GetStringSlice("email.recipients")
	c.JSON(http.StatusOK, gin.H{"recipients": recipients})
}

func AddRecipientHandler(c *gin.Context) {
	var payload recipientPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}

	email := strings.TrimSpace(payload.Email)
	if email == "" {
		respondError(c, http.StatusBadRequest, "email cannot be empty", "")
		return
	}

	recipients := viper.GetStringSlice("email.recipients")
	for _, r := range recipients {
		if strings.EqualFold(r, email) {
			respondError(c, http.StatusConflict, "recipient already exists", "")
			return
		}
	}

	recipients = append(recipients, email)
	if err := saveRecipients(recipients); err != nil {
		respondError(c, http.StatusInternalServerError, "failed to save configuration", err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"recipients": recipients})
}

func DeleteRecipientHandler(c *gin.Context) {
	var payload recipientPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, "invalid payload", err.Error())
		return
	}

	email := strings.TrimSpace(payload.Email)
	if email == "" {
		respondError(c, http.StatusBadRequest, "email cannot be empty", "")
		return
	}

	recipients := viper.GetStringSlice("email.recipients")
	filtered := make([]string, 0, len(recipients))
	found := false
	for _, r := range recipients {
		if strings.EqualFold(r, email) {
			found = true
			continue
		}
		filtered = append(filtered, r)
	}

	if !found {
		respondError(c, http.StatusNotFound, "recipient not found", "")
		return
	}

	if err := saveRecipients(filtered); err != nil {
		respondError(c, http.StatusInternalServerError, "failed to save configuration", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"recipients": filtered})
}

func saveRecipients(recipients []string) error {
	viper.Set("email.recipients", recipients)
	if err := viper.WriteConfig(); err != nil {
		configPath := viper.ConfigFileUsed()
		if configPath == "" {
			return err
		}
		return viper.WriteConfigAs(configPath)
	}
	return nil
}

func GetConfigHandler(c *gin.Context) {
	address := viper.GetString("email.address")
	smtpHost := viper.GetString("email.smtp_host")
	smtpPort := viper.GetInt("email.smtp_port")
	recipients := viper.GetStringSlice("email.recipients")

	hasCredentials := false
	if address != "" {
		if _, err := keyring.GetPassword(address); err == nil {
			hasCredentials = true
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"address":         address,
		"smtp_host":       smtpHost,
		"smtp_port":       smtpPort,
		"recipients":      recipients,
		"has_credentials": hasCredentials,
	})
}

func PostTestEmailHandler(c *gin.Context) {
	cfg, ok := loadEmailConfig(c)
	if !ok {
		return
	}

	message := mail.NewMsg()
	if err := message.From(cfg.From); err != nil {
		respondError(c, http.StatusInternalServerError, "failed to prepare test email", err.Error())
		return
	}
	if err := message.To(cfg.Recipients...); err != nil {
		respondError(c, http.StatusInternalServerError, "failed to prepare test email", err.Error())
		return
	}
	message.Subject("cogmoteGO test email")
	message.SetBodyString(mail.TypeTextHTML, "<p>This is a test email from <b>cogmoteGO</b>. Your email configuration is working correctly.</p>")

	if !deliverEmail(c, cfg, message) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "test email sent successfully",
		"message_id": message.GetMessageID(),
	})
}

func RegisterRoutes(r gin.IRouter) {
	r.POST("/email", PostEmailHandler)
	r.GET("/email/config", GetConfigHandler)
	r.POST("/email/test", PostTestEmailHandler)
	r.GET("/email/recipients", GetRecipientsHandler)
	r.POST("/email/recipients", AddRecipientHandler)
	r.DELETE("/email/recipients", DeleteRecipientHandler)
}
