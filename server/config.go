package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
)

// getTicketMentionUsers returns the users that should be mentioned in ticket
func (p *Plugin) getTicketMentionUsers(teamName string, channelId string) []string {
	var teamMembers []string

	config := p.API.GetConfig()
	if config != nil && config.PluginSettings.Plugins[pluginID] != nil {
		teamMembersConfig := config.PluginSettings.Plugins[pluginID]["ticketmentionconfig"]

		if teamMembersConfig != nil {
			teamMembersStr := teamMembersConfig.(string)
			if teamMembersStr != "" {
				var teamMembersMap map[string][]string
				if err := json.Unmarshal([]byte(teamMembersStr), &teamMembersMap); err == nil {
					// Check one channel teamName
					channelName := p.getChannelName(channelId)
					configName := fmt.Sprintf("%s__%s", channelName, teamName)
					if members, exists := teamMembersMap[configName]; exists {
						teamMembers = append(teamMembers, members...)
					}

					// Always add channelId members if they exist
					if allMembers, exists := teamMembersMap[channelName]; exists {
						teamMembers = append(teamMembers, allMembers...)
					}
				} else {
					p.API.LogError("Failed to parse team members config", "error", err.Error(), "rawConfig", teamMembersStr, "configBytes", []byte(teamMembersStr))
				}
			}
		}
	}

	return teamMembers
}

// getTeamOptions returns the team options from configuration or falls back to defaults
func (p *Plugin) getTeamOptions() []*model.PostActionOptions {
	config := p.API.GetConfig()
	if config != nil && config.PluginSettings.Plugins[pluginID] != nil {
		teamOptionsConfig := config.PluginSettings.Plugins[pluginID]["teamoptionsconfig"]

		if teamOptionsConfig != nil {
			raw := teamOptionsConfig.(string)
			if raw != "" {
				type option struct {
					Text  string `json:"Text"`
					Value string `json:"Value"`
				}
				var parsed []option
				if err := json.Unmarshal([]byte(raw), &parsed); err == nil {
					var opts []*model.PostActionOptions
					for _, o := range parsed {
						if o.Text == "" || o.Value == "" {
							continue
						}
						opts = append(opts, &model.PostActionOptions{Text: o.Text, Value: o.Value})
					}
					if len(opts) > 0 {
						return opts
					}
				} else {
					p.API.LogError("Failed to parse team options config", "error", err.Error(), "rawConfig", raw)
				}
			}
		}
	}

	return teamOptions
}

// getProjectOptions returns the project options from configuration or falls back to defaults
func (p *Plugin) getProjectOptions() []*model.PostActionOptions {
	config := p.API.GetConfig()
	if config != nil && config.PluginSettings.Plugins[pluginID] != nil {
		projectOptionsConfig := config.PluginSettings.Plugins[pluginID]["projectoptionsconfig"]

		if projectOptionsConfig != nil {
			raw := projectOptionsConfig.(string)
			if raw != "" {
				type option struct {
					Text  string `json:"Text"`
					Value string `json:"Value"`
				}
				var parsed []option
				if err := json.Unmarshal([]byte(raw), &parsed); err == nil {
					var opts []*model.PostActionOptions
					for _, o := range parsed {
						if o.Text == "" || o.Value == "" {
							continue
						}
						opts = append(opts, &model.PostActionOptions{Text: o.Text, Value: o.Value})
					}
					if len(opts) > 0 {
						return opts
					}
				} else {
					p.API.LogError("Failed to parse project options config", "error", err.Error(), "rawConfig", raw)
				}
			}
		}
	}

	return projectOptions
}

// getAllowedChannels returns the list of channel names that are allowed to use the plugin.
// The list is parsed from the `DefaultChannel` plugin setting which now accepts
// a comma-separated or newline-separated list of channel names. If empty, all
// channels are allowed.
func (p *Plugin) getAllowedChannels() []string {
	config := p.API.GetConfig()
	if config == nil || config.PluginSettings.Plugins[pluginID] == nil {
		return nil
	}

	rawSetting := config.PluginSettings.Plugins[pluginID]["defaultchannel"]
	raw, ok := rawSetting.(string)
	if !ok || strings.TrimSpace(raw) == "" {
		return nil
	}

	normalized := strings.ReplaceAll(raw, "\n", ",")
	parts := strings.Split(normalized, ",")
	var channels []string
	for _, part := range parts {
		name := strings.TrimSpace(part)
		if name == "" {
			continue
		}
		channels = append(channels, name)
	}
	if len(channels) == 0 {
		return nil
	}
	return channels
}

// validateChannel checks if the given channel ID matches the default channel
func (p *Plugin) validateChannel(channelId string) bool {
	allowed := p.getAllowedChannels()
	if len(allowed) == 0 {
		return true
	}

	channelName := p.getChannelName(channelId)
	for _, name := range allowed {
		if strings.EqualFold(channelName, name) {
			return true
		}
	}
	return false
}

// Get channel name from it's ID
func (p *Plugin) getChannelName(channelId string) string {
	channel, err := p.API.GetChannel(channelId)
	if err != nil {
		p.API.LogError("Failed to get channel for validation", "error", err.Error(), "channel_id", channelId)
		return ""
	}
	return channel.Name
}
