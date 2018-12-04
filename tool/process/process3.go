package process

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"jmeter-kubernetes/tool/util"
	"os"
	"os/exec"
	"strings"
)

// Execjmeter 処理番号3 : jmeter起動
func Execjmeter() {
	c := make(chan string, 2)

	fmt.Println("--------------------select the masterPod--------------------")
	go getPods(c)                 // Pod一覧取得処理
	util.Kurukuru("Pod一覧を取得中", c) // 実行処理演出

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
	jmxFile := util.StrStdin()

	go copyjmx(jmxFile, kubePods[podNumber-1], c) // jmxファイルコピー処理
	util.Kurukuru("jmxファイルをPodへコピーしています", c)      // 実行処理演出
	<-c                                           // 処理を止める

	fmt.Println("--------------------start jmeter--------------------")
	go func(chan string) {
		outputByte, err := exec.Command("kubectl", "get", "endpoints", "-o=jsonpath=\"{range .items[*]}{range .subsets[*]}{range .addresses[*]}{.ip}{'\\n'}{end}\"", "-n", jmeterSlave).CombinedOutput()
		util.ExecAfterProcess(outputByte, err, c)
	}(c)
	util.Kurukuru("PodのIPを取得しています", c) // 実行処理演出
	ipListString := <-c

	ipList := util.GetSliceNotBlank(strings.Split(ipListString, "\n"))

	for i, ip := range ipList {
		ipList[i] = ip + ":1099"
	}

	startjmeter(jmxFile, strings.Join(ipList, ","), kubePods[podNumber-1], c) // jmeter開始
	<-c

}

// getPods Pod一覧を取得
func getPods(c chan string) {
	// Pod一覧を取得(byte配列)
	outputByte, err := exec.Command("kubectl", "get", "pods", "-o", "custom-columns=:metadata.name", "-n", jmeterMaster).CombinedOutput()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}

// copyjmx jmxファイルをコンテナへコピー
func copyjmx(jmxFile string, kubePod string, c chan string) {
	jmxPath := jmxFile + ".jmx"
	kubejmPath := kubePod + ":/jmeter/bin/"

	// jmxファイルをコンテナへコピー
	outputByte, err := exec.Command("kubectl", "cp", jmxPath, kubejmPath, "-n", jmeterMaster).CombinedOutput()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}

// startjmeter jmeter開始
func startjmeter(jmxFile string, ipList string, kubePod string, c chan string) {
	jmeterCmd := "\"jmeter -n -t /jmeter/bin/" + jmxFile + ".jmx -l /jmeter/bin/" + jmxFile + ".jtl -R" + ipList + "\""

	// jmeterコマンド
	cmd := exec.Command("kubectl", "exec", "-n", jmeterMaster, "-i", kubePod, "--", "/bin/sh", "-c", jmeterCmd)
	// jmeter開始
	runCommand(cmd)
}

func runCommand(cmd *exec.Cmd) {
	outReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}
	errReader, err := cmd.StderrPipe()
	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}

	var bufout, buferr bytes.Buffer
	outReader2 := io.TeeReader(outReader, &bufout)
	errReader2 := io.TeeReader(errReader, &buferr)

	if err = cmd.Start(); err != nil {
		fmt.Print(err)
		os.Exit(0)
	}

	go printOutputWithHeader(outReader2)
	go printOutputWithHeader(errReader2)

	err = cmd.Wait()
	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}

	fmt.Println("処理が正常終了しました")
}

func printOutputWithHeader(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Text())
	}
}
