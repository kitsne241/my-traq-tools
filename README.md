# my-traq-tools

[traQ](https://github.com/traPtitech/traQ) 向けの [Go WebSocket Bot ライブラリ](https://github.com/traPtitech/traq-ws-bot) を利用する際に、チャンネルやスタンプ、ユーザーなどの各種 ID を簡易的に扱うためのユーティリティ集です。

直近の情報をあらかじめメモリ上にキャッシュすることで、パス名や名前からすぐに UUID を引けるようにする機能を提供します。

```sh
go get github.com/kitsne241/my-traq-tools@latest
```

## パッケージ

それぞれのパッケージは独立しており、必要なものだけを選んで使うことができます。

- **qchannels**: チャンネルのパス と ID の相互変換、および親子関係の取得を行います。
- **qstamps**: スタンプの名前と ID の相互変換、およびメッセージに対する連続したスタンプの付与をサポートします。
- **qusers**: ユーザーの Display ID と ID の相互変換を行います。
- **qutils**: そのほかいろいろとお役立ちな関数を用意する予定です。

## 注意

複数のチャンネルにまたがって長期的に駆動する Bot に qchannels, qstamps, qusers を利用する場合は、定期的に Refresh を実行して最新の情報に追随させることが推奨されます。ただし、Refresh 処理はやや重い処理なので、秒単位などの高い頻度で繰り返し実行することは推奨されません。

qutils に用意されている埋め込み検出の実装は [オリジナル](https://github.com/traPtitech/traQ_S-UI/blob/master/src/lib/markdown/detector.ts) とはやや異なります（が、ほとんどの場合には問題なく機能するはずです）。

## 用例

### qchannels

```go
package main

import (
	"fmt"
	"github.com/kitsne241/my-traq-tools/qchannels"
	wsbot "github.com/traPtitech/traq-ws-bot"
)

func main() {
	bot, _ := wsbot.NewBot(&wsbot.Options{AccessToken: "ACCESS_TOKEN"})
	qc, _ := qchannels.New(bot) // 内部で Refresh を実行

	channel, _ := qc.GetChannelByPath("times/24/kitsnegra")
	json, _ := channel.MarshalJSON()
	fmt.Println(string(json))

	qc.Refresh() // 情報を更新
}
```

### qstamps

```go
package main

import (
	"fmt"
	"github.com/kitsne241/my-traq-tools/qstamps"
	wsbot "github.com/traPtitech/traq-ws-bot"
)

func main() {
	bot, _ := wsbot.NewBot(&wsbot.Options{AccessToken: "ACCESS_TOKEN"})
	qs, _ := qstamps.New(bot) // 内部で Refresh を実行

	id, _ := qs.GetStampID("LGTM")
	fmt.Println("スタンプ :LGTM: の UUID : ", id)

    // 指定のスタンプをつける
	qs.Stamp("MESSAGE_ID", "oisu-1", "oisu-2", "oisu-3", "oisu-4yoko")

	qs.Refresh() // 情報を更新
}
```

### qusers

```go
package main

import (
	"fmt"
	"github.com/kitsne241/my-traq-tools/qusers"
	wsbot "github.com/traPtitech/traq-ws-bot"
)

func main() {
	bot, _ := wsbot.NewBot(&wsbot.Options{AccessToken: "ACCESS_TOKEN"})
	qu, _ := qusers.New(bot) // 内部で Refresh を実行

	id, _ := qu.GetUserID("kitsne")
	fmt.Println("ユーザー @kitsne の UUID : ", id)

	qu.Refresh() // 情報を更新
}
```

## 依存関係

- [github.com/traPtitech/traq-ws-bot](https://github.com/traPtitech/traq-ws-bot)
- [github.com/traPtitech/go-traq](https://github.com/traPtitech/go-traq)
