package config

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
	DirPerms   uint32 `yaml:"dir-perms"`
}
