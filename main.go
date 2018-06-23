package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/oschwald/geoip2-golang"
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
		db, err := geoip2.Open("GeoLite2-City.mmdb")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		var ip []string
		var text string
		re, _ := regexp.Compile("(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])")
		r.ParseForm()

		text = strings.Join(r.Form["text"], " ")
		ip = re.FindAllString(text, -1)
		sort.Strings(ip)
		ip = unique(ip)

		for _, v := range ip {
			fmt.Fprint(w, v)
			fmt.Fprint(w, "\n")
			i := net.ParseIP(v)
			record, err := db.City(i)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprint(w, "City name: ", record.City.Names["en"])
			//fmt.Fprint(w, "\nSubdivision name: ", record.Subdivisions[0].Names["en"])
			fmt.Fprint(w, "\nCountry name: ", record.Country.Names["en"])
			fmt.Fprint(w, "\nISO country code: ", record.Country.IsoCode)
			fmt.Fprint(w, "\nTime zone: ", record.Location.TimeZone)
			fmt.Fprint(w, "\nCoordinates: ", record.Location.Latitude, record.Location.Longitude)
			fmt.Fprint(w, "\n\n")
		}
	}
}

func unique(s []string) []string {
	m := make(map[string]bool)
	for _, item := range s {
		if _, ok := m[item]; ok {
			// duplicate item
			fmt.Println(item, "is a duplicate")
		} else {
			m[item] = true
		}
	}

	var result []string
	for item, _ := range m {
		result = append(result, item)
	}
	return result
}
