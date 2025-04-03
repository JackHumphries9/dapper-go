package actions

import (
	"github.com/JackHumphries9/dapper-go/discord"
	"github.com/JackHumphries9/dapper-go/discord/interaction_callback_type"
)

type Modal struct {
	Modal      discord.ModalCallback
	Properties ActionOptions
	OnSubmit   InteractionHandler
}

func (m Modal) CustomID() string {
	return m.Modal.CustomId
}

func (m Modal) Options() ActionOptions {
	return m.Properties
}

func (m Modal) Type() ActionType {
	return ActionTypeModal
}

func (m Modal) Handler(itc *InteractionContext) {
	m.OnSubmit(itc)
}

func (m Modal) AssociatedActions() []Action {
	return []Action{}
}

func (m *Modal) GetModalResponse() discord.InteractionResponse {
	return discord.InteractionResponse{
		Type: interaction_callback_type.Modal,
		Data: m.Modal,
	}
}
