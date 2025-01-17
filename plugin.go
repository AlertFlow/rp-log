package main

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/AlertFlow/runner/pkg/executions"
	"github.com/AlertFlow/runner/pkg/models"
	"github.com/tidwall/gjson"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type LogPlugin struct{}

func (p *LogPlugin) Init() models.Plugin {
	return models.Plugin{
		Name:    "Log",
		Type:    "action",
		Version: "1.0.6",
		Creator: "JustNZ",
	}
}

func (p *LogPlugin) Details() models.PluginDetails {
	params := []models.Param{
		{
			Key:         "AdditionalMessage",
			Type:        "text",
			Default:     "",
			Required:    false,
			Description: "Additional message to log. To use payload data, use 'payload.<key>'",
		},
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		log.Error(err)
	}

	return models.PluginDetails{
		Action: models.ActionDetails{
			ID:          "log",
			Name:        "Log Message",
			Description: "Prints a Log Message on Runner stdout",
			Icon:        "solar:clipboard-list-broken",
			Type:        "log",
			Category:    "Utility",
			Function:    p.Execute,
			Params:      json.RawMessage(paramsJSON),
		},
	}
}

func (p *LogPlugin) Execute(execution models.Execution, flow models.Flows, payload models.Payload, steps []models.ExecutionSteps, step models.ExecutionSteps, action models.Actions) (data map[string]interface{}, finished bool, canceled bool, no_pattern_match bool, failed bool) {
	additionalMessage := ""

	for _, param := range action.Params {
		if param.Key == "AdditionalMessage" {
			additionalMessage = param.Value
		}
	}

	if strings.Contains(additionalMessage, "payload.") {
		// convert payload to string
		payloadBytes, err := json.Marshal(payload.Payload)
		if err != nil {
			log.Error("Error converting payload to JSON:", err)
			return nil, false, false, false, true
		}
		payloadString := string(payloadBytes)

		additionalMessage = gjson.Get(payloadString, strings.Replace(additionalMessage, "payload.", "", 1)).String()
	}

	log.WithFields(log.Fields{
		"Execution":          execution.ID,
		"StepID":             step.ID,
		"Additional Message": additionalMessage,
	}).Info("Log Action triggered")

	err := executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
		ID:       step.ID,
		ActionID: action.ID.String(),
		ActionMessages: []string{
			additionalMessage,
			"Log Action finished",
		},
		Pending:    false,
		Finished:   true,
		StartedAt:  time.Now(),
		FinishedAt: time.Now(),
	})
	if err != nil {
		return nil, false, false, false, true
	}
	return nil, true, false, false, false
}

func (p *LogPlugin) Handle(context *gin.Context) {}

// Exported symbol for the plugin
var Plugin LogPlugin
