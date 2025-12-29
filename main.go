package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type AuthRoundTripper struct {
	APIKey string
	Base   http.RoundTripper
}

type Point struct {
	Timestamp int64   `json:"ts"`
	Value     float64 `json:"value"`
}

func (rt *AuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+rt.APIKey)
	return rt.Base.RoundTrip(req)
}

func QueryInstant(ctx context.Context, endpoint, promql string) ([]Point, error) {
	client, err := api.NewClient(api.Config{
		Address: endpoint,
		RoundTripper: &AuthRoundTripper{
			APIKey: os.Getenv("THANOS_API_KEY"),
			Base:   http.DefaultTransport,
		},
	})
	if err != nil {
		return nil, err
	}

	apiV1 := v1.NewAPI(client)

	r := v1.Range{
		Start: time.Now().Add(-5 * time.Minute),
		End:   time.Now(),
		Step:  time.Minute,
	}

	result, _, err := apiV1.QueryRange(ctx,
		`rate(prometheus_http_requests_total[1m])`,
		r,
	)

	matrix := result.(model.Matrix)

	points := []Point{}
	for _, stream := range matrix {
		for _, sample := range stream.Values {
			points = append(points, Point{
				Timestamp: int64(sample.Timestamp), // unix seconds
				Value:     float64(sample.Value),
			})
		}
	}

	return points, nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := QueryInstant(
		ctx,
		"http://localhost:9090",
		`rate(prometheus_http_requests_total[1m])`,
	)
	if err != nil {
		panic(err)
	}

	x, _ := json.Marshal(res)
	fmt.Println(string(x))
}
