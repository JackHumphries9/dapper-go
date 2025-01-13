package discord

import "github.com/JackHumphries9/dapper-go/discord/command_option_type"

type ApplicationCommandDataOption struct {
	Name    string                                `json:"name"`
	Type    command_option_type.CommandOptionType `json:"type"`
	Value   interface{}                           `json:"value,omitempty"`
	Options []ApplicationCommandDataOption        `json:"options,omitempty"`
	Focused bool                                  `json:"focused,omitempty"`
}
