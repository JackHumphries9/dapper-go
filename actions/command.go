package actions

import (
	"github.com/JackHumphries9/dapper-go/client"
)

type Command struct {
	Command    client.CreateApplicationCommand
	Actions    []Action
	Properties ActionOptions
	OnInvoke   InteractionHandler
}

func (c Command) CustomID() string {
	return c.Command.Name
}

func (c Command) Options() ActionOptions {
	return c.Properties
}

func (c Command) Type() ActionType {
	return ActionTypeCommand
}

func (c Command) Handler(itc *InteractionContext) {
	c.OnInvoke(itc)
}

func (c Command) AssociatedActions() []Action {
	return c.Actions
}
