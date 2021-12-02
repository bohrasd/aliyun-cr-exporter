package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bohrasd/aliyun-cr-exporter/aliyun"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	exporter := NewAliyunExporter()
	prometheus.MustRegister(&exporter)

	// serve metrics
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Aliyun ACR Exporter</title></head>
             <body>
             <h1>Aliyun ACR Exporter</h1>
             <p><a href="/metrics">Metrics</a></p>
             <h2>Build</h2>
             </body>
             </html>`))
	})
	http.HandleFunc("/-/healthy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})
	http.HandleFunc("/-/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	log.Fatal(http.ListenAndServe(":9101", nil))
}

type MetricMap map[string]*prometheus.Desc

type AliyunExporter struct {
	instance   string
	metric_map MetricMap
	// mu         sync.Mutex
}

type Datas struct {
	namespaces map[string]aliyun.Namespace
	repos      map[int64]aliyun.Repo
	repoBuilds map[int64][]aliyun.Build
	repoTotal  map[int64]int
}

var data = Datas{
	namespaces: make(map[string]aliyun.Namespace),
	repos:      make(map[int64]aliyun.Repo),
	repoBuilds: make(map[int64][]aliyun.Build),
	repoTotal:  make(map[int64]int),
}

func NewAliyunExporter() AliyunExporter {

	ae := AliyunExporter{

		metric_map: MetricMap{
			"aliyun_acr_namespace_info": prometheus.NewDesc("aliyun_acr_namespace_info", "Aliyun ACR Namespace",
				[]string{"namespaceStatus", "namespace", "authorizeType"}, nil),
		},
	}

	ae.metric_map.initRepoDisc()
	ae.metric_map.initBuildDisc()

	go func() {
		ac := aliyun.NewAliyunClientMutex()
		namespaces := ac.NamespacesList()
		limit := make(chan struct{}, 10)

		multi_chan := map[string]chan aliyun.Repo{}
		for _, namespace := range namespaces {

			data.namespaces[namespace.Namespace] = namespace
			multi_chan[namespace.Namespace] = make(chan aliyun.Repo, 100)

			go func() {
				ac.ReposList(namespace, multi_chan[namespace.Namespace])
			}()

			for repo := range multi_chan[namespace.Namespace] {
				data.repos[repo.RepoId] = repo
				if repo.RepoBuildType == "AUTO_BUILD" {

					limit <- struct{}{}
					data.repoBuilds[repo.RepoId], data.repoTotal[repo.RepoId] = ac.GetLastestBuilds(namespace.Namespace, repo.RepoName)

					<-limit
				}

			}
		}

		<-time.Tick(time.Second * 30)
	}()

	return ae

}

//implement Collector
func (exporter *AliyunExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range exporter.metric_map {
		ch <- desc
	}

}

//implement Collector
func (exporter *AliyunExporter) Collect(prom_ch chan<- prometheus.Metric) {
	// exporter.mu.Lock()
	// defer exporter.mu.Unlock()

	namespaces := data.namespaces

	for _, namespace := range namespaces {

		prom_ch <- prometheus.MustNewConstMetric(exporter.metric_map["aliyun_acr_namespace_info"], prometheus.GaugeValue, 1,
			namespace.Namespace, namespace.NamespaceStatus, namespace.AuthorizeType)

	}
	for _, repo := range data.repos {
		if repo.RepoBuildType == "AUTO_BUILD" {
			exporter.metric_map.collectBuildTotal(prom_ch, repo, data.repoTotal[repo.RepoId])
			exporter.metric_map.collectBuildMetric(prom_ch, data.repoBuilds[repo.RepoId])

		}

		exporter.metric_map.collectRepoMetric(prom_ch, repo)
	}
}
