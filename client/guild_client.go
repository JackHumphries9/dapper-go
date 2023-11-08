package client

import (
	"github.com/tgwaffles/gladis/discord"
	"net/http"
	"net/url"
	"strconv"
)

type GuildClient struct {
	GuildId discord.Snowflake
	Bot     *BotClient
}

func (guildClient *GuildClient) MakeRequest(discordRequest DiscordRequest) (response *http.Response, err error) {
	discordRequest.ValidateEndpoint()
	discordRequest.Endpoint = "/guilds/" + guildClient.GuildId.String() + discordRequest.Endpoint

	return guildClient.Bot.MakeRequest(discordRequest)
}

func (guildClient *GuildClient) FetchGuild() (*discord.Guild, error) {
	guild := &discord.Guild{}
	_, err := guildClient.MakeRequest(DiscordRequest{
		Method:         "GET",
		Endpoint:       "",
		Body:           nil,
		ExpectedStatus: 200,
		UnmarshalTo:    guild,
	})

	if err != nil {
		return nil, err
	}

	return guild, nil
}

func (guildClient *GuildClient) GetMemberClient(memberId discord.Snowflake) *MemberClient {
	return &MemberClient{
		MemberId:    memberId,
		GuildClient: guildClient,
	}
}

type ActiveThreadsResponse struct {
	Threads       []discord.Channel      `json:"threads"`
	ThreadMembers []discord.ThreadMember `json:"members"`
}

func (guildClient *GuildClient) GetActiveThreads() (ActiveThreadsResponse, error) {
	response := ActiveThreadsResponse{}
	_, err := guildClient.MakeRequest(DiscordRequest{
		Method:         "GET",
		Endpoint:       "/threads/active",
		Body:           nil,
		ExpectedStatus: 200,
		UnmarshalTo:    &response,
	})

	if err != nil {
		return ActiveThreadsResponse{}, err
	}

	return response, nil
}

type ListMembersRequest struct {
	// The last member fetched
	After *discord.Snowflake
	// Max members to fetch in one request
	Limit *int
}

func (guildClient *GuildClient) ListMembers(request ListMembersRequest) ([]discord.Member, error) {
	members := make([]discord.Member, 0)

	query := make(url.Values)
	if request.After != nil {
		query.Add("after", request.After.String())
	}
	if request.Limit != nil {
		query.Add("limit", strconv.Itoa(*request.Limit))
	}
	endpoint := "/members"
	encodedQuery := query.Encode()
	if len(encodedQuery) > 0 {
		endpoint += "?" + encodedQuery
	}

	req := DiscordRequest{
		ExpectedStatus: 200,
		Method:         "GET",
		Endpoint:       endpoint,
		Body:           nil,
		UnmarshalTo:    &members,
	}

	_, err := guildClient.MakeRequest(req)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (guildClient *GuildClient) GetChannels() ([]discord.Channel, error) {
	channels := make([]discord.Channel, 0)
	_, err := guildClient.MakeRequest(DiscordRequest{
		Method:         "GET",
		Endpoint:       "/channels",
		Body:           nil,
		ExpectedStatus: 200,
		UnmarshalTo:    &channels,
	})

	if err != nil {
		return nil, err
	}

	return channels, nil
}
