package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/couchbase/go-blip"
	"encoding/json"
	"log"
)

// https://github.com/couchbase/couchbase-lite-core/blob/master/modules/docs/pages/replication-protocol.adoc#changes
var changesCmd = &cobra.Command{
	Use:   "changes [sync gateway url]",
	Short: "Send changes",
	Long:  `Send changes.`,
	Run:   changesRun,
	Args:  cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.AddCommand(changesCmd)
}

func changesRun(cmd *cobra.Command, args []string) {

	serverUrl := args[0]

	sgBlipContext, err := NewSgBlipContext(serverUrl)
	if err != nil {
		panic(fmt.Sprintf("Error creeating sgblip context: %v", err))
	}
	defer sgBlipContext.BlipSender.Close()

	// TODO: this should be passed via the command line, either as raw data from stdio, or via file
	changesStr := `
		[
			[99, "docID", "revID", false],
			[100, "docID2", "revID2", false]
		]`

	changes := []interface{}{}
	err = json.Unmarshal([]byte(changesStr), &changes)
	if err != nil {
		panic(fmt.Sprintf("Error unmarshalling changes: %v", err))
	}

	changesRequest := blip.NewRequest()
	changesRequest.SetProfile("changes")
	changesRequest.SetJSONBody(changes)
	changesRequest.SetCompressed(false)

	sent := sgBlipContext.BlipSender.Send(changesRequest)
	if !sent {
			panic(fmt.Sprintf("Unable to subscribe to changes."))
	}

	// read the response
	response := changesRequest.Response()
	log.Printf("response properties: %+v", response.Properties)
	body, err := response.Body()
	if err != nil {
		panic(fmt.Sprintf("Error reading response body: %v", err))
	}
	log.Printf("response body: %s", body)


}
