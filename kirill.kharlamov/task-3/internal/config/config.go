package config

type Config struct {
	InputFile      string  `yaml:"input-file"`
	OutputFile     string  `yaml:"output-file"`
	DirPermissions *uint32 `yaml:"dir-permissions,omitempty"`
}
