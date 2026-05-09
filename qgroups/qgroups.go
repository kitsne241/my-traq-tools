package qgroups

import (
	"context"

	"github.com/traPtitech/go-traq"
	wsbot "github.com/traPtitech/traq-ws-bot"
)

type QGroups struct {
	bot         *wsbot.Bot
	groupNameID map[string]string
	groupIDData map[string]traq.UserGroup
}

// 引数の Bot をもとにインスタンスを生成
func New(bot *wsbot.Bot) (*QGroups, error) {
	q := &QGroups{bot: bot}
	err := q.Refresh()
	if err != nil {
		return nil, err
	}
	return q, nil
}

// traQ の直近の全グループの名前と ID の対応表を取得
func (q *QGroups) Refresh() error {
	groups, _, err := q.bot.API().GroupAPI.GetUserGroups(context.Background()).Execute()
	if err != nil {
		return err
	}

	groupNameID := map[string]string{}
	groupIDData := map[string]traq.UserGroup{}

	for _, group := range groups {
		groupNameID[group.Name] = group.Id
		groupIDData[group.Id] = group
	}

	q.groupNameID = groupNameID
	q.groupIDData = groupIDData

	return nil
}

// 引数の ID をもつグループの現在のデータを取得
func (q *QGroups) GetGroup(id string) (traq.UserGroup, bool) {
	group, ok := q.groupIDData[id]
	return group, ok
}

// 引数の ID をもつグループの名前を取得
func (q *QGroups) GetGroupName(id string) (string, bool) {
	group, ok := q.GetGroup(id)
	if !ok {
		return "", false
	}
	return group.Name, true
}

// 引数の名前をもつグループの ID を取得
func (q *QGroups) GetGroupIDByName(name string) (string, bool) {
	groupID, ok := q.groupNameID[name]
	return groupID, ok
}

// 引数の名前を持つグループの現在のデータを取得
func (q *QGroups) GetGroupByName(name string) (traq.UserGroup, bool) {
	groupID, ok := q.GetGroupIDByName(name)
	if !ok {
		return traq.UserGroup{}, false
	}
	return q.GetGroup(groupID)
}
