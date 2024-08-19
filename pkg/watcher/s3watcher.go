package watcher

import "github.com/gorilla/websocket"

type S3Watcher struct {
}

func newS3Watcher() *S3Watcher {
	return &S3Watcher{}
}

func (w *S3Watcher) Watch(dirPath []string, clients *map[*websocket.Conn]bool) {

}

func (w *S3Watcher) notifyClients(filename string, clients *map[*websocket.Conn]bool) {

}
