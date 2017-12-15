package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packr"
)

type holder struct {
	dir    string
	title  string
	slides []slide
	styles string
}

type slide struct {
	content string
	image   string
}

func main() {
	port := flag.String("port", "8080", "http port the server is starting on")
	rootDir := flag.String("dir", "example", "root dir of your presentation")
	title := flag.String("title", "Slide", "html title")
	flag.Parse()

	if !dirExist(*rootDir) {
		log.Fatal("cannot find root directory :(")
	}

	h := holder{
		dir:   *rootDir,
		title: *title,
	}

	h.parse()

	http.HandleFunc("/", h.handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(*rootDir))))
	http.HandleFunc("/favicon.ico", h.na)
	fmt.Println("starting on port: " + *port + " for directory " + *rootDir)
	http.ListenAndServe(":"+*port, nil)
}

type SlideContent struct {
	Slides string
	Title  string
	Styles template.CSS
}

func (h *holder) na(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func (h *holder) parse() {
	filepath.Walk(h.dir, func(path string, info os.FileInfo, err error) error {
		if isDir(path) {
			return nil
		}
		content, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println("cannot read... skipping ", path)
			return nil
		}

		switch filepath.Ext(path) {
		case ".css":
			h.styles += string(content)
		case ".md":
			h.generateSlides(string(content))
			// case ".jpg", ".png", ".gif":
			// 	styles += addStyleRule(path)
		}

		return nil
	})
}

func (h *holder) handler(w http.ResponseWriter, r *http.Request) {
	box := packr.NewBox("./www")
	t, _ := template.New("slide").Parse(box.String("slide.html"))

	slides := ""
	styles := h.styles

	for i, s := range h.slides {
		slides += s.content
		slides += "\n"

		if s.image != "" {
			styles += "\n"
			styles += addStyleRule(s.image, i)
		}

	}

	s := SlideContent{
		Slides: slides,
		Styles: template.CSS(styles),
		Title:  h.title,
	}
	t.Execute(w, s)
}

func addStyleRule(filename string, slideNumber int) string {

	imgURL := "/static/" + filename

	css := fmt.Sprintf(`.slide-%d {
		background: url("%s");
		background-repeat: no-repeat;
		background-size: contain;
		background-position: center;
	}
	`, slideNumber, imgURL)

	return css
}

func (h *holder) generateSlides(content string) {
	cleanup := strings.Trim(content, "\n \t")

	scanner := bufio.NewScanner(strings.NewReader(cleanup))
	s := slide{}
	for scanner.Scan() {
		tmp := strings.Trim(scanner.Text(), "\t ")
		if tmp == "" {
			h.slides = append(h.slides, s)
			s = slide{}
		} else {
			if strings.HasPrefix(tmp, "@img") {
				s.image = strings.Replace(tmp, "@img", "", -1)
			} else {
				s.content += "\t" + tmp + "\n"
			}
		}
	}
	h.slides = append(h.slides, s)

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
