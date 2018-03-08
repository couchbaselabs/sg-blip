package cmd

import (
	"fmt"
	"net/url"
	"github.com/couchbase/go-blip"
	"log"
	"golang.org/x/net/websocket"
)

type SgBlipContext struct {
	BlipContext *blip.Context
	BlipSender *blip.Sender
}

func NewSgBlipContext(serverUrl string) (*SgBlipContext, error) {

	// Construct URL to connect to blipsync target endpoint
	destUrl := fmt.Sprintf("%s/_blipsync", serverUrl)
	u, err := url.Parse(destUrl)
	if err != nil {
		return nil, fmt.Errorf("Error parsing url: %v", destUrl)
	}
	u.Scheme = "ws"

	// Make BLIP/Websocket connection
	blipContext := blip.NewContext(BlipCBMobileReplication)
	blipContext.Logger = func(eventType blip.LogEventType, fmt string, params ...interface{}) {
		log.Printf(fmt, params...)
	}

	origin := "http://localhost" // TODO: what should be used here?

	config, err := websocket.NewConfig(u.String(), origin)
	if err != nil {
		return nil, fmt.Errorf("Error creating websocket config.  Error: %v", err)
	}

	sender, err := blipContext.DialConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to Sync Gateway.  Error: %v", err)
	}

	return &SgBlipContext{
		BlipContext: blipContext,
		BlipSender: sender,
	}, nil

}