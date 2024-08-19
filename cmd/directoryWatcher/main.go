package main

import (
	watcher2 "directoryWatcher/pkg/watcher"
	"flag"
	"log"
)

func main() {
	// getting cli arguments, example - go run .\main.go --type local --path ./ --server 127.0.0.1:8097
	wType := flag.String("type", "local", "valid watcher types are s3, local and hdfs")
	dirPath := flag.String("path", "./", "enter directory path to watch")
	serverIp := flag.String("server", "127.0.0.1:8097", "enter server ip address")
	flag.Parse()

	watcher, err := watcher2.NewWatcher(*wType, *dirPath)
	if err != nil {
		log.Fatalf("error while initializing the watcher...")
	}

	watcher.Start(*serverIp)
}
