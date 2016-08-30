package files

import (
	"errors"
	"net/url"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/storage"
	"gopkg.in/mgo.v2"
)

type gridStoreAdapter struct {
	gfs *mgo.GridFS
}

func newGridStoreAdapter() *gridStoreAdapter {
	g := &gridStoreAdapter{}
	g.gfs = storage.TomatoDB.MongoDatabase.GridFS("fs")
	return g
}

func (g *gridStoreAdapter) createFile(filename string, data []byte, contentType string) error {
	file, err := g.gfs.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	n, err := file.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return errors.New("createFile failed")
	}

	if contentType != "" {
		file.SetContentType(contentType)
	}

	return nil
}

func (g *gridStoreAdapter) deleteFile(filename string) error {
	return g.gfs.Remove(filename)
}

func (g *gridStoreAdapter) getFileData(filename string) []byte {
	file, err := g.gfs.Open(filename)
	if err != nil {
		return []byte{}
	}
	defer file.Close()

	data := []byte{}
	buf := make([]byte, 1024)
	for {
		n, _ := file.Read(buf)
		if n == 0 {
			break
		}
		data = append(data, buf[:n]...)
	}

	return data
}

func (g *gridStoreAdapter) getFileLocation(filename string) string {
	return config.TConfig.ServerURL + "/files/" + config.TConfig.AppID + "/" + url.QueryEscape(filename)
}

func (g *gridStoreAdapter) getFileStream(filename string) (FileStream, error) {
	file, err := g.gfs.Open(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}
