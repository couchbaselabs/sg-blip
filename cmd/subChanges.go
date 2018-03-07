package cmd

import (
	"fmt"

	"encoding/json"
	"log"
	"net/url"

	"github.com/couchbase/go-blip"
	"github.com/spf13/cobra"
	"golang.org/x/net/websocket"
)

const (
	BlipCBMobileReplication = "CBMobile_2"
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

	// Construct URL to connect to blipsync target endpoint
	destUrl := fmt.Sprintf("%s/_blipsync", serverUrl)
	u, err := url.Parse(destUrl)
	if err != nil {
		panic(fmt.Errorf("Error parsing url: %v", destUrl))
	}
	u.Scheme = "ws"

	// Make BLIP/Websocket connection
	blipContext := blip.NewContext(BlipCBMobileReplication)
	blipContext.Logger = func(eventType blip.LogEventType, fmt string, params ...interface{}) {
		log.Printf(fmt, params...)
	}
	// blipContext.LogMessages = true
	// blipContext.LogFrames = true

	origin := "http://localhost" // TODO: what should be used here?

	config, err := websocket.NewConfig(u.String(), origin)
	if err != nil {
		panic(fmt.Errorf("Error creating websocket config.  Error: %v", err))

	}

	// TODO -- take username and param
	//if len(spec.connectingUsername) > 0 {
	//	config.Header = http.Header{
	//		"Authorization": {"Basic " + base64.StdEncoding.EncodeToString([]byte(spec.connectingUsername+":"+spec.connectingPassword))},
	//	}
	//}

	sender, err := blipContext.DialConfig(config)
	if err != nil {
		panic(fmt.Errorf("Error connecting to Sync Gateway.  Error: %v", err))
	}

	// When this test sends subChanges, Sync Gateway will send a changes request that must be handled
	blipContext.HandlerForProfile["changes"] = func(request *blip.Message) {

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

	sent := sender.Send(subChangesRequest)
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
