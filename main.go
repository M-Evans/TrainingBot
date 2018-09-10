package main

import (
	//"encoding/json"
	"fmt"
	//"log"
	"net/http"
	"strings"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message
	w.Write([]byte(message))
	/*
		json, err := json.MarshalIndent(r, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "\n%s\n", string(json))
		// */
	fmt.Fprintf(w, "PostForm: %+v\n", (*r).PostForm)
}

func main() {
	http.HandleFunc("/slack/event", sayHello)
	if err := http.ListenAndServe(":8090", nil); err != nil {
		panic(err)
	}
}
