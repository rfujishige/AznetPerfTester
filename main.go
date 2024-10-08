package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type swanctlConf struct {
	Local_addrs       string
	Local_publicAddrs string
	Remote_addrs0     string
	Remote_addrs1     string
	Psk               string
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("/opt/AznetPerfTester/htmlTemplates/index.html")
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
		t, _ := template.ParseFiles("/opt/AznetPerfTester/htmlTemplates/addNeighbor.html")
		t.Execute(w, nil)
	}

	if r.Method == "POST" {
		t, _ := template.ParseFiles("/opt/AznetPerfTester/htmlTemplates/addNeighbor.html")
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
		t, _ := template.ParseFiles("/opt/AznetPerfTester/htmlTemplates/addRoute.html")
		t.Execute(w, nil)
	}

	if r.Method == "POST" {
		t, _ := template.ParseFiles("/opt/AznetPerfTester/htmlTemplates/addRoute.html")
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
		t, _ := template.ParseFiles("/opt/AznetPerfTester/htmlTemplates/addRouteTakusan.html")
		t.Execute(w, nil)
	}

	if r.Method == "POST" {
		t, _ := template.ParseFiles("/opt/AznetPerfTester/htmlTemplates/addRouteTakusan.html")
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

func addVpnConnection(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("/opt/AznetPerfTester/htmlTemplates/addVpnConnection.html")
		t.Execute(w, nil)
	}

	if r.Method == "POST" {

		f, err := os.Create("/etc/swanctl/swanctl.conf")
		if err != nil {
			log.Fatal(err)
			return
		}
		defer f.Close()

		t, _ := template.New("swanctl.conf").ParseFiles("/opt/AznetPerfTester/configTemplates/swanctl.conf")
		r.ParseForm()
		data := r.Form

		swanconf := swanctlConf{
			Local_addrs:       data["Local_addrs"][0],
			Local_publicAddrs: data["Local_publicAddrs"][0],
			Remote_addrs0:     data["Remote_addrs0"][0],
			Remote_addrs1:     data["Remote_addrs1"][0],
			Psk:               data["Psk"][0],
		}
		t.Execute(f, swanconf)

		/*
			err = exec.Command("systemctl", "restart", "ipsec").Run()
			if err != nil {
				log.Fatal(err)
				return
			}
			time.Sleep(1)
		*/
		err = exec.Command("swanctl", "--load-all").Run()
		if err != nil {
			log.Fatal(err)
			return
		}
		time.Sleep(1)

		err = exec.Command("swanctl", "--initiate", "--ike", "vng0", "--child", "s2s0").Start()
		if err != nil {
			log.Fatal(err)
			return
		}
		time.Sleep(1)

		err = exec.Command("swanctl", "--initiate", "--ike", "vng1", "--child", "s2s0").Start()
		if err != nil {
			log.Fatal(err)
			return
		}
		time.Sleep(1)

		out, err := exec.Command("ipsec", "status").Output()
		if err != nil {
			log.Fatal(err)
			return
		}
		time.Sleep(1)

		fmt.Fprintf(w, string(out))

	}
}

func chechVPNstatus(w http.ResponseWriter, r *http.Request) {
	out, err := exec.Command("ipsec", "statusall").Output()
	if err != nil {
		log.Fatal(err)
		return
	}
	/*
		TemplateFuncs := map[string]interface{}{
			// Replaces newlines with <br>
			"nl2br": func(text string) template.HTML {
				return template.HTML(strings.Replace(template.HTMLEscapeString(text), "\n", "<br>", -1))
			},
		}

		t, _ := template.New("default.html").Funcs(TemplateFuncs).ParseFiles("./htmlTemplates/default.html")

		t.Execute(w,
			struct {
				Message string
			}{
				Message: string(out),
			},
		)
	*/
	fmt.Fprintf(w, string(out))
}

func main() {
	http.HandleFunc("/checkRoutes", checkRoutesHandler)
	http.HandleFunc("/checkNeighbor", checkNeighborHandler)
	http.HandleFunc("/addNeighbor", addNeighborHandler)
	http.HandleFunc("/addRoute", addRouteHandler)
	http.HandleFunc("/addRouteTakusan", addRouteTakusanHandler)
	http.HandleFunc("/addVpnConnection", addVpnConnection)
	http.HandleFunc("/chechVPNstatus", chechVPNstatus)
	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe(":8080", nil)
}
