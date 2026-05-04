package qchannels

import (
	"context"

	traqwsbot "github.com/traPtitech/traq-ws-bot"
)

type QChannels struct {
	bot           *traqwsbot.Bot
	idTree        map[string]string
	channelPathID map[string]string
	channelIDPath map[string]string
}

// 引数の Bot をもとにインスタンスを生成
func New(bot *traqwsbot.Bot) (*QChannels, error) {
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
	channelIDName := map[string]string{}
	channelPathID := map[string]string{}
	channelIDPath := map[string]string{}

	for _, channel := range channels.Public {
		channelIDName[channel.Id] = channel.Name
		parentID := channel.ParentId.Get()
		if parentID != nil {
			idTree[channel.Id] = *parentID
		}
	}
	// ID の木構造とそれぞれの ID をもつチャンネルの名称を取得。この木からそれぞれのパスを作る

	for _, channel := range channels.Public {
		currentID := channel.Id
		exists := false
		path := channelIDName[currentID]
		for {
			currentID, exists = idTree[currentID]
			if !exists {
				break
			}
			path = channelIDName[currentID] + "/" + path
		}
		channelPathID[path] = channel.Id
		channelIDPath[channel.Id] = path
	}

	q.idTree = idTree
	q.channelPathID = channelPathID
	q.channelIDPath = channelIDPath

	return nil
}

// 引数のパスをもつチャンネルの現在の ID を取得
func (q *QChannels) GetChannelID(path string) (string, bool) {
	channelID, ok := q.channelPathID[path]
	return channelID, ok
}

// 引数の ID をもつチャンネルの現在のパスを取得
func (q *QChannels) GetChannelPath(id string) (string, bool) {
	channelPath, ok := q.channelIDPath[id]
	return channelPath, ok
}

// 引数の ID をもつチャンネルの親チャンネルの ID を取得
func (q *QChannels) GetParent(channelID string) (string, bool) {
	parentID, ok := q.idTree[channelID]
	return parentID, ok
}

// 引数の ID をもつチャンネルの子チャンネルの ID の配列を取得
func (q *QChannels) GetChildren(channelID string) []string {
	children := []string{}
	for childID, parentID := range q.idTree {
		if parentID == channelID {
			children = append(children, childID)
		}
	}
	return children
}
