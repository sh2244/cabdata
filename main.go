package main

import (
	httptransport "github.com/go-kit/kit/transport/http"
	"log"
	"net/http"
)

func main() {
	connectDB()
	tripsService := NewTrips()

	flushHandler := httptransport.NewServer(
		makeMedallionsFlushEndpoint(tripsService),
		decodeFlushRequest,
		encodeResponse,
	)
	http.Handle("/flush", flushHandler)

	countHandler := httptransport.NewServer(
		makeMedallionsCountEndpoint(tripsService),
		decodeCountRequest,
		encodeResponse,
	)
	http.Handle("/count", countHandler)

	// listen port would be normally be provided externally eg via service discovery (consul, etcd)
	log.Fatal(http.ListenAndServe(":12345", nil))
}
