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
	go DeleteDeploymentAndService(util.JmeterSlave, c) // jmslave削除
	util.Kurukuru("deployment/service削除中", c)          // 実行処理演出
	<-c                                                // 処理を止める
	go DeleteNamespace(util.JmeterSlave, c)            // namespaceを削除
	util.Kurukuru("namespace削除中", c)                   // 実行処理演出
	<-c                                                // 処理を止める

	fmt.Println("--------------------jmmaster削除--------------------")
	go DeleteDeploymentAndService(util.JmeterMaster, c) // jmmaster削除
	util.Kurukuru("deployment/service削除中", c)           // 実行処理演出
	<-c                                                 // 処理を止める
	go DeleteNamespace(util.JmeterMaster, c)            // namespaceを削除
	util.Kurukuru("namespace削除中", c)                    // 実行処理演出
	<-c                                                 // 処理を止める

	fmt.Println("--------------------jminfluxdb削除--------------------")
	go DeleteDeploymentAndService(util.JmeterInfluxdb, c) // jminfluxdb削除
	util.Kurukuru("deployment/service削除中", c)             // 実行処理演出
	<-c                                                   // 処理を止める
	go DeleteNamespace(util.JmeterInfluxdb, c)            // namespaceを削除
	util.Kurukuru("namespace削除中", c)                      // 実行処理演出
	<-c                                                   // 処理を止める

	fmt.Println("--------------------jmgrafana削除--------------------")
	go DeleteDeploymentAndService(util.JmeterGrafana, c) // jmgrafana削除
	util.Kurukuru("deployment/service削除中", c)            // 実行処理演出
	<-c                                                  // 処理を止める
	go DeleteNamespace(util.JmeterGrafana, c)            // namespaceを削除
	util.Kurukuru("namespace削除中", c)                     // 実行処理演出
	<-c                                                  // 処理を止める

	fmt.Println("--------------------cluster削除--------------------")
	go deleteCluster(c)         // クラスタ削除
	util.Kurukuru("クラスタ削除中", c) // 実行処理演出

	close(c) // channelを閉じる
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
