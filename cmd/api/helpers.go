package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) readIdParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

type envelope map[string]any

func (app *application) writeJson(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	res, err := json.Marshal(data)
	if err != nil {
		return err
	}

	res = append(res, '\n')

	w.Header().Set("Content-Type", "application/json")
	for k, v := range headers {
		w.Header()[k] = v
	}

	w.WriteHeader(status)
	w.Write(res)

	return nil
}
