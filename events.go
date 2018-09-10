package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

const TOKEN = "xoxb-222707315955-395033242117-xXN44aHB2Ahl3r0MxaVgMOnQ"

// You more than likely want your "Bot User OAuth Access Token" which starts with "xoxb-"
var api = slack.New(TOKEN)

func main() {
	http.HandleFunc("/slack/event", func(w http.ResponseWriter, r *http.Request) {
		log.Println("handling event")
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		log.Println("body:")
		log.Println(body)

		eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{"hLzQ6ZOhVRC4LaN6yxJ0SRah"}))
		if e != nil {
			log.Println("internal server error - parseEvent")
			w.WriteHeader(http.StatusInternalServerError)
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				log.Println("internal server error - unmarshal request json")
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "text")
			w.Write([]byte(r.Challenge))
		}
		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			postParams := slack.PostMessageParameters{}
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				api.PostMessage(ev.Channel, "Yes, hello.", postParams)
			}
		}
	})
	fmt.Println("[INFO] Server listening")
	if err := http.ListenAndServe(":710", nil); err != nil {
		panic(err)
	}
}
