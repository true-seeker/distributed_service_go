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

func generateTaskView(w http.ResponseWriter, r *http.Request) {
	GenerateTask(5000)
}

func StartWebServerListener() {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/generateTask", generateTaskView)
	fmt.Println("http listener started")
	err := http.ListenAndServe(":3333", nil)
	FailOnError(err, "Cant start http listener")
}
