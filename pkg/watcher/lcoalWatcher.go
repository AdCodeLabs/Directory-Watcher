package watcher

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type LocalWatcher struct {
	mu *sync.Mutex
}

func newLocalWatcher() *LocalWatcher {
	var mu sync.Mutex
	return &LocalWatcher{
		mu: &mu,
	}
}

func (w *LocalWatcher) Watch(dirPath []string, clients *map[*websocket.Conn]bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {

		}
	}(watcher)

	for _, dir := range dirPath {
		err = watcher.Add(dir)
		if err != nil {
			log.Fatal(err)
		}
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				log.Println("New file created:", event.Name)
				w.notifyClients(fmt.Sprintf("New file created: %s", event.Name), clients)
			} else if event.Op&fsnotify.Remove == fsnotify.Remove {
				log.Println("A file removed:", event.Name)
				w.notifyClients(fmt.Sprintf("A file removed: %s", event.Name), clients)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error:", err)
		}
	}
}

func (w *LocalWatcher) notifyClients(message string, clients *map[*websocket.Conn]bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for client := range *clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("write error: %v", err)
			err := client.Close()
			if err != nil {
				return
			}
			delete(*clients, client)
		}
	}
}
