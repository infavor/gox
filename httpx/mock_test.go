package httpx_test

import (
	"encoding/json"
	"github.com/hetianyi/gox/httpx"
	"github.com/hetianyi/gox/logger"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func init() {
	logger.Init(nil)
}

type Response struct {
	Detail string `json:"detail"`
}

func TestMockGet1(t *testing.T) {
	result, _, err := httpx.Mock().URL("https://hub.docker.com/v2/user/").Get().Success(&Response{}, 200, 401).Error(func(status int, response []byte) {
		logrus.Error("status ", status, ", response: ", string(response))
	}).Do()
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Info(result.(*Response).Detail)
	}
}

func TestMockGet2(t *testing.T) {
	result, _, err := httpx.Mock().URL("https://hub.docker.com/v2/user/").Get().Success("", 200, 401).Error(func(status int, response []byte) {
		logrus.Error("status ", status, ", response: ", string(response))
	}).Do()
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Info(result)
	}
}

type Edition struct {
	Label string `json:"label"`
	Name  string `json:"name"`
}

type Response1 struct {
	Editions []Edition
}

func TestMockGet3(t *testing.T) {
	result, _, err := httpx.Mock().URL("https://hub.docker.com/api/content/v1/platforms").Get().Success(&Response1{}).Error(func(status int, response []byte) {
		logrus.Error("status ", status, ", response: ", string(response))
	}).Do()
	if err != nil {
		logrus.Error(err)
	} else {
		bs, _ := json.MarshalIndent(result, "", " ")
		logrus.Info(string(bs))
	}
}

func TestMockGet4(t *testing.T) {
	result, _, err := httpx.Mock().URL("https://hub.docker.com/v2/repositories/library/redis/tags/").
		Parameter("page_size", "25").
		Parameter("page", "1").
		Get().Success(new(map[string]interface{})).Error(func(status int, response []byte) {
		logrus.Error("status ", status, ", response: ", string(response))
	}).Do()
	if err != nil {
		logrus.Error(err)
	} else {
		bs, _ := json.MarshalIndent(result, "", " ")
		logrus.Info(string(bs))
	}
}

func TestMockGet5(t *testing.T) {
	result, _, err := httpx.Mock().URL("https://hub.docker.com/api/content/v1/products/images/redis").
		Get().Success("").Error(func(status int, response []byte) {
		logrus.Error("status ", status, ", response: ", string(response))
	}).Do()
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Info(result)
	}
}

func TestMockPost1(t *testing.T) {
	logrus.Info("start ", time.Now().UTC())
	result, _, err := httpx.Mock().URL("http://cdn.yifuls.com/api/cdn/detail").
		Post().
		Body(map[string][]string{"id": {"10004"}}).
		Success(new(map[string]interface{})).Error(func(status int, response []byte) {
		logrus.Error("status ", status, ", response: ", string(response))
	}).Do()
	if err != nil {
		logrus.Error(err)
	} else {
		bs, _ := json.MarshalIndent(result, "", " ")
		logrus.Info(string(bs))
	}
	logrus.Info("end ", time.Now().UTC())
}
