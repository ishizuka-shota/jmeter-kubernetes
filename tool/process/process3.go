package process

import (
	"fmt"
	"jmeter-kubernetes/tool/util"
	"os/exec"
	"strings"
)

// Execjmeter 処理番号3 : jmeter起動
func Execjmeter() {
	c := make(chan string, 2)

	fmt.Println("--------------------masterPodの選択--------------------")
	// Pod一覧取得処理
	go getPods(c)
	// 実行処理演出
	util.Kurukuru("Pod一覧を取得中", c)

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
	go copyjmx(jmxPath, kubePods[podNumber-1], c)
	util.Kurukuru("jmxファイルをPodへコピーしています", c)
	// 処理を止める
	<-c

	fmt.Println("--------------------jmeter起動--------------------")
	// jmeter起動

}

// getPods Pod一覧を取得
func getPods(c chan string) {
	// Pod一覧を取得(byte配列)
	outputByte, err := exec.Command("kubectl", "get", "pods", "-o", "custom-columns=:metadata.name", "-n", jmeterSlave).Output()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}

// copyjmx jmxファイルをコンテナへコピー
func copyjmx(path string, kubePod string, c chan string) {
	kubejmPath := kubePod + ":/jmeter/bin/"

	// jmxファイルをコンテナへコピー
	outputByte, err := exec.Command("kubectl", "cp", path, kubejmPath, "-n", jmeterSlave).Output()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}

// startjmeter jmeter開始
func startjmeter(c chan string) {
	// jmeter開始
	outputByte, err := exec.Command("kubectl", "exec", "-it", "/jmeter/bin/jmeter -n -t /jmeter/bin/lavoro.jmx -l /jmeter/bin/lavoro.jtl -R").Output()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}
