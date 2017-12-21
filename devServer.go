package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func (h *holder) ws(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer connection.Close()

	// add connection to pool
	h.connection = connection

	for {
		mt, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = connection.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func debounce(interval time.Duration, input chan string, cb func(arg string)) {
	var item string
	timer := time.NewTimer(interval)
	for {
		select {
		case item = <-input:
			timer.Reset(interval)
		case <-timer.C:
			if item != "" {
				cb(item)
			}
		}
	}
}

func (h *holder) startFileWatcher(dir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	eventChan := make(chan string)
	go debounce(time.Second, eventChan, func(name string) {
		fmt.Println("reloading... ", name)
<<<<<<< HEAD
		if h.connection != nil {
			h.connection.WriteMessage(websocket.TextMessage, []byte("reload!"))
=======
		for _, wsConnection := range h.connections {
			wsConnection.WriteMessage(websocket.TextMessage, []byte("reload!"))
>>>>>>> ebb473eb4f1963f2282ab9020ba8ce2827de76c1
		}
	})

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				eventChan <- event.Name
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
