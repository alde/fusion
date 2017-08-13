package server

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/alde/fusion/db"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type errorDocument struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Route enforces the structure of a route
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.Handler
}

// Routes is a collection of route structs
type Routes []Route

// NewRouter creates a new http router
func NewRouter(sql *db.FusionDAO) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes(sql) {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			//Handler(prometheus.InstrumentHandler(route.Name, route.Handler))
			Handler(route.Handler)
	}
	router.NotFoundHandler = http.HandlerFunc(notFound())
	return router
}

func routes(sql *db.FusionDAO) Routes {
	return Routes{
		Route{
			Name:    "news",
			Method:  "GET",
			Pattern: "/api/v1/news",
			Handler: news(sql),
		},
	}
}

func news(sql *db.FusionDAO) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		offset := parseQuery(q, "offset", 0)
		limit := parseQuery(q, "limit", 50)
		logrus.WithFields(logrus.Fields{
			"limit":  limit,
			"offset": offset,
		}).Debug("Fetching news")
		news, err := sql.News(offset, limit)

		if err != nil {
			handleError(err, w, "Unable to read news from the database")
			return
		}

		writeJSON(200, news, w)
	}
}

func handleError(err error, w http.ResponseWriter, message string) {
	if err == nil {
		return
	}

	errorMessage := errorDocument{
		err.Error(),
		message,
	}

	if err = writeJSON(422, errorMessage, w); err != nil {
		logrus.WithError(err).WithField("message", message).Panic("Unable to respond")
	}
}

func writeJSON(status int, data interface{}, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func notFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(nil)
	}
}

func parseQuery(query url.Values, key string, defaultValue int) int {
	val := query.Get(key)
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defaultValue
	}
	return int(i)
}
