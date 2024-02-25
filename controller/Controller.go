package controller

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"service/model"
)

/*
Type Controller implements:
- Name : This is the identifier the yaml uses to parse it into an actual endpoint
- http.Handler interface : This is the function that tells the controller what to do
- Model : Handles any interaction witht the database
*/
type Controller struct {
	Name     string
	Model    *model.Model
	Fallback []byte
	http.Handler
}

/* Constructor for the controller, outside of package used like this 'Controller.Create(x,y)' */
func Create(name string, datamodel *model.Model, fallback []byte) Controller {
	return Controller{Name: name, Model: datamodel, Fallback: fallback}
}

/* logic is the function to fulfill the http.Handler interface. */
func (c Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if c.Model == nil {
		_, err := w.Write(c.Fallback)
		if err != nil {
			log.Err(err).Msg("Something went wrong with Fallback")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Err(err).Msg("Something went wrong with reading the request body")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		query, err := c.Model.Querybuilder(body)
		if err != nil {
			log.Err(err).Msg("Something went wrong with building query")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		result, err := c.Model.Query(query)
		if err != nil {
			log.Err(err).Msg("Something went wrong with querying database")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer result.Close()

		// Convert the result to a slice of maps
		var results []map[string]interface{}
		columns, err := result.Columns()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		for result.Next() {
			row := make(map[string]interface{})
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range columns {
				valuePtrs[i] = &values[i]
			}
			if err := result.Scan(valuePtrs...); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			for i, col := range columns {
				val := values[i]
				row[col] = val
			}
			results = append(results, row)
		}

		// Marshal the results into JSON
		resp, err := json.Marshal(results)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Send the JSON response
		_, err = w.Write(resp)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
