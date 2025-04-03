package actions

import (
	"github.com/JackHumphries9/dapper-go/discord"
	"github.com/JackHumphries9/dapper-go/helpers"
)

type Button struct {
	Button     *discord.Button
	Properties ActionOptions
	OnPress    InteractionHandler
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

func (b *Button) CreateButtonInstance(opts helpers.ButtonInstanceOptions) discord.MessageComponent {
	return helpers.CreateButtonInstance(b.Button, opts)
}
