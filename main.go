package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"strconv"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./htmlTemplates/index.html")
	t.Execute(w, nil)
}

func checkNeighborHandler(w http.ResponseWriter, r *http.Request) {
	out, err := exec.Command("gobgp", "neighbor").Output()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Fprintf(w, string(out))
}

func checkRoutesHandler(w http.ResponseWriter, r *http.Request) {
	out, err := exec.Command("gobgp", "global", "rib").Output()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Fprintf(w, string(out))
}

func addNeighborHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./htmlTemplates/addNeighbor.html")
		t.Execute(w, nil)
	}

	if r.Method == "POST" {
		t, _ := template.ParseFiles("./htmlTemplates/addNeighbor.html")
		t.Execute(w, nil)
		r.ParseForm()
		neighbor := r.Form

		// duplicate ip causes error
		err := exec.Command("gobgp", "neighbor", "add", neighbor["neighborIp"][0], "as", neighbor["asn"][0]).Run()
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func addRouteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./htmlTemplates/addRoute.html")
		t.Execute(w, nil)
	}

	if r.Method == "POST" {
		t, _ := template.ParseFiles("./htmlTemplates/addRoute.html")
		t.Execute(w, nil)
		r.ParseForm()
		route := r.Form

		err := exec.Command("gobgp", "global", "rib", "add", route["network"][0], "origin", "igp", "aspath", route["aspath"][0]).Run()

		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func addRouteTakusanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./htmlTemplates/addRouteTakusan.html")
		t.Execute(w, nil)
	}

	if r.Method == "POST" {
		t, _ := template.ParseFiles("./htmlTemplates/addRouteTakusan.html")
		t.Execute(w, nil)
		r.ParseForm()
		numberMap := r.Form["number"][0]
		number, err := strconv.Atoi(numberMap)
		if err != nil {
			log.Fatal(err)
			return
		}

		for i := 0; i < number; i++ {
			network := "192.168." + strconv.Itoa(i) + ".0/24"
			println(network)
			err := exec.Command("gobgp", "global", "rib", "add", network, "origin", "igp").Run()

			if err != nil {
				log.Fatal(err)
				return
			}
		}
	}
}

func main() {
	http.HandleFunc("/checkRoutes", checkRoutesHandler)
	http.HandleFunc("/checkNeighbor", checkNeighborHandler)
	http.HandleFunc("/addNeighbor", addNeighborHandler)
	http.HandleFunc("/addRoute", addRouteHandler)
	http.HandleFunc("/addRouteTakusan", addRouteTakusanHandler)
	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe(":8080", nil)
}
