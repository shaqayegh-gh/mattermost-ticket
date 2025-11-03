package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

// handleTicketCommand handles the /ticket slash command
func (p *Plugin) handleTicketCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	if !p.validateChannel(args.ChannelId) {
		allowed := p.getAllowedChannels()
		var where string
		if len(allowed) > 0 {
			where = "the following channels: `" + strings.Join(allowed, "`, `") + "`"
		} else {
			where = "allowed channels"
		}
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("❌ This command can only be used in %s.", where),
		}, nil
	}

	dialog := model.OpenDialogRequest{
		TriggerId: args.TriggerId,
		URL:       fmt.Sprintf("/plugins/%s/api/v1/dialog", pluginID),
		Dialog: model.Dialog{
			Title:            "Create New Ticket",
			IntroductionText: "Please fill in the details for your ticket:",
			Elements: []model.DialogElement{
				{
					DisplayName: "Team Name",
					Name:        "team_name",
					Type:        "select",
					Placeholder: "Select your team",
					Options:     p.getTeamOptions(),
				},
				{
					DisplayName: "Project Name",
					Name:        "project_name",
					Type:        "select",
					Placeholder: "Select issue project",
					Options:     p.getProjectOptions(),
				},
				{
					DisplayName: "Environment",
					Name:        "environment",
					Type:        "select",
					Placeholder: "Select environment",
					Options:     environmentOptions,
				},
				{
					DisplayName: "Message Priority",
					Name:        "priority",
					Type:        "select",
					Placeholder: "Select priority",
					Options:     priorityOptions,
					Default:     "standard",
				},
				{
					DisplayName: "Issue Description",
					Name:        "description",
					Type:        "textarea",
					Placeholder: "Describe the issue in detail...",
					MaxLength:   2000,
				},
			},
			SubmitLabel:    "Create Ticket",
			NotifyOnCancel: true,
		},
	}

	if err := p.API.OpenInteractiveDialog(dialog); err != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Failed to open ticket dialog: " + err.Error(),
		}, nil
	}

	return &model.CommandResponse{}, nil
}

// handleResolveCommand handles the /resolve slash command
func (p *Plugin) handleResolveCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	if !p.validateChannel(args.ChannelId) {
		allowed := p.getAllowedChannels()
		var where string
		if len(allowed) > 0 {
			where = "the following channels: `" + strings.Join(allowed, "`, `") + "`"
		} else {
			where = "allowed channels"
		}
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("❌ This command can only be used in %s.", where),
		}, nil
	}

	parts := strings.Fields(args.Command)
	if len(parts) < 2 {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Usage: /resolve <post_id>\nExample: /resolve abc123def456",
		}, nil
	}

	postId := parts[1]

	post, err := p.API.GetPost(postId)
	if err != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Failed to find post: " + err.Error(),
		}, nil
	}

	if !strings.Contains(post.Message, "🎫 **New Ticket Created**") {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "This post is not a ticket. Please use the post ID of a ticket.",
		}, nil
	}

	if strings.Contains(post.Message, "✅ **Status:** Resolved") {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "This ticket is already resolved.",
		}, nil
	}

	// Update post to resolved status
	updatePost := post.Clone()
	updatePost.Message = strings.Replace(updatePost.Message, "**Status:** Open", "✅ **Status:** Resolved", 1)
	updatePost.Message = strings.Replace(updatePost.Message, "💡 **To mark as resolved:** Use `/resolve "+postId+"`", "", 1)

	// Add reopen button
	p.attachReopenButton(updatePost, postId, post.ChannelId)

	if _, err := p.API.UpdatePost(updatePost); err != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Failed to update post: " + err.Error(),
		}, nil
	}

	replyPost := &model.Post{
		ChannelId: post.ChannelId,
		UserId:    args.UserId,
		Message:   "✅ Resolved",
		Type:      model.PostTypeDefault,
		RootId:    postId,
	}

	if _, err := p.API.CreatePost(replyPost); err != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Failed to create reply: " + err.Error(),
		}, nil
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         "Ticket marked as resolved successfully!",
	}, nil
}
