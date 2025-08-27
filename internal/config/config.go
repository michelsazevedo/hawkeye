package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"

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

//go:embed config.yaml
var embeddedConfig []byte

// NewConfig ...
func NewConfig() (*Config, error) {
	data := []byte(os.ExpandEnv(string(embeddedConfig)))

	configs := &Config{}

	if err := yaml.Unmarshal(data, configs); err != nil {
		return nil, err
	}

	return configs, nil
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
