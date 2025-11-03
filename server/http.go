package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

// ServeHTTP handles HTTP requests
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.API.LogInfo("ServeHTTP request received", "path", r.URL.Path, "method", r.Method, "raw_query", r.URL.RawQuery, "host", r.Host, "remote_addr", r.RemoteAddr)

	if r.URL.Path == "/api/v1/dialog" {
		p.handleDialogSubmit(w, r)
		return
	}

	if r.URL.Path == "/api/v1/runresolve" || strings.Contains(r.URL.Path, "/runresolve") {
		p.API.LogInfo("Button integration path matched", "path", r.URL.Path)
		p.handleRunResolve(w, r)
		return
	}

	if r.URL.Path == "/api/v1/runreopen" || strings.Contains(r.URL.Path, "/runreopen") {
		p.API.LogInfo("Reopen button integration path matched", "path", r.URL.Path)
		p.handleRunReopen(w, r)
		return
	}

	if !strings.Contains(r.URL.Path, "/api/v1/") {
		p.API.LogWarn("Unhandled path in ServeHTTP", "path", r.URL.Path)
	}
	http.NotFound(w, r)
}

// handleDialogSubmit processes the ticket creation dialog submission
func (p *Plugin) handleDialogSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request model.SubmitDialogRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var ticketData TicketDialog
	if teamNameVal, ok := request.Submission["team_name"].(string); ok {
		ticketData.TeamName = teamNameVal
	}
	if projectNameVal, ok := request.Submission["project_name"].(string); ok {
		ticketData.ProjectName = projectNameVal
	}
	if environmentVal, ok := request.Submission["environment"].(string); ok {
		ticketData.Environment = environmentVal
	}
	if priorityVal, ok := request.Submission["priority"].(string); ok {
		ticketData.Priority = priorityVal
	}
	if descriptionVal, ok := request.Submission["description"].(string); ok {
		ticketData.Description = descriptionVal
	}

	if ticketData.TeamName == "" || ticketData.ProjectName == "" || ticketData.Environment == "" || ticketData.Description == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status": "OK"}`)); err != nil {
			p.API.LogError("failed to write response body", "error", err.Error())
		}
		return
	}

	// Validate that the dialog was submitted from an allowed channel
	if !p.validateChannel(request.ChannelId) {
		allowed := p.getAllowedChannels()
		where := "allowed channels"
		if len(allowed) > 0 {
			where = "the following channels: `" + strings.Join(allowed, "`, `") + "`"
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if _, err := fmt.Fprintf(w, `{"error": "Tickets can only be created in %s"}`, where); err != nil {
			p.API.LogError("failed to write error response body", "error", err.Error())
		}
		return
	}

	// Create ticket
	if err := p.createTicket(ticketData, request.ChannelId, request.UserId); err != nil {
		http.Error(w, "Failed to create ticket", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{"status": "OK"}`)); err != nil {
		p.API.LogError("failed to write response body", "error", err.Error())
	}
}

// handleRunResolve executes /resolve command when button is clicked
func (p *Plugin) handleRunResolve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		p.API.LogError("Invalid method for runresolve", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.PostActionIntegrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		p.API.LogError("Failed to decode run resolve request", "error", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	postID, ok := req.Context["post_id"].(string)
	if !ok || postID == "" {
		p.API.LogError("Missing or invalid post_id in context")
		http.Error(w, "Missing post_id", http.StatusBadRequest)
		return
	}

	channelID, ok := req.Context["channel_id"].(string)
	if !ok || channelID == "" {
		p.API.LogError("Missing or invalid channel_id in context")
		http.Error(w, "Missing channel_id", http.StatusBadRequest)
		return
	}

	args := &model.CommandArgs{
		UserId:    req.UserId,
		ChannelId: channelID,
		Command:   fmt.Sprintf("/resolve %s", postID),
	}

	resp, _ := p.handleResolveCommand(nil, args)

	integrationResp := &model.PostActionIntegrationResponse{
		EphemeralText: resp.Text,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(integrationResp); err != nil {
		p.API.LogError("failed to encode integration response", "error", err.Error())
	}
}

// handleRunReopen executes reopen action when button is clicked
func (p *Plugin) handleRunReopen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		p.API.LogError("Invalid method for runreopen", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.PostActionIntegrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		p.API.LogError("Failed to decode run reopen request", "error", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	postID, ok := req.Context["post_id"].(string)
	if !ok || postID == "" {
		p.API.LogError("Missing or invalid post_id in context")
		http.Error(w, "Missing post_id", http.StatusBadRequest)
		return
	}

	channelID, ok := req.Context["channel_id"].(string)
	if !ok || channelID == "" {
		p.API.LogError("Missing or invalid channel_id in context")
		http.Error(w, "Missing channel_id", http.StatusBadRequest)
		return
	}

	post, err := p.API.GetPost(postID)
	if err != nil {
		p.API.LogError("Failed to get post for reopen", "error", err.Error())
		http.Error(w, "Failed to get post", http.StatusInternalServerError)
		return
	}

	if !strings.Contains(post.Message, "✅ **Status:** Resolved") {
		integrationResp := &model.PostActionIntegrationResponse{
			EphemeralText: "This ticket is not resolved.",
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(integrationResp); err != nil {
			p.API.LogError("failed to encode integration response", "error", err.Error())
		}
		return
	}

	// Update post to open status
	updatePost := post.Clone()
	updatePost.Message = strings.Replace(updatePost.Message, "✅ **Status:** Resolved", "**Status:** Open", 1)
	p.attachResolveButton(updatePost, postID, channelID)

	if _, err := p.API.UpdatePost(updatePost); err != nil {
		p.API.LogError("Failed to update post for reopen", "error", err.Error())
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	replyPost := &model.Post{
		ChannelId: channelID,
		UserId:    req.UserId,
		Message:   "🔄 Reopened",
		Type:      model.PostTypeDefault,
		RootId:    postID,
	}

	if _, err := p.API.CreatePost(replyPost); err != nil {
		p.API.LogError("Failed to create reopen reply", "error", err.Error())
	}

	integrationResp := &model.PostActionIntegrationResponse{
		EphemeralText: "Ticket reopened successfully!",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(integrationResp); err != nil {
		p.API.LogError("failed to encode integration response", "error", err.Error())
	}
}
