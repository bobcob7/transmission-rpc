/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// sessionStatsCmd represents the sessionStats command
var getTorrentsCmd = &cobra.Command{
	Use:   "torrents",
	Short: "Get torrent information from transmission server",
	RunE: func(cmd *cobra.Command, args []string) error {
		ids := make([]int, len(args))
		for i, arg := range args {
			id, err := strconv.ParseInt(arg, 10, 64)
			if err != nil {
				return fmt.Errorf("failed parsing int argument")
			}
			ids[i] = int(id)
		}
		torrents, err := tr.GetTorrents(cmd.Context(), ids...)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to get torrents:", err)
			os.Exit(1)
		}
		if err := json.NewEncoder(os.Stdout).Encode(torrents); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to marshal output:", err)
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	getCmd.AddCommand(getTorrentsCmd)
}
