package main

import (
	"time"

	"gitlab.justlab.xyz/alertflow-public/runner/pkg/executions"
	"gitlab.justlab.xyz/alertflow-public/runner/pkg/models"

	log "github.com/sirupsen/logrus"
)

type LogPlugin struct{}

func (p *LogPlugin) Init() models.Plugin {
	return models.Plugin{
		Name:    "Log",
		Type:    "action",
		Version: "1.0.1",
		Creator: "JustNZ",
	}
}

func (p *LogPlugin) Details() models.ActionDetails {
	return models.ActionDetails{
		ID:          "log",
		Name:        "Log Message",
		Description: "Prints a Log Message on Runner stdout",
		Icon:        "solar:clipboard-list-broken",
		Type:        "log",
		Category:    "Utility",
		Function:    p.Execute,
		Params:      nil,
	}
}

func (p *LogPlugin) Execute(execution models.Execution, flow models.Flows, payload models.Payload, steps []models.ExecutionSteps, step models.ExecutionSteps, action models.Actions) (data map[string]interface{}, finished bool, canceled bool, no_pattern_match bool, failed bool) {
	log.WithFields(log.Fields{
		"Execution": execution.ID,
		"StepID":    step.ID,
	}).Info("Log Action triggered")

	err := executions.UpdateStep(execution.ID.String(), models.ExecutionSteps{
		ID:             step.ID,
		ActionID:       action.ID.String(),
		ActionMessages: []string{"Log Action finished"},
		Pending:        false,
		Finished:       true,
		StartedAt:      time.Now(),
		FinishedAt:     time.Now(),
	})
	if err != nil {
		return nil, false, false, false, true
	}
	return nil, true, false, false, false
}

// Exported symbol for the plugin
var Plugin LogPlugin
