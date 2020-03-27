package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func (h *holder) ws(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("could not start dev-websocket:", err)
		return
	}
	defer connection.Close()

	h.devCon = connection

	for {
		// mt, message, err := connection.ReadMessage()
		_, message, err := connection.ReadMessage()
		if err != nil {
			// log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		// err = connection.WriteMessage(mt, message)
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }
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

func findChangedSlide(h *holder) string {
	prev := h.slides
	h.parse()

	for i, s := range prev {
		if len(h.slides)-1 >= i {
			if s.hash != h.slides[i].hash {
				return strconv.Itoa(i)
			}
		}
	}
	// no changes found so far, if list length has changed,
	// we predict one slide is added to the end
	// (otherweise a hash check should find the change)
	if len(h.slides) != len(prev) {
		return strconv.Itoa(len(h.slides) - 1)
	}
	return "-"
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
		if h.devCon != nil {

			jsonPayload := make(map[string]string)
			jsonPayload["do"] = "reload"
			jsonPayload["slide"] = findChangedSlide(h)
			fmt.Println("changed slide: ", jsonPayload["slide"])

			h.devCon.WriteJSON(jsonPayload)
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
