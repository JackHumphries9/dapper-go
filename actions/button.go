package actions

import (
	"github.com/JackHumphries9/dapper-go/discord"
	"github.com/JackHumphries9/dapper-go/discord/button_style"
)

type Button struct {
	Button     *discord.Button
	Properties ActionOptions
	OnPress    InteractionHandler
}

type ButtonInstanceOptions struct {
	ID       *string
	Disabled *bool
	Style    *button_style.ButtonStyle
	Emoji    *discord.Emoji
}

func (b Button) CustomID() string {
	if b.Button == nil || b.Button.CustomId == nil {
		return ""
	}

	return *b.Button.CustomId

}

func (b Button) Options() ActionOptions {
	return b.Properties
}

func (b Button) Type() ActionType {
	return ActionTypeButton
}

func (b Button) Handler(itc *InteractionContext) {
	b.OnPress(itc)
}

func (b Button) AssociatedActions() []Action {
	return []Action{}
}

func (b *Button) CreateButtonInstance(opts ButtonInstanceOptions) discord.MessageComponent {
	if b.Button == nil {
		return nil
	}

	newOpts := ButtonInstanceOptions{
		ID:       b.Button.CustomId,
		Disabled: b.Button.Disabled,
		Style:    &b.Button.Style,
		Emoji:    b.Button.Emoji,
	}

	if opts.Disabled != nil {
		newOpts.Disabled = opts.Disabled
	}
	if opts.Emoji != nil {
		newOpts.Emoji = opts.Emoji
	}
	if opts.Emoji != nil {
		newOpts.Emoji = opts.Emoji
	}
	if opts.ID != nil {
		id := b.CustomID() + ":" + *opts.ID
		newOpts.ID = &id
	}

	return &discord.Button{
		Style:      *newOpts.Style,
		Label:      b.Button.Label,
		Emoji:      newOpts.Emoji,
		Url:        b.Button.Url,
		Disabled:   newOpts.Disabled,
		ButtonType: b.Button.ButtonType,
		CustomId:   newOpts.ID,
	}
}
