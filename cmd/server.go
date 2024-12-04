package main

import (
	"fmt"

	"github.com/Vergangenheit/kafka-scratch/server"
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	host string
	port string
)

func init() {
	// Add flags specific to the server command
	serverStartCmd.Flags().StringVar(&port, "port", "6379", "port to run server from")
	serverStartCmd.Flags().StringVar(&host, "host", "localhost", "host address to run server from")

	// Bind flags to Viper
	viper.BindPFlag("port", serverStartCmd.Flags().Lookup("port"))
	viper.BindPFlag("host", serverStartCmd.Flags().Lookup("host"))
}

var serverStartCmd = &cobra.Command{
	Use:   "server-start",
	Short: "Start the server",
	Long:  `Start the server with the specified configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve values from Viper
		port := viper.GetString("port")
		host := viper.GetString("host")

		fmt.Printf("Starting server on address %s...\n", host)
		fmt.Printf("Starting server on port %s...\n", port)

		// set up logger
		logger := hclog.New(&hclog.LoggerOptions{
			Name:  "redis-server",
			Level: hclog.LevelFromString("INFO"),
		})
		server := server.NewServer(host, port, logger)
		err := server.Start()
		if err != nil {
			logger.Error("Failed to run server", "error", err)
			return
		}
	},
}
