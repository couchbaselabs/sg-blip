package cmd

import (
	"fmt"

	"encoding/json"
	"log"

	"github.com/couchbase/go-blip"
	"github.com/spf13/cobra"
)

// subChangesCmd represents the subChanges command
var subChangesCmd = &cobra.Command{
	Use:   "subChanges [sync gateway url]",
	Short: "Subscribe to changes",
	Long:  `This will print out any changes that are received from Sync Gateway.  For example, when docs are added, updated, or deleted.`,
	Run:   subChangesRun,
	Args:  cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.AddCommand(subChangesCmd)
}

func subChangesRun(cmd *cobra.Command, args []string) {

	fmt.Printf("subChanges called.... args %v\n", args)

	continuous := true // TODO: set via flag

	serverUrl := args[0]

	sgBlipContext, err := NewSgBlipContext(serverUrl)
	if err != nil {
		panic(fmt.Sprintf("Error creeating sgblip context: %v", err))
	}
	defer sgBlipContext.BlipSender.Close()

	// When this test sends subChanges, Sync Gateway will send a changes request that must be handled
	sgBlipContext.BlipContext.HandlerForProfile["changes"] = func(request *blip.Message) {

		requestBody, err := request.Body()
		if err != nil {
			panic(fmt.Errorf("Error getting request body.  Error: %v", err))

		}
		log.Printf("Got change: %s", requestBody)

		if !request.NoReply() {

			// Send an empty response to avoid the Sync: Invalid response to 'changes' message
			response := request.Response()
			emptyResponseVal := []interface{}{}
			emptyResponseValBytes, err := json.Marshal(emptyResponseVal)
			if err != nil {
				panic(fmt.Sprintf("Error marshalling response: %v", err))
			}
			response.SetBody(emptyResponseValBytes)
		}

		// request.Sender.CloseAbruptly()
		request.Sender.Close()

	}

	// Send subChanges to subscribe to changes, which will cause the "changes" profile handler above to be called back
	subChangesRequest := blip.NewRequest()
	subChangesRequest.SetProfile("subChanges")
	switch continuous {
	case true:
		subChangesRequest.Properties["continuous"] = "true"
	default:
		subChangesRequest.Properties["continuous"] = "false"
	}
	subChangesRequest.Properties["foo"] = "bar"
	subChangesRequest.SetCompressed(false)

	sent := sgBlipContext.BlipSender.Send(subChangesRequest)
	if !sent {
		panic(fmt.Sprintf("Unable to subscribe to changes."))
	}
	subChangesResponse := subChangesRequest.Response()
	if subChangesResponse.SerialNumber() != subChangesRequest.SerialNumber() {
		panic(fmt.Sprintf("subChangesResponse.SerialNumber() != subChangesRequest.SerialNumber().  %v != %v", subChangesResponse.SerialNumber(), subChangesRequest.SerialNumber()))
	}

	// Block until the user cancels
	select {}

}
