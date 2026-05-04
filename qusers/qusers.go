package qusers

import (
	"context"

	traqwsbot "github.com/traPtitech/traq-ws-bot"
)

type QUsers struct {
	bot        *traqwsbot.Bot
	userNameID map[string]string
	userIDName map[string]string
}

// 引数の Bot をもとにインスタンスを生成
func New(bot *traqwsbot.Bot, use []string) *QUsers {
	q := &QUsers{bot: bot}
	q.RefreshBimap()
	return q
}

// traQ の直近の全ユーザーの Display ID と ID の対応表を取得
func (q *QUsers) RefreshBimap() error {
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
func (q *QUsers) GetUserID(name string) *string {
	userID, exists := q.userNameID[name]
	if exists {
		return &userID
	} else {
		return nil
	}
}

// 引数の ID をもつユーザーの Display ID を取得
func (q *QUsers) GetUserName(id string) *string {
	userName, exists := q.userIDName[id]
	if exists {
		return &userName
	} else {
		return nil
	}
}
