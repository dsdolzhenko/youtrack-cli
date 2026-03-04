package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	serverURL  string
	token      string
	jsonOutput bool
)

var rootCmd = &cobra.Command{
	Use:          "yt",
	Short:        "YouTrack CLI — access YouTrack from command line and agents",
	SilenceErrors: true,
}

func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Silence usage for runtime errors; args validation errors still show it
		// because PersistentPreRunE runs after args validation but before RunE.
		cmd.SilenceUsage = true
		loadConfig()
		return nil
	}

	rootCmd.PersistentFlags().StringVar(&serverURL, "url", "", "YouTrack server URL (env: YOUTRACK_URL)")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "YouTrack API token (env: YOUTRACK_TOKEN)")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
}

func writeJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func requireArgs(names ...string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == len(names) {
			return nil
		}
		if len(args) > len(names) {
			return fmt.Errorf("accepts %d argument(s), received %d", len(names), len(args))
		}
		missing := make([]string, len(names)-len(args))
		for i, n := range names[len(args):] {
			missing[i] = "<" + n + ">"
		}
		if len(missing) == 1 {
			return fmt.Errorf("missing required argument %s", missing[0])
		}
		return fmt.Errorf("missing required arguments %s", strings.Join(missing, " and "))
	}
}

func checkConfig() error {
	if serverURL == "" {
		return fmt.Errorf("server URL is required (--url or YOUTRACK_URL)")
	}
	if token == "" {
		return fmt.Errorf("token is required (--token or YOUTRACK_TOKEN)")
	}
	return nil
}

// loadConfig applies config file values for unset flags, then env vars.
// Priority: flag > env var > config file.
func loadConfig() {
	cfg := readConfigFile()

	if serverURL == "" {
		if v := os.Getenv("YOUTRACK_URL"); v != "" {
			serverURL = v
		} else if v := cfg["url"]; v != "" {
			serverURL = v
		}
	}

	if token == "" {
		if v := os.Getenv("YOUTRACK_TOKEN"); v != "" {
			token = v
		} else if v := cfg["token"]; v != "" {
			token = v
		}
	}
}

func readConfigFile() map[string]string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	path := filepath.Join(home, ".config", "youtrack", "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var cfg map[string]string
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil
	}
	return cfg
}
