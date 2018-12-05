package util

import (
	"fmt"
	"os"
	"os/exec"
	"unsafe"
)

// ExecAfterProcess プロセス実行後処理
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

// ExecProcess プロセス実行処理
func ExecProcess(c chan string, arg ...string) {
	// コマンド作成
	cmd := exec.Cmd{
		Args: arg,
	}

	// プロセス実行
	outputByte, err := cmd.CombinedOutput()

	// 実行後処理
	ExecAfterProcess(outputByte, err, c)
}
