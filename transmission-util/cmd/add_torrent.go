/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// sessionStatsCmd represents the sessionStats command
var addTorrentsCmd = &cobra.Command{
	Use:   "torrent",
	Short: "Add torrent information from transmission server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("Missing magnet link")
		}
		magnetLink := args[0]
		id, err := tr.AddMagnetLink(cmd.Context(), magnetLink)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to add torrent:", err)
			os.Exit(1)
		}
		fmt.Println("Add torrent:", id)
		return nil
	},
}

func init() {
	addCmd.AddCommand(addTorrentsCmd)
}
