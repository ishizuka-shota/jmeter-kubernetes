package util

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
)

// GetHomePath ホームディレクトリパス取得
func GetHomePath() (home string) {
	homepath, err := homedir.Dir()
	if err != nil {
		fmt.Println("ホームディレクトリの取得に失敗しました")
		os.Exit(0)
	}

	home, err = homedir.Expand(homepath)
	if err != nil {
		fmt.Println("ホームディレクトリの取得に失敗しました")
		os.Exit(0)
	}

	return
}
