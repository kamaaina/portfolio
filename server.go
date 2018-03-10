package main

import ("net/http"
	"html/template")

type Data struct {
	Accounts []string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	accts := make([]string, 0)
	accts = append(accts, "V")
	accts = append(accts, "F")
	accts = append(accts, "RM")
	accts = append(accts, "RC")
	accts = append(accts, "IRA")
	
	data := Data{Accounts: accts}
	t, _ := template.ParseFiles("portfolioTemplate.html")
	
	err := t.Execute(w, data)
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	http.ListenAndServe(":8080", nil)
}
