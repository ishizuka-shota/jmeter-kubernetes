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
	// kubectl exec -ti -n jmgrafana jmeter-grafana-dep-678c9cdb7c-4hl2g -- curl "http://admin:admin@127.0.0.1:3000/api/datasources" -X POST -H "Content-Type: application/json;charset=UTF-8" --data-binary "{\"name\":\"jmeterdb\",\"type\":\"influxdb\",\"url\":\"http://jmeter-influxdb:8086\",\"access\":\"proxy\",\"isDefault\":true,\"database\":\"jmeter\",\"user\":\"admin\",\"password\":\"admin\"}"
}
