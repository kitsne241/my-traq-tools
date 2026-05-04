package qutils

import (
	"encoding/json"
	"slices"
)

type Embed struct {
	EmbedContent
	Start int // 埋め込み前のテキストにおける埋め込みの開始位置
	End   int // 埋め込み前のテキストにおける埋め込みの終了位置
}

type EmbedContent struct {
	Type string `json:"type"` // user, group, channel のいずれか
	Raw  string `json:"raw"`  // 任意文字列で埋め込みの内容
	ID   string `json:"id"`   // UUID
}

// メッセージ本文を埋め込みのないもとの Markdown の形式に変換し、埋め込みの内容を返す
func Unembed(text string) (string, []Embed) {
	textRune := []rune(text)
	inEmbed := false
	embeds := []Embed{}
	var data Embed

	// !{ ... } の形式を最短一致で検出し、中身を JSON としてパース
	for i := 0; i < len(textRune); i++ {
		if inEmbed {
			if textRune[i] == '}' {
				inEmbed = false
				data.End = i + 1

				var content EmbedContent
				err := json.Unmarshal([]byte(string(textRune[data.Start+1:i+1])), &content)
				if err == nil {
					data.EmbedContent = content
					embeds = append(embeds, data)
				}
			}
		} else {
			if (i < len(textRune)-1) && textRune[i] == '!' && textRune[i+1] == '{' {
				inEmbed = true
				data = Embed{Start: i}
			}
		}
	}

	// 得られた embed を後ろから順に置き換えて埋め込みを解消する
	slices.Reverse(embeds)
	for _, data := range embeds {
		tempRune := append([]rune(data.Raw), textRune[data.End:]...)
		textRune = append(textRune[:data.Start], tempRune...)
	}

	slices.Reverse(embeds)
	return string(textRune), embeds
}

// 埋め込みは type, raw, id の 3 つのキーのみから構成される JSON 文字列 !{ ... } である
// type は user, group, channel のいずれかで、id は UUID 形式の文字列であり、いずれも波括弧を含まないことが保証されている
// raw は自由にかけるので波括弧を含められるが、ふつうの利用においては含まないとしてよいはず
