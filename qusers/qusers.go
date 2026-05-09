package qusers

import (
	"context"
	"strings"

	"github.com/traPtitech/go-traq"
	wsbot "github.com/traPtitech/traq-ws-bot"
)

type QUsers struct {
	bot             *wsbot.Bot
	userLowerNameID map[string]string
	userNameID      map[string]string
	userIDData      map[string]traq.User
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

// traQ の直近の全ユーザーの名前と ID の対応表を取得
func (q *QUsers) Refresh() error {
	users, _, err := q.bot.API().UserAPI.GetUsers(context.Background()).IncludeSuspended(true).Execute()
	if err != nil {
		return err
	}

	userLowerNameID := map[string]string{}
	userNameID := map[string]string{}
	userIDData := map[string]traq.User{}

	for _, user := range users {
		userLowerNameID[strings.ToLower(user.Name)] = user.Id
		userNameID[user.Name] = user.Id
		userIDData[user.Id] = user
	}

	q.userLowerNameID = userLowerNameID
	q.userNameID = userNameID
	q.userIDData = userIDData

	return nil
}

// 引数の名前（完全一致）をもつユーザーの ID を取得
func (q *QUsers) GetUserID(name string) (string, bool) {
	userID, ok := q.userNameID[name]
	return userID, ok
}

// 引数の名前（Case を無視）をもつユーザーの ID を取得
func (q *QUsers) GetUserIDCaseInsentive(name string) (string, bool) {
	userID, ok := q.userLowerNameID[strings.ToLower(name)]
	return userID, ok
}

// 引数の ID をもつユーザーの現在のデータを取得
func (q *QUsers) GetUser(id string) (traq.User, bool) {
	user, ok := q.userIDData[id]
	return user, ok
}

// 引数の ID をもつユーザーの名前を取得
func (q *QUsers) GetUserName(id string) (string, bool) {
	user, ok := q.GetUser(id)
	if !ok {
		return "", false
	}
	return user.Name, true
}

// 引数の名前（完全一致）を持つユーザーの現在のデータを取得
func (q *QUsers) GetUserByName(name string) (traq.User, bool) {
	userID, ok := q.GetUserID(name)
	if !ok {
		return traq.User{}, false
	}
	return q.GetUser(userID)
}

// 引数の名前（Case を無視）を持つユーザーの現在のデータを取得
func (q *QUsers) GetUserByCaseInsentiveName(name string) (traq.User, bool) {
	userID, ok := q.GetUserIDCaseInsentive(name)
	if !ok {
		return traq.User{}, false
	}
	return q.GetUser(userID)
}
