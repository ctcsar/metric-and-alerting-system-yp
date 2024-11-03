package files

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/storage"
)

type MyFile struct {
	File    *os.File
	Path    string
	Content *storage.Storage
}

func NewFile() *MyFile {
	return &MyFile{}
}

func (f *MyFile) WriteFile(m *storage.Storage, filePath string) {
	f.Path = filePath
	f.Content = m

	defer f.File.Close()

	file, err := os.OpenFile(f.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)

	if err != nil {
		fmt.Println(err)
	}
	data, err := json.MarshalIndent(f.Content, "", "  ")

	if err != nil {
		fmt.Println(err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		fmt.Println(err)
	}

	if _, err := file.Write(data); err != nil {
		fmt.Println(err)
	}

}

func (f *MyFile) ReadFromFile(filePath string, metrics *storage.Storage) {

	f.Path = filePath

	file, err := os.OpenFile(f.Path, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	data, err := os.ReadFile(f.Path)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, metrics)
	if err != nil {
		fmt.Println(err)
	}
}
