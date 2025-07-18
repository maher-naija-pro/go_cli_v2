package config

import (
    "errors"
    "fmt"
    "log"
    "os"
    "gopkg.in/yaml.v3"
)

type Context struct {
    SystemPrompt string `yaml:"system_prompt"`
}

type Config struct {
    OpenAIAPIKey string                        `yaml:"openai_api_key"`
    Model        string                        `yaml:"model"`
    BaseURL      string                        `yaml:"base_url"`
    Commands     map[string]map[string]Context `yaml:"commands"`
}

// Load reads the configuration from the given path, applies environment variable overrides,
// and performs validation. Logs and returns on error.
func Load(path string) Config {
    var cfg Config

    // Check if file exists
    fileInfo, err := os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            log.Fatalf("Config file not found: %s", path)
        } else {
            log.Fatalf("Error checking config file: %v", err)
        }
    }
    if fileInfo.IsDir() {
        log.Fatalf("Config path is a directory, not a file: %s", path)
    }

    // Read file
    data, err := os.ReadFile(path)
    if err != nil {
        log.Fatalf("Failed to read config file: %v", err)
    }

    // Parse YAML
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        log.Fatalf("Failed to parse YAML config: %v", err)
    }

    // Environment variable overrides
    if v := os.Getenv("OPENAI_API_KEY"); v != "" {
        log.Printf("Overriding openai_api_key from environment variable")
        cfg.OpenAIAPIKey = v
    }
    if v := os.Getenv("OPENAI_MODEL"); v != "" {
        log.Printf("Overriding model from environment variable")
        cfg.Model = v
    }
    if v := os.Getenv("OPENAI_BASE_URL"); v != "" {
        log.Printf("Overriding base_url from environment variable")
        cfg.BaseURL = v
    }

    // Validation and defaults
    if cfg.OpenAIAPIKey == "" {
        log.Fatal("OPENAI_API_KEY is required but not set in config or environment")
    }
    if cfg.Model == "" {
        log.Printf("Model not set, defaulting to gpt-3.5-turbo")
        cfg.Model = "gpt-3.5-turbo"
    }
    if cfg.BaseURL == "" {
        log.Printf("Base URL not set, defaulting to https://api.openai.com/v1")
        cfg.BaseURL = "https://api.openai.com/v1"
    }
    if cfg.Commands == nil {
        log.Printf("No commands found in configuration")
        cfg.Commands = make(map[string]map[string]Context)
    }
    // Validate commands structure
    for mainCmd, subCmds := range cfg.Commands {
        if subCmds == nil {
            log.Printf("Main command '%s' has no subcommands", mainCmd)
            continue
        }
        for subCmd, ctx := range subCmds {
            if ctx.SystemPrompt == "" {
                log.Printf("Warning: system_prompt is empty for command '%s/%s'", mainCmd, subCmd)
            }
        }
    }

    return cfg
}

// WriteDefault writes a default config file to the given path.
// Returns an error if writing fails.
func WriteDefault(path string) error {
    cfg := Config{
        OpenAIAPIKey: "",
        Model:        "gpt-3.5-turbo",
        BaseURL:      "https://api.openai.com/v1",
        Commands: map[string]map[string]Context{
            "example": {
                "hello": {SystemPrompt: "Say hello in a fun way."},
            },
        },
    }
    out, err := yaml.Marshal(cfg)
    if err != nil {
        log.Printf("Failed to marshal default config: %v", err)
        return err
    }

    // Check if file already exists
    if _, err := os.Stat(path); err == nil {
        log.Printf("Config file already exists at %s, not overwriting", path)
        return errors.New(fmt.Sprintf("config file already exists at %s", path))
    } else if !os.IsNotExist(err) {
        log.Printf("Error checking config file: %v", err)
        return err
    }

    err = os.WriteFile(path, out, 0644)
    if err != nil {
        log.Printf("Failed to write default config: %v", err)
        return err
    }
    log.Printf("Default config written to %s", path)
    return nil
}
