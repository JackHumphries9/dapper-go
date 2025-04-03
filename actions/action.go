package actions

type InteractionHandler func(itc *InteractionContext)

type ActionType string

const (
	ActionTypeCommand         ActionType = "command"
	ActionTypeSubcommand                 = "subcommand"
	ActionTypeSubcommandGroup            = "subcommand_group"
	ActionTypeButton                     = "button"
	ActionTypeSelect                     = "select"
	ActionTypeModal                      = "modal"
)

type ActionOptions struct {
	CancelDefer bool
	Ephemeral   bool
}

type Action interface {
	CustomID() string
	Handler(itc *InteractionContext)
	AssociatedActions() []Action
	Options() ActionOptions
	Type() ActionType
}
