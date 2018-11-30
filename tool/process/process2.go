package process

import (
	"jmeter-kubernetes/tool/util"
	"os/exec"
)

const (
	jmeterSlave  = "jmslave"
	jmeterMaster = "jmmaster"
)

// GetPods Pod一覧を取得
func GetPods(c chan string) {
	// Pod一覧を取得(byte配列)
	outputByte, err := exec.Command("kubectl", "get", "pods", "-o", "custom-columns=:metadata.name", "-n", jmeterSlave).Output()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}

// Copyjmx jmxファイルをコンテナへコピー
func Copyjmx(path string, kubePod string, c chan string) {
	kubejmPath := kubePod + ":/jmeter/bin/"

	// jmxファイルをコンテナへコピー
	outputByte, err := exec.Command("kubectl", "cp", path, kubejmPath, "-n", jmeterSlave).Output()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}
