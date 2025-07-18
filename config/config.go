package config

import (
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

func Load(path string) Config {
    var cfg Config
    if _, err := os.Stat(path); os.IsNotExist(err) {
        log.Fatalf("❌ config file not found: %s", path)
    }
    data, err := os.ReadFile(path)
    if err != nil {
        log.Fatalf("❌ failed to read config: %v", err)
    }
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        log.Fatalf("❌ failed to parse yaml: %v", err)
    }
    if v := os.Getenv("OPENAI_API_KEY"); v != "" {
        cfg.OpenAIAPIKey = v
    }
    if v := os.Getenv("OPENAI_MODEL"); v != "" {
        cfg.Model = v
    }
    if v := os.Getenv("OPENAI_BASE_URL"); v != "" {
        cfg.BaseURL = v
    }
    if cfg.OpenAIAPIKey == "" {
        log.Fatal("❌ OPENAI_API_KEY is required")
    }
    if cfg.Model == "" {
        cfg.Model = "gpt-3.5-turbo"
    }
    if cfg.BaseURL == "" {
        cfg.BaseURL = "https://api.openai.com/v1"
    }
    return cfg
}

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
        return err
    }
    return os.WriteFile(path, out, 0644)
}
