package watcher

import "github.com/gorilla/websocket"

type HDFSWatcher struct {
}

func newHDFSWatcher() *HDFSWatcher {
	return &HDFSWatcher{}
}

func (w *HDFSWatcher) Watch(dirPath []string, clients *map[*websocket.Conn]bool) {

}

func (w *HDFSWatcher) notifyClients(filename string, clients *map[*websocket.Conn]bool) {

}
