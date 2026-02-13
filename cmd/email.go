/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
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
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")

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

		if err := keyring.SaveCredentials(emailAddress, password); err != nil {
			fmt.Fprintf(os.Stderr, "failed to store credentials: %v\n", err)
			return
		}

		viper.Set("email.address", emailAddress)
		if host != "" {
			viper.Set("email.smtp_host", host)
		}
		if port > 0 {
			viper.Set("email.smtp_port", port)
		}

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

func init() {
	rootCmd.AddCommand(emailCmd)
	emailCmd.AddCommand(emailSetCmd)
	emailCmd.AddCommand(emailShowCmd)
	emailCmd.AddCommand(emailDeleteCmd)

	emailSetCmd.Flags().String("host", "", "SMTP host address")
	emailSetCmd.Flags().Int("port", 0, "SMTP port")
}
