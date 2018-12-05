package process

import (
	"fmt"
	"io"
	"jmeter-kubernetes/tool/util"
	"os/exec"
)

// DeleteKubernetesExecEnv 処理番号2 : クラスタ削除
func DeleteKubernetesExecEnv() {
	c := make(chan string, 2)

	fmt.Println("--------------------jmslave削除--------------------")
	go deleteDeploymentAndService(util.JmeterSlave, c) // jmslave削除
	util.Kurukuru("deployment/service削除中", c)          // 実行処理演出
	<-c                                                // 処理を止める
	go deleteNamespace(util.JmeterSlave, c)            // namespaceを削除
	util.Kurukuru("namespace削除中", c)                   // 実行処理演出
	<-c                                                // 処理を止める

	fmt.Println("--------------------jmmaster削除--------------------")
	go deleteDeploymentAndService(util.JmeterMaster, c) // jmmaster削除
	util.Kurukuru("deployment/service削除中", c)           // 実行処理演出
	<-c                                                 // 処理を止める
	go deleteNamespace(util.JmeterMaster, c)            // namespaceを削除
	util.Kurukuru("namespace削除中", c)                    // 実行処理演出
	<-c                                                 // 処理を止める

	fmt.Println("--------------------jmgrafana削除--------------------")
	go deleteDeploymentAndService(util.JmeterGrafana, c) // jmgrafana削除
	util.Kurukuru("deployment/service削除中", c)            // 実行処理演出
	<-c                                                  // 処理を止める
	go deleteNamespace(util.JmeterGrafana, c)            // namespaceを削除
	util.Kurukuru("namespace削除中", c)                     // 実行処理演出
	<-c                                                  // 処理を止める

	fmt.Println("--------------------jminfluxdb削除--------------------")
	go deleteDeploymentAndService(util.JmeterInfluxdb, c) // jminfluxdb削除
	util.Kurukuru("deployment/service削除中", c)             // 実行処理演出
	<-c                                                   // 処理を止める
	go deleteNamespace(util.JmeterInfluxdb, c)            // namespaceを削除
	util.Kurukuru("namespace削除中", c)                      // 実行処理演出
	<-c                                                   // 処理を止める

	fmt.Println("--------------------cluster削除--------------------")
	go deleteCluster(c)         // クラスタ削除
	util.Kurukuru("クラスタ削除中", c) // 実行処理演出
}

// deleteNamespace namespace削除
func deleteNamespace(namespace string, c chan string) {
	// namespace削除
	outputByte, err := exec.Command("kubectl", "delete", "namespace", namespace).CombinedOutput()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}

// deleteDeploymentAndService デプロイメント・サービス削除
func deleteDeploymentAndService(fileName string, c chan string) {
	filePath := "../" + fileName + ".yaml"

	// デプロイメント・サービス削除
	outputByte, err := exec.Command("kubectl", "delete", "-f", filePath).CombinedOutput()

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

	outputByte, err := cmd.CombinedOutput()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)

}
