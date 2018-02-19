package main

import ("net/http"
	"html/template")

type Person struct {
	Firstname string
	Lastname string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	person := Person{Firstname: "Mike", Lastname: "White"}
	t, _ := template.ParseFiles("personTemplate.html")
	t.Execute(w, person)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}
