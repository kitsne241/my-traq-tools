package qutils

import (
	"context"

	"github.com/traPtitech/go-traq"
	wsbot "github.com/traPtitech/traq-ws-bot"
)

// メッセージの投稿だけをしたいときのためのショートカット
func Post(bot *wsbot.Bot, text string, channelID string) error {
	_, _, err := bot.API().MessageAPI.PostMessage(context.Background(), channelID).
		PostMessageRequest(traq.PostMessageRequest{Content: text}).Execute()
	return err
}
