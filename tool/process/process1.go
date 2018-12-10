package process

import (
	"fmt"
	"jmeter-kubernetes/tool/util"
	"os/exec"
	"strings"
)

// CreateKubernetesExecEnv 処理番号1 : kubernetes環境作成
func CreateKubernetesExecEnv() {

	c := make(chan string, 2)

	fmt.Println("--------------------cluster作成--------------------")
	go createCluster(c)         // クラスタ作成
	util.Kurukuru("クラスタ作成中", c) // 実行処理演出
	<-c                         // 処理を止める

	fmt.Println("--------------------cluster認証--------------------")
	// コンテナの認証情報を取得
	go func(chan string) {
		outputByte, err := exec.Command("gcloud", "container", "clusters", "get-credentials", "jmeter").CombinedOutput()
		util.ExecAfterProcess(outputByte, err, c)
	}(c)
	util.Kurukuru("クラスタ認証中", c) // 実行処理演出
	<-c                         // 処理を止める

	fmt.Println("--------------------jmslave作成--------------------")
	go CreateNamespace(util.JmeterSlave, c) // namespaceを作成
	util.Kurukuru("namespace作成中", c)        // 実行処理演出
	<-c                                     // 処理を止める

	fmt.Print("slaveサイズ入力 >> ")
	podSize := util.StrStdin()
	go CreateDepoymentAndService(util.JmeterSlave, c) // jmslave作成
	util.Kurukuru("deployment/service作成中", c)         // 実行処理演出
	<-c                                               // 処理を止める
	go scalePodJmslave(podSize, c)                    // podサイズ調整
	util.Kurukuru("podサイズ調整中", c)                     // 実行処理演出
	<-c                                               // 処理を止める

	fmt.Println("--------------------jmmaster作成--------------------")
	go CreateNamespace(util.JmeterMaster, c)           // namespaceを作成
	util.Kurukuru("namespace作成中", c)                   // 実行処理演出
	<-c                                                // 処理を止める
	go CreateDepoymentAndService(util.JmeterMaster, c) // jmmaster作成
	util.Kurukuru("deployment/service作成中", c)          // 実行処理演出
	<-c                                                // 処理を止める

	fmt.Println("--------------------jminfluxdb作成--------------------")
	go CreateNamespace(util.JmeterInfluxdb, c)                     // namespaceを作成
	util.Kurukuru("namespace作成中", c)                               // 実行処理演出
	<-c                                                            // 処理を止める
	go CreateDepoymentAndService(util.JmeterInfluxdb, c)           // jminfluxdb作成
	util.Kurukuru("deployment/service作成中", c)                      // 実行処理演出
	<-c                                                            // 処理を止める
	go GetPods(util.JmeterInfluxdb, c)                             //influxdbのpod取得
	util.Kurukuru("InfluxdbのPodを取得中", c)                           // 実行処理演出
	influxdbPod := util.GetSliceNotBlank(strings.Split(<-c, "\n")) // Pod名取得
	go createDATABASE(influxdbPod[0], c)                           // データベース構築
	util.Kurukuru("DATABASE構築中", c)                                // 実行処理演出
	<-c                                                            // 処理を止める

	fmt.Println("--------------------jmgrafana作成--------------------")
	go CreateNamespace(util.JmeterGrafana, c)           // namespaceを作成
	util.Kurukuru("namespace作成中", c)                    // 実行処理演出
	<-c                                                 // 処理を止める
	go CreateDepoymentAndService(util.JmeterGrafana, c) // jmgrafana作成
	util.Kurukuru("deployment/service作成中", c)           // 実行処理演出
	<-c                                                 // 処理を止める

	close(c) // channelを閉じる
}

// createCluster クラスタ作成
func createCluster(c chan string) {
	// jmxファイルをコンテナへコピー
	outputByte, err := exec.Command("gcloud", "container", "clusters", "create", "jmeter", "--machine-type=n1-highmem-4", "--num-nodes=1",
		"--disk-size=10", "--zone=asia-northeast1-a", "--enable-basic-auth", "--issue-client-certificate", "--no-enable-ip-alias", "--metadata",
		"disable-legacy-endpoints=true").CombinedOutput()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}

// createDATABASE データベース作成
func createDATABASE(influxdbPod string, c chan string) {
	createDatabaseCmd := "CREATE DATABASE jmeter"

	// データベース作成
	outputByte, err := exec.Command("kubectl", "exec", "-n", util.JmeterInfluxdb, "-i", influxdbPod, "--", "influx", "-execute", createDatabaseCmd).CombinedOutput()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}

// scalePodJmslave jmslaveのポッドサイズのスケール
func scalePodJmslave(podSize string, c chan string) {
	replicas := fmt.Sprintf("--replicas=%s", podSize)

	// デプロイメント・サービス作成
	outputByte, err := exec.Command("kubectl", "scale", replicas, "deployment/jmeter-slaves-dep", "-n", util.JmeterSlave).CombinedOutput()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}
