package config

import (
	"encoding/json"
	"log"
	"os"
)

// ConfigFile json file with configuration
const ConfigFile = "gologdog.json"

// GetDefaultConfig get default config
func GetDefaultConfig() Config {
	return Config{
		UseTemp: true,
		Format:  "%h:%m:%s.%n %x %i %t",
		Target: Target{
			LogdogFormat: "logdogformat",
			Logdog:       []string{"NRVC_SRV"},
		},
		Key: Key{
			Installed: "privatekey2048_GL.pem",
			Factory:   "logenc_factorykey.pem",
		},
	}
}

// GetConfig get configs from file
func GetConfig(filename string) (config Config, err error) {
	config = GetDefaultConfig()
	file, err := os.Open(filename)
	if err != nil {
		return
	}

	dec := json.NewDecoder(file)
	if err = dec.Decode(&config); err != nil {
		log.Fatal(err)
	}

	return config, err
}
