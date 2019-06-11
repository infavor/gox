package httpx_test

import (
	"encoding/json"
	"github.com/hetianyi/gox/file"
	"github.com/hetianyi/gox/httpx"
	"github.com/hetianyi/gox/logger"
	"github.com/sirupsen/logrus"
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

func TestMockPost2(t *testing.T) {
	logrus.Info("start ", time.Now().UTC())
	result, _, err := httpx.Mock().URL("http://cdn.yifuls.com/api/cdn/detail").
		Post().
		Body(map[string][]string{"id": {"10004"}}).
		Success(nil).Error(func(status int, response []byte) {
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

func TestMockPost3(t *testing.T) {
	out, err := file.CreateFile("D:\\tmp\\response.pdf")
	if err != nil {
		logrus.Fatal(err)
	}
	defer out.Close()
	logrus.Info("start ", time.Now().UTC())
	result, _, err := httpx.Mock().URL("https://doctool.cbim.org.cn/docgen/warehouse/2631c7fca887dc9c097188b40232f1e8.pdf?filename=%E4%B8%AD%E5%9B%BD%E5%86%9C%E4%B8%9A%E7%A7%91%E6%8A%80%E5%9B%BD%E9%99%85%E4%BA%A4%E6%B5%81%E4%B8%AD%E5%BF%83.pdf").
		Get().
		Success(out).Error(func(status int, response []byte) {
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
		logrus.Fatal(err)
	}
	req.Header.Add("Content-Type", m.FormDataContentType())

	resp, err := httpClient.Do(req)
	if err != nil {
		logrus.Fatal(err)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info(string(bs))
}

func TestMockUpload(t *testing.T) {
	httpx.SetTTL(0)
	logrus.Info("start ", time.Now().UTC())
	result, _, err := httpx.Mock().URL("http://localhost:8001/upload").
		Success(nil).
		Error(func(status int, response []byte) {
			logrus.Error("status ", status, ", response: ", string(response))
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
		logrus.Error(err)
	} else {
		bs, _ := json.MarshalIndent(result, "", " ")
		logrus.Info(string(bs))
	}
	logrus.Info("end ", time.Now().UTC())
}
