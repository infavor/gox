package httpx_test

import (
	"encoding/json"
	"github.com/infavor/gox/file"
	"github.com/infavor/gox/httpx"
	"github.com/infavor/gox/logger"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
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
		logger.Error("status ", status, ", response: ", string(response))
	}).Do()
	if err != nil {
		logger.Error(err)
	} else {
		logger.Info(result.(*Response).Detail)
	}
}

func TestMockGet2(t *testing.T) {
	result, _, err := httpx.Mock().URL("https://hub.docker.com/v2/user/").Get().Success("", 200, 401).Error(func(status int, response []byte) {
		logger.Error("status ", status, ", response: ", string(response))
	}).Do()
	if err != nil {
		logger.Error(err)
	} else {
		logger.Info(result)
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
		logger.Error("status ", status, ", response: ", string(response))
	}).Do()
	if err != nil {
		logger.Error(err)
	} else {
		bs, _ := json.MarshalIndent(result, "", " ")
		logger.Info(string(bs))
	}
}

func TestMockGet4(t *testing.T) {
	result, _, err := httpx.Mock().URL("https://hub.docker.com/v2/repositories/library/redix/tags/").
		Parameter("page_size", "25").
		Parameter("page", "1").
		Get().Success(new(map[string]interface{})).Error(func(status int, response []byte) {
		logger.Error("status ", status, ", response: ", string(response))
	}).Do()
	if err != nil {
		logger.Error(err)
	} else {
		bs, _ := json.MarshalIndent(result, "", " ")
		logger.Info(string(bs))
	}
}

func TestMockGet5(t *testing.T) {
	result, _, err := httpx.Mock().URL("https://hub.docker.com/api/content/v1/products/images/redix").
		Get().Success("").Error(func(status int, response []byte) {
		logger.Error("status ", status, ", response: ", string(response))
	}).Do()
	if err != nil {
		logger.Error(err)
	} else {
		logger.Info(result)
	}
}

func TestMockPost1(t *testing.T) {
	logger.Info("start ", time.Now().UTC())
	result, _, err := httpx.Mock().URL("http://cdn.yifuls.com/api/cdn/detail").
		Post().
		Body(map[string][]string{"id": {"10004"}}).
		Success(new(map[string]interface{})).Error(func(status int, response []byte) {
		logger.Error("status ", status, ", response: ", string(response))
	}).Do()
	if err != nil {
		logger.Error(err)
	} else {
		bs, _ := json.MarshalIndent(result, "", " ")
		logger.Info(string(bs))
	}
	logger.Info("end ", time.Now().UTC())
}

func TestMockPost2(t *testing.T) {
	logger.Info("start ", time.Now().UTC())
	result, _, err := httpx.Mock().URL("http://cdn.yifuls.com/api/cdn/detail").
		Post().
		Body(map[string][]string{"id": {"10004"}}).
		Success(nil).Error(func(status int, response []byte) {
		logger.Error("status ", status, ", response: ", string(response))
	}).Do()
	if err != nil {
		logger.Error(err)
	} else {
		bs, _ := json.MarshalIndent(result, "", " ")
		logger.Info(string(bs))
	}
	logger.Info("end ", time.Now().UTC())
}

func TestMockPost3(t *testing.T) {
	out, err := file.CreateFile("D:\\tmp\\轻松一刻语音版-所谓塑料同事-下班缘分就尽了.mp3")
	if err != nil {
		logger.Fatal(err)
	}
	defer out.Close()
	logger.Info("start ", time.Now().UTC())
	result, _, err := httpx.Mock().URL("http://mobilepics.ws.126.net/UIg9y7iEZIrCxoikpo3HeNeikjLMqkV7%3D%3DFTRPUTL6.mp3").
		Get().
		Success(out).Error(func(status int, response []byte) {
		logger.Error("status ", status, ", response: ", string(response))
	}).Do()
	if err != nil {
		logger.Error(err)
	} else {
		bs, _ := json.MarshalIndent(result, "", " ")
		logger.Info(string(bs))
	}
	logger.Info("end ", time.Now().UTC())
}

func TestMultipart(t *testing.T) {

	httpClient := &http.Client{
		Timeout: time.Second * 20,
	}

	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()

		o, _ := m.CreateFormField("Name")
		o.Write([]byte("zhangsan"))
		o, _ = m.CreateFormField("Name")
		o.Write([]byte("lisi"))

		fi, _ := file.GetFile("E:\\godfs-storage\\123.zip")
		defer fi.Close()
		o, _ = m.CreateFormFile("secrets", "123.zip")
		io.Copy(o, fi)
	}()

	req, err := http.NewRequest("POST", "http://localhost:8001/upload", r)
	if err != nil {
		logger.Fatal(err)
	}
	req.Header.Add("Content-Type", m.FormDataContentType())

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Fatal(err)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info(string(bs))
}

func TestMockUpload(t *testing.T) {
	httpx.SetTTL(0)
	logger.Info("start ", time.Now().UTC())
	result, _, err := httpx.Mock().URL("http://localhost:8001/upload").
		Success(nil).
		Error(func(status int, response []byte) {
			logger.Error("status ", status, ", response: ", string(response))
		}).
		Multipart(func(writer *multipart.Writer) {
			o, _ := writer.CreateFormField("Name")
			o.Write([]byte("zhangsan"))
			o, _ = writer.CreateFormField("Name")
			o.Write([]byte("lisi"))

			fi, _ := file.GetFile("F:\\Software\\fastdfs_client_v1.24.jar")
			defer fi.Close()
			o, _ = writer.CreateFormFile("secrets", filepath.Base("F:\\Software\\fastdfs_client_v1.24.jar"))
			io.Copy(o, fi)
		}).
		Do()
	if err != nil {
		logger.Error(err)
	} else {
		bs, _ := json.MarshalIndent(result, "", " ")
		logger.Info(string(bs))
	}
	logger.Info("end ", time.Now().UTC())
}
