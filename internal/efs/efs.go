package efs

import (
	"embed"
	"fmt"
)

//go:embed all:*
var FS embed.FS

func GetIcon() []byte {
	b, err := FS.ReadFile("favicon.ico")
	if err != nil {
		fmt.Print(err)
	}
	return b
}

func FileExists(path string) bool {
	_, err := FS.ReadFile(path)
	return err == nil
}
