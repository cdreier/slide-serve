package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/markbates/pkger"
)

type presenterMsg struct {
	Type string `json:"type,omitempty"`
}

func (h *holder) presenterSocket(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("could not start dev-websocket:", err)
		return
	}
	defer connection.Close()

	for {
		// mt, message, err := connection.ReadMessage()
		_, message, err := connection.ReadMessage()
		if err != nil {
			// log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		msg := presenterMsg{}
		json.Unmarshal(message, &msg)

		switch msg.Type {
		case "next":
			if h.presenterCon != nil {
				h.presenterCon.WriteJSON(presenterMsg{
					Type: "requestNext",
				})
			} else {
				log.Println("presenter socker is null")
			}
			break
		case "prev":
			if h.presenterCon != nil {
				h.presenterCon.WriteJSON(presenterMsg{
					Type: "requestPrev",
				})
			}
			break
		case "presentation:join":
			h.presenterCon = connection
			log.Println("presentation joined")
			break
		}
		// err = connection.WriteMessage(mt, message)
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }
	}
}

func (h *holder) presenterHandler(w http.ResponseWriter, r *http.Request) {
	slideFile, _ := pkger.Open("/www/presenter.html")
	t, _ := template.New("slide").Parse(mustFileToString(slideFile))

	slides := ""
	notes := ""
	styles := h.styles

	for i, s := range h.slides {
		slides += renderSlide(s, i, h.codeTheme)
		notes += fmt.Sprintf("<div class=\"note\">%s</div>", s.notes)

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

	s := slideContent{
		Slides:     template.HTML(slides),
		Notes:      template.HTML(notes),
		Styles:     template.CSS(styles),
		PrintStyle: template.CSS(""),
		Title:      h.title,
		SlideRatio: h.slideRatio,
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
