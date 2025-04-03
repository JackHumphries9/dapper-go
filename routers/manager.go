package routers

import (
	"fmt"

	"github.com/JackHumphries9/dapper-go/actions"
	"github.com/JackHumphries9/dapper-go/client"
	"github.com/JackHumphries9/dapper-go/discord"
	"github.com/JackHumphries9/dapper-go/discord/interaction_callback_type"
	"github.com/JackHumphries9/dapper-go/discord/interaction_type"
	"github.com/JackHumphries9/dapper-go/helpers"
)

type InteractionRouter struct {
	actions map[string]actions.Action
}

func NewInteractionRouter() InteractionRouter {
	return InteractionRouter{
		actions: make(map[string]actions.Action, 0),
	}
}

func (ir *InteractionRouter) RouteInteraction(interaction *discord.Interaction) (discord.InteractionResponse, error) {
	var interactionCustomId string

	if interaction.Type == interaction_type.ApplicationCommand {
		interactionCustomId = interaction.Data.(*discord.ApplicationCommandData).Name

	} else if interaction.Type == interaction_type.MessageComponent {
		interactionCustomId = interaction.Data.(*discord.MessageComponentData).CustomId

	} else if interaction.Type == interaction_type.ModalSubmit {
		interactionCustomId = interaction.Data.(*discord.ModalSubmitData).CustomId

	} else {
		return discord.InteractionResponse{}, fmt.Errorf("invalid interaction type")
	}

	// Find associated action
	if action, ok := ir.actions[interactionCustomId]; ok {
		deferralChan := make(chan *discord.InteractionResponse)

		itc := actions.NewInteractionContext(interaction, deferralChan, action.Options().CancelDefer)

		if action.Options().Ephemeral {
			itc.SetEphemeral(true)
		}

		go action.Handler(&itc)

		if !action.Options().CancelDefer {
			response := <-deferralChan

			return *response, nil
		}

		return discord.InteractionResponse{
			Type: interaction_callback_type.DeferredUpdateMessage,
			Data: &discord.MessageCallbackData{
				Flags: helpers.Ptr(int(itc.GetMessageFlags())),
			},
		}, nil
	}

	return discord.InteractionResponse{}, fmt.Errorf("Cannot find interaction: %s", interactionCustomId)
}

func (ir *InteractionRouter) bindAction(action actions.Action) {
	if _, ok := ir.actions[action.CustomID()]; ok {
		panic("action already exists")
	}

	fmt.Printf("binding action %s\n", action.CustomID())

	ir.actions[action.CustomID()] = action
}

func (ir *InteractionRouter) RegisterAction(action actions.Action) {
	ir.bindAction(action)

	// Bind associated actions
	for _, act := range action.AssociatedActions() {
		ir.bindAction(act)
	}
}

func (ir *InteractionRouter) RegisterCommandsWithDiscord(appId discord.Snowflake, botClient *client.BotClient) error {
	discordCommands := make([]client.CreateApplicationCommand, 0)

	for _, cmd := range ir.actions {
		if cmd.Type() == actions.ActionTypeCommand {
			discordCommands = append(discordCommands, cmd.(actions.Command).Command)
		}
	}

	return botClient.GetApplicationClient(appId).RegisterCommands(discordCommands)
}
