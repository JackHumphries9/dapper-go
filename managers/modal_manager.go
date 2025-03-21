package managers

import (
	"fmt"

	"github.com/JackHumphries9/dapper-go/discord"
	"github.com/JackHumphries9/dapper-go/interactable"
)

type ModalManager struct {
	modals map[string]interactable.Modal
}

func (dmm *ModalManager) RouteInteraction(itx *discord.Interaction) (discord.InteractionResponse, error) {
	submitData := itx.Data.(*discord.ModalSubmitData)

	if modal, ok := dmm.modals[submitData.CustomId]; ok {
		itc := interactable.InteractionContext{
			Interaction:  itx,
			DeferChannel: make(chan *discord.InteractionResponse),
		}

		go modal.OnSubmit(&itc)

		response := <-itc.DeferChannel

		return *response, nil
	}

	return discord.InteractionResponse{}, fmt.Errorf("Cannot find interaction %s", submitData.CustomId)
}

func (dmm *ModalManager) Register(modal interactable.Modal, prefix string) {
	customId := modal.Modal.CustomId

	if customId == "" {
		panic("cannot register a modal with no id")
	}

	id := prefix + "." + customId

	dmm.modals[id] = modal
}

func NewDapperModalManager() ModalManager {
	return ModalManager{
		modals: make(map[string]interactable.Modal, 0),
	}
}
