package aliyun

import (
	"sync"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
)

type AliyunClientMutex struct {
	Mu     *sync.Mutex
	Client *cr.Client
}

func NewAliyunClientMutex() AliyunClientMutex {
	// collect metrics info
	client, err := cr.NewClientWithAccessKey("cn-hangzhou", "ABC", "DEF")
	c_mutex := AliyunClientMutex{Client: client}

	if err != nil {
		panic(err)
	}

	return c_mutex
}
