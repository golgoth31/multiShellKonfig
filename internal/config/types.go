package config

type Konfigs struct {
	Konfigs []Konfig `yaml:"kubeConfigs"`
}

type Konfig struct {
	Path string `yaml:"path"`
	ID   string `yaml:"id"`
}
