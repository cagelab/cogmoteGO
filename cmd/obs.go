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

var obsCmd = &cobra.Command{
	Use:   "obs",
	Short: "Manage OBS configuration",
	Long:  "Manage OBS scene/source settings and password stored in the system keyring.",
}

var obsSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set OBS configuration interactively",
	Long:  "Interactively set scene name, source name, and password (password stored in keyring).",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		currentSceneName := viper.GetString("obs.scene_name")
		currentSourceName := viper.GetString("obs.source_name")

		fmt.Printf("Enter scene name [%s]: ", currentSceneName)
		sceneName, _ := reader.ReadString('\n')
		sceneName = strings.TrimSpace(sceneName)
		if sceneName == "" {
			sceneName = currentSceneName
		}

		fmt.Printf("Enter source name [%s]: ", currentSourceName)
		sourceName, _ := reader.ReadString('\n')
		sourceName = strings.TrimSpace(sourceName)
		if sourceName == "" {
			sourceName = currentSourceName
		}

		if sceneName == sourceName {
			fmt.Fprintln(os.Stderr, "error: scene_name and source_name cannot be the same (OBS limitation)")
			return
		}

		fmt.Print("Enter password (will be stored in keyring): ")
		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read password: %v\n", err)
			return
		}
		password := strings.TrimSpace(string(passwordBytes))

		viper.Set("obs.scene_name", sceneName)
		viper.Set("obs.source_name", sourceName)

		if err := writeConfigWithFallback(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to save configuration: %v\n", err)
			return
		}

		if password != "" {
			if err := keyring.SaveObsPassword(password); err != nil {
				fmt.Fprintf(os.Stderr, "failed to store password in keyring: %v\n", err)
				return
			}
		}

		fmt.Println("OBS configuration saved")
	},
}

var obsShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current OBS configuration",
	Long:  "Display the current OBS configuration. Password is not shown, only whether it is configured.",
	Run: func(cmd *cobra.Command, args []string) {
		sceneName := viper.GetString("obs.scene_name")
		sourceName := viper.GetString("obs.source_name")

		fmt.Println("OBS Configuration:")
		fmt.Printf("  Scene name  : %s\n", valueOrNotSet(sceneName))
		fmt.Printf("  Source name : %s\n", valueOrNotSet(sourceName))

		if _, err := keyring.GetObsPassword(); err == nil {
			fmt.Printf("  Password    : configured\n")
		} else {
			fmt.Printf("  Password    : not configured\n")
		}
	},
}

var obsDeletePasswordCmd = &cobra.Command{
	Use:   "delete-password",
	Short: "Delete OBS password from keyring",
	Long:  "Remove the OBS password from the system keyring. Scene and source settings are preserved.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := keyring.DeleteObsPassword(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to delete password: %v\n", err)
			return
		}
		fmt.Println("OBS password deleted from system keyring")
	},
}

func init() {
	rootCmd.AddCommand(obsCmd)
	obsCmd.AddCommand(obsSetCmd)
	obsCmd.AddCommand(obsShowCmd)
	obsCmd.AddCommand(obsDeletePasswordCmd)
}
