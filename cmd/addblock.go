/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/avirup-ghosal/VaultChain/core"
	"github.com/spf13/cobra"
)

var data string

// addblockCmd represents the addblock command
var addblockCmd = &cobra.Command{
	Use:   "addblock",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if data == "" {
			fmt.Println("Please provide transaction data using --data flag")
			os.Exit(1)
		}
		newblockchain := core.NewBlockchain()
		err := newblockchain.AddBlock(data)
		if err != nil {
			log.Fatalf("error adding block: %w", err)
		}

		fmt.Println("addblock called")
	},
}

func init() {
	rootCmd.AddCommand(addblockCmd)

	addblockCmd.Flags().StringVar(&data, "data", "", "Transaction data to be added to the block")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addblockCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addblockCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
