package main

import "net/http"
import f "fmt"

func main() {

	http.HandleFunc("/", rootHandler)
	http.ListenAndServe(":11180", nil)
}

var globalIntger *int = new(int)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	*globalIntger++
	f.Println(w)
	f.Println("アクセスしたURL↓")
	f.Println(r.URL.Path)
	f.Println("アクセスしたURL↑")
	f.Fprint(w, "Hello World")
	f.Fprintf(w, "グローバル変数の値%d", *globalIntger)
}
