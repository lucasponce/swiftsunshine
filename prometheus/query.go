package prometheus

import (
	"net/http"
	"fmt"
	"html"
	"github.com/lucasponce/swiftsunshine/version"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"log"
	"context"
	"time"
)

const reqsFmt = "sum(rate(istio_request_count[%s])) by (source_service, destination_service, source_version, destination_version)"
const emptyFilter = " > 0"

type queryHandler struct {
	prometheusAddr string
	result string
}

func NewQueryHandler(prometheusAddr string) http.Handler {
	return &queryHandler{prometheusAddr, ""}
}

func (h *queryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Sunshine, %q \n", html.EscapeString(r.URL.Path))
	fmt.Fprintf(w, "Version [%q] \n", version.String())

	timeHorizon := r.URL.Query().Get("time_horizon")
	if timeHorizon == "" {
		timeHorizon = "5m"
	}
	filterEmpty := false
	filterEmptyStr := r.URL.Query().Get("filter_empty")
	if filterEmptyStr == "true" {
		filterEmpty = true
	}

	// validate time_horizon
	if _, err := model.ParseDuration(timeHorizon); err != nil {
		writeError(w, fmt.Errorf("could not parse time_horizon: %v", err))
		return
	}

	client, err := api.NewClient(api.Config{Address: h.prometheusAddr})
	if err != nil {
		writeError(w, fmt.Errorf("could not create a prometheus client: %v", err))
		return
	}

	api := v1.NewAPI(client)
	query := fmt.Sprintf(reqsFmt, timeHorizon)
	if filterEmpty {
		query += emptyFilter
		fmt.Println(query)
	}
	fmt.Fprintf(w, "Query: [%q] \n", query)
	val, err := api.Query(context.Background(), query, time.Now())
	if err != nil {
		writeError(w, fmt.Errorf("could not perform a prometheus query: %v", err))
		return
	}

	switch val.Type() {
	case model.ValVector:
		matrix := val.(model.Vector)
		for _, sample := range matrix {
			metric := sample.Metric
			value := sample.Value
			fmt.Fprintf(w, "Metric [%q] Value [%q] \n", metric.String(), value.String())
		}
	default:
		fmt.Fprintf(w, "No Prometheus data found \n")
		return
	}
}

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, writeErr := w.Write([]byte(err.Error()))
	log.Print(writeErr)
}