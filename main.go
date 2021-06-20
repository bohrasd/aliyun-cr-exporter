package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
)

type Namespaces struct {
	Namespaces []Namespace
}

type Namespace struct {
	NamespaceStatus string
	Namespace       string
	AuthorizeType   string
}

type Repos struct {
	Total    int
	Page     int
	PageSize int
	Repos    []Repo
}

type Repo struct {
	Downloads         int64
	GmtCreate         int64
	GmtModified       int64
	Logo              string
	RegionId          string
	RepoAuthorizeType string
	RepoBuildType     string
	RepoDomainList    struct {
		Internal string
		Public   string
		Vpc      string
	}
	RepoId         int64
	RepoName       string
	RepoNamespace  string
	RepoOriginType string
	RepoStatus     string
	RepoType       string
	Stars          int64
	Summary        string
}

func main() {

	// declare metrics
	namespacesGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aliyun_acr_namespace_info",
			Help: "Aliyun ACR Namespace",
		},
		[]string{"namespaceStatus", "namespace", "authorizeType"},
	)

	reposGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aliyun_acr_repo_info",
			Help: "Aliyun ACR Repository",
		},
		[]string{ //"downloads", "gmtCreate", "gmtModified", "logo",
			"regionId", "repoAuthorizeType", "repoBuildType", //"repoDomainList",
			"repoId", "repoName", "repoNamespace", "repoStatus", "repoType", "summary"},
	)

	prometheus.MustRegister(namespacesGauge, reposGauge)

	// collect metrics info
	client, err := cr.NewClientWithAccessKey("cn-hangzhou", "ABC", "DEF")
	if err != nil {
		// Handle exceptions
		panic(err)
	}

	getNSResp, err := client.GetNamespaceList(cr.CreateGetNamespaceListRequest())
	if err != nil {
		// Handle exceptions
		panic(err)
	}
	if getNSResp.GetHttpStatus() == 200 {

		namespaces := struct{ Data Namespaces }{}

		json.Unmarshal(getNSResp.GetHttpContentBytes(), &namespaces)

		for _, namespace := range namespaces.Data.Namespaces {
			namespacesGauge.With(prometheus.Labels{
				"namespace":       namespace.Namespace,
				"namespaceStatus": namespace.NamespaceStatus,
				"authorizeType":   namespace.AuthorizeType,
			}).Set(1)

			go func(namespace Namespace) {
				repoListReq := cr.CreateGetRepoListByNamespaceRequest()
				repoListReq.PathParams["RepoNamespace"] = namespace.Namespace
				repoListResp, err := client.GetRepoListByNamespace(repoListReq)
				if err != nil {
					// Handle exceptions
					panic(err)
				}

				repos := struct{ Data Repos }{}
				json.Unmarshal(repoListResp.GetHttpContentBytes(), &repos)

				dt := repos.Data
				maxPage := dt.Total/dt.PageSize + 1
				for i := 2; i <= maxPage; i++ {

					repoListReq.Page = requests.NewInteger(i)
					repoListResp, err = client.GetRepoListByNamespace(repoListReq)

					if err != nil {
						// Handle exceptions
						panic(err)
					}

					repos = struct{ Data Repos }{}

					json.Unmarshal(repoListResp.GetHttpContentBytes(), &repos)
					for _, repo := range repos.Data.Repos {

						reposGauge.With(prometheus.Labels{
							// "downloads": repo.Downloads,
							// "gmtCreate": repo.GmtCreate,
							// "gmtModified": repo.GmtModified,
							// "logo": repo.Logo,
							"regionId":          repo.RegionId,
							"repoAuthorizeType": repo.RepoAuthorizeType,
							"repoBuildType":     repo.RepoBuildType,
							// "repoDomainList": repo.RepoDomainList,
							"repoId":        strconv.FormatInt(repo.RepoId, 10),
							"repoName":      repo.RepoName,
							"repoNamespace": repo.RepoNamespace,
							// "repoOriginType": repo.RepoOriginType,
							"repoStatus": repo.RepoStatus,
							"repoType":   repo.RepoType,
							// "stars": repo.Stars,
							"summary": repo.Summary,
						})
					}

				}
			}(namespace)

		}

	}

	// serve metrics
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))
}
