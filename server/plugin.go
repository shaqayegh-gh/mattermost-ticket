package main

import (
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/pkg/errors"
)

type Plugin struct {
	plugin.MattermostPlugin
}

// OnActivate is called when the plugin is activated
func (p *Plugin) OnActivate() error {
	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          "ticket",
		DisplayName:      "Create Ticket",
		Description:      "Create a new ticket",
		AutoComplete:     true,
		AutoCompleteDesc: "Create a new ticket",
		AutoCompleteHint: "",
	}); err != nil {
		return errors.Wrap(err, "failed to register command")
	}

	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          "resolve",
		DisplayName:      "Resolve Ticket",
		Description:      "Mark a ticket as resolved",
		AutoComplete:     true,
		AutoCompleteDesc: "Mark a ticket as resolved",
		AutoCompleteHint: "<post_id>",
	}); err != nil {
		return errors.Wrap(err, "failed to register resolve command")
	}

	return nil
}

// ExecuteCommand handles the /ticket and /resolve slash commands
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	if strings.HasPrefix(args.Command, "/ticket") {
		return p.handleTicketCommand(c, args)
	} else if strings.HasPrefix(args.Command, "/resolve") {
		return p.handleResolveCommand(c, args)
	}

	return &model.CommandResponse{}, nil
}
