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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type configField struct {
	key          string
	fieldType    string // "int" or "string"
	defaultValue any
	validate     func(string) error
	description  string
	group        string
}

var configFields = []configField{
	{
		key:          "port",
		fieldType:    "int",
		defaultValue: 9012,
		validate:     validatePort,
		description:  "Server listening port",
		group:        "Server",
	},
	{
		key:          "instance_id",
		fieldType:    "string",
		defaultValue: "",
		validate: func(v string) error {
			if v == "" {
				return fmt.Errorf("instance_id cannot be empty")
			}
			return nil
		},
		description: "Unique instance identifier",
		group:       "Server",
	},
	{
		key:          "proxy.handshake_timeout",
		fieldType:    "int",
		defaultValue: 5000,
		validate:     validatePositiveInt,
		description:  "WebSocket handshake timeout (ms)",
		group:        "Proxy",
	},
	{
		key:          "proxy.msg_timeout",
		fieldType:    "int",
		defaultValue: 5000,
		validate:     validatePositiveInt,
		description:  "Message timeout (ms)",
		group:        "Proxy",
	},
	{
		key:          "proxy.max_retries",
		fieldType:    "int",
		defaultValue: 3,
		validate:     validatePositiveInt,
		description:  "Maximum retry attempts",
		group:        "Proxy",
	},
	{
		key:          "proxy.retry_interval",
		fieldType:    "int",
		defaultValue: 200,
		validate:     validatePositiveInt,
		description:  "Retry interval (ms)",
		group:        "Proxy",
	},
	{
		key:          "obs.scene_name",
		fieldType:    "string",
		defaultValue: "cagelab",
		validate: func(v string) error {
			if v == "" {
				return fmt.Errorf("scene_name cannot be empty")
			}
			return nil
		},
		description: "OBS scene name",
		group:       "OBS",
	},
	{
		key:          "obs.source_name",
		fieldType:    "string",
		defaultValue: "cogmoteGO",
		validate: func(v string) error {
			if v == "" {
				return fmt.Errorf("source_name cannot be empty")
			}
			return nil
		},
		description: "OBS text source name",
		group:       "OBS",
	},
}

func findConfigField(key string) *configField {
	for i := range configFields {
		if configFields[i].key == key {
			return &configFields[i]
		}
	}
	return nil
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage application configuration",
	Long:  "View, modify, and reset application configuration settings.",
}

var configShowCmd = &cobra.Command{
	Use:   "show [key]",
	Short: "Show configuration values",
	Long:  "Display all configuration values or a specific key.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			key := args[0]
			field := findConfigField(key)
			if field == nil {
				fmt.Fprintf(os.Stderr, "unknown configuration key: %s\n", key)
				return
			}
			fmt.Printf("%s = %v\n", key, viper.Get(key))
			return
		}

		// Show all, grouped
		currentGroup := ""
		for _, f := range configFields {
			if f.group != currentGroup {
				if currentGroup != "" {
					fmt.Println()
				}
				fmt.Printf("%s:\n", f.group)
				currentGroup = f.group
			}
			// Compute display label from key
			label := f.key
			fmt.Printf("  %-27s : %v\n", label, viper.Get(f.key))
		}
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long:  "Set the value of a configuration key. The value is validated before saving.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		field := findConfigField(key)
		if field == nil {
			fmt.Fprintf(os.Stderr, "unknown configuration key: %s\n", key)
			return
		}

		if err := field.validate(value); err != nil {
			fmt.Fprintf(os.Stderr, "invalid value for %s: %v\n", key, err)
			return
		}

		// Confirm before changing instance_id
		if key == "instance_id" {
			fmt.Print("Changing instance_id may affect identification. Continue? [y/N]: ")
			if !confirmAction() {
				fmt.Println("aborted")
				return
			}
		}

		switch field.fieldType {
		case "int":
			n, _ := strconv.Atoi(value) // already validated
			viper.Set(key, n)
		case "string":
			viper.Set(key, value)
		}

		if err := writeConfigWithFallback(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to save configuration: %v\n", err)
			return
		}

		fmt.Printf("%s set to %s\n", key, value)
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset [key]",
	Short: "Reset configuration to default values",
	Long:  "Reset a specific key or all non-email configuration to default values.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			key := args[0]
			field := findConfigField(key)
			if field == nil {
				fmt.Fprintf(os.Stderr, "unknown configuration key: %s\n", key)
				return
			}
			viper.Set(key, field.defaultValue)
			if err := writeConfigWithFallback(); err != nil {
				fmt.Fprintf(os.Stderr, "failed to save configuration: %v\n", err)
				return
			}
			fmt.Printf("%s reset to %v\n", key, field.defaultValue)
			return
		}

		// Reset all
		fmt.Print("Reset all configuration to defaults (email settings will not be affected)? [y/N]: ")
		if !confirmAction() {
			fmt.Println("aborted")
			return
		}

		for _, f := range configFields {
			viper.Set(f.key, f.defaultValue)
		}

		if err := writeConfigWithFallback(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to save configuration: %v\n", err)
			return
		}
		fmt.Println("all configuration reset to defaults")
	},
}

func writeConfigWithFallback() error {
	if err := viper.WriteConfig(); err != nil {
		configPath := viper.ConfigFileUsed()
		if configPath == "" {
			return err
		}
		if writeErr := viper.WriteConfigAs(configPath); writeErr != nil {
			return writeErr
		}
	}
	return nil
}

func validatePort(v string) error {
	n, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("must be an integer")
	}
	if n < 1 || n > 65535 {
		return fmt.Errorf("must be between 1 and 65535")
	}
	return nil
}

func validatePositiveInt(v string) error {
	n, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("must be an integer")
	}
	if n <= 0 {
		return fmt.Errorf("must be greater than 0")
	}
	return nil
}

func confirmAction() bool {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configResetCmd)
}
