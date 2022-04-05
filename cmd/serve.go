/*
Package cmd
Copyright © 2022 Arthur Isac isacartur@gmail.com

*/
package cmd

import (
	"github.com/spf13/cobra"
)

var port string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs a memstore server",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")
	serveCmd.LocalFlags().StringVarP(&port, "port", "p", "8080", "Defines port for memstore server. Default value is 8080")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
