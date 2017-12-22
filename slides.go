package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packr"
)

type SlideContent struct {
	Slides  string
	Title   string
	Styles  template.CSS
	DevMode template.HTML
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
		tmp := strings.Trim(scanner.Text(), "\t ")
		if tmp == "" {
			h.slides = append(h.slides, s)
			s = slide{}
		} else {
			if strings.HasPrefix(tmp, "@img") {
				s.image = strings.Replace(tmp, "@img", "", -1)
			} else if strings.HasPrefix(tmp, "@css") {
				// TODO reading file? - need to add slide index
				s.styles = strings.Replace(tmp, "@css", "", -1)
			} else {
				s.content += "\t" + tmp + "\n"
			}
		}
	}
	h.slides = append(h.slides, s)

}
