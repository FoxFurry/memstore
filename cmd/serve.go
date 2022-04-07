/*
Package cmd

*/
package cmd

import (
	"context"
	"fmt"
	"github.com/FoxFurry/memstore/internal/api/server"
	"github.com/spf13/cobra"
)

var port string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs a memstore server",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if port == "" {
			fmt.Println("Cannot use empty port")
			return
		}

		ctx, cancel := context.WithCancel(context.Background())

		storeServer := server.New(ctx)

		if err := storeServer.Start(port); err != nil {
			fmt.Printf("Unexpected error while running server: %v", err)
		}
		cancel()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().StringVarP(&port, "port", "p", "8080", "Defines port for memstore server. Default value is 8080")
}
