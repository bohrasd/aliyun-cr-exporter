package aliyun

import (
	"encoding/json"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
)

type Repos struct {
	Total    int
	Page     int
	PageSize int
	Repos    []Repo
}

type Repo struct {
	Downloads int64
	// GmtCreate         int64
	// GmtModified       int64
	Logo              string
	RegionId          string
	RepoAuthorizeType string
	RepoBuildType     string
	// RepoDomainList    struct {
	// Internal string
	// Public   string
	// Vpc      string
	// }
	RepoId        int64
	RepoName      string
	RepoNamespace string
	// RepoOriginType string
	RepoStatus string
	RepoType   string
	Stars      int64
	Summary    string
}

// make sure the lock tiny
// or it will be hold by channel block
func (c_mutex *AliyunClientMutex) repoRequest(ns string, page int, repos *struct{ Data Repos }) {

	c_mutex.Mu.Lock()
	defer c_mutex.Mu.Unlock()

	req := cr.CreateGetRepoListByNamespaceRequest()
	req.PathParams["RepoNamespace"] = ns
	req.Page = requests.NewInteger(page)
	resp, err := c_mutex.Client.GetRepoListByNamespace(req)
	if err != nil {
		// Handle exceptions
		panic(err)
	}

	json.Unmarshal(resp.GetHttpContentBytes(), &repos)
}

// the repos could be too long for your memory
// so lets iterate with channel
func (c_mutex *AliyunClientMutex) ReposList(namespace Namespace, ch chan Repo) {

	defer close(ch)

	repos := struct{ Data Repos }{}
	c_mutex.repoRequest(namespace.Namespace, 1, &repos)

	for _, repo := range repos.Data.Repos {
		ch <- repo
	}
	dt := repos.Data
	maxPage := dt.Total/dt.PageSize + 1
	for i := 2; i <= maxPage; i++ {

		repos = struct{ Data Repos }{}

		c_mutex.repoRequest(namespace.Namespace, i, &repos)
		for _, repo := range repos.Data.Repos {
			ch <- repo
		}
	}
}
