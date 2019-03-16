package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"net/http"
)

// TripService is an interface allowing query of token counts and cache flushing
type TripService interface {
	CountByMedallions(medallions []string, date string) (medallionCounts []MedallionsCount)
	CountByMedallionsBypass(medallions []string, date string) (medallionCounts []MedallionsCount)
	FlushCache()
}

// decode a "/flush" ignoring any additional parameters
func decodeFlushRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	var i interface{}
	return i, nil
}

// setup the flush endpoint; patrickmn/go-cache Flush() doesn't return a result,
// so assume flushes always work
func makeMedallionsFlushEndpoint(tripservice TripService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		tripservice.FlushCache()
		return `{result:ok}`, nil
	}
}

// MedallionsCountRequest is a json request on  "/count" for querying token counts
type MedallionsCountRequest struct {
	Medallions []string `json:"medallions"` // medallions to search for
	Date       string   `json:"date"`       // date to search for medallions
	Fresh      bool     `json:"fresh"`      // if true bypass cache
}

// decode a "/count" request and it's json body
func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request MedallionsCountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// makeMedallionsCountEndpoint is the endpoint for "/count"
func makeMedallionsCountEndpoint(tripservice TripService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(MedallionsCountRequest)
		if req.Fresh {
			return tripservice.CountByMedallionsBypass(req.Medallions, req.Date), nil
		}
		return tripservice.CountByMedallions(req.Medallions, req.Date), nil
	}
}

// encodeResponse encodes reponses to "/flush" and "/count"
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
