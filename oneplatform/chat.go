package oneplatform

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/patcharp/go_swth/requests"
	"net/http"
)

const ChatProductionEndpoint = "https://chat-public.one.th:8034/api/v1"

type Chat struct {
	BotId       string
	Token       string
	TokenType   string
	ApiEndpoint string
}

type ChatFriend struct {
	OneEmail    string `json:"one_email"`
	UserId      string `json:"user_id"`
	AccountId   string `json:"one_id"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"`
}

type Choice struct {
	Label string `json:"label"`
	Type  string `json:"type"`
	Url   string `json:"url"`
	Size  string `json:"size"`
}
type Elements struct {
	Image   string   `json:"image"`
	Title   string   `json:"title"`
	Detail  string   `json:"detail"`
	Choices []Choice `json:"choice"`
}

type QuickReply struct {
	Label   string      `json:"label"`
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

func NewChatBot(botId string, token string, tokenType string) Chat {
	return Chat{
		BotId:       botId,
		Token:       token,
		TokenType:   tokenType,
		ApiEndpoint: ChatProductionEndpoint,
	}
}

func (c *Chat) FindOneChatFriend(keyword string) (ChatFriend, error) {
	var chatFriend ChatFriend
	msg := struct {
		BotId   string `json:"bot_id"`
		Keyword string `json:"key_search"`
	}{
		BotId:   c.BotId,
		Keyword: keyword,
	}
	body, _ := json.Marshal(&msg)
	r, err := c.send(http.MethodPost, c.url("/searchfriend"), body)
	if err != nil {
		return chatFriend, nil
	}
	chatFriendResult := struct {
		Status string     `json:"status"`
		Friend ChatFriend `json:"friend"`
	}{}
	if err := json.Unmarshal(r.Body, &chatFriendResult); err != nil {
		return chatFriend, err
	}
	chatFriend = chatFriendResult.Friend
	return chatFriend, nil
}

func (c *Chat) PushTextMessage(to string, msg string, customNotify *string) error {
	pushMessage := struct {
		To           string `json:"to"`
		BotId        string `json:"bot_id"`
		Type         string `json:"type"`
		Message      string `json:"message"`
		CustomNotify string `json:"custom_notification,omitempty"`
	}{
		To:      to,
		BotId:   c.BotId,
		Type:    "text",
		Message: msg,
	}
	if customNotify != nil {
		pushMessage.CustomNotify = *customNotify
	}
	body, _ := json.Marshal(&pushMessage)
	_, err := c.send(http.MethodPost, c.url("/push_message"), body)
	return err
}

func (c *Chat) PushWebView(to string, label string, path string, img string, title string, detail string, customNotify *string) error {
	pushMessage := struct {
		To           string     `json:"to"`
		BotId        string     `json:"bot_id"`
		Type         string     `json:"type"`
		CustomNotify string     `json:"custom_notification,omitempty"`
		Elements     []Elements `json:"elements"`
	}{
		To:    to,
		BotId: c.BotId,
		Type:  "template",
		Elements: []Elements{
			{
				Image:  img,
				Title:  title,
				Detail: detail,
				Choices: []Choice{
					{
						Label: label,
						Type:  "webview",
						Url:   path,
						Size:  "full",
					},
				},
			},
		},
	}

	if customNotify != nil {
		pushMessage.CustomNotify = *customNotify
	}
	body, _ := json.Marshal(&pushMessage)
	r, err := c.send(http.MethodPost, c.url("/push_message"), body)
	if err != nil {
		return err
	}
	if r.Code != 200 {
		return errors.New(fmt.Sprintf("server return error with http code %d : %s", r.Code, string(r.Body)))
	}
	return nil
}

func (c *Chat) PushLink(to string, label string, path string, img string, title string, detail string, customNotify *string) error {
	pushMessage := struct {
		To           string     `json:"to"`
		BotId        string     `json:"bot_id"`
		Type         string     `json:"type"`
		CustomNotify string     `json:"custom_notification,omitempty"`
		Elements     []Elements `json:"elements"`
	}{
		To:    to,
		BotId: c.BotId,
		Type:  "template",
		Elements: []Elements{
			{
				Image:  img,
				Title:  title,
				Detail: detail,
				Choices: []Choice{
					{
						Label: label,
						Type:  "link",
						Url:   path,
					},
				},
			},
		},
	}

	if customNotify != nil {
		pushMessage.CustomNotify = *customNotify
	}
	body, _ := json.Marshal(&pushMessage)
	r, err := c.send(http.MethodPost, c.url("/push_message"), body)
	if err != nil {
		return err
	}
	if r.Code != 200 {
		return errors.New(fmt.Sprintf("server return error with http code %d : %s", r.Code, string(r.Body)))
	}
	return nil
}

func (c *Chat) PushQuickReply(to string, message string, quickReply []QuickReply) error {
	pushQuickReply := struct {
		To         string       `json:"to"`
		BotId      string       `json:"bot_id"`
		Message    string       `json:"message"`
		QuickReply []QuickReply `json:"quick_reply"`
	}{
		To:         to,
		BotId:      c.BotId,
		Message:    message,
		QuickReply: quickReply,
	}
	body, _ := json.Marshal(&pushQuickReply)
	r, err := c.send(http.MethodPost, c.url("/push_quickreply"), body)
	if err != nil {
		return err
	}
	if r.Code != 200 {
		return errors.New(fmt.Sprintf("server return error with http code %d : %s", r.Code, string(r.Body)))
	}
	return nil
}

func (c *Chat) send(method string, url string, body []byte) (requests.Response, error) {
	headers := map[string]string{
		echo.HeaderContentType:   "application/json",
		echo.HeaderAuthorization: fmt.Sprintf("%s %s", c.TokenType, c.Token),
	}
	r, err := requests.Request(method, url, headers, bytes.NewBuffer(body), 0)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (c *Chat) url(path string) string {
	return fmt.Sprintf("%s%s", c.ApiEndpoint, path)
}
