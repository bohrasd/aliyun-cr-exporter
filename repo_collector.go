package main

import (
	"strconv"

	"github.com/bohrasd/aliyun-cr-exporter/aliyun"
	"github.com/prometheus/client_golang/prometheus"
)

// let's not make it too complicated
func (metric_map MetricMap) initRepoDisc() {

	metric_map["aliyun_acr_repo_info"] = prometheus.NewDesc(
		"aliyun_acr_repo_info",
		"Aliyun ACR repository downloads",
		[]string{ //"downloads", "gmtCreate", "gmtModified", "logo",
			"regionId", "repoAuthorizeType", "repoBuildType", //"repoDomainList",
			"repoId", "repoName", "repoNamespace", "repoStatus", "repoType", "summary"}, nil)

	metric_map["aliyun_acr_repo_downloads"] = prometheus.NewDesc(
		"aliyun_acr_repo_downloads",
		"Aliyun ACR repository downloads",
		[]string{"repoName", "repoNamespace", "regionId"}, nil)

	// "aliyun_acr_repo_stars": prometheus.NewDesc(
	// "aliyun_acr_repo_stars",
	// "Aliyun ACR repository stars",
	// []string{"repoId", "repoNamespace", "regionId"}, nil),
}

func (metric_map MetricMap) collectRepoMetric(prom_ch chan<- prometheus.Metric, repo aliyun.Repo) {

	prom_ch <- prometheus.MustNewConstMetric(metric_map["aliyun_acr_repo_info"], prometheus.GaugeValue,
		1,
		repo.RegionId,
		repo.RepoAuthorizeType,
		repo.RepoBuildType,
		strconv.FormatInt(repo.RepoId, 10),
		repo.RepoName,
		repo.RepoNamespace,
		repo.RepoStatus,
		repo.RepoType,
		repo.Summary,
	)
	prom_ch <- prometheus.MustNewConstMetric(metric_map["aliyun_acr_repo_downloads"], prometheus.GaugeValue,
		float64(repo.Downloads),
		repo.RegionId,
		repo.RepoName,
		repo.RepoNamespace,
	)

	// prom_ch <- prometheus.MustNewConstMetric(exporter.metric_map["aliyun_acr_repo_stars"], prometheus.GaugeValue,
	// float64(repo.Stars),
	// repo.RegionId,
	// strconv.FormatInt(repo.RepoId, 10),
	// repo.RepoNamespace,
	// )
}
