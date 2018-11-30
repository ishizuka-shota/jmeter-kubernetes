package process

import (
	"io"
	"jmeter-kubernetes/tool/util"
	"os/exec"
)

// DeleteDeployment デプロイメント削除
func DeleteDeployment(c chan string) {
	// デプロイメント削除
	outputByte, err := exec.Command("kubectl", "delete", "-f", "../jmeter-slave.yaml").Output()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}

// DeleteCluster クラスタ削除
func DeleteCluster(c chan string) {
	// jmxファイルをコンテナへコピー
	cmd := exec.Command("gcloud", "container", "clusters", "delete", "jmeter")

	// 標準入力
	stdin, _ := cmd.StdinPipe()
	io.WriteString(stdin, "y")
	stdin.Close()

	outputByte, err := cmd.Output()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)

}
