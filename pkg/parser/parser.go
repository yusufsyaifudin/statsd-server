package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type MetricType string

var validMetricType = map[string]MetricType{
	"c":  Counter,
	"g":  Gauge,
	"ms": Timing,
	"s":  Set,
}

const (
	Counter MetricType = "c"
	Gauge   MetricType = "g"
	Timing  MetricType = "ms"
	Set     MetricType = "s"
)

type Metric struct {
	RawPayload string     `json:"raw_payload"`
	BucketName string     `json:"bucket_name,omitempty"`
	Value      float64    `json:"value,omitempty"`
	MetricType MetricType `json:"metric_type,omitempty"`
}

// Parse handle multi-packet, separated with newline
// https://github.com/statsd/statsd/blob/8e6e29ea1f00062c9be27eebb122fa8c427bc74b/docs/metric_types.md#multi-metric-packets
func Parse(payload string) (metrics []Metric, err error) {
	metrics = make([]Metric, 0)
	strMetrics := strings.Split(payload, "\n")
	for _, strMetric := range strMetrics {
		strMetric = strings.TrimSpace(strMetric)
		if strMetric == "" {
			continue
		}

		var metric Metric
		metric, err = ParseMetric(strMetric)
		if err != nil {
			err = fmt.Errorf("cannot parse metric '%s': %w", strMetric, err)
			return
		}

		metrics = append(metrics, metric)
	}

	return
}

// ParseMetric parsing the message payload as statsd Metric
// See how to parse the message here:
// https://github.com/statsd/statsd/blob/9cf77d87855bcb69a8663135f59aa23825db9797/stats.js#L241-L326
// https://github.com/statsd/statsd/blob/9cf77d87855bcb69a8663135f59aa23825db9797/docs/metric_types.md
// TODO: support tags and sampler rate
func ParseMetric(payload string) (metric Metric, err error) {
	metric = Metric{
		RawPayload: payload,
	}

	// ** Parse bucket name
	metricPartStr := ""
	bucketParts := strings.Split(payload, ":")
	if len(bucketParts) >= 2 {
		metric.BucketName = bucketParts[0]
		metricPartStr = strings.Join(bucketParts[1:], "")
	}

	metricParts := strings.Split(metricPartStr, "|")
	metricPartsLen := len(metricParts)

	if metricPartsLen >= 1 {
		metric.Value, _ = strconv.ParseFloat(metricParts[0], 64)
	}

	if metricPartsLen >= 2 {
		metricTypePayload := metricParts[1]
		metricType, ok := validMetricType[metricTypePayload]
		if !ok {
			err = fmt.Errorf("not valid metric type '%s'", metricTypePayload)
			return
		}

		metric.MetricType = metricType
	}

	return
}
