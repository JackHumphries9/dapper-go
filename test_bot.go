package main

import (
	"fmt"
	"os"

	"github.com/JackHumphries9/dapper-go/client"
	"github.com/JackHumphries9/dapper-go/dapper"
	"github.com/JackHumphries9/dapper-go/discord"
	"github.com/JackHumphries9/dapper-go/discord/button_style"
	"github.com/JackHumphries9/dapper-go/server"
	"github.com/icza/gox/gox"
	"github.com/joho/godotenv"
)

var comp = dapper.DapperButton{
	Component: &discord.Button{
		Style: button_style.Primary,
		Label: gox.Ptr("Next"),
		Emoji: &discord.Emoji{
			Name: gox.Ptr("➡️"),
		},
		CustomId: gox.Ptr("hello-next1"),
	},
	OnPress: func(itx *discord.Interaction) {
		err := itx.EditResponse(discord.ResponseEditData{
			Embeds: []discord.Embed{
				{
					Title: "I'm the second Embed",
				},
			},
		})

		if err != nil {
			fmt.Printf("failed to send edit response")
		}
	},
}

func main() {
	godotenv.Load()
	botServer := server.NewInteractionServer(os.Getenv("PUBLIC_KEY"))
	botClient := client.NewBot(os.Getenv("BOT_TOKEN"))
	appId, err := discord.GetSnowflake(os.Getenv("APP_ID"))

	if err != nil {
		panic("Heyo you messed up")
	}

	cmdTest := dapper.DapperCommand{
		Command: client.CreateApplicationCommand{
			Name:        "hello",
			Description: gox.Ptr("Test Description"),
		}, CommandOptions: dapper.DapperCommandOptions{
			Ephemeral: false,
		},
		OnCommand: func(itx *discord.Interaction) {
			err := itx.EditResponse(discord.ResponseEditData{
				Embeds: []discord.Embed{
					{
						Title: "Hello From Dapper",
					},
				},
				Components: []discord.MessageComponent{
					&discord.ActionRow{
						Components: []discord.MessageComponent{
							comp.Component,
						},
					},
				},
			})

			if err != nil {
				fmt.Printf("failed to send edit response")
			}

		},
	}

	botServer.RegisterCommand(cmdTest)

	botServer.RegisterCommand(dapper.DapperCommand{
		Command: client.CreateApplicationCommand{
			Name:        "world",
			Description: gox.Ptr("The World!"),
		}, CommandOptions: dapper.DapperCommandOptions{
			Ephemeral: false,
		},
		OnCommand: func(itx *discord.Interaction) {
			err := itx.EditResponse(discord.ResponseEditData{
				Embeds: []discord.Embed{
					{
						Title:       "Hello World!",
						Description: "This is some more stuff haha lol",
					},
				},
			})

			if err != nil {
				fmt.Printf("failed to send edit response")
			}

		},
	})

	botServer.RegisterComponent(comp)

	err = botServer.RegisterCommandsWithDiscord(appId, botClient)

	if err != nil {
		fmt.Printf("Failed to register discord commands: %v\n", err)
	}

	botServer.Listen(3000)
}
