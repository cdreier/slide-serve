package main

import (
	"encoding/json"
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
	wwwDir := pkger.Dir("/www")
	slideFile, _ := wwwDir.Open("presenter.html")
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

	s := slideContent{
		Slides:     template.HTML(slides),
		Styles:     template.CSS(styles),
		PrintStyle: template.CSS(""),
		Title:      h.title,
		SlideRatio: h.slideRatio,
	}

	t.Execute(w, s)
}
