package router

import (
	"hw6coursera/service"
	"net/http"
	"regexp"
)

type RequestProcessor interface {
	getRecords(w http.ResponseWriter, r *http.Request)
	insertRecord(w http.ResponseWriter, r *http.Request)
	getSingleRecord(w http.ResponseWriter, r *http.Request)
	updateRecord(w http.ResponseWriter, r *http.Request)
	deleteRecord(w http.ResponseWriter, r *http.Request)
	getAllTables(w http.ResponseWriter, r *http.Request)
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
			router.getRecords(w, r)
		case "PUT":
			router.insertRecord(w, r)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	case router.tableAndIdPattern.MatchString(r.RequestURI):
		switch r.Method {
		case "GET":
			router.getSingleRecord(w, r)
		case "POST":
			router.updateRecord(w, r)
		case "DELETE":
			router.deleteRecord(w, r)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		router.getAllTables(w, r)
	}
}
