package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gobuffalo/packr"
)

type slide struct {
	content string
	image   string
	styles  string
	hash    string
}

func (s *slide) buildHash() {
	hasher := md5.New()
	hasher.Write([]byte(s.content + s.styles + s.image))
	s.hash = hex.EncodeToString(hasher.Sum(nil))
}

type slideContent struct {
	Slides  string
	Title   string
	Styles  template.CSS
	DevMode template.HTML
}

func (h *holder) parse() {
	if h.demo {
		// example presentation
		exampleBox := packr.NewBox("./example")
		fmt.Println("serving example presentation")
		all := exampleBox.List()
		sort.Strings(all)
		for _, path := range all {
			if filepath.Base(path) == "styles.css" {
				h.styles += exampleBox.String(path)
			}

			if filepath.Ext(path) == ".md" {
				h.generateSlides(exampleBox.String(path))
			}
		}

	} else {
		// user presentation
		h.slides = make([]slide, 0)
		filepath.Walk(h.dir, func(path string, info os.FileInfo, err error) error {
			if info == nil || info.IsDir() {
				return filepath.SkipDir
			}
			// reading all the files, check file ext before reading?
			content, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Println("cannot read... skipping ", path)
				return nil
			}

			if filepath.Base(path) == "styles.css" {
				h.styles += string(content)
			}

			if filepath.Ext(path) == ".md" {
				h.generateSlides(string(content))
			}

			return nil
		})
	}

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

		if s.styles != "" {
			styles += "\n"
			slideStyle := strings.Replace(s.styles, "SLIDENUMBER", strconv.Itoa(i), -1)
			styles += slideStyle
		}

	}

	s := slideContent{
		Slides: slides,
		Styles: template.CSS(styles),
		Title:  h.title,
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
		tmp := strings.TrimRight(scanner.Text(), "\t ")
		if tmp == "" {
			s.buildHash()
			h.slides = append(h.slides, s)
			s = slide{}
		} else {
			if strings.HasPrefix(tmp, "@img") {
				s.image = strings.Replace(tmp, "@img", "", -1)
			} else if strings.HasPrefix(tmp, "@css") {
				filename := strings.Replace(tmp, "@css", "", -1)
				data, err := ioutil.ReadFile(h.dir + filename)
				if err == nil {
					s.styles = string(data)
				}

			} else {
				s.content += "\t" + tmp + "\n"
			}
		}
	}
	s.buildHash()
	h.slides = append(h.slides, s)

}
