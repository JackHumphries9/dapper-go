package discord

import (
	"github.com/JackHumphries9/dapper-go/discord/sku_flags"
	sku_type "github.com/JackHumphries9/dapper-go/discord/sku_types"
)

type SKU struct {
	ID            Snowflake          `json:"id"`
	Type          sku_type.SkuType   `json:"type"`
	ApplicationID Snowflake          `json:"application_id"`
	Name          string             `json:"name"`
	Slug          string             `json:"slug"`
	Flags         sku_flags.SKUFlags `json:"flags"`
}
