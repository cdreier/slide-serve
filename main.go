package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type holder struct {
	dir        string
	demo       bool
	title      string
	slides     []slide
	styles     string
	dev        bool
	connection *websocket.Conn
}

func main() {
	port := flag.String("port", "8080", "http port the server is starting on")
	rootDir := flag.String("dir", "example", "root dir of your presentation")
	title := flag.String("title", "Slide", "html title")
	devMode := flag.Bool("dev", false, "dev true to start a filewatcher and reload the edited slide")
	flag.Parse()

	isDemo := false

	if *rootDir == "example" && !dirExist(*rootDir) {
		isDemo = true
		*devMode = false
		*title = "Slide"
	} else if !dirExist(*rootDir) {
		log.Fatal("cannot find root directory :(")
	}

	h := holder{
		dir:   *rootDir,
		title: *title,
		dev:   *devMode,
		demo:  isDemo,
	}

	h.parse()

	if *devMode {
		http.HandleFunc("/ws", h.ws)
		go h.startFileWatcher(*rootDir)
	}

	http.HandleFunc("/", h.handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(*rootDir))))
	http.HandleFunc("/favicon.ico", h.na)
	fmt.Println("starting on port: " + *port + " for directory " + *rootDir)
	http.ListenAndServe(":"+*port, nil)
}

func (h *holder) na(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func isDir(dir string) bool {
	stat, _ := os.Stat(dir)
	return stat.IsDir()
}

func dirExist(dir string) bool {
	_, err := os.Stat(dir)
	if err == nil {
		return true
	}
	return !os.IsNotExist(err)
}
