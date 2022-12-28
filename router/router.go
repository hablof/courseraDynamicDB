package router

import (
	"hw6coursera/service"
	"net/http"
	"regexp"
)

type RequestProcessor interface {
	GetRecords(w http.ResponseWriter, r *http.Request)
	InsertRecord(w http.ResponseWriter, r *http.Request)
	GetSingleRecord(w http.ResponseWriter, r *http.Request)
	UpdateRecord(w http.ResponseWriter, r *http.Request)
	DeleteRecord(w http.ResponseWriter, r *http.Request)
	GetAllTables(w http.ResponseWriter, r *http.Request)
}

type Router struct {
	tableAndIdPattern *regexp.Regexp
	tablePattern      *regexp.Regexp

	RequestProcessor
}

func NewRouter(s *service.Service) *Router {
	tableAndIdPattern := regexp.MustCompile(`\A\/\w+\/\d+\/?\z`)
	tablePattern := regexp.MustCompile(`\A\/\w+(?:\?\w+=\w+)?(?:&\w+=\w+)?\/?\z`)
	return &Router{
		RequestProcessor:  newRequectProcessor(s),
		tableAndIdPattern: tableAndIdPattern,
		tablePattern:      tablePattern,
	}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case router.tablePattern.MatchString(r.RequestURI):
		switch r.Method {
		case "GET":
			router.GetRecords(w, r)
		case "PUT":
			router.InsertRecord(w, r)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	case router.tableAndIdPattern.MatchString(r.RequestURI):
		switch r.Method {
		case "GET":
			router.GetSingleRecord(w, r)
		case "PUT":
			router.UpdateRecord(w, r)
		case "DELETE":
			router.DeleteRecord(w, r)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		router.GetAllTables(w, r)
	}
}
