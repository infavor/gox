package httpx_test

import (
	"github.com/hetianyi/gox/httpx"
	"github.com/hetianyi/gox/logger"
	"github.com/sirupsen/logrus"
	"testing"
)

func init() {
	logger.Init(nil)
}

type Response struct {
	Detail string `json:"detail"`
}

func TestMockGet(t *testing.T) {
	var resp = &Response{}
	err := httpx.Mock().URL("https://hub.docker.com/v2/user/").Get().Success(resp).Error(func(status int, response []byte) {
		logrus.Error("status ", status, ", response: ", string(response))
	}).Do()
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Info(resp.Detail)
	}
}
