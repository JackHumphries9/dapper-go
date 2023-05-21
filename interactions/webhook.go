package interactions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tgwaffles/gladis/components"
	"github.com/tgwaffles/gladis/discord"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	webhookUrl = apiUrl + "/webhooks/%d/%s"
)

type WebhookRequest struct {
	Content         *string                       `json:"content,omitempty"`
	Username        string                        `json:"username,omitempty"`
	AvatarUrl       string                        `json:"avatar_url,omitempty"`
	TTS             *bool                         `json:"tts,omitempty"`
	Embeds          []discord.Embed               `json:"embeds,omitempty"`
	AllowedMentions *discord.AllowedMentions      `json:"allowed_mentions,omitempty"`
	Flags           *int                          `json:"flags,omitempty"`
	Components      []components.MessageComponent `json:"components,omitempty"`
	Attachments     []discord.Attachment          `json:"attachments,omitempty"`
	ThreadName      string                        `json:"thread_name,omitempty"`
}

type WebhookMessageResponse struct {
}

type Webhook struct {
	Id            discord.Snowflake  `json:"id"`
	Type          uint8              `json:"type"`
	GuildId       *discord.Snowflake `json:"guild_id,omitempty"`
	ChannelId     *discord.Snowflake `json:"channel_id,omitempty"`
	User          *discord.User      `json:"user,omitempty"`
	Name          *string            `json:"name,omitempty"`
	Avatar        *string            `json:"avatar,omitempty"`
	Token         *string            `json:"token,omitempty"`
	ApplicationId *discord.Snowflake `json:"application_id,omitempty"`
	SourceGuild   *discord.Guild     `json:"source_guild,omitempty"`
	SourceChannel *discord.Channel   `json:"source_channel,omitempty"`
	Url           *string            `json:"url,omitempty"`
}

func WebhookFromUrl(url string) (webhook *Webhook, err error) {
	// Webhook is in the format https://discord.com/api/webhooks/<id>/<token>
	splitUrl := strings.Split(url, "webhooks/")
	if len(splitUrl) != 2 {
		return nil, fmt.Errorf("invalid webhook url - missing 'webhooks/' (given: %s)", url)
	}

	splitUrl = strings.Split(splitUrl[1], "/")
	if len(splitUrl) != 2 {
		return nil, fmt.Errorf("invalid webhook url - missing id or token (given: %s)", url)
	}

	webhook = &Webhook{}
	webhookId, err := strconv.ParseUint(splitUrl[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid webhook url - invalid id (given: %s), error: %w", splitUrl[0], err)
	}
	webhook.Id = discord.Snowflake(webhookId)
	token := strings.Trim(splitUrl[1], "/ \n\r\t") // Trim trailing slashes and whitespace
	webhook.Token = &token
	webhook.Url = &url

	return webhook, nil
}

func (hook *Webhook) GetUrl() string {
	if hook.Url == nil {
		url := fmt.Sprintf(webhookUrl, hook.Id, *hook.Token)
		hook.Url = &url
	}

	return *hook.Url
}

func (hook *Webhook) Send(req WebhookRequest) (err error) {
	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error marshaling data to JSON: %w", err)
	}
	request, err := http.NewRequest("POST", hook.GetUrl(), bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %w", err)
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("expected status code 204, got %d", resp.StatusCode)
	}

	return nil
}

type WebhookGetMessageRequest struct {
	// String so it can be "@original"
	MessageId string             `json:"-"` // Not sent in request body
	ThreadId  *discord.Snowflake `json:"thread_id,omitempty"`
}

func (hook *Webhook) GetMessage(req WebhookGetMessageRequest) (message *discord.Message, err error) {
	var body io.Reader
	if req.ThreadId != nil {
		data, err := json.Marshal(req)
		if err != nil {
			return nil, fmt.Errorf("error marshaling data to JSON: %w", err)
		}
		body = bytes.NewReader(data)
	}

	request, err := http.NewRequest("GET", hook.GetUrl()+fmt.Sprintf("/messages/%s", req.MessageId), body)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("expected status code 200, got %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&message)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}

	return message, nil
}

func (hook *Webhook) EditMessage(messageId string, data ResponseEditData) error {
	err := data.Verify()
	if err != nil {
		return fmt.Errorf("error verifying edit data: %w", err)
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling data to JSON: %w", err)
	}

	request, err := http.NewRequest("PATCH", hook.GetUrl()+fmt.Sprintf("/messages/%s", messageId), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading HTTP response body: %w", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("expected status code 200, got %d. Response body: %s", resp.StatusCode, string(responseBody))
	}

	return nil
}

func (hook *Webhook) DeleteMessage(messageId string) error {
	request, err := http.NewRequest("DELETE", hook.GetUrl()+fmt.Sprintf("/messages/%s", messageId), nil)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading HTTP response body: %w", err)
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("expected status code 204, got %d. Response body: %s", resp.StatusCode, string(responseBody))
	}

	return nil
}
