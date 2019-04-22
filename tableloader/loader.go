package tableloader

import (
	"io/ioutil"
	"path"
	"strings"
)

var DataDir = "data"

func Load(name string) []byte {
	data, err := ioutil.ReadFile(path.Join(DataDir, strings.ToLower(name) + ".bytes"))
	if err != nil {
		panic(err)
	}
	return data
}
