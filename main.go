package main

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/websocket"
	"github.com/urfave/cli"
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
	slideRatio string
	connection *websocket.Conn
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
		dir:        rootDir,
		title:      title,
		dev:        devMode,
		demo:       isDemo,
		codeTheme:  c.String("syntaxhl"),
		pdfPrint:   c.Bool("pdf"),
		slideRatio: c.String("ratio"),
	}

	h.parse()

	if devMode {
		http.HandleFunc("/ws", h.ws)
		go h.startFileWatcher(rootDir)
	}

	http.HandleFunc("/", h.handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(rootDir))))
	http.HandleFunc("/favicon.ico", http.NotFound)
	port := c.String("port")
	log.Println("starting on port: " + port + " for directory " + rootDir)
	return http.ListenAndServe(":"+port, nil)
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
		SlideRatio: h.slideRatio,
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
