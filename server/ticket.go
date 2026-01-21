package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
)

// buildPostPriority constructs a PostPriority object from a plain string value.
func buildPostPriority(priority string) *model.PostPriority {
	p := priority
	requestedAck := false
	persistent := false
	return &model.PostPriority{
		Priority:                &p,
		RequestedAck:            &requestedAck,
		PersistentNotifications: &persistent,
	}
}

// createTicket creates a new ticket post with the provided data
func (p *Plugin) createTicket(ticketData TicketDialog, channelId, userId string) error {
	ticketMentions := p.getTicketMentionUsers(ticketData.TeamName, channelId)

	priority := ticketData.Priority
	if priority == "" {
		priority = "standard"
	}

	// Create ticket post
	ticketPost := &model.Post{
		ChannelId: channelId,
		UserId:    userId,
		Message: fmt.Sprintf("🎫 **New Ticket Created**\n\n"+
			"**Ticket Details:**\n\n"+
			"• Team: **%s**\n"+
			"• Project: **%s**\n"+
			"• Environment: **%s**\n\n\n"+
			"**Status:** Open\n\n"+
			"💡 **To mark as resolved:** Use `/resolve %s`",
			ticketData.TeamName,
			ticketData.ProjectName,
			ticketData.Environment,
			"placeholder"),
		Type: model.PostTypeDefault,
	}

	ticketPost.Metadata = &model.PostMetadata{Priority: buildPostPriority(priority)}

	for _, member := range ticketMentions {
		ticketPost.Message += fmt.Sprintf(" @%s", member)
	}

	// Create the post
	firstPost, err := p.API.CreatePost(ticketPost)
	if err != nil {
		p.API.LogError("Failed to create ticket post", "error", err.Error())
		return err
	}

	updatePost := firstPost.Clone()
	updatePost.Message = strings.Replace(updatePost.Message, "placeholder", firstPost.Id, 1)
	if updatePost.Metadata == nil {
		updatePost.Metadata = &model.PostMetadata{}
	}
	updatePost.Metadata.Priority = buildPostPriority(priority)
	p.attachResolveButton(updatePost, firstPost.Id, firstPost.ChannelId)

	if _, err := p.API.UpdatePost(updatePost); err != nil {
		p.API.LogError("Failed to update post with button", "error", err.Error())
		return err
	}

	descriptionPost := &model.Post{
		ChannelId: channelId,
		UserId:    userId,
		Message:   ticketData.Description,
		Type:      model.PostTypeDefault,
		RootId:    firstPost.Id,
	}

	if _, err := p.API.CreatePost(descriptionPost); err != nil {
		p.API.LogError("Failed to create description post", "error", err.Error())
		return err
	}

	return nil
}

// attachResolveButton adds a resolve button to the post
func (p *Plugin) attachResolveButton(post *model.Post, postID, channelID string) {
	integrationURL := fmt.Sprintf("/plugins/%s/api/v1/runresolve", pluginID)

	attachment := &model.SlackAttachment{
		Text:     "Click below to resolve this ticket",
		Fallback: "Resolve Ticket",
		Color:    "#e40504",
		Actions: []*model.PostAction{
			{
				Id:   "runresolve",
				Type: model.PostActionTypeButton,
				Name: "Resolve Ticket",
				Integration: &model.PostActionIntegration{
					URL: integrationURL,
					Context: map[string]interface{}{
						"post_id":    postID,
						"channel_id": channelID,
					},
				},
			},
		},
	}

	if post.Props == nil {
		post.Props = make(model.StringInterface)
	}
	post.Props["attachments"] = []*model.SlackAttachment{attachment}
}

// attachReopenButton adds a reopen button to the post
func (p *Plugin) attachReopenButton(post *model.Post, postID, channelID string) {
	integrationURL := fmt.Sprintf("/plugins/%s/api/v1/runreopen", pluginID)

	attachment := &model.SlackAttachment{
		Text:     "Click below to reopen this ticket",
		Fallback: "Reopen Ticket",
		Color:    "#0ddb21",
		Actions: []*model.PostAction{
			{
				Id:   "runreopen",
				Type: model.PostActionTypeButton,
				Name: "Reopen Ticket",
				Integration: &model.PostActionIntegration{
					URL: integrationURL,
					Context: map[string]interface{}{
						"post_id":    postID,
						"channel_id": channelID,
					},
				},
			},
		},
	}

	if post.Props == nil {
		post.Props = make(model.StringInterface)
	}
	post.Props["attachments"] = []*model.SlackAttachment{attachment}
}
