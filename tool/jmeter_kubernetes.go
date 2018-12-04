package main

import (
	"fmt"
	"jmeter-kubernetes/tool/process"
	"jmeter-kubernetes/tool/util"
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

	switch util.StrStdin() {
	case "1":
		process.CreateKubernetesExecEnv()
	case "2":
		process.DeleteKubernetesExecEnv()
	case "3":
		process.Execjmeter()
	default:
		fmt.Println("タスク番号を入力してください")
	}
}
