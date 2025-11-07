package config

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
	DirPerms   int    `yaml:"dir-perms"`
}
