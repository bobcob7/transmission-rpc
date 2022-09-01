/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// sessionStatsCmd represents the sessionStats command
var sessionStatsCmd = &cobra.Command{
	Use:   "session-stats",
	Short: "Get session statistics from transmission server",
	Run: func(cmd *cobra.Command, args []string) {
		session, err := tr.GetSession(cmd.Context())
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to get session:", err)
			os.Exit(1)
		}
		if err := json.NewEncoder(os.Stdout).Encode(session); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to marshal output:", err)
			os.Exit(1)
		}
	},
}

func init() {
	getCmd.AddCommand(sessionStatsCmd)
}
