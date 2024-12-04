package main

import (
	"fmt"

	"github.com/Vergangenheit/kafka-scratch/client"
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	hostPort string
	request  string
)

func init() {
	// Add flags specific to the server command
	sendRequest.Flags().StringVar(&hostPort, "hostport", "localhost:9092", "host port of the server")
	sendRequest.Flags().StringVar(&request, "request", "00000023001200046f7fc66100096b61666b612d636c69000a6b61666b612d636c6904302e3100", "request in string format")

	viper.BindPFlag("hostport", sendRequest.Flags().Lookup("hostport"))
	viper.BindPFlag("request", sendRequest.Flags().Lookup("request"))

}

var sendRequest = &cobra.Command{
	Use:   "send",
	Short: "Send generic request to the server",
	Long:  `send a request to the server with a string`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := hclog.New(&hclog.LoggerOptions{
			Name:  "client",
			Level: hclog.LevelFromString("INFO"),
		})
		// Retrieve values from Viper
		hostPort := viper.GetString("hostport")
		request := viper.GetString("request")

		logger.Info(fmt.Sprintf("Sending request to server on %s...\n", hostPort))

		// send command
		cl, err := client.NewKafkaClient(hostPort)
		if err != nil {
			logger.Error("error instantiating client", err)
			return
		}
		defer cl.Close()

		req := []byte(request)
		err = cl.Send(req)
		if err != nil {
			logger.Error("error sending request", err)
			return
		}

	},
}
