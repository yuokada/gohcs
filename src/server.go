package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	port       int
	docroot    string
	checklist  string
	targetList []TargetServer
)

type TargetServer struct {
	network, address string
}

func checkServer(targets *[]TargetServer) (bool, error) {
	for _, t := range *targets {
		_, err := net.Dial(t.network, t.address)
		if err != nil {
			err = errors.Errorf("connection refused error: %s:%s",
				t.network, t.address)
			return false, err
		}
	}
	return true, nil
}

func CheckAndServerHandler(f http.Handler,
	cf func(*[]TargetServer) (bool, error),
	targets []TargetServer) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Now Checking!!")

		result, err := cf(&targets)

		if result {
			f.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, err.Error())
			return
		}
	}
	return http.HandlerFunc(fn)
}

func init() {
	flag.IntVar(&port, "port", 8000, "listen port")
	flag.StringVar(&docroot, "docroot", "/var/www/html", "Document Root")
	flag.StringVar(&checklist, "checklist", "", "check list file by json")
	flag.Parse()

	if checklist != "" {
		var targetObjcts [][]string
		b, err := ioutil.ReadFile(checklist)
		if err != nil {
			log.Fatal(err)
		}
		jerr := json.Unmarshal(b, &targetObjcts)
		if jerr != nil {
			jerr = errors.Wrapf(jerr, "Cant't Parse %s", checklist)
			log.Fatal(jerr)
		}
		targetList = make([]TargetServer, 0, len(targetObjcts))
		for _, vv := range targetObjcts {
			obj := TargetServer{
				vv[0],
				vv[1],
			}
			targetList = append(targetList, obj)
		}
	} else {
		log.Printf("Please pass checklist filepath\n")
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	mux := http.NewServeMux()
	// static file handler
	func() {
		paths := []string{"/favicon.ico"}
		for _, p := range paths {
			mux.Handle(p, http.FileServer(http.Dir(docroot)))
		}
	}()

	// serve /status.html
	handler := CheckAndServerHandler(
		http.FileServer(http.Dir(docroot)), checkServer, targetList)
	mux.Handle("/status.html", handler)

	fmt.Printf("Listing port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Fatal(err)
	}
}
