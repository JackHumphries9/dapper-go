package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/JackHumphries9/dapper-go/actions"
	"github.com/JackHumphries9/dapper-go/client"
	"github.com/JackHumphries9/dapper-go/discord"
	"github.com/JackHumphries9/dapper-go/discord/command_option_type"
	"github.com/JackHumphries9/dapper-go/helpers"
	"github.com/JackHumphries9/dapper-go/server"
)

const FILENAME = "./examples/env.json"

type Env struct {
	PublicKey string `json:"PUBLIC_KEY"`
	BotToken  string `json:"BOT_TOKEN"`
	AppId     string `json:"APP_ID"`
}

func LoadJSONEnv() Env {
	plan, err := os.ReadFile(FILENAME)

	if err != nil {
		panic("no env file")
	}

	var data Env
	err = json.Unmarshal(plan, &data)

	if err != nil {
		panic("cannot unmarshal")
	}

	return data
}

var command = actions.Command{
	Command: client.CreateApplicationCommand{
		Name:        "user",
		Description: helpers.Ptr("testing"),
	},
	Actions: []actions.Action{
		actions.Subcommand{
			Subcommand: discord.ApplicationCommandOption{
				Type: command_option_type.SubCommand,
				Name: "get",
				Options: []discord.ApplicationCommandOption{
					{
						Type:        command_option_type.User,
						Name:        "user",
						Description: "testing",
					},
				},
			},
			OnInvoke: func(itc *actions.InteractionContext) {
				itc.SetEphemeral(true)
				itc.Defer()

				user, err := itc.GetUserCommandOption("user")

				if err != nil {
					panic("woahhh")
				}

				err = itc.Respond(discord.ResponseEditData{

					Embeds: []discord.Embed{
						{
							Title:       "Got a user!",
							Description: fmt.Sprintf("Got user: %s", user.MentionUserString()),
						},
					},
				})

				if err != nil {
					fmt.Printf("cannot respond to message %v", err)
				}
			},
		},
	},
}

func main() {
	var env = LoadJSONEnv()

	botServer := server.NewInteractionHandler(env.PublicKey)
	botClient := client.NewBot(env.BotToken)
	appId, err := discord.GetSnowflake(env.AppId)

	if err != nil {
		panic("Heyo you messed up")
	}

	botServer.RegisterAction(command)

	botServer.RegisterCommandsWithDiscord(appId, botClient)

	http.HandleFunc("/", botServer.Handle)

	// Start the server on port 8080
	fmt.Println("Starting server on :3000")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
