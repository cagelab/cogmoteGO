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
	Long:  "Interactively set scene name, source name, and install method.",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		currentSceneName := viper.GetString("obs.scene_name")
		currentSourceName := viper.GetString("obs.source_name")
		currentInstallMethod := viper.GetString("obs.install_method")

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

		fmt.Printf("Enter install method (system/flatpak) [%s]: ", currentInstallMethod)
		installMethod, _ := reader.ReadString('\n')
		installMethod = strings.TrimSpace(installMethod)
		if installMethod == "" {
			installMethod = currentInstallMethod
		}
		if installMethod != "system" && installMethod != "flatpak" {
			fmt.Fprintf(os.Stderr, "error: install_method must be 'system' or 'flatpak'\n")
			return
		}

		viper.Set("obs.scene_name", sceneName)
		viper.Set("obs.source_name", sourceName)
		viper.Set("obs.install_method", installMethod)

		if err := writeConfigWithFallback(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to save configuration: %v\n", err)
			return
		}

		fmt.Println("OBS configuration saved")
	},
}

var obsSetPasswordCmd = &cobra.Command{
	Use:   "set-password",
	Short: "Set OBS password",
	Long:  "Set OBS websocket password (stored in system keyring).",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter password: ")
		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read password: %v\n", err)
			return
		}
		password := strings.TrimSpace(string(passwordBytes))

		if password == "" {
			fmt.Fprintln(os.Stderr, "error: password cannot be empty")
			return
		}

		if err := keyring.SaveObsPassword(password); err != nil {
			fmt.Fprintf(os.Stderr, "failed to store password in keyring: %v\n", err)
			return
		}

		fmt.Println("Password saved")
	},
}

var obsShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current OBS configuration",
	Long:  "Display the current OBS configuration. Password is not shown, only whether it is configured.",
	Run: func(cmd *cobra.Command, args []string) {
		sceneName := viper.GetString("obs.scene_name")
		sourceName := viper.GetString("obs.source_name")
		installMethod := viper.GetString("obs.install_method")

		fmt.Println("OBS Configuration:")
		fmt.Printf("  Scene name     : %s\n", valueOrNotSet(sceneName))
		fmt.Printf("  Source name    : %s\n", valueOrNotSet(sourceName))
		fmt.Printf("  Install method : %s\n", valueOrNotSet(installMethod))

		if _, err := keyring.GetObsPassword(); err == nil {
			fmt.Printf("  Password       : configured\n")
		} else {
			fmt.Printf("  Password       : not configured\n")
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
	obsCmd.AddCommand(obsSetPasswordCmd)
	obsCmd.AddCommand(obsShowCmd)
	obsCmd.AddCommand(obsDeletePasswordCmd)
}
