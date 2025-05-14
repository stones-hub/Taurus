package app

import (
	"Taurus/config"
	"Taurus/pkg/util"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// initialize calls the initialization functions of all modules
func initialize(configPath string, env string) {
	// initialize environment variables, if empty, do not load
	err := godotenv.Load(env)
	if err != nil {
		log.Printf("Error loading .env file: %v\n", err.Error())
	}

	// load application configuration file
	log.Printf("Loading application configuration file: %s", configPath)
	loadConfig(configPath)

	// print application configuration
	if config.Core.PrintEnable {
		log.Println("Configuration:", util.ToJsonString(config.Core))
	}

	InitialzeLog()
	InitializeDB()
	InitializeRedis()
	InitializeTemplates()
	InitializeCron()
	InitializeInjector()
	InitializeWebsocket()
	InitializeMCP()
}

// loadConfig reads and parses configuration files from a directory or a single file
func loadConfig(path string) {
	info, err := os.Stat(path)
	if err != nil {
		log.Fatalf("Failed to access config path: %v\n", err)
	}

	if info.IsDir() {
		// Recursively load all configuration files in the directory
		err := filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				log.Printf("Error accessing file %s: %v\n", filePath, err)
				return nil
			}

			// Skip directories
			if fileInfo.IsDir() {
				return nil
			}

			loadConfigFile(filePath)
			return nil
		})
		if err != nil {
			log.Fatalf("Failed to walk through config directory: %v\n", err)
		}
	} else {
		// Load a single configuration file
		loadConfigFile(path)
	}

	log.Println("Configuration loaded successfully")
}

// loadConfigFile loads a single configuration file based on its extension
func loadConfigFile(filePath string) {
	ext := filepath.Ext(filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Failed to open config file: %v\n", err)
		return
	}
	// Replace placeholders with environment variables
	content := replacePlaceholders(string(data))

	switch ext {
	case ".json":
		err = json.Unmarshal([]byte(content), &config.Core)
		if err != nil {
			log.Printf("Failed to parse JSON config file: %s; error: %v\n", filePath, err)
		}
	case ".yaml", ".yml":
		err = yaml.Unmarshal([]byte(content), &config.Core)
		if err != nil {
			log.Printf("Failed to parse YAML config file: %s; error: %v\n", filePath, err)
		}
	case ".toml":
		_, err = toml.Decode(content, &config.Core)
		if err != nil {
			log.Printf("Failed to parse TOML config file: %s; error: %v\n", filePath, err)
		}
	default:
		log.Printf("Unsupported config file format: %s\n", filePath)
	}
}

// replacePlaceholders replaces placeholders in the config content with environment variables
func replacePlaceholders(content string) string {
	re := regexp.MustCompile(`\$\{(\w+):([^}]+)\}`)
	return re.ReplaceAllStringFunc(content, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) == 3 {
			envVar := parts[1]
			defaultValue := parts[2]
			if value, exists := os.LookupEnv(envVar); exists {
				return value
			}
			return defaultValue
		}
		return match
	})
}
