package main

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/markbates/pkger"
	"github.com/urfave/cli"
)

type holder struct {
	dir          string
	demo         bool
	title        string
	slides       []slide
	styles       string
	dev          bool
	pdfPrint     bool
	clickEnabled bool
	codeTheme    string
	slideRatio   string
	devCon       *websocket.Conn
	presenterCon *websocket.Conn
}

func main() {

	app := cli.NewApp()
	app.Name = "slide-serve"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "dir",
			Value: "example",
			Usage: "root dir of your presentation",
		},
		cli.StringFlag{
			Name:  "port",
			Value: "8080",
			Usage: "`PORT` to start the http server on",
		},
		cli.StringFlag{
			Name:  "title",
			Value: "Slide",
			Usage: "html title",
		},
		cli.StringFlag{
			Name:  "syntaxhl, hl",
			Value: "monokai",
			Usage: "code highlighter theme",
		},
		cli.StringFlag{
			Name:  "ratio",
			Value: "16x9",
			Usage: "ratio of your slides, 4x3, 16x9 or 16x10",
		},
		cli.BoolFlag{
			Name:  "pdf",
			Usage: "printing a pdf",
		},
		cli.BoolFlag{
			Name:  "dev",
			Usage: "dev true to start a filewatcher and reload the edited slide",
		},
		cli.BoolFlag{
			Name:  "enableClick, click, c",
			Usage: "on default you only navigate with arrow keys, this enabled 'next slide' on click",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal("cannot start slide-serve server! ", err.Error())
	}

}

func run(c *cli.Context) error {
	isDemo := false
	rootDir := c.String("dir")
	devMode := c.Bool("dev")
	title := c.String("title")

	if rootDir == "example" && !dirExist(rootDir) {
		isDemo = true
		devMode = false
		title = "Slide"
	} else if !dirExist(rootDir) {
		return errors.New("cannot find root directory :(")
	}

	h := holder{
		dir:          rootDir,
		title:        title,
		dev:          devMode,
		demo:         isDemo,
		codeTheme:    c.String("syntaxhl"),
		pdfPrint:     c.Bool("pdf"),
		slideRatio:   c.String("ratio"),
		clickEnabled: c.Bool("enableClick"),
	}

	h.parse()

	if devMode {
		http.HandleFunc("/ws", h.ws)
		go h.startFileWatcher(rootDir)
	}

	http.HandleFunc("/", h.handler)
	http.HandleFunc("/presenter", h.presenterHandler)
	http.HandleFunc("/presenterws", h.presenterSocket)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(rootDir))))
	http.HandleFunc("/favicon.ico", http.NotFound)
	port := c.String("port")
	log.Println("starting on port: " + port + " for directory " + rootDir)
	return http.ListenAndServe(":"+port, nil)
}

func (h *holder) handler(w http.ResponseWriter, r *http.Request) {
	slideFile, _ := pkger.Open("/www/slide.html")
	t, _ := template.New("slide").Parse(mustFileToString(slideFile))

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

	cssFile, _ := pkger.Open("/www/summaryStyle.css")
	if h.pdfPrint {
		cssFile, _ = pkger.Open("/www/pdfStyle.css")
	}

	s := slideContent{
		Slides:        template.HTML(slides),
		Styles:        template.CSS(styles),
		PrintStyle:    template.CSS(mustFileToString(cssFile)),
		Title:         h.title,
		SlideRatio:    h.slideRatio,
		ClickListener: "",
	}

	if h.clickEnabled {
		s.ClickListener = "window.onclick = next;"
	}

	if h.dev {
		devModeFile, _ := pkger.Open("/www/devMode.html")
		js, _ := template.New("devmode").Parse(mustFileToString(devModeFile))
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

func mustFileToString(f http.File) string {
	content, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("must read file to string failed: ", err.Error())
	}
	return string(content)
}
