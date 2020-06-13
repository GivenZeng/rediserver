package rediserver

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func ParseYAMLFile(filePath string, conf interface{}) error {
	byts, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(byts, conf)
}
