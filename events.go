package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

// You more than likely want your "Bot User OAuth Access Token"
func getOauthToken() string {
	dat, err := ioutil.ReadFile("./secrets/oauth")
	if (err != nil) { panic(err) }
	var ret = string(dat)
	return ret[:len(ret)-1]
}

func getVerificationToken() string {
	dat, err := ioutil.ReadFile("./secrets/slack")
	if (err != nil) { panic(err) }
	var ret = string(dat)
	return ret[:len(ret)-1]
}

func main() {
	var oauthToken = getOauthToken();
	var verificationToken = getVerificationToken();
	fmt.Printf("oauthToken: '%s'", oauthToken)
	fmt.Printf("verificationToken: '%s'", verificationToken)
	var api = slack.New(oauthToken);

	http.HandleFunc("/slack/event", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handling event")
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		fmt.Println("body:")
		fmt.Println(body)

		eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body),
					slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: verificationToken}))
		if e != nil {
			fmt.Println(e)
			fmt.Println("internal server error - parseEvent")
			w.WriteHeader(http.StatusInternalServerError)
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				fmt.Println("internal server error - unmarshal request json")
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "text")
			w.Write([]byte(r.Challenge))
		}
		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
				case *slackevents.AppMentionEvent:
					api.PostMessage(ev.Channel, slack.MsgOptionText("At your service.", false))
			}
		}
	})
	fmt.Println("[INFO] Server listening")
	if err := http.ListenAndServe(":710", nil); err != nil {
		panic(err)
	}
}
