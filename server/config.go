package main

import (
	"encoding/json"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
)

// getTeamMembers returns the members of a specific team from configuration
func (p *Plugin) getTeamMembers(teamName string) []string {
	var teamMembers []string

	config := p.API.GetConfig()
	if config != nil && config.PluginSettings.Plugins[pluginID] != nil {
		teamMembersConfig := config.PluginSettings.Plugins[pluginID]["teammembersconfig"]

		if teamMembersConfig != nil {
			teamMembersStr := teamMembersConfig.(string)
			if teamMembersStr != "" {
				var teamMembersMap map[string][]string
				if err := json.Unmarshal([]byte(teamMembersStr), &teamMembersMap); err == nil {
					if members, exists := teamMembersMap[teamName]; exists {
						teamMembers = append(teamMembers, members...)
					}

					// Always add "all" members if they exist
					if allMembers, exists := teamMembersMap["all"]; exists {
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

// getDefaultChannel returns the default channel name from configuration
func (p *Plugin) getDefaultChannel() string {
	config := p.API.GetConfig()
	if config != nil && config.PluginSettings.Plugins[pluginID] != nil {
		if defaultChannelConfig := config.PluginSettings.Plugins[pluginID]["defaultchannel"]; defaultChannelConfig != nil {
			if defaultChannel, ok := defaultChannelConfig.(string); ok && defaultChannel != "" {
				return defaultChannel
			}
		}
	}
	return ""
}

// validateChannel checks if the given channel ID matches the default channel
func (p *Plugin) validateChannel(channelId string) bool {
	defaultChannelName := p.getDefaultChannel()
	if defaultChannelName == "" {
		return true
	}

	channel, err := p.API.GetChannel(channelId)
	if err != nil {
		p.API.LogError("Failed to get channel for validation", "error", err.Error(), "channel_id", channelId)
		return false
	}

	return strings.EqualFold(channel.Name, defaultChannelName)
}
