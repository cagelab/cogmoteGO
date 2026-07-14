package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Ccccraz/cogmoteGO/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Manage backup trusted roots",
}

var backupRootsCmd = &cobra.Command{
	Use:   "roots",
	Short: "Manage backup source and Samba roots",
}

var backupRootsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured backup roots",
	Run: func(cmd *cobra.Command, args []string) {
		for _, rootType := range []string{"source", "samba"} {
			roots, err := loadBackupRoots(rootType)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to load %s roots: %v\n", rootType, err)
				return
			}
			fmt.Printf("%s roots:\n", rootType)
			for _, root := range roots {
				fmt.Printf("  %s: %s\n", root.ID, root.Path)
			}
		}
	},
}

var backupRootsAddCmd = &cobra.Command{
	Use:   "add <source|samba> <id> <absolute-path>",
	Short: "Add a trusted backup root",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		rootType, id, path := args[0], args[1], args[2]
		if !validRootType(rootType) {
			fmt.Fprintln(os.Stderr, "root type must be source or samba")
			return
		}
		if !config.IsValidBackupRootID(id) {
			fmt.Fprintln(os.Stderr, "root id may contain only letters, numbers, underscores, and hyphens")
			return
		}
		if !filepath.IsAbs(path) {
			fmt.Fprintln(os.Stderr, "root path must be absolute")
			return
		}
		info, err := os.Stat(path)
		if err != nil || !info.IsDir() {
			fmt.Fprintf(os.Stderr, "root path must be an existing directory: %s\n", path)
			return
		}
		roots, err := loadBackupRoots(rootType)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load roots: %v\n", err)
			return
		}
		for _, root := range roots {
			if root.ID == id {
				fmt.Fprintf(os.Stderr, "root id already exists: %s\n", id)
				return
			}
		}
		viper.Set(backupRootKey(rootType), append(roots, config.BackupRoot{ID: id, Path: path}))
		if err := writeConfigWithFallback(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to save backup root: %v\n", err)
			return
		}
		fmt.Printf("%s root %s added\n", rootType, id)
	},
}

var backupRootsRemoveCmd = &cobra.Command{
	Use:   "remove <source|samba> <id>",
	Short: "Remove a trusted backup root",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		rootType, id := args[0], args[1]
		if !validRootType(rootType) {
			fmt.Fprintln(os.Stderr, "root type must be source or samba")
			return
		}
		roots, err := loadBackupRoots(rootType)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load roots: %v\n", err)
			return
		}
		filtered := make([]config.BackupRoot, 0, len(roots))
		found := false
		for _, root := range roots {
			if root.ID == id {
				found = true
				continue
			}
			filtered = append(filtered, root)
		}
		if !found {
			fmt.Fprintf(os.Stderr, "root id not found: %s\n", id)
			return
		}
		viper.Set(backupRootKey(rootType), filtered)
		if err := writeConfigWithFallback(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to save backup root: %v\n", err)
			return
		}
		fmt.Printf("%s root %s removed\n", rootType, id)
	},
}

func loadBackupRoots(rootType string) ([]config.BackupRoot, error) {
	var roots []config.BackupRoot
	if err := viper.UnmarshalKey(backupRootKey(rootType), &roots); err != nil {
		return nil, err
	}
	return roots, nil
}

func backupRootKey(rootType string) string {
	return "backup." + rootType + "_roots"
}

func validRootType(rootType string) bool {
	return rootType == "source" || rootType == "samba"
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.AddCommand(backupRootsCmd)
	backupRootsCmd.AddCommand(backupRootsListCmd, backupRootsAddCmd, backupRootsRemoveCmd)
}
