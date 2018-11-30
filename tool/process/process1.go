package process

import (
	"jmeter-kubernetes/tool/util"
	"os/exec"
)

// CreateCluster クラスタ作成
func CreateCluster(c chan string) {
	// jmxファイルをコンテナへコピー
	outputByte, err := exec.Command("gcloud", "container", "clusters", "create", "jmeter", "--preemptible", "--machine-type=g1-small", "--num-nodes=3",
		"--disk-size=10", "--zone=asia-northeast1-a", "--enable-basic-auth", "--issue-client-certificate", "--no-enable-ip-alias", "--metadata",
		"disable-legacy-endpoints=true").Output()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)

}

// CreateDeployment デプロイメント作成
func CreateDeployment(c chan string) {
	// デプロイメント作成
	outputByte, err := exec.Command("kubectl", "apply", "-f", "../jmeter-slave.yaml").Output()

	// 実行後処理
	util.ExecAfterProcess(outputByte, err, c)

}
