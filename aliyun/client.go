package aliyun

import (
	"os"
	"sync"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/joho/godotenv"
)

type AliyunClientMutex struct {
	Mu     sync.Mutex
	Client *cr.Client
}

func NewAliyunClientMutex() *AliyunClientMutex {
	godotenv.Load()
	// collect metrics info
	client, err := cr.NewClientWithAccessKey(os.Getenv("ALIYUN_REGION"), os.Getenv("ALIYUN_AK"), os.Getenv("ALIYUN_SK"))
	c_mutex := AliyunClientMutex{Client: client}

	if err != nil {
		panic(err)
	}

	return &c_mutex
}
