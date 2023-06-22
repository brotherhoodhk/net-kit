package config

import (
	"encoding/json"
	"encoding/xml"

	"github.com/go-yaml/yaml"
)

type Host struct {
	Address string   `json:"address" xml:"address" yaml:"address"`
	Port    int      `json:"port" xml:"port" yaml:"port"`
	Name    string   `json:"name" xml:"name" yaml:"name"`
	Id      string   `json:"id" xml:"id" yaml:"id"`
	Tags    []string `json:"tags" xml:"tags" yaml:"tags"`
}

var Parse = make(map[string]func([]byte, *Host) error)
var Format = make(map[string]func(*Host) ([]byte, error))

func yamlparse(data []byte, dst *Host) error {
	return yaml.Unmarshal(data, dst)
}
func xmlparse(data []byte, dst *Host) error {
	return xml.Unmarshal(data, dst)
}
func jsonparse(data []byte, dst *Host) error {
	return json.Unmarshal(data, dst)
}
func yamlformat(src *Host) ([]byte, error) {
	return yaml.Marshal(src)
}

func xmlformat(src *Host) ([]byte, error) {
	return xml.MarshalIndent(src, "", "\t")
}
func jsonformat(src *Host) ([]byte, error) {
	return json.MarshalIndent(src, "", "\t")
}
