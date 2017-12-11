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
	dir   string
	title string
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

	http.HandleFunc("/", h.handler)
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
func (h *holder) handler(w http.ResponseWriter, r *http.Request) {
	box := packr.NewBox("www")
	t, _ := template.New("slide").Parse(box.String("slide.html"))

	slides, styles := getSlides(h.dir)

	s := SlideContent{
		Slides: slides,
		Styles: template.CSS(styles),
		Title:  h.title,
	}
	t.Execute(w, s)
}

func getSlides(dir string) (string, string) {
	slides := ""
	styles := ""
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if isDir(path) {
			return nil
		}
		content, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println("cannot read... skipping ", path)
			return nil
		}

		if filepath.Base(path) == "styles.css" {
			styles = string(content)
		} else {
			slides += getSlideContent(string(content))
		}

		return nil
	})

	return slides, styles
}

func getSlideContent(content string) string {
	cleanup := strings.TrimLeft(content, "\n \t")
	cleanup = strings.TrimRight(cleanup, "\n \t")
	finalSlide := ""

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		finalSlide += "\t" + scanner.Text() + "\n"
	}

	finalSlide += "\n"
	return finalSlide
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