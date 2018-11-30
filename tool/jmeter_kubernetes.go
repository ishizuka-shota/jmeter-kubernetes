package main

import (
	"fmt"
	"github.com/tj/go-spin"
	"io"
	"jmeter-kubernetes/tool/process"
	"jmeter-kubernetes/tool/util"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	jmeterSlave  = "jmslave"
	jmeterMaster = "jmmaster"
)

func main() {
	// メニュー表示
	fmt.Println("実施したいタスクを選択")
	fmt.Println("cluster作成　[1]")
	fmt.Println("cluster削除　[2]")
	fmt.Println("jmeter実行　 [3]")
	fmt.Print(">> ")

	input := util.StrStdin()

	switch input {
	case "1":
		createKubernetesExecEnv()
	case "2":
		deleteKubernetesExecEnv()
	case "3":
		execjmeter()
	default:
		print("タスク番号を入力してください")
	}
}

// createCluster 処理番号1 : クラスタ作成
func createKubernetesExecEnv() {
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
	go process.CreateCluster(c)
	// 実行処理演出
	kurukuru("クラスタを作成中", c)
	// 処理を止める
	<-c

	fmt.Println("--------------------クラスタ認証--------------------")
	// コンテナの認証情報を取得
	go func(chan string) {
		outputByte, err := exec.Command("gcloud", "container", "clusters", "get-credentials", "jmeter").Output()
		util.ExecAfterProcess(outputByte, err, c)
	}(c)
	// 実行処理演出
	kurukuru("クラスタ認証中", c)
	// 処理を止める
	<-c

	fmt.Println("--------------------namespace作成--------------------")
	// namespaceを作成
	go func(chan string) {
		outputByte, err := exec.Command("kubectl", "create", "namespace", "jmslave").Output()
		util.ExecAfterProcess(outputByte, err, c)
	}(c)
	// 実行処理演出
	kurukuru("namespace作成中", c)
	// 処理を止める
	<-c

	fmt.Println("--------------------deployment作成--------------------")
	// デプロイメント作成
	go process.CreateDeployment(c)
	// 実行処理演出
	kurukuru("デプロイメントを作成中", c)
}

// deleteCluster 処理番号2 : クラスタ削除
func deleteKubernetesExecEnv() {
	c := make(chan string, 2)

	fmt.Println("--------------------deployment削除--------------------")
	// デプロイメント削除
	go process.DeleteDeployment(c)
	// 実行処理演出
	kurukuru("デプロイメントを削除中", c)
	// 処理を止める
	<-c

	fmt.Println("--------------------cluster削除--------------------")
	// クラスタ削除
	go process.DeleteCluster(c)
	// 実行処理演出
	kurukuru("クラスタを削除中", c)
}

// execjmeter 処理番号3 : jmeter起動
func execjmeter() {
	c := make(chan string, 2)

	fmt.Println("--------------------masterPodの選択--------------------")
	// Pod一覧取得処理
	go process.GetPods(c)
	// 実行処理演出
	kurukuru("Pod一覧を取得中", c)

	// 受信
	kubePodsString := <-c

	// 改行でスライスし、配列の中身がブランクのものをすべて取り除く
	kubePods := util.GetSliceNotBlank(strings.Split(kubePodsString, "\n"))

	// Pod一覧出力
	for i, kubePod := range kubePods {
		fmt.Printf("%s  [%d]\n", kubePod, i+1)
	}

	fmt.Println("masterとして使用するPodを選択")
	podNumber := util.IntStdin()

	fmt.Println("--------------------使用するjmxファイルのコピー--------------------")
	fmt.Println("jmxファイルのパスを入力")
	jmxPath := util.StrStdin()

	// jmxファイルコピー処理
	go process.Copyjmx(jmxPath, kubePods[podNumber-1], c)
	kurukuru("jmxファイルをPodへコピーしています", c)
	// 処理を止める
	<-c

	fmt.Println("--------------------jmeter起動--------------------")
	// jmeter起動

}

// kurukuru コマンドライン実行中くるくる
func kurukuru(text string, c chan string) {

	s := spin.New()
	for {
		// channnelから値が送信されたかどうかを判定
		if len(c) > 0 {
			fmt.Printf("\r%s  完了", text)
			fmt.Println()
			fmt.Println()
			break
		}

		// くるくる
		fmt.Printf("\r%s %s", text, s.Next())
		time.Sleep(100 * time.Millisecond)
	}

}

// startjmeter jmeter開始
func startjmeter(c chan string) {
	// jmeter開始
	outputByte, err := exec.Command("kubectl", "exec", "-it", "/jmeter/bin/jmeter -n -t /jmeter/bin/lavoro.jmx -l /jmeter/bin/lavoro.jtl -R").Output()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}
