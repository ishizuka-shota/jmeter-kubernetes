package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// StrStdin 標準入力から文字列取得
func StrStdin() (stringInput string) {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	stringInput = scanner.Text()

	stringInput = strings.TrimSpace(stringInput)
	return
}

// IntStdin 標準入力から整数取得
func IntStdin() (intInput int) {
	stringInput := StrStdin()
	intInput, err := strconv.Atoi(strings.TrimSpace(stringInput))
	if err != nil {
		fmt.Println("標準入力でエラーが発生しました。")
		os.Exit(0)
	}
	return
}
