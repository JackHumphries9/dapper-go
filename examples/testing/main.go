package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/JackHumphries9/dapper-go/actions"
	"github.com/JackHumphries9/dapper-go/client"
	"github.com/JackHumphries9/dapper-go/discord"
	"github.com/JackHumphries9/dapper-go/discord/button_style"
	"github.com/JackHumphries9/dapper-go/discord/text_input_style"
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

var button = actions.Button{
	Button: &discord.Button{
		Style:    button_style.Primary,
		Label:    helpers.Ptr("Test"),
		CustomId: helpers.Ptr("test-btn"),
	},
	OnPress: func(itc *actions.InteractionContext) {
		itc.SetEphemeral(true)

		_ = itc.ShowModal(actions.Modal{
			Modal: discord.ModalCallback{
				CustomId: "test",
				Title:    "Testing Modal",
				Components: helpers.CreateActionRow(&discord.TextInput{
					CustomId:    "t",
					Style:       text_input_style.Short,
					Label:       "Test Label",
					Required:    true,
					Placeholder: "some string",
				}),
			},
		})
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

	botServer.RegisterAction(actions.Command{
		Command: client.CreateApplicationCommand{
			Name:        "fruit",
			Description: helpers.Ptr("Ping Pong"),
		},
		Actions: []actions.Action{button},
		OnInvoke: func(itc *actions.InteractionContext) {
			itc.SetEphemeral(true)
			// itc.Defer()

			// fileBytes, err := os.ReadFile("./examples/testing/test-image.png")
			// if err != nil {
			// 	fmt.Printf("Error reading file: %v\n", err)
			// 	return
			// }

			err = itc.Respond(discord.ResponseEditData{
				Embeds: []discord.Embed{
					{
						Title:       "Hello World!",
						Description: "Hello World!",
						// Image: &discord.EmbedImage{
						// 	URL: "attachment://test-image.png",
						// },
					},
				},
				Components: helpers.CreateActionRow(button.Button),
				// Attachments: []discord.MessageAttachment{
				// 	discord.NewBytesAttachment(fileBytes, "test-image.png", "image/png"),
				// },
			})

			if err != nil {
				fmt.Printf("cannot respond to message %v", err)
			}
		},
	})

	botServer.RegisterCommandsWithDiscord(appId, botClient)

	http.HandleFunc("/", botServer.Handle)

	// Start the server on port 8080
	fmt.Println("Starting server on :3000")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
