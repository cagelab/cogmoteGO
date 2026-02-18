/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Ccccraz/cogmoteGO/internal/keyring"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var emailCmd = &cobra.Command{
	Use:   "email",
	Short: "Manage email credentials and SMTP configuration",
	Long:  "Manage email credentials stored in the system keyring and SMTP settings in the config file.",
}

var emailSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set email credentials and SMTP configuration",
	Long:  "Interactively set email address and password (stored in keyring), and SMTP host/port (stored in config).",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter email address: ")
		emailAddress, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read email address: %v\n", err)
			return
		}
		emailAddress = strings.TrimSpace(emailAddress)
		if emailAddress == "" {
			fmt.Fprintln(os.Stderr, "email address cannot be empty")
			return
		}

		fmt.Print("Enter email password: ")
		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read password: %v\n", err)
			return
		}
		password := strings.TrimSpace(string(passwordBytes))
		if password == "" {
			fmt.Fprintln(os.Stderr, "password cannot be empty")
			return
		}

		fmt.Print("Enter SMTP host: ")
		host, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read SMTP host: %v\n", err)
			return
		}
		host = strings.TrimSpace(host)
		if host == "" {
			fmt.Fprintln(os.Stderr, "SMTP host cannot be empty")
			return
		}

		fmt.Print("Enter SMTP port: ")
		portStr, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read SMTP port: %v\n", err)
			return
		}
		portStr = strings.TrimSpace(portStr)
		if portStr == "" {
			fmt.Fprintln(os.Stderr, "SMTP port cannot be empty")
			return
		}
		port, err := strconv.Atoi(portStr)
		if err != nil || port <= 0 {
			fmt.Fprintln(os.Stderr, "invalid SMTP port")
			return
		}

		if err := keyring.SaveCredentials(emailAddress, password); err != nil {
			fmt.Fprintf(os.Stderr, "failed to store credentials: %v\n", err)
			return
		}

		viper.Set("email.address", emailAddress)
		viper.Set("email.smtp_host", host)
		viper.Set("email.smtp_port", port)

		if err := viper.WriteConfig(); err != nil {
			configPath := viper.ConfigFileUsed()
			if configPath == "" {
				fmt.Fprintf(os.Stderr, "failed to save configuration: %v\n", err)
				return
			}
			if writeErr := viper.WriteConfigAs(configPath); writeErr != nil {
				fmt.Fprintf(os.Stderr, "failed to save configuration: %v\n", writeErr)
				return
			}
		}

		fmt.Println("email credentials saved to system keyring")
		fmt.Println("email configuration updated")
	},
}

var emailShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current email configuration",
	Long:  "Display the current email configuration. Passwords are not shown, only whether credentials are configured.",
	Run: func(cmd *cobra.Command, args []string) {
		sendEmail := viper.GetString("email.address")
		smtpHost := viper.GetString("email.smtp_host")
		smtpPort := viper.GetInt("email.smtp_port")
		recipients := viper.GetStringSlice("email.recipients")

		fmt.Println("Email Configuration:")
		fmt.Printf("  Email address : %s\n", valueOrNotSet(sendEmail))
		fmt.Printf("  SMTP host     : %s\n", valueOrNotSet(smtpHost))
		if smtpPort > 0 {
			fmt.Printf("  SMTP port     : %d\n", smtpPort)
		} else {
			fmt.Printf("  SMTP port     : (not set)\n")
		}
		if len(recipients) > 0 {
			fmt.Printf("  Recipients    : %s\n", strings.Join(recipients, ", "))
		} else {
			fmt.Printf("  Recipients    : (not set)\n")
		}

		if sendEmail != "" {
			if _, err := keyring.GetPassword(sendEmail); err == nil {
				fmt.Printf("  Credentials   : configured\n")
			} else {
				fmt.Printf("  Credentials   : not configured\n")
			}
		} else {
			fmt.Printf("  Credentials   : not configured\n")
		}
	},
}

var emailDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete email credentials and configuration",
	Long:  "Remove email credentials from the system keyring and clear email configuration from the config file.",
	Run: func(cmd *cobra.Command, args []string) {
		sendEmail := viper.GetString("email.address")

		if sendEmail != "" {
			if err := keyring.DeleteCredentials(sendEmail); err != nil {
				fmt.Fprintf(os.Stderr, "failed to delete credentials: %v\n", err)
			} else {
				fmt.Println("email credentials deleted from system keyring")
			}
		}

		viper.Set("email.address", "")
		viper.Set("email.smtp_host", "")
		viper.Set("email.smtp_port", 0)

		if err := viper.WriteConfig(); err != nil {
			configPath := viper.ConfigFileUsed()
			if configPath == "" {
				fmt.Fprintf(os.Stderr, "failed to save configuration: %v\n", err)
				return
			}
			if writeErr := viper.WriteConfigAs(configPath); writeErr != nil {
				fmt.Fprintf(os.Stderr, "failed to save configuration: %v\n", writeErr)
				return
			}
		}

		fmt.Println("email configuration cleared")
	},
}

func valueOrNotSet(s string) string {
	if s == "" {
		return "(not set)"
	}
	return s
}

var emailRecipientsCmd = &cobra.Command{
	Use:   "recipients",
	Short: "Manage email recipients (subscribers)",
	Long:  "Manage the list of email recipients who will receive notifications.",
}

var emailRecipientsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all email recipients",
	Long:  "Display all configured email recipients.",
	Run: func(cmd *cobra.Command, args []string) {
		recipients := viper.GetStringSlice("email.recipients")
		if len(recipients) == 0 {
			fmt.Println("No recipients configured")
			return
		}
		fmt.Println("Email Recipients:")
		for i, r := range recipients {
			fmt.Printf("  %d. %s\n", i+1, r)
		}
	},
}

var emailRecipientsAddCmd = &cobra.Command{
	Use:   "add [email]",
	Short: "Add an email recipient",
	Long:  "Add a new email address to the recipients list. If email is not provided, it will prompt interactively.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var email string
		if len(args) > 0 {
			email = strings.TrimSpace(args[0])
		} else {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter recipient email: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to read email: %v\n", err)
				return
			}
			email = strings.TrimSpace(input)
		}

		if email == "" {
			fmt.Fprintln(os.Stderr, "email cannot be empty")
			return
		}

		recipients := viper.GetStringSlice("email.recipients")
		for _, r := range recipients {
			if strings.EqualFold(r, email) {
				fmt.Fprintln(os.Stderr, "recipient already exists")
				return
			}
		}

		recipients = append(recipients, email)
		if err := saveRecipients(recipients); err != nil {
			fmt.Fprintf(os.Stderr, "failed to save recipients: %v\n", err)
			return
		}

		fmt.Printf("added recipient: %s\n", email)
	},
}

var emailRecipientsRemoveCmd = &cobra.Command{
	Use:   "remove [email]",
	Short: "Remove an email recipient",
	Long:  "Remove an email address from the recipients list. If email is not provided, it will prompt interactively.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var email string
		if len(args) > 0 {
			email = strings.TrimSpace(args[0])
		} else {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter recipient email to remove: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to read email: %v\n", err)
				return
			}
			email = strings.TrimSpace(input)
		}

		if email == "" {
			fmt.Fprintln(os.Stderr, "email cannot be empty")
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
			fmt.Fprintf(os.Stderr, "recipient not found: %s\n", email)
			return
		}

		if err := saveRecipients(filtered); err != nil {
			fmt.Fprintf(os.Stderr, "failed to save recipients: %v\n", err)
			return
		}

		fmt.Printf("removed recipient: %s\n", email)
	},
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

func init() {
	rootCmd.AddCommand(emailCmd)
	emailCmd.AddCommand(emailSetCmd)
	emailCmd.AddCommand(emailShowCmd)
	emailCmd.AddCommand(emailDeleteCmd)
	emailCmd.AddCommand(emailRecipientsCmd)

	emailRecipientsCmd.AddCommand(emailRecipientsListCmd)
	emailRecipientsCmd.AddCommand(emailRecipientsAddCmd)
	emailRecipientsCmd.AddCommand(emailRecipientsRemoveCmd)
}
