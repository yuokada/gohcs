package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/pkg/errors"
)

var (
	port       int
	docroot    string
	checklist  string
	targetList []TargetServer
)

var logger = logrus.New()

type TargetServer struct {
	network, address string
}

const (
	HEADER_SERVER = "gohcs"
)

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
		logger.Traceln("Now Checking!!")

		result, err := cf(&targets)
		w.Header().Add("Server", HEADER_SERVER)
		if result {
			f.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
			_, err = fmt.Fprintln(w, err.Error())
			if err != nil {
				logger.Error(err.Error())
			}
			return
		}
	}
	return http.HandlerFunc(fn)
}

func init() {
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Out = os.Stderr

	flag.IntVar(&port, "port", 8000, "listen port")
	flag.StringVar(&docroot, "docroot", "/var/www/html", "Document Root")
	flag.StringVar(&checklist, "checklist", "", "check list file by json")
	loglevel := flag.String("loglevel", "warn", "set log level")
	flag.Parse()

	logLevel, err := logrus.ParseLevel(*loglevel)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Level = logLevel

	if checklist != "" {
		var targetObjcts [][]string
		b, err := ioutil.ReadFile(checklist)
		if err != nil {
			logger.Fatal(err)
		}
		jerr := json.Unmarshal(b, &targetObjcts)
		if jerr != nil {
			jerr = errors.Wrapf(jerr, "Cant't Parse %s", checklist)
			logger.Fatal(jerr)
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
		logger.Error("Please pass checklist filepath")
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

	handler := CheckAndServerHandler(
		http.FileServer(http.Dir(docroot)), checkServer, targetList)
	// mux.Handle("/status.html", handler)
	mux.Handle("/", handler)

	fmt.Printf("Listing port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		logger.Fatal(err)
	}
}
