package util

import (
	"bufio"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"strings"
)

var server, token string

// CreateKubernetesClient kubernetesAPI実行時クライアント作成
func CreateKubernetesClient() *kubernetes.Clientset {
	filePath := GetHomePath()

	f, err := os.Open(filePath + "/.kube/config")
	if err != nil {
		fmt.Fprintf(os.Stderr, "File %s could not read: %v\n", filePath, err)
		os.Exit(1)
	}

	// 関数終了時にファイル閉じる
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "server:") {
			first := strings.Index(scanner.Text(), "server:")
			server = scanner.Text()[first+8:]
		}
		if strings.Contains(scanner.Text(), "access-token:") {
			first := strings.Index(scanner.Text(), "access-token:")
			token = scanner.Text()[first+14:]
		}
	}
	if serr := scanner.Err(); serr != nil {
		fmt.Fprintf(os.Stderr, "File %s scan error: %v\n", filePath, err)
	}

	config, err := clientcmd.BuildConfigFromFlags(server, "")
	if err != nil {
		fmt.Print("サーバーが見つかりません。")
		os.Exit(0)
	}
	config.Insecure = true
	config.BearerToken = token
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Print("kubeクライアントの作成に失敗しました。")
		os.Exit(0)
	}
	return client
}
