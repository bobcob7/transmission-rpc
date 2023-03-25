/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/bobcob7/transmission-rpc"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Start transmission interface
	var err error
	tr, err = transmission.New(context.Background(), address)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to server:", err)
		os.Exit(1)
	}
	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var address string
var tr *transmission.Client

func init() {
	rootCmd.PersistentFlags().StringVar(&address, "base-url", "https://transmission.bobcob7.com", "URL to transmission server")
}
