package watcher

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// DirWatcher - an interface for directory watchers
type DirWatcher interface {
	Watch(dirPath []string, clients *map[*websocket.Conn]bool)
	notifyClients(filename string, clients *map[*websocket.Conn]bool)
}

// Watcher - Main manager in application
type Watcher struct {
	server      *http.Server
	upgrader    *websocket.Upgrader
	dirWatcher  DirWatcher
	clients     map[*websocket.Conn]bool
	dirsToWatch []string

	mu *sync.Mutex
}

func NewWatcher(watcherType string, dirPath string) (*Watcher, error) {
	var dirWatcher DirWatcher
	var mu sync.Mutex
	dirsToWatch := make([]string, 0)

	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	switch strings.ToLower(watcherType) {
	case "s3":
		dirWatcher = newS3Watcher()
	case "local":
		dirWatcher = newLocalWatcher()

		err := filepath.Walk(dirPath,
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
	case "hdfs":
		dirWatcher = newHDFSWatcher()
	}

	return &Watcher{
		upgrader:    upgrader,
		dirWatcher:  dirWatcher,
		clients:     make(map[*websocket.Conn]bool),
		mu:          &mu,
		server:      nil,
		dirsToWatch: dirsToWatch,
	}, nil
}

func (wc *Watcher) Start(serverIp string) {
	go wc.dirWatcher.Watch(wc.dirsToWatch, &wc.clients)
	log.Println("started watching the directory...")

	http.HandleFunc("/", wc.handleFunc)

	err := http.ListenAndServe(serverIp, nil)
	log.Println("started the web server...")
	if err != nil {
		log.Println("error while starting the server...")
	}
}

func (wc *Watcher) handleFunc(w http.ResponseWriter, r *http.Request) {
	conn, err := wc.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error while reading message...")
	} else {
		log.Println("started the handler...")
	}

	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("error while closing the connection")
		}
	}(conn)

	wc.mu.Lock()
	wc.clients[conn] = true
	wc.mu.Unlock()
	wc.websocketHandler(conn)
}

func (wc *Watcher) websocketHandler(conn *websocket.Conn) {
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("client disconnected...", err)
			wc.mu.Lock()
			delete(wc.clients, conn)
			wc.mu.Unlock()
			break
		}
	}
}
