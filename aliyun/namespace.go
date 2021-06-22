package aliyun

import (
	"encoding/json"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
)

type Namespace struct {
	NamespaceStatus string
	Namespace       string
	AuthorizeType   string
}

func (c_mutex *AliyunClientMutex) NamespacesList() []Namespace {
	c_mutex.Mu.Lock()
	defer c_mutex.Mu.Unlock()

	getNSResp, err := c_mutex.Client.GetNamespaceList(cr.CreateGetNamespaceListRequest())
	if err != nil {
		// Handle exceptions
		panic(err)
	}

	namespaces := struct {
		Data struct{ Namespaces []Namespace }
	}{}

	json.Unmarshal(getNSResp.GetHttpContentBytes(), &namespaces)
	return namespaces.Data.Namespaces
}
