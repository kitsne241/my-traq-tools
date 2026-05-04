package qusers

import (
	"context"

	wsbot "github.com/traPtitech/traq-ws-bot"
)

type QUsers struct {
	bot        *wsbot.Bot
	userNameID map[string]string
	userIDName map[string]string
}

// 引数の Bot をもとにインスタンスを生成
func New(bot *wsbot.Bot) (*QUsers, error) {
	q := &QUsers{bot: bot}
	err := q.Refresh()
	if err != nil {
		return nil, err
	}
	return q, nil
}

// traQ の直近の全ユーザーの Display ID と ID の対応表を取得
func (q *QUsers) Refresh() error {
	users, _, err := q.bot.API().UserAPI.GetUsers(context.Background()).IncludeSuspended(true).Execute()
	if err != nil {
		return err
	}

	userNameID := map[string]string{}
	userIDName := map[string]string{}

	for _, user := range users {
		userNameID[user.Name] = user.Id
		userIDName[user.Id] = user.Name
	}

	q.userNameID = userNameID
	q.userIDName = userIDName

	return nil
}

// 引数の Display ID をもつユーザーの ID を取得
func (q *QUsers) GetUserID(name string) (string, bool) {
	userID, ok := q.userNameID[name]
	return userID, ok
}

// 引数の ID をもつユーザーの Display ID を取得
func (q *QUsers) GetUserName(id string) (string, bool) {
	userName, ok := q.userIDName[id]
	return userName, ok
}
