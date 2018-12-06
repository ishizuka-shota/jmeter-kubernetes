package util

import (
	"fmt"
	"github.com/tj/go-spin"
	"time"
)

// Kurukuru コマンドライン実行中くるくる
func Kurukuru(text string, c chan string) {

	s := spin.New()
	for {
		// channnelから値が送信されたかどうかを判定
		if len(c) > 0 {
			fmt.Printf("\r%s  完了", text)
			fmt.Println()
			fmt.Println()
			break
		}

		// // くるくる
		fmt.Printf("\r%s %s", text, s.Next())
		time.Sleep(100 * time.Millisecond)
	}

}
