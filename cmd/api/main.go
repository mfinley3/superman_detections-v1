package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/oschwald/geoip2-golang"

	detection "github.com/mfinley3/superman_detections-v1/internal/detections/service"
	loginrepository "github.com/mfinley3/superman_detections-v1/internal/detections/sqlite"
	transportHTTP "github.com/mfinley3/superman_detections-v1/internal/detections/transport/http"
)

func main() {
	db, err := loginrepository.ConnectAndMigrateDB("./resources/superman_detections.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	geobase, err := geoip2.Open("./resources/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer geobase.Close()

	lr := loginrepository.New(db)
	ls := detection.New(lr, geobase)

	errs := make(chan error, 1)

	go startServer(ls, errs)

	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-sig)
	}()

	err = <-errs
	log.Fatal("superman_detections-v1 terminated: " + err.Error())
}

func startServer(detectionsvc detection.Service, errs chan error) {
	addr := fmt.Sprintf("0.0.0.0:8080")

	mux := chi.NewRouter()
	mux.Route("/superman_detections/v1", func(mux chi.Router) {
		mux.Mount("/", transportHTTP.Handler(detectionsvc))
	})

	log.Printf("listening on: %s\n", addr)
	errs <- http.ListenAndServe(addr, mux)
}
