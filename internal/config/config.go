package config

import (
	"os"

	"bytes"

	"github.com/BurntSushi/toml"
)

type Config struct {
	APIURL           string   `toml:"apiURL"`
	CloudflareZoneID string   `toml:"cloudflareZoneID"`
	APIToken         string   `toml:"apiToken"`
	RecordNames      []string `toml:"recordNames"`
	UpdateInterval   int      `toml:"updateInterval"`
}

func Load(filename string) (Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if _, err := toml.Decode(string(file), &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func Save(filename string, config Config) error {
	var buf bytes.Buffer
	writer := toml.NewEncoder(&buf)
	if err := writer.Encode(config); err != nil {
		return err
	}
	file := buf.Bytes()

	return os.WriteFile(filename, file, 0644)
}
