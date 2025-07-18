package config

import (
    "errors"
    "fmt"
    "os"
    "ai/logger"
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

    // Read file
    data, err := os.ReadFile(path)
    if err != nil {
        logger.Warnf("Failed to read config file: %v", err)
    }

    // Parse YAML
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        logger.Warnf("Failed to parse YAML config: %v", err)
    }

    // Environment variable overrides
    if v := os.Getenv("OPENAI_API_KEY"); v != "" {
        logger.Infof("Overriding openai_api_key from environment variable")
        cfg.OpenAIAPIKey = v
    }
    if v := os.Getenv("OPENAI_MODEL"); v != "" {
        logger.Infof("Overriding model from environment variable")
        cfg.Model = v
    }
    if v := os.Getenv("OPENAI_BASE_URL"); v != "" {
        logger.Infof("Overriding base_url from environment variable")
        cfg.BaseURL = v
    }

    // Validation and defaults
    if cfg.OpenAIAPIKey == "" {
        logger.Warnf("OPENAI_API_KEY is required but not set in config or environment")
    }
    if cfg.Model == "" {
        logger.Infof("Model not set, defaulting to gpt-3.5-turbo")
        cfg.Model = "gpt-3.5-turbo"
    }
    if cfg.BaseURL == "" {
        logger.Infof("Base URL not set, defaulting to https://api.openai.com/v1")
        cfg.BaseURL = "https://api.openai.com/v1"
    }
    if cfg.Commands == nil {
        logger.Warnf("No commands found in configuration")
        cfg.Commands = make(map[string]map[string]Context)
    }
    // Validate commands structure
    for mainCmd, subCmds := range cfg.Commands {
        if subCmds == nil {
            logger.Warnf("Main command '%s' has no subcommands", mainCmd)
            continue
        }
        for subCmd, ctx := range subCmds {
            if ctx.SystemPrompt == "" {
                logger.Warnf("Warning: system_prompt is empty for command '%s/%s'", mainCmd, subCmd)
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
            "dev": {
                "logs":     {SystemPrompt: "Explain how to view and interpret application logs."},
                "explain":  {SystemPrompt: "Explain what this code or command does."},
                "cmd":      {SystemPrompt: "Generate a shell command for the described task."},
                "generate": {SystemPrompt: "Generate code or configuration for the described requirement."},
                "debug":    {SystemPrompt: "Suggest debugging steps for the described issue."},
                "review":   {SystemPrompt: "Review the following code for bugs and improvements."},
            },
            "kubernetes": {
                "troubleshoot": {SystemPrompt: "Troubleshoot the described Kubernetes issue."},
            },
            "ai": {
                "prompt_engineer": {SystemPrompt: "Suggest improvements to the following prompt."},
            },
        },
    }
    out, err := yaml.Marshal(cfg)
    if err != nil {
        logger.Infof("Failed to marshal default config: %v", err)
        return err
    }

    // Check if file already exists
    if _, err := os.Stat(path); err == nil {
        logger.Warnf("Config file already exists at %s, not overwriting", path)
        return errors.New(fmt.Sprintf("config file already exists at %s", path))
    } else if !os.IsNotExist(err) {
        logger.Infof("Error checking config file: %v", err)
        return err
    }

    err = os.WriteFile(path, out, 0644)
    if err != nil {
        logger.Infof("Failed to write default config: %v", err)
        return err
    }
    logger.Infof("Default config written to %s", path)
    return nil
}
