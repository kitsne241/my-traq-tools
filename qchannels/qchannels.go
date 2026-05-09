package qchannels

import (
	"context"

	"github.com/traPtitech/go-traq"
	wsbot "github.com/traPtitech/traq-ws-bot"
)

type QChannels struct {
	bot           *wsbot.Bot
	idTree        map[string]string
	channelPathID map[string]string
	channelIDPath map[string]string
	channelIDData map[string]traq.Channel
}

// 引数の Bot をもとにインスタンスを生成
func New(bot *wsbot.Bot) (*QChannels, error) {
	q := &QChannels{bot: bot}
	err := q.Refresh()
	if err != nil {
		return nil, err
	}
	return q, nil
}

// traQ の直近の全チャンネルのパスと ID の対応表を取得
func (q *QChannels) Refresh() error {
	channels, _, err := q.bot.API().ChannelAPI.GetChannels(context.Background()).IncludeDm(false).Execute()
	if err != nil {
		return err
	}

	idTree := map[string]string{} // 子チャンネルの ID をキー、親チャンネルの ID を値にもつ
	channelIDData := map[string]traq.Channel{}
	channelPathID := map[string]string{}
	channelIDPath := map[string]string{}

	for _, channel := range channels.Public {
		channelIDData[channel.Id] = channel
		parentID := channel.ParentId.Get()
		if parentID != nil {
			idTree[channel.Id] = *parentID
		}
	}
	// ID の木構造とそれぞれの ID をもつチャンネルの名称を取得。この木からそれぞれのパスを作る

	for _, channel := range channels.Public {
		currentID := channel.Id
		exists := false
		path := channelIDData[currentID].Name
		for {
			currentID, exists = idTree[currentID]
			if !exists {
				break
			}
			path = channelIDData[currentID].Name + "/" + path
		}
		channelPathID[path] = channel.Id
		channelIDPath[channel.Id] = path
	}

	q.idTree = idTree
	q.channelIDData = channelIDData
	q.channelPathID = channelPathID
	q.channelIDPath = channelIDPath

	return nil
}

// 引数の ID をもつチャンネルの現在のデータを取得
func (q *QChannels) GetChannel(id string) (traq.Channel, bool) {
	channel, ok := q.channelIDData[id]
	return channel, ok
}

// 引数の ID をもつチャンネルの現在のパスを取得
func (q *QChannels) GetChannelPath(id string) (string, bool) {
	channelPath, ok := q.channelIDPath[id]
	return channelPath, ok
}

// 引数のパスをもつチャンネルの現在の ID を取得
func (q *QChannels) GetChannelIDByPath(path string) (string, bool) {
	channelID, ok := q.channelPathID[path]
	return channelID, ok
}

// 引数のパスをもつチャンネルの現在のデータを取得
func (q *QChannels) GetChannelByPath(path string) (traq.Channel, bool) {
	channelID, ok := q.GetChannelIDByPath(path)
	if !ok {
		return traq.Channel{}, false
	}
	return q.GetChannel(channelID)
}
