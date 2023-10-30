package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var DefaultConfigFilename = "config.json"

type Config struct {
	DebugLevel     string `json:"debugLevel,omitempty"` // Logging level for all subsystems {trace, debug, info, warn, error, critical}
	LogDir         string `json:"logDir,omitempty"`
	MaxLogFiles    int    `json:"maxLogFiles,omitempty"`    // Maximum logfiles to keep (0 for no rotation)
	MaxLogFileSize int    `json:"maxLogFileSize,omitempty"` // Maximum logfile size in MB

	MerkleRootAPI string `json:"merkleRootAPIhttp,omitempty"`

	AutherKeystore string `json:"autherKeystore,omitempty"`
	AutherPassword string `json:"autherPassword,omitempty"`

	MerkleRootListenerInterval int64 `json:"merkleRootListenerInterval,omitempty"` //seconds

	ChainRPC    string `json:"chainRPC,omitempty"`
	CredaOracle string `json:"credaOracle,omitempty"`
}

func DefaultConfig() Config {
	return Config{
		DebugLevel:                 "Info",
		LogDir:                     "",
		MaxLogFiles:                1,
		MaxLogFileSize:             100,
		MerkleRootAPI:              "",
		AutherKeystore:             "auther.keystore",
		AutherPassword:             "",
		MerkleRootListenerInterval: 60,
		ChainRPC:                   "",
		CredaOracle:                "",
	}
}

func LoadConfig() (*Config, error) {
	preCfg := DefaultConfig()

	file, err := ioutil.ReadFile(DefaultConfigFilename)
	if err != nil {
		return nil, err
	}
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	err = json.Unmarshal(file, &preCfg)
	if err != nil {
		return nil, err
	}

	if err := preCfg.ValidateConfig(); err != nil {
		return nil, err
	}
	return &preCfg, nil
}

func (cfg *Config) ValidateConfig() error {
	cfg.LogDir = CleanAndExpandPath(cfg.LogDir)

	if cfg.ChainRPC == "" {
		return errors.New("ChainRPC is empty")
	}

	if cfg.CredaOracle == "" {
		return errors.New("CredaOracle is empty")
	}
	return nil
}

// CleanAndExpandPath expands environment variables and leading ~ in the
// passed path, cleans the result, and returns it.
// This function is taken from https://github.com/btcsuite/btcd
func CleanAndExpandPath(path string) string {
	if path == "" {
		return ""
	}

	// Expand initial ~ to OS specific home directory.
	if strings.HasPrefix(path, "~") {
		var homeDir string
		u, err := user.Current()
		if err == nil {
			homeDir = u.HomeDir
		} else {
			homeDir = os.Getenv("HOME")
		}

		path = strings.Replace(path, "~", homeDir, 1)
	}

	// NOTE: The os.ExpandEnv doesn't work with Windows-style %VARIABLE%,
	// but the variables can still be expanded via POSIX-style $VARIABLE.
	return filepath.Clean(os.ExpandEnv(path))
}
