package process

import (
	"fmt"
	"jmeter-kubernetes/tool/util"
	"os/exec"
)

// CreateDepoymentAndService デプロイメント・サービス作成
func CreateDepoymentAndService(fileName string, c chan string) {
	filePath := fmt.Sprintf("../%s.yaml", fileName)

	// デプロイメント・サービス作成
	outputByte, err := exec.Command("kubectl", "apply", "-f", filePath).CombinedOutput()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}

// CreateNamespace namespace作成
func CreateNamespace(namespace string, c chan string) {
	// namespace作成
	outputByte, err := exec.Command("kubectl", "create", "namespace", namespace).CombinedOutput()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}

// DeleteDeploymentAndService デプロイメント・サービス削除
func DeleteDeploymentAndService(fileName string, c chan string) {
	filePath := fmt.Sprintf("../%s.yaml", fileName)

	// デプロイメント・サービス削除
	outputByte, err := exec.Command("kubectl", "delete", "-f", filePath).CombinedOutput()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}

// DeleteNamespace namespace削除
func DeleteNamespace(namespace string, c chan string) {
	// namespace削除
	outputByte, err := exec.Command("kubectl", "delete", "namespace", namespace).CombinedOutput()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}

// GetPods Pod一覧を取得
func GetPods(namespace string, c chan string) {
	// Pod一覧を取得(byte配列)
	outputByte, err := exec.Command("kubectl", "get", "pods", "-o", "custom-columns=:metadata.name", "-n", namespace).CombinedOutput()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}

// GetPodsIP PodのIP一覧を取得
func GetPodsIP(namespace string, c chan string) {
	// PodのIP一覧を取得(byte配列)
	outputByte, err := exec.Command("kubectl", "get", "endpoints", "-o=jsonpath=\"{range .items[*]}{range .subsets[*]}{range .addresses[*]}{.ip}{'\\n'}{end}\"", "-n", namespace).CombinedOutput()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)
}
