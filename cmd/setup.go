package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

//go:embed skills
var skillsFS embed.FS

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Install the Claude Code skill into ~/.claude/skills/youtrack/",
	Long: `Installs a Claude Code skill that lets Claude agents invoke the yt CLI
on your behalf from any project. Run this once after installing yt.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return installSkill()
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func installSkill() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dest := filepath.Join(home, ".claude", "skills", "youtrack")
	const src = "skills/youtrack"

	return fs.WalkDir(skillsFS, src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dest, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		data, err := skillsFS.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(target, data, 0644); err != nil {
			return err
		}
		fmt.Printf("installed %s\n", target)
		return nil
	})
}
