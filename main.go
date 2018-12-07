package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nlopes/slack/slackevents"
)

func getVerificationToken() string {
	data, err := ioutil.ReadFile("./secrets/slack")
	if (err != nil) {
		panic(err)
	}
	var ret = string(data)
	return ret[:len(ret)-1]
}

func main() {
	var verificationToken = getVerificationToken();

	http.HandleFunc("/slack/event", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handling event")
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		fmt.Println("body:")
		fmt.Println(body)

		eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body),
					slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: verificationToken}))
		if err != nil {
			fmt.Println(err)
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
			handle(eventsAPIEvent.InnerEvent);
		}
	})

	fmt.Println("[INFO] Server listening")
	if err := http.ListenAndServe(":710", nil); err != nil {
		panic(err)
	}
}
