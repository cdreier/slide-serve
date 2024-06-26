package main

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

type slide struct {
	content    string
	code       string
	image      string
	styles     string
	classes    string
	javascript string
	hash       string
	notes      string
}

func (s *slide) buildHash() {
	s.hash = md5Hash(s.content + s.styles + s.image + s.javascript)
}

type slideContent struct {
	Slides         template.HTML
	Notes          template.HTML
	Title          string
	SlideRatio     string
	Styles         template.CSS
	PrintStyle     template.CSS
	DevMode        template.HTML
	ClickListener  template.JS
	SocketCode     template.JS
	SocketExecuter template.JS
}

func (h *holder) parse() {
	if h.demo {
		// example presentation
		fmt.Println("serving example presentation")

		// TODO remove pkger in favor of embed

		// pkger.Walk("/example", func(path string, info os.FileInfo, err error) error {

		// 	if filepath.Base(path) == "styles.css" {
		// 		f, _ := pkger.Open(path)
		// 		h.styles += mustFileToString(f)
		// 	}

		// 	if filepath.Ext(path) == ".md" {
		// 		f, _ := pkger.Open(path)
		// 		h.generateSlides(mustFileToString(f))
		// 	}

		// 	return nil
		// })

	} else {
		// user presentation
		h.slides = make([]slide, 0)
		h.styles = ""
		filepath.Walk(h.dir, func(path string, info os.FileInfo, _ error) error {
			if info == nil || info.IsDir() {
				return nil
			}
			// reading all the files, check file ext before reading?
			content, err := os.ReadFile(path)
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

func addStyleRule(filename string, slideNumber int, imgRootUrl string) string {

	imgURL := imgRootUrl + filename

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

	cleanup := strings.Trim(content, "\n\t")

	hideSlide := false
	scanner := bufio.NewScanner(strings.NewReader(cleanup))
	s := slide{}
	for scanner.Scan() {
		tmp := strings.TrimRight(scanner.Text(), "\t")
		// empty line marks new slide
		if tmp == "" {
			if !hideSlide {
				s.buildHash()
				h.slides = append(h.slides, s)
			}
			s = slide{}
			hideSlide = false
		} else {

			if strings.HasPrefix(tmp, "@img") {
				s.image = strings.Replace(tmp, "@img", "", -1)
			} else if strings.HasPrefix(tmp, "@css") {
				filename := strings.Replace(tmp, "@css", "", -1)
				data, err := os.ReadFile(h.dir + filename)
				if err == nil {
					s.styles = string(data)
				}

			} else if strings.HasPrefix(tmp, "@js") {
				filename := strings.Replace(tmp, "@js", "", -1)
				data, err := os.ReadFile(h.dir + filename)
				if err == nil {
					s.javascript = string(data)
				}

			} else if strings.HasPrefix(tmp, "@code/") {
				s.code = strings.Replace(tmp, "@code/", "", -1)
			} else if strings.HasPrefix(tmp, "@classes/") {
				s.classes = strings.Replace(tmp, "@classes/", "", -1)
			} else if strings.HasPrefix(tmp, "@append") {
				prevSlide := h.slides[len(h.slides)-1]
				buf := s.content
				s = prevSlide
				s.content += buf
			} else if strings.HasPrefix(tmp, "@hidden") {
				hideSlide = true
			} else if strings.HasPrefix(tmp, "@note") {
				s.notes += strings.Replace(tmp, "@note", "", -1)
				s.notes += "<br/>"
			} else {
				s.content += tmp + "\n"
			}
		}
	}
	if !hideSlide {
		s.buildHash()
		h.slides = append(h.slides, s)
	}

}
