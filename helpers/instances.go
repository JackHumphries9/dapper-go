package helpers

import (
	"github.com/JackHumphries9/dapper-go/discord"
	"github.com/JackHumphries9/dapper-go/discord/button_style"
)

type ButtonInstanceOptions struct {
	State    *string
	CustomID *string // Setting this overrides state
	Disabled *bool
	Style    *button_style.ButtonStyle
	Emoji    *discord.Emoji
}

func CreateButtonInstance(button *discord.Button, opts ButtonInstanceOptions) discord.MessageComponent {
	buttonInstance := button

	if opts.Disabled != nil {
		buttonInstance.Disabled = opts.Disabled
	}
	if opts.Emoji != nil {
		buttonInstance.Emoji = opts.Emoji
	}
	if opts.Emoji != nil {
		buttonInstance.Emoji = opts.Emoji
	}
	if opts.State != nil {
		id := *button.CustomId + ":" + *opts.State
		buttonInstance.CustomId = &id
	}
	if opts.CustomID != nil {
		buttonInstance.CustomId = opts.CustomID
	}

	return buttonInstance
}
