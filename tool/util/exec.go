package util

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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

// ExecProcessRealTimeLog プロセス実行処理(リアルタイムなログ)
func ExecProcessRealTimeLog(cmd *exec.Cmd) {
	outReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}
	errReader, err := cmd.StderrPipe()
	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}

	var bufout, buferr bytes.Buffer
	outReader2 := io.TeeReader(outReader, &bufout)
	errReader2 := io.TeeReader(errReader, &buferr)

	if err = cmd.Start(); err != nil {
		fmt.Print(err)
		os.Exit(0)
	}

	go printOutputWithHeader(outReader2)
	go printOutputWithHeader(errReader2)

	err = cmd.Wait()
	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}
}

func printOutputWithHeader(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Text())
	}
}
