package files

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"strings"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/files/sinastorage"
)

type sinaAdapter struct {
	bucket string
	url    string
	scs    *sinastorage.SCS
}

func newSinaAdapter() *sinaAdapter {
	url := strings.Replace(config.TConfig.SinaDomain, "http://", "", -1)
	url = strings.Replace(url, "/", "", -1)
	s := &sinaAdapter{
		bucket: config.TConfig.SinaBucket,
		url:    url,
	}
	s.scs = &sinastorage.SCS{
		Accessk: config.TConfig.SinaAccessKey,
		Secretk: config.TConfig.SinaSecretKey,
		URI:     url,
	}
	return s
}

func (s *sinaAdapter) createFile(filename string, data []byte, contentType string) error {
	code, _, err := s.scs.PutObjectData(s.bucket, filename, data, "public-read")
	if code != 200 || err != nil {
		return errs.E(errs.FileSaveError, "createFile failed.")
	}
	return nil
}

func (s *sinaAdapter) deleteFile(filename string) error {
	code, _, err := s.scs.DeleteObjectData(s.bucket, filename)
	if code != 204 || err != nil {
		return errs.E(errs.FileDeleteError, "deleteFile failed.")
	}
	return nil
}

func (s *sinaAdapter) getFileData(filename string) ([]byte, error) {
	return s.download(filename)
}

func (s *sinaAdapter) getFileLocation(filename string) string {
	if config.TConfig.FileDirectAccess {
		return fmt.Sprintf("http://%s/%s/%s?formatter=json", s.url, s.bucket, url.QueryEscape(filename))
	}
	return config.TConfig.ServerURL + "/files/" + config.TConfig.AppID + "/" + url.QueryEscape(filename)
}

func (s *sinaAdapter) getFileStream(filename string) (FileStream, error) {
	return nil, errs.E(errs.FileReadError, "no such file or directory")
}

func (s *sinaAdapter) getAdapterName() string {
	return "sinaAdapter"
}

func (s *sinaAdapter) download(filename string) ([]byte, error) {
	path := fmt.Sprintf("http://%s/%s/%s?formatter=json", s.url, s.bucket, url.QueryEscape(filename))
	request, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errs.E(errs.FileReadError, "no such file or directory")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
