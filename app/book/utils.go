package coinbook

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/GeertJohan/go.rice"
)

func ExecuteUnpack(datadir string, box *rice.Box) error {
	return box.Walk("", func(name string, info os.FileInfo, err error) error {
		if name == "" {
			dirname := path.Join(datadir, box.Name())
			return os.MkdirAll(dirname, 0777)
		}
		if info.IsDir() == true {
			dirname := path.Join(datadir, box.Name(), name)
			return os.MkdirAll(dirname, 0777)
		}
		data, err := box.Bytes(name)
		if err != nil {
			return err
		}
		filename := path.Join(datadir, box.Name(), name)
		return ioutil.WriteFile(filename, data, 0777)
	})
}
