package server

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/JackHumphries9/dapper-go/actions"
	"github.com/JackHumphries9/dapper-go/client"
	"github.com/JackHumphries9/dapper-go/discord"

	"github.com/JackHumphries9/dapper-go/routers"
	"github.com/JackHumphries9/dapper-go/verification"
)

type InteractionServerOptions struct {
	PublicKey      ed25519.PublicKey
	DapperLogger   *DapperLogger
	StateDelimiter string
}

var defaultConfig = InteractionServerOptions{
	PublicKey: ed25519.PublicKey(""),
}

type InteractionHandler struct {
	opts              InteractionServerOptions
	interactionRouter routers.InteractionRouter
	logger            *DapperLogger
}

func (ih *InteractionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		ih.logger.Error("Only POST method is supported")
		return
	}

	verify := verification.Verify(r, ih.opts.PublicKey)

	if !verify {
		ih.logger.Error("Recieved an invalid request")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rawBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		ih.logger.Error("Failed to read body")
		return
	}

	interaction, err := discord.ParseInteraction(string(rawBody))

	if err != nil {
		ih.logger.Error(fmt.Sprintf("Failed to parse interaction: %v\n", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ih.logger.OnInteractionRecieved(interaction)

	if interaction.IsPing() {
		discord.CreatePongResponse().ToHttpResponse().WriteResponse(w)
		return
	}

	interactionResponse, err := ih.interactionRouter.RouteInteraction(interaction)

	if err != nil {
		ih.logger.Error(fmt.Sprintf("failed to route the interaction: %+v\n", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: Handle attachments here

	body, err := json.Marshal(interactionResponse)
	if err != nil {
		ih.logger.Error("An error occured while responding to interaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)

	return
}

func (ih *InteractionHandler) RegisterAction(action actions.Action) {
	ih.interactionRouter.RegisterAction(action)
}

func (ih *InteractionHandler) RegisterCommandsWithDiscord(appId discord.Snowflake, client *client.BotClient) error {
	err := ih.interactionRouter.RegisterCommandsWithDiscord(appId, client)

	if err != nil {
		ih.logger.Error(fmt.Sprintf("Failed to register discord commands: %v\n", err))
	} else {
		ih.logger.Info("Successfully registered discord commands")
	}

	return err
}
func NewInteractionHandler(publicKey string) InteractionHandler {
	key, err := hex.DecodeString(publicKey)

	if err != nil {
		panic("Invalid public key")
	}

	return NewInteractionHandlerWithOptions(InteractionServerOptions{
		PublicKey:      ed25519.PublicKey(key),
		DapperLogger:   &DefaultLogger,
		StateDelimiter: "/",
	})
}

func NewInteractionHandlerWithOptions(iso InteractionServerOptions) InteractionHandler {
	if iso.DapperLogger == nil {
		iso.DapperLogger = &DefaultLogger
	}

	return InteractionHandler{
		opts:              iso,
		interactionRouter: routers.NewInteractionRouter(iso.StateDelimiter),
		logger:            iso.DapperLogger,
	}
}
