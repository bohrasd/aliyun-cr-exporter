package main

import (
	"strconv"

	"github.com/bohrasd/aliyun-cr-exporter/aliyun"
	"github.com/prometheus/client_golang/prometheus"
)

func (metric_map MetricMap) initBuildDisc() {

	metric_map["aliyun_acr_build_status"] = prometheus.NewDesc(
		"aliyun_acr_build_status",
		"Aliyun ACR build status",
		[]string{"repoName", "repoNamespace", "tag", "buildId", "startTime"}, nil)

	metric_map["aliyun_acr_build_succeeded_duration_seconds"] = prometheus.NewDesc(
		"aliyun_acr_build_succeeded_duration_seconds",
		"Aliyun ACR suceeded builds duration",
		[]string{"repoName", "repoNamespace", "tag", "buildId", "startTime"}, nil)
}

var buildStatusMap = map[string]float64{
	"SUCCESS":  1,
	"BUILDING": 2,
	"CANCELED": 3,
	"FAILED":   4,
}

// latest several builds could be enough
// to make sure some of those builds wont slip
// through your collect time gap
func (metric_map MetricMap) collectBuildMetric(prom_ch chan<- prometheus.Metric, builds []aliyun.Build) {
	for _, build := range builds {

		status, exists := buildStatusMap[build.BuildStatus]
		if !exists {
			status = 0
		}
		prom_ch <- prometheus.MustNewConstMetric(metric_map["aliyun_acr_build_status"], prometheus.GaugeValue,
			status,
			build.Image.RepoName,
			build.Image.RepoNamespace,
			build.Image.Tag,
			build.BuildId,
			strconv.FormatInt(build.StartTime, 10),
		)

		if build.BuildStatus == "SUCCESS" {
			prom_ch <- prometheus.MustNewConstMetric(metric_map["aliyun_acr_build_succeeded_duration_seconds"], prometheus.GaugeValue,
				float64(build.EndTime-build.StartTime),
				build.Image.RepoName,
				build.Image.RepoNamespace,
				build.Image.Tag,
				build.BuildId,
				strconv.FormatInt(build.StartTime, 10),
			)
		}
	}
}
