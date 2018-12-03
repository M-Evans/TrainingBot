package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"io/ioutil"
	"net/http"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

// You more than likely want your "Bot User OAuth Access Token"
func getOauthToken() string {
	dat, err := ioutil.ReadFile("./secrets/oauth")
	if (err != nil) { panic(err) }
	return string(dat)
}

func getVerificationToken() string {
	dat, err := ioutil.ReadFile("./secrets/slack")
	if (err != nil) { panic(err) }
	return string(dat)
}

func main() {
	var oauthToken = getOauthToken();
	var verificationToken = getVerificationToken();
	var api = slack.New(oauthToken);

	http.HandleFunc("/slack/event", func(w http.ResponseWriter, r *http.Request) {
		log.Println("handling event")
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		log.Println("body:")
		log.Println(body)

		eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body),
					slackevents.OptionVerifyToken(&slackevents.TokenComparator{verificationToken}))
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
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
				case *slackevents.AppMentionEvent:
					api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			}
		}
	})
	fmt.Println("[INFO] Server listening")
	if err := http.ListenAndServe(":710", nil); err != nil {
		panic(err)
	}
}
