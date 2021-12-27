package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
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
	addWebsocket bool
	imageRootUrl string
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
	app.Commands = []cli.Command{
		{
			Name:   "export",
			Usage:  "export slides",
			Action: export,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dest",
					Value: "export",
					Usage: "destination folder",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal("cannot start slide-serve server! ", err.Error())
	}

}
func copyAllFileToFolderNotIncludeExtension(folder string, dest string, ext string) error {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return err
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ext) {
			continue
		}
		src := folder + "/" + f.Name()
		dst := dest + "/" + f.Name()
		if err := copyFile(src, dst); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dest string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		log.Fatal(err)
	}
	return ioutil.WriteFile(dest, data, 0644)
}

func export(c *cli.Context) error {
	rootDir := c.GlobalString("dir")
	dest := c.String("dest")
	log.Println("start export", c.Args())
	if !dirExist(rootDir) {
		return errors.New("cannot find root directory :(")
	}
	h := holder{
		dir:          rootDir,
		title:        c.GlobalString("title"),
		dev:          false,
		demo:         false,
		codeTheme:    c.GlobalString("syntaxhl"),
		pdfPrint:     c.GlobalBool("pdf"),
		slideRatio:   c.GlobalString("ratio"),
		clickEnabled: c.GlobalBoolT("enableClick"),
		addWebsocket: false,
		imageRootUrl: "static",
	}

	h.parse()
	os.Mkdir(dest, 0755)
	os.Mkdir(fmt.Sprintf("%s/static", dest), 0755)
	f, err := os.Create(fmt.Sprintf("%s/Side.html", dest))
	copyAllFileToFolderNotIncludeExtension(rootDir, fmt.Sprintf("%s/static", dest), ".md")
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	h.handle(w, nil)
	return nil
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
		addWebsocket: true,
		imageRootUrl: "/static",
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
func (h *holder) handle(wr io.Writer, r *http.Request) {
	slideFile, _ := pkger.Open("/www/slide.html")
	t, _ := template.New("slide").Parse(mustFileToString(slideFile))

	slides := ""
	styles := h.styles

	for i, s := range h.slides {
		slides += renderSlide(s, i, h.codeTheme)

		if s.image != "" {
			styles += "\n"
			styles += addStyleRule(s.image, i, h.imageRootUrl)
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
		SocketCode:    "",
	}

	if h.clickEnabled {
		s.ClickListener = "window.onclick = next;"
	}

	if h.addWebsocket {
		s.SocketExecuter = "startSocket()"
		s.SocketCode = `
		function startSocket(){
			ws = new WebSocket("ws://"+location.host+"/presenterws");
			ws.onopen = function(){
				ws.send(JSON.stringify({
					type: "presentation:join",
				}))
			}
			ws.onmessage = function(data){
				var msg = JSON.parse(data.data)
				switch(msg.type){
					case "requestNext":
						next();
						break;
					case "requestPrev":
						prev();
						break;
				}
			}
		}
		`
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

	t.Execute(wr, s)
}

func (h *holder) handler(w http.ResponseWriter, r *http.Request) {
	h.handle(w, r)
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
