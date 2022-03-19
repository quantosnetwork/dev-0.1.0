/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"quantos/0.1.0/blockchain"

	"github.com/spf13/cobra"
)

// testsCmd represents the tests command
var testsCmd = &cobra.Command{
	Use:   "tests",
	Short: "call various Quantos tests",
	Long:  `calling various Quantos tests utilities`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tests called")
		bc := blockchain.GetBlockProxy()
		bc.NewBlock()

	},
}

func init() {
	rootCmd.AddCommand(testsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
