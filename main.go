package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

func main() {
	http.HandleFunc("/parse", parse)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func parse(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("input.gtpl")
		t.Execute(w, nil)
	} else {
		var ip []string
		var text string
		re, _ := regexp.Compile("(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])")
		r.ParseForm()

		text = strings.Join(r.Form["text"], " ")
		ip = re.FindAllString(text, -1)
		sort.Strings(ip)

		for _, v := range ip {
			fmt.Fprint(w, v)
			fmt.Fprint(w, "\n")
		}
	}
}
