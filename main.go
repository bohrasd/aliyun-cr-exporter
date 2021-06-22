package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bohrasd/aliyun-cr-exporter/aliyun"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// declare metrics
	// namespacesGauge := prometheus.NewGaugeVec(
	// prometheus.GaugeOpts{
	// Name: "aliyun_acr_namespace_info",
	// Help: "Aliyun ACR Namespace",
	// },
	// []string{"namespaceStatus", "namespace", "authorizeType"},
	// )

	// repoDownloadsGauge := prometheus.NewGaugeVec(
	// prometheus.GaugeOpts{
	// Name: "aliyun_acr_repo_downloads",
	// Help: "Aliyun ACR Repository Downloads",
	// },
	// []string{ //"downloads", "gmtCreate", "gmtModified", "logo",
	// "regionId", "repoAuthorizeType", "repoBuildType", //"repoDomainList",
	// "repoId", "repoName", "repoNamespace", "repoStatus", "repoType", "summary"},
	// )

	// prometheus.MustRegister(namespacesGauge, repoDownloadsGauge)
	exporter := NewAliyunExporter()
	prometheus.MustRegister(&exporter)

	// c_mutex := aliyun.NewAliyunClientMutex()

	// namespaces := c_mutex.NamespacesList()

	// for _, namespace := range namespaces {

	//     namespacesGauge.With(prometheus.Labels{
	//         "namespace":       namespace.Namespace,
	//         "namespaceStatus": namespace.NamespaceStatus,
	//         "authorizeType":   namespace.AuthorizeType,
	//     }).Set(1)

	//     ch := make(chan aliyun.Repo, 100)
	//     go c_mutex.ReposList(namespace, ch)

	//     go func(namespace aliyun.Namespace) {
	//         for repo := range ch {
	//             fmt.Println(repo)

	//             repoDownloadsGauge.With(prometheus.Labels{
	//                 // "downloads": repo.Downloads,
	//                 // "gmtCreate": repo.GmtCreate,
	//                 // "gmtModified": repo.GmtModified,
	//                 // "logo": repo.Logo,
	//                 "regionId":          repo.RegionId,
	//                 "repoAuthorizeType": repo.RepoAuthorizeType,
	//                 "repoBuildType":     repo.RepoBuildType,
	//                 // "repoDomainList": repo.RepoDomainList,
	//                 "repoId":        strconv.FormatInt(repo.RepoId, 10),
	//                 "repoName":      repo.RepoName,
	//                 "repoNamespace": repo.RepoNamespace,
	//                 // "repoOriginType": repo.RepoOriginType,
	//                 "repoStatus": repo.RepoStatus,
	//                 "repoType":   repo.RepoType,
	//                 // "stars": repo.Stars,
	//                 "summary": repo.Summary,
	//             }).Set(float64(repo.Downloads))

	//         }
	//     }(namespace)

	// }

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
	ac         *aliyun.AliyunClientMutex
	metric_map MetricMap
	// mu         sync.Mutex
}

func NewAliyunExporter() AliyunExporter {

	ae := AliyunExporter{

		ac: aliyun.NewAliyunClientMutex(),

		metric_map: MetricMap{
			"aliyun_acr_namespace_info": prometheus.NewDesc("aliyun_acr_namespace_info", "Aliyun ACR Namespace",
				[]string{"namespaceStatus", "namespace", "authorizeType"}, nil),
		},
	}

	ae.metric_map.initRepoDisc()
	ae.metric_map.initBuildDisc()

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

	namespaces := exporter.ac.NamespacesList()

	multi_chan := map[string]chan aliyun.Repo{}
	for _, namespace := range namespaces {

		prom_ch <- prometheus.MustNewConstMetric(exporter.metric_map["aliyun_acr_namespace_info"], prometheus.GaugeValue, 1,
			namespace.Namespace, namespace.NamespaceStatus, namespace.AuthorizeType)

		multi_chan[namespace.Namespace] = make(chan aliyun.Repo, 100)

		go func(namespace aliyun.Namespace) {
			exporter.ac.ReposList(namespace, multi_chan[namespace.Namespace])
		}(namespace)

		for repo := range multi_chan[namespace.Namespace] {
			if repo.RepoBuildType == "AUTO_BUILD" {

				build, total := exporter.ac.GetLastestBuilds(namespace.Namespace, repo.RepoName)
				exporter.metric_map.collectBuildTotal(prom_ch, repo, total)
				exporter.metric_map.collectBuildMetric(prom_ch, build)

			}

			exporter.metric_map.collectRepoMetric(prom_ch, repo)
		}
	}
}
