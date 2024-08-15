package main

import (
	watcher2 "directoryWatcher/pkg/watcher"
	"flag"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// getting cli arguments, example - go run .\main.go --type local --path ./ --server 127.0.0.1:8097
	wType := flag.String("type", "local", "valid watcher types are s3, local and hdfs")
	dirPath := flag.String("path", "./", "enter directory path to watch")
	serverIp := flag.String("server", "127.0.0.1:8097", "enter server ip address")
	flag.Parse()

	watcher, err := watcher2.NewWatcher(*wType)
	if err != nil {
		log.Fatalf("error while initializing the watcher...")
	}

	dirsToWatch := make([]string, 0)

	err = filepath.Walk(*dirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				dirsToWatch = append(dirsToWatch, path)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	watcher.Start(*serverIp, dirsToWatch)
}
