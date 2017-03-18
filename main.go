package main

import (
	"bufio"
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"strings"
)

var wsaddr = flag.String("ws", ":8080", "websocket address")

func main() {
	flag.Parse()

	hub := newHub()
	go hub.run()
	go ws(hub)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		hub.broadcast <- scanner.Bytes()
	}

}

func ws(h *Hub) {

	http.HandleFunc("/stream/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/stream/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(h, w, r)
	})

	http.HandleFunc("/stream/thisiamoausodmusdojads",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
			text := r.Form.Get("text")
			fmt.Printf("Outputting: %s\n", text)
			for k, v := range r.Form {
				fmt.Println("key:", k)
				fmt.Println("val:", strings.Join(v, ""))
			}
			h.broadcast <- []byte(text)
		},
	)
	err := http.ListenAndServe(*wsaddr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
