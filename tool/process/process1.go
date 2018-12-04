package process

import (
	"fmt"
	"jmeter-kubernetes/tool/util"
	"os"
	"os/exec"
)

// CreateKubernetesExecEnv 処理番号1 : kubernetes環境作成
func CreateKubernetesExecEnv() {

	c := make(chan string, 2)

	fmt.Println("gcloudのゾーン設定をasia-northeast1-a(東京リージョン)に変更する必要があります。")
	fmt.Print("変更しますか？ (y/n) >> ")
	yesOrNo(c)

	fmt.Println("--------------------cluster作成--------------------")
	fmt.Print("クラスタサイズ入力 >> ")
	clusterSize := util.StrStdin()
	go createCluster(clusterSize, c) // クラスタ作成
	util.Kurukuru("クラスタ作成中", c)      // 実行処理演出
	<-c                              // 処理を止める

	fmt.Println("--------------------cluster認証--------------------")
	// コンテナの認証情報を取得
	go func(chan string) {
		outputByte, err := exec.Command("gcloud", "container", "clusters", "get-credentials", "jmeter").Output()
		util.ExecAfterProcess(outputByte, err, c)
	}(c)
	util.Kurukuru("クラスタ認証中", c) // 実行処理演出
	<-c                         // 処理を止める

	fmt.Println("--------------------jmslave作成--------------------")
	// namespaceを作成
	go func(chan string) {
		outputByte, err := exec.Command("kubectl", "create", "namespace", "jmslave").Output()
		util.ExecAfterProcess(outputByte, err, c)
	}(c)
	util.Kurukuru("namespace作成中", c) // 実行処理演出
	<-c                              // 処理を止める

	go createJmslave(c)                       // デプロイメント・サービス作成
	util.Kurukuru("deployment/service作成中", c) // 実行処理演出
	<-c                                       // 処理を止める

	fmt.Println("--------------------jmmaster作成--------------------")
	// namespaceを作成
	go func(chan string) {
		outputByte, err := exec.Command("kubectl", "create", "namespace", "jmmaster").Output()
		util.ExecAfterProcess(outputByte, err, c)
	}(c)
	util.Kurukuru("namespace作成中", c) // 実行処理演出
	<-c                              // 処理を止める

	go createJmmaster(c)                      // デプロイメント・サービス作成
	util.Kurukuru("deployment/service作成中", c) // 実行処理演出
}

// createCluster クラスタ作成
func createCluster(clusterSize string, c chan string) {
	numNodes := "--num-nodes=" + clusterSize

	// jmxファイルをコンテナへコピー
	outputByte, err := exec.Command("gcloud", "container", "clusters", "create", "jmeter", "--preemptible", "--machine-type=g1-small", numNodes,
		"--disk-size=10", "--zone=asia-northeast1-a", "--enable-basic-auth", "--issue-client-certificate", "--no-enable-ip-alias", "--metadata",
		"disable-legacy-endpoints=true").Output()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}

// createJmslave jmslave関係作成
func createJmslave(c chan string) {
	// jmslaveのデプロイメント・サービス作成
	outputByte, err := exec.Command("kubectl", "apply", "-f", "../jmeter-slave.yaml").Output()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}

// createJmmaster jmmaster関係作成
func createJmmaster(c chan string) {
	// jmslaveのデプロイメント・サービス作成
	outputByte, err := exec.Command("kubectl", "apply", "-f", "../jmeter-master.yaml").Output()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}

func yesOrNo(c chan string) {
	switch util.StrStdin() {
	case "y":
		go func(chan string) {
			outputByte, err := exec.Command("gcloud", "config", "set", "compute/zone", "asia-northeast1-a").Output()
			util.ExecAfterProcess(outputByte, err, c)
		}(c)
		util.Kurukuru("設定の変更中", c) // 実行処理演出
		<-c                        // 処理を止める
	case "n":
		fmt.Println("処理を中断します。")
		os.Exit(0)
	default:
		fmt.Print("yもしくはnを入力してください >> ")
		yesOrNo(c)
	}
}
