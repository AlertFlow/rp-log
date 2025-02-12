package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlertFlow/runner/pkg/models"
	"github.com/AlertFlow/runner/pkg/protocol"
)

func main() {
	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		os.Exit(0)
	}()

	// Process requests
	for {
		var req protocol.Request
		if err := decoder.Decode(&req); err != nil {
			os.Exit(1)
		}

		// Handle the request
		resp := handle(req)

		if err := encoder.Encode(resp); err != nil {
			os.Exit(1)
		}
	}
}

func Details() models.Plugin {
	return models.Plugin{
		Name:    "Log",
		Type:    "action",
		Version: "1.0.7",
		Author:  "JustNZ",
		Action: models.ActionDetails{
			ID:          "log",
			Name:        "Log Message",
			Description: "Prints a Log Message on Runner stdout",
			Icon:        "solar:clipboard-list-broken",
			Category:    "Utility",
			Params: []models.Param{
				{
					Key:         "AdditionalMessage",
					Type:        "text",
					Default:     "",
					Required:    false,
					Description: "Additional message to log. To use payload data, use 'payload.<key>'",
				},
			},
		},
	}
}

func handle(req protocol.Request) protocol.Response {
	// Plugin-specific logic here
	return protocol.Response{
		Success: true,
		Data: map[string]interface{}{
			"result": "processed",
		},
	}
}
