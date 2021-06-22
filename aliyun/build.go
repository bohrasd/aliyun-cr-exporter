package aliyun

import (
	"encoding/json"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
)

type Build struct {
	Image struct {
		RepoNamespace string
		RepoName      string
		Tag           string
	}
	StartTime   int64
	EndTime     int64
	BuildId     string
	BuildStatus string
}

// you wont need all of them right
func (c_mutex *AliyunClientMutex) GetLastestBuilds(ns string, repo string) []Build {
	c_mutex.Mu.Lock()
	defer c_mutex.Mu.Unlock()

	req := cr.CreateGetRepoBuildListRequest()
	req.PathParams["RepoNamespace"] = ns
	req.PathParams["RepoName"] = repo
	resp, err := c_mutex.Client.GetRepoBuildList(req)
	if err != nil {
		// Handle exceptions
		panic(err)
	}

	builds := struct {
		Data struct{ Builds []Build }
	}{}

	json.Unmarshal(resp.GetHttpContentBytes(), &builds)
	return builds.Data.Builds
}
