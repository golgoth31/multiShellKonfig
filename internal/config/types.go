package config

type Konfigs struct {
	Konfigs []Konfig `yaml:"kubeConfigs"`
}

type Konfig struct {
	Path string `yaml:"path"`
	// ID is a UUID v4 that uniquely identifies this kubeconfig entry and
	// namespaces its on-disk state under contextsPath/<ID>/.
	ID string `yaml:"id"`
}
