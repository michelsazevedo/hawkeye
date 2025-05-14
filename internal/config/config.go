package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Settings
type Settings struct {
	Secret          string `yaml:"secret"`
	ApplicationName string `yaml:"application_name"`
	Server          struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	}
}

type Observability struct {
	JaegerEndpoint string `yaml:"jaeger_endpoint"`
}

type EsIndex struct {
	Name       string
	Properties []struct {
		Name     string `yaml:"name"`
		Type     string `yaml:"type"`
		Analyzer string `yaml:"analyzer,omitempty"`
		Index    *bool  `yaml:"index,omitempty"`
	} `yaml:"properties"`
}

// Config ...
type Config struct {
	Elasticsearch struct {
		URL     string    `yaml:"url"`
		Indices []EsIndex `yaml:"indices"`
	} `yaml:"elasticsearch"`

	NATS struct {
		URL string `yaml:"url"`
	} `yaml:"nats"`

	Observability struct {
		JaegerEndpoint string `yaml:"jaeger_endpoint"`
	} `yaml:"observability"`

	Settings Settings `yaml:"settings"`
}

// NewConfig ...
func NewConfig() (*Config, error) {
	absPath, _ := filepath.Abs(".")
	file, err := ExpandEnv(filepath.Join(path.Clean(absPath), "internal/config", "config.yaml"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	yd := yaml.NewDecoder(file)
	err = yd.Decode(cfg)

	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) GetElasticsearchUrls() []string {
	return []string{c.Elasticsearch.URL}
}

func (c *Config) GetElasticSearchIndex(name string) ([]byte, error) {
	for _, index := range c.Elasticsearch.Indices {
		if index.Name == name {
			props := make(map[string]any)

			for _, p := range index.Properties {
				field := map[string]any{
					"type": p.Type,
				}
				if p.Analyzer != "" {
					field["analyzer"] = p.Analyzer
				}
				if p.Index != nil {
					field["index"] = *p.Index
				}
				props[p.Name] = field
			}

			mapping := map[string]any{
				"mappings": map[string]any{
					"properties": props,
				},
			}

			return json.Marshal(mapping)
		}
	}
	return nil, fmt.Errorf("index %s not found", name)
}

func ExpandEnv(configs string) (io.Reader, error) {
	file, err := os.Open(configs)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = file.Close()
	}()

	bufferConfigs := new(bytes.Buffer)
	_, err = bufferConfigs.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	bytesConfigs := []byte(os.ExpandEnv(bufferConfigs.String()))
	return bytes.NewReader(bytesConfigs), nil
}
