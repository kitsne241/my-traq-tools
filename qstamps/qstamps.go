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
	stampIDData map[string]traq.Stamp
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
	stampsWithThumbnail, _, err := q.bot.API().StampAPI.GetStamps(context.Background()).Execute()
	if err != nil {
		return err
	}

	// スタンプには謎に型が Stamp と StampWithThumbnail の 2 つある
	// とりあえずシンプルな方を使ってみて、hasThumbnail が思ったより重要だと分かった暁には乗り換えてみる

	stamps := make([]traq.Stamp, len(stampsWithThumbnail))

	for i, stampWithThumbnail := range stampsWithThumbnail {
		stamps[i] = traq.Stamp{
			Id:        stampWithThumbnail.Id,
			Name:      stampWithThumbnail.Name,
			CreatorId: stampWithThumbnail.CreatorId,
			CreatedAt: stampWithThumbnail.CreatedAt,
			UpdatedAt: stampWithThumbnail.UpdatedAt,
			FileId:    stampWithThumbnail.FileId,
			IsUnicode: stampWithThumbnail.IsUnicode,
		}
	}

	stampNameID := map[string]string{}
	stampIDData := map[string]traq.Stamp{}

	for _, stamp := range stamps {
		stampIDData[stamp.Id] = stamp
		stampNameID[stamp.Name] = stamp.Id
	}

	q.stampNameID = stampNameID
	q.stampIDData = stampIDData

	return nil
}

// 引数の名前をもつスタンプの現在の ID を取得
func (q *QStamps) GetStampID(name string) (string, bool) {
	stampID, ok := q.stampNameID[name]
	return stampID, ok
}

// 引数の ID をもつスタンプの現在のデータを取得
func (q *QStamps) GetStamp(id string) (traq.Stamp, bool) {
	stamp, ok := q.stampIDData[id]
	return stamp, ok
}

// 引数の ID をもつスタンプの現在の名前を取得
func (q *QStamps) GetStampName(id string) (string, bool) {
	stamp, ok := q.GetStamp(id)
	if !ok {
		return "", false
	}
	return stamp.Name, true
}

// 引数の名前をもつスタンプの現在のデータを取得
func (q *QStamps) GetStampByName(name string) (traq.Stamp, bool) {
	stampID, ok := q.GetStampID(name)
	if !ok {
		return traq.Stamp{}, false
	}
	return q.GetStamp(stampID)
}

// メッセージに引数の名前のスタンプを順番につける。返り値は成功したスタンプの名前の配列
func (q *QStamps) Stamp(messageID string, stampNames ...string) ([]string, error) {
	var err error = nil
	successfulStamps := []string{}

	for _, name := range stampNames {
		stampID, ok := q.GetStampID(name)
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
