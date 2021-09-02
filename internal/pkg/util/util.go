package util

import (
	"fmt"
	"os"
	"path/filepath"
)

func ParseArryInterface(data interface{}, i int) {
	dataType := data.([]interface{})
	switch dataType[i].(type) {
	case int:
		fmt.Println("test ok")
	}

}

func GetFileDirSize(path string) (int64, error) {
	var size int64
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})

	return size, err
}
