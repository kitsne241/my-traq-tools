package qstamps

import (
	"context"

	traq "github.com/traPtitech/go-traq"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
)

type QStamps struct {
	bot         *traqwsbot.Bot
	stampNameID map[string]string
	stampIDName map[string]string
}

// 引数の Bot をもとにインスタンスを生成
func New(bot *traqwsbot.Bot, use []string) *QStamps {
	q := &QStamps{bot: bot}
	q.RefreshBimap()
	return q
}

// traQ の直近の全スタンプの名前と ID の対応表を取得
func (q *QStamps) RefreshBimap() error {
	stamps, _, err := q.bot.API().StampAPI.GetStamps(context.Background()).Execute()
	if err != nil {
		return err
	}

	stampNameID := map[string]string{}
	stampIDName := map[string]string{}

	for _, stamp := range stamps {
		stampIDName[stamp.Id] = stamp.Name
		stampNameID[stamp.Name] = stamp.Id
	}

	q.stampNameID = stampNameID
	q.stampIDName = stampIDName

	return nil
}

// 引数の名前をもつスタンプの現在の ID を取得
func (q *QStamps) GetStampID(name string) *string {
	stampID, exists := q.stampNameID[name]
	if exists {
		return &stampID
	} else {
		return nil
	}
}

// 引数の ID をもつスタンプの現在の名前を取得
func (q *QStamps) GetStampName(id string) *string {
	stampName, exists := q.stampIDName[id]
	if exists {
		return &stampName
	} else {
		return nil
	}
}

// メッセージに引数の名前のスタンプを順番につける
func (q *QStamps) Stamp(messageID string, stampNames ...string) error {
	for _, name := range stampNames {
		stampID, exists := q.stampNameID[name]
		if !exists {
			continue // 存在しないスタンプは無視
		}

		_, err := q.bot.API().MessageAPI.AddMessageStamp(context.Background(), messageID, stampID).
			PostMessageStampRequest(*traq.NewPostMessageStampRequestWithDefaults()).Execute()

		if err != nil {
			return err
		}
	}

	return nil
}
