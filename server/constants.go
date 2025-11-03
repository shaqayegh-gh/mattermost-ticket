package main

import "github.com/mattermost/mattermost/server/public/model"

const pluginID string = "com.github.mattermost-ticket-plugin"

// Default team options
var teamOptions = []*model.PostActionOptions{
	{Text: "Develop", Value: "develop"},
	{Text: "Design", Value: "design"},
	{Text: "Scrum Master", Value: "scrum-master"},
	{Text: "QA", Value: "qa"},
	{Text: "Product", Value: "product"},
	{Text: "Marketing", Value: "marketing"},
	{Text: "Support", Value: "support"},
	{Text: "DevOps", Value: "devops"},
	{Text: "HR", Value: "hr"},
	{Text: "Others...", Value: "others"},
}

// Default project options
var projectOptions = []*model.PostActionOptions{
	{Text: "Backend", Value: "backend"},
	{Text: "Frontend", Value: "frontend"},
	{Text: "Others...", Value: "others"},
}

// Default environment options
var environmentOptions = []*model.PostActionOptions{
	{Text: "Development", Value: "develop"},
	{Text: "Stage", Value: "stage"},
	{Text: "Production", Value: "production"},
	{Text: "Others...", Value: "others"},
}

// Message priority options (using Mattermost native priority)
var priorityOptions = []*model.PostActionOptions{
	{Text: "Standard", Value: "standard"},
	{Text: "Important", Value: "important"},
	{Text: "Urgent", Value: "urgent"},
}
