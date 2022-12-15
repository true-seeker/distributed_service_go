package services

import (
	"fmt"
	"io"
	"net/http"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "ok")
}

func StartWebServerListener() {
	http.HandleFunc("/", getRoot)
	fmt.Println("http listener started")
	err := http.ListenAndServe(":3333", nil)
	FailOnError(err, "Cant start http listener")
}
