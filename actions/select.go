package actions

import (
	"github.com/JackHumphries9/dapper-go/discord"
)

type Select struct {
	Select     *discord.SelectMenu
	Properties ActionOptions
	OnSelect   InteractionHandler
}

func (s Select) CustomID() string {
	if s.Select == nil {
		return ""
	}

	return s.Select.CustomId
}

func (s Select) Options() ActionOptions {
	return s.Properties
}

func (s Select) Type() ActionType {
	return ActionTypeSelect
}

func (s Select) Handler(itc *InteractionContext) {
	s.OnSelect(itc)
}

func (s Select) AssociatedActions() []Action {
	return []Action{}
}

// func (s *Select) CreateComponentInstance(opts Sele) discord.MessageComponent {
// 	return &discord.SelectMenu{
// 		MenuType:     s.Component.MenuType,
// 		Options:      s.Component.Options,
// 		ChannelTypes: db.Component.ChannelTypes,
// 		Placeholder:  db.Component.Placeholder,
// 		MinValues:    db.Component.MinValues,
// 		MaxValues:    db.Component.MaxValues,
// 		Disabled:     &opts.Disabled,
// 		CustomId:     db.Component.CustomId + ":" + opts.ID,
// 	}
// }
