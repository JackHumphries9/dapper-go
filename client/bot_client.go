package client

import (
	"encoding/json"
	"fmt"
	"github.com/tgwaffles/gladis/client/errors"
	"github.com/tgwaffles/gladis/discord"
	"net/http"
)

type BotClient struct {
	Token  string
	Client *http.Client
}

func NewBot(token string) *BotClient {
	return &BotClient{
		Token:  token,
		Client: http.DefaultClient,
	}
}

func getUserAgent() string {
	return "DiscordBot (https://github.com/TGWaffles/GLaDIs, v0.2) Interactions HTTP Client"
}

func (botClient *BotClient) MakeRequest(discordRequest DiscordRequest) (response *http.Response, err error) {
	discordRequest.ValidateEndpoint()
	request, err := http.NewRequest(discordRequest.Method, discordRequest.GetUrl(), discordRequest.getBodyAsReader())
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	if !discordRequest.DisableAuth {
		request.Header.Set("Authorization", "Bot "+botClient.Token)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", getUserAgent())
	request.Header.Set("Accept", "application/json")

	for key, value := range discordRequest.AdditionalHeaders {
		request.Header.Set(key, value)
	}

	response, err = botClient.Client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %w", err)
	}

	if !discordRequest.DisableStatusCheck && response.StatusCode != discordRequest.ExpectedStatus {
		return nil, errors.StatusError{
			Code:     errors.StatusErrorCode(response.StatusCode),
			Response: response,
		}
	}

	if discordRequest.UnmarshalTo != nil {
		err = json.NewDecoder(response.Body).Decode(discordRequest.UnmarshalTo)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling response: %w", err)
		}
		return nil, nil
	}

	return response, nil
}

func (botClient *BotClient) GetGuildClient(guildId discord.Snowflake) *GuildClient {
	return &GuildClient{
		GuildId: guildId,
		Bot:     botClient,
	}
}

func (botClient *BotClient) GetChannelClient(channelId discord.Snowflake) *ChannelClient {
	return &ChannelClient{
		ChannelId: channelId,
		Bot:       botClient,
	}
}

func (botClient *BotClient) GetUserClient(userId discord.Snowflake) *UserClient {
	if userId == 0 {
		return nil
	}
	return &UserClient{
		UserId: userId,
		Bot:    botClient,
	}
}

func (botClient *BotClient) GetSelfUserClient() *UserClient {
	return &UserClient{
		UserId: discord.Snowflake(0),
		Bot:    botClient,
	}
}
