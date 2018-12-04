package util

import (
	"fmt"
	"os"
	"unsafe"
)

// ExecAfterProcess 実行後処理
func ExecAfterProcess(outputByte []byte, err error, c chan string) {
	if err != nil {
		fmt.Println("Error.")
		fmt.Println(err)
		os.Exit(0)
	}

	// byte配列を文字列型に変換
	output := *(*string)(unsafe.Pointer(&outputByte))

	// 受信側に文字列を送信
	c <- output
}
