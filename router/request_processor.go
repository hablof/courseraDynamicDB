package router

import (
	"errors"
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
func (rp *requestProcessor) deleteRecord(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	tableName := path[0]
	id, err := strconv.Atoi(path[1])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	switch err := rp.service.DeleteById(tableName, id); {
	case err == service.ErrRecordNotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("record not found"))
		return
	case err == service.ErrTableNotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("unknown table"))
		return
	case err != nil:
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("deleted record id %d", id)))
}

// GetAllTables implements RequestProcessor
func (rp *requestProcessor) getAllTables(w http.ResponseWriter, r *http.Request) {
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
func (rp *requestProcessor) getRecords(w http.ResponseWriter, r *http.Request) {
	tableName := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/"), "/")
	limit := getIntFieldOrDefault(r, limitField, service.DefaultLimit)
	offset := getIntFieldOrDefault(r, offsetField, service.DefaultOffset)
	data, err := rp.service.GetAllRecords(tableName, limit, offset)
	switch {
	case err == service.ErrTableNotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("unknown table"))
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to get records"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application-json")
	w.Write(data)
}

// GetSingleRecord implements RequestProcessor
func (rp *requestProcessor) getSingleRecord(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	tableName := path[0]
	id, err := strconv.Atoi(path[1])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := rp.service.GetById(tableName, id)
	switch {
	case err == service.ErrRecordNotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("record not found"))
		return
	case err == service.ErrTableNotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("unknown table"))
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to service"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application-json")
	w.Write(data)
}

// InsertRecord implements RequestProcessor
func (rp *requestProcessor) insertRecord(w http.ResponseWriter, r *http.Request) {
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
	switch {
	case err == service.ErrRecordNotFound || err == service.ErrTableNotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	case errors.As(err, &service.ErrType{}) || errors.As(err, &service.ErrCannotBeNull{}):
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to insert record"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("last insert id %d", lastInsertedId)))
}

// UpdateRecord implements RequestProcessor
func (rp *requestProcessor) updateRecord(w http.ResponseWriter, r *http.Request) {
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

	switch err := rp.service.UpdateById(tableName, id, unit); {
	case err == service.ErrRecordNotFound || err == service.ErrTableNotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	case err == service.ErrMissingUpdData || errors.As(err, &service.ErrType{}) || errors.As(err, &service.ErrCannotBeNull{}):
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to update record"))
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

func getIntFieldOrDefault(r *http.Request, field string, defaultValue int) int {
	valueStr := r.URL.Query().Get(field)
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
