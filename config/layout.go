package config

// Config main struct
type Config struct {
	UseTemp bool   `json:"use_temp"`
	Format  string `json:"format"`
	Target  Target `json:"target"`
	Key     Key    `json:"key"`
}

// Target all target to parse
type Target struct {
	LogdogFormat string   `json:"logdog_format"`
	Logdog       []string `json:"logdog"`
}

// Key keys for decryption
type Key struct {
	Installed string `json:"installed"`
	Factory   string `json:"factory"`
}
