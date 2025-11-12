package config

type Config struct {
	InputFile      string `yaml:"input-file"`
	OutputFile     string `yaml:"output-file"`
	DirPermissions *int   `yaml:"dir-permissions,omitempty"`
}
