package router

import (
	"fmt"
	"hw6coursera/service"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	limitField  = "limit"
	offsetField = "offset"
)

type requestProcessor struct {
	service *service.Service
}

// DeleteRecord implements RequestProcessor
func (rp *requestProcessor) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	tableName := path[0]
	id, err := strconv.Atoi(path[1])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := rp.service.DeleteById(tableName, id); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("deleted record id %d", id)))
}

// GetAllTables implements RequestProcessor
func (rp *requestProcessor) GetAllTables(w http.ResponseWriter, r *http.Request) {
	data, err := rp.service.RecordService.GetAllTables()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to get tables"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application-json")
	w.Write(data)
}

// GetRecords implements RequestProcessor
func (rp *requestProcessor) GetRecords(w http.ResponseWriter, r *http.Request) {
	tableName := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/"), "/")
	limit := getIntFieldOrDefault(r, limitField, service.DefaultLimit)
	offset := getIntFieldOrDefault(r, offsetField, service.DefaultOffset)
	data, err := rp.service.GetAllRecords(tableName, limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to get records"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application-json")
	w.Write(data)
}

func getIntFieldOrDefault(r *http.Request, field string, defaultValue int) int {
	valueStr := r.URL.Query().Get("limit")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		// Просто логгируемся, не крашимся
		log.Printf("got error parsing to int %s value (%s): %+v\n", field, valueStr, err)
		return defaultValue
	}
	return value
}

// GetSingleRecord implements RequestProcessor
func (rp *requestProcessor) GetSingleRecord(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	tableName := path[0]
	id, err := strconv.Atoi(path[1])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := rp.service.GetById(tableName, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to get tables"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application-json")
	w.Write(data)

}

// InsertRecord implements RequestProcessor
func (rp *requestProcessor) InsertRecord(w http.ResponseWriter, r *http.Request) {
	tableName := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/"), "/")

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	urlVals := r.PostForm
	unit := make(map[string]string)
	for k := range urlVals {
		unit[k] = urlVals.Get(k)
	}

	lastInsertedId, err := rp.service.Create(tableName, unit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to insert record"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("last insert id %d", lastInsertedId)))
}

// UpdateRecord implements RequestProcessor
func (rp *requestProcessor) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	tableName := path[0]
	id, err := strconv.Atoi(path[1])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	urlVals := r.PostForm
	unit := make(map[string]string)
	for k := range urlVals {
		unit[k] = urlVals.Get(k)
	}

	if err := rp.service.UpdateById(tableName, id, unit); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to insert record"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("updated record id %d", id)))
}

func newRequectProcessor(s *service.Service) *requestProcessor {
	return &requestProcessor{
		service: s,
	}
}
