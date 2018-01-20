package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/websocket"
)

type holder struct {
	dir        string
	demo       bool
	title      string
	slides     []slide
	styles     string
	dev        bool
	pdfPrint   bool
	codeTheme  string
	connection *websocket.Conn
}

func main() {
	port := flag.String("port", "8080", "http port the server is starting on")
	rootDir := flag.String("dir", "example", "root dir of your presentation")
	title := flag.String("title", "Slide", "html title")
	devMode := flag.Bool("dev", false, "dev true to start a filewatcher and reload the edited slide")
	codeTheme := flag.String("syntaxhl", "monokai", "code highlighter theme")
	pdfPrint := flag.Bool("pdf", false, "printing a pdf")
	// control := flag.Bool("control", false, "attach controller with peer to peer ")
	flag.Parse()

	isDemo := false

	// if *control {
	// 	qrterminal.GenerateHalfBlock("http://drailing.net", qrterminal.L, os.Stdout)
	// }

	if *rootDir == "example" && !dirExist(*rootDir) {
		isDemo = true
		*devMode = false
		*title = "Slide"
	} else if !dirExist(*rootDir) {
		log.Fatal("cannot find root directory :(")
	}

	h := holder{
		dir:       *rootDir,
		title:     *title,
		dev:       *devMode,
		demo:      isDemo,
		codeTheme: *codeTheme,
		pdfPrint:  *pdfPrint,
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

func (h *holder) handler(w http.ResponseWriter, r *http.Request) {
	box := packr.NewBox("./www")
	t, _ := template.New("slide").Parse(box.String("slide.html"))

	slides := ""
	styles := h.styles

	for i, s := range h.slides {
		slides += renderSlide(s, i, h.codeTheme)

		if s.image != "" {
			styles += "\n"
			styles += addStyleRule(s.image, i)
		}

		if s.styles != "" {
			styles += "\n"
			slideStyle := strings.Replace(s.styles, "SLIDENUMBER", strconv.Itoa(i), -1)
			styles += slideStyle
		}

	}

	printStyles := box.String("summaryStyle.css")
	if h.pdfPrint {
		printStyles = box.String("pdfStyle.css")
	}

	s := slideContent{
		Slides:     template.HTML(slides),
		Styles:     template.CSS(styles),
		PrintStyle: template.CSS(printStyles),
		Title:      h.title,
	}

	if h.dev {
		js, _ := template.New("devmode").Parse(box.String("devMode.html"))
		var buf bytes.Buffer
		data := make(map[string]string)
		data["url"] = "ws://" + r.Host + "/ws"
		js.Execute(&buf, data)
		s.DevMode = template.HTML(buf.String())
	}

	t.Execute(w, s)
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
