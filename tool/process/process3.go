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

	fmt.Println("--------------------select the masterPod--------------------")
	go GetPods(util.JmeterMaster, c)    // MasterPod一覧取得処理
	util.Kurukuru("MasterPod一覧を取得中", c) // 実行処理演出

	kubePodsString := <-c // 受信

	// 改行でスライスし、配列の中身がブランクのものをすべて取り除く
	kubePods := util.GetSliceNotBlank(strings.Split(kubePodsString, "\n"))

	// Pod一覧出力
	for i, kubePod := range kubePods {
		fmt.Printf("%s  [%d]\n", kubePod, i+1)
	}

	fmt.Println("masterとして使用するPodを選択")
	fmt.Print(" >> ")
	podNumber := util.IntStdin()

	fmt.Println("--------------------copy jmx file--------------------")
	fmt.Print("jmxファイル名を入力(拡張子不要) >> ")
	file := util.StrStdin()

	go copyjmx(file, kubePods[podNumber-1], c) // jmxファイルコピー処理
	util.Kurukuru("jmxファイルをPodへコピー", c)        // 実行処理演出
	<-c                                        // 処理を止める

	fmt.Println("--------------------start jmeter--------------------")
	go GetPodsIP(util.JmeterSlave, c)     //SlavePodのIP一覧取得
	util.Kurukuru("SlavePodのIP一覧を取得中", c) // 実行処理演出
	ipListString := <-c

	ipList := util.GetSliceNotBlank(strings.Split(ipListString, "\n"))

	for i, ip := range ipList {
		ipList[i] = fmt.Sprintf("%s:1099", ip)
	}

	fmt.Println("jmeterによる負荷テストを開始")
	startjmeter(file, strings.Join(ipList, ","), kubePods[podNumber-1]) // jmeter開始
	fmt.Println("処理が正常終了しました")

	fmt.Println("--------------------get jtl file--------------------")
	go getjtl(file, kubePods[podNumber-1], c)
	util.Kurukuru("jtlファイルをローカルへコピー", c)

	close(c) // channelを閉じる
}

// copyjmx jmxファイルをコンテナへコピー
func copyjmx(jmxFile string, kubePod string, c chan string) {
	jmxPath := fmt.Sprintf("../jmx/%s.jmx", jmxFile)
	kubejmPath := fmt.Sprintf("%s:/jmeter/bin/", kubePod)

	// jmxファイルをコンテナへコピー
	outputByte, err := exec.Command("kubectl", "cp", jmxPath, kubejmPath, "-n", util.JmeterMaster).CombinedOutput()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}

// startjmeter jmeter開始
func startjmeter(jmxFile string, ipList string, kubePod string) {
	jmeterCmd := fmt.Sprintf("jmeter -n -t /jmeter/bin/%s.jmx -l /jmeter/bin/%s.jtl -R%s", jmxFile, jmxFile, ipList)

	// jmeterコマンド
	cmd := exec.Command("kubectl", "exec", "-n", util.JmeterMaster, "-i", kubePod, "--", "/bin/ash", "-c", jmeterCmd)

	// jmeter開始
	util.ExecProcessRealTimeLog(cmd)
}

// getjtl jtlファイル取得
func getjtl(file string, kubePod string, c chan string) {
	kubejmPath := fmt.Sprintf("%s:/jmeter/bin/%s.jtl", kubePod, file)

	// jmxファイルをコンテナへコピー
	outputByte, err := exec.Command("kubectl", "cp", kubejmPath, "../jtl", "-n", util.JmeterMaster).CombinedOutput()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)

}
