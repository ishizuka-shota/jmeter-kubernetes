package process

import (
	"fmt"
	"jmeter-kubernetes/tool/util"
	"os"
	"os/exec"
)

// CreateKubernetesExecEnv 処理番号1 : kubernetes環境作成
func CreateKubernetesExecEnv() {
	fmt.Println("gcloudのゾーン設定をasia-northeast1-a(東京リージョン)に変更する必要があります。")
	fmt.Println("変更しますか？ (y/n)")

	switch util.StrStdin() {
	case "y":
		fmt.Println("設定を変更しています...")
		_, err := exec.Command("gcloud", "config", "set", "compute/zone", "asia-northeast1-a").Output()
		if err != nil {
			fmt.Println("処理を中断します。")
			os.Exit(0)
		}
		fmt.Println("設定の変更が完了")
	case "n":
		fmt.Println("処理を中断します。")
		os.Exit(0)
	default:
		fmt.Println("yもしくはnを入力してください")
	}

	c := make(chan string, 2)

	fmt.Println("--------------------cluster作成--------------------")
	// クラスタ作成
	go createCluster(c)
	// 実行処理演出
	util.Kurukuru("クラスタを作成中", c)
	// 処理を止める
	<-c

	fmt.Println("--------------------クラスタ認証--------------------")
	// コンテナの認証情報を取得
	go func(chan string) {
		outputByte, err := exec.Command("gcloud", "container", "clusters", "get-credentials", "jmeter").Output()
		util.ExecAfterProcess(outputByte, err, c)
	}(c)
	// 実行処理演出
	util.Kurukuru("クラスタ認証中", c)
	// 処理を止める
	<-c

	fmt.Println("--------------------namespace作成--------------------")
	// namespaceを作成
	go func(chan string) {
		outputByte, err := exec.Command("kubectl", "create", "namespace", "jmslave").Output()
		util.ExecAfterProcess(outputByte, err, c)
	}(c)
	// 実行処理演出
	util.Kurukuru("namespace作成中", c)
	// 処理を止める
	<-c

	fmt.Println("--------------------deployment・service作成--------------------")
	// デプロイメント・サービス作成
	go createDeploymentAndService(c)
	// 実行処理演出
	util.Kurukuru("デプロイメント・サービスを作成中", c)
}

// createCluster クラスタ作成
func createCluster(c chan string) {
	// jmxファイルをコンテナへコピー
	outputByte, err := exec.Command("gcloud", "container", "clusters", "create", "jmeter", "--preemptible", "--machine-type=g1-small", "--num-nodes=3",
		"--disk-size=10", "--zone=asia-northeast1-a", "--enable-basic-auth", "--issue-client-certificate", "--no-enable-ip-alias", "--metadata",
		"disable-legacy-endpoints=true").Output()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}

// createDeploymentAndService デプロイメント・サービス作成
func createDeploymentAndService(c chan string) {
	// デプロイメント・サービス作成
	outputByte, err := exec.Command("kubectl", "apply", "-f", "../jmeter-slave.yaml").Output()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}
