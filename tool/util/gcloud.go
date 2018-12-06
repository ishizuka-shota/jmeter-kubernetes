package util

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ChangeComputeZone ゾーン設定取得、必要であれば変更
func ChangeComputeZone(c chan string) {
	// gcloudのゾーン設定を取得
	outputByte, err := exec.Command("gcloud", "config", "get-value", "compute/zone").CombinedOutput()
	ExecAfterProcess(outputByte, err, c)
	zone := strings.TrimRight(<-c, "\n\r")

	if zone != Tokyo1a {
		fmt.Println("gcloudのゾーン設定をasia-northeast1-a(東京リージョン)に変更する必要があります。")
		fmt.Println("現在のゾーン設定は" + zone + "です。")
		fmt.Print("変更しますか？ (y/n) >> ")

		switch StrStdin() {
		case "y":
			go func(chan string) {
				outputByte, err := exec.Command("gcloud", "config", "set", "compute/zone", Tokyo1a).CombinedOutput()
				ExecAfterProcess(outputByte, err, c)
			}(c)
			Kurukuru("設定の変更中", c) // 実行処理演出
			<-c                   // 処理を止める
		default:
			fmt.Println("処理を中断します。")
			os.Exit(0)
		}
	}
}
