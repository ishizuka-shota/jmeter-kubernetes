package process

import (
	"fmt"
	"io"
	"jmeter-kubernetes/tool/util"
	"os/exec"
)

const (
	jmeterSlave  = "jmslave"
	jmeterMaster = "jmmaster"
)

// DeleteKubernetesExecEnv 処理番号2 : クラスタ削除
func DeleteKubernetesExecEnv() {
	c := make(chan string, 2)

	fmt.Println("--------------------deployment・service削除--------------------")
	// デプロイメント削除
	go deleteDeploymentAndService(c)
	// 実行処理演出
	util.Kurukuru("デプロイメント・サービスを削除中", c)
	// 処理を止める
	<-c

	fmt.Println("--------------------cluster削除--------------------")
	// クラスタ削除
	go deleteCluster(c)
	// 実行処理演出
	util.Kurukuru("クラスタを削除中", c)
}

// deleteDeploymentAndService デプロイメント・サービス削除
func deleteDeploymentAndService(c chan string) {
	// デプロイメント・サービス削除
	outputByte, err := exec.Command("kubectl", "delete", "-f", "../jmeter-slave.yaml").Output()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}

// deleteCluster クラスタ削除
func deleteCluster(c chan string) {
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
