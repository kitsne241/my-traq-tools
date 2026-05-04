package qstamps

import (
	"context"
	"fmt"

	traq "github.com/traPtitech/go-traq"
	wsbot "github.com/traPtitech/traq-ws-bot"
)

type QStamps struct {
	bot         *wsbot.Bot
	stampNameID map[string]string
	stampIDName map[string]string
}

// 引数の Bot をもとにインスタンスを生成
func New(bot *wsbot.Bot) (*QStamps, error) {
	q := &QStamps{bot: bot}
	err := q.Refresh()
	if err != nil {
		return nil, err
	}
	return q, nil
}

// traQ の直近の全スタンプの名前と ID の対応表を取得
func (q *QStamps) Refresh() error {
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
func (q *QStamps) GetStampID(name string) (string, bool) {
	stampID, ok := q.stampNameID[name]
	return stampID, ok
}

// 引数の ID をもつスタンプの現在の名前を取得
func (q *QStamps) GetStampName(id string) (string, bool) {
	stampName, exists := q.stampIDName[id]
	return stampName, exists
}

// メッセージに引数の名前のスタンプを順番につける。返り値は成功したスタンプの名前の配列
func (q *QStamps) Stamp(messageID string, stampNames ...string) ([]string, error) {
	var err error = nil
	successfulStamps := []string{}

	for _, name := range stampNames {
		stampID, ok := q.stampNameID[name]
		if !ok {
			err = fmt.Errorf("スタンプ :%s: はキャッシュに存在しません", name)
			break
		}

		_, err = q.bot.API().MessageAPI.AddMessageStamp(context.Background(), messageID, stampID).
			PostMessageStampRequest(*traq.NewPostMessageStampRequestWithDefaults()).Execute()

		if err != nil {
			break
		}

		successfulStamps = append(successfulStamps, name)
	}

	return successfulStamps, err
}
