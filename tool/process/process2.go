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

	fmt.Println("--------------------jmslave削除--------------------")
	go deleteJmslave(c)                       // jmslave削除
	util.Kurukuru("deployment/service削除中", c) // 実行処理演出
	<-c                                       // 処理を止める

	fmt.Println("--------------------jmmaster削除--------------------")
	go deleteJmmaster(c)                      // jmmaster削除
	util.Kurukuru("deployment/service削除中", c) // 実行処理演出
	<-c                                       // 処理を止める

	fmt.Println("--------------------cluster削除--------------------")
	go deleteCluster(c)         // クラスタ削除
	util.Kurukuru("クラスタ削除中", c) // 実行処理演出
}

// deleteJmslave jmslave関係削除
func deleteJmslave(c chan string) {
	// jmslaveのデプロイメント・サービス削除
	outputByte, err := exec.Command("kubectl", "delete", "-f", "../jmeter-slave.yaml").Output()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}

// deleteJmmaster jmmaster関係削除
func deleteJmmaster(c chan string) {
	// jmmasterのデプロイメント・サービス削除
	outputByte, err := exec.Command("kubectl", "delete", "-f", "../jmeter-master.yaml").Output()

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
