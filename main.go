package main

import (
    "os"
    "ai/cmd"
    "ai/config"
    "ai/openai"
    "ai/logger"
    "github.com/urfave/cli/v2"
)

func main() {
    // Initialize logger to stderr with INFO level
    logger.InitLogger(os.Stderr, logger.INFO)

    configPath := os.Getenv("AI_CONFIG_PATH")
    if configPath == "" {
        // Default to ~/.ai/config.yaml if AI_CONFIG_PATH is not set
        home, err := os.UserHomeDir()
        if err != nil {
            logger.Infof("Failed to get user home directory: %v", err)
        }
        configPath = home + string(os.PathSeparator) + ".ai" + string(os.PathSeparator) + "config.yaml"
        logger.Infof("AI_CONFIG_PATH not set, using default: %s", configPath)
    } else {
        logger.Infof("Using config path from AI_CONFIG_PATH: %s", configPath)
    }

    // Check if config file exists
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        logger.Warnf("Config file not found at path: %s", configPath)
    }

    cfg := config.Load(configPath)
    client := openai.New(cfg.OpenAIAPIKey, cfg.Model, cfg.BaseURL)
    if client == nil {
        logger.Infof("Failed to initialize OpenAI client. Exiting.")
    }

    commands := cmd.Load(cfg, client)
    if len(commands) == 0 {
        logger.Infof("No commands loaded from configuration.")
    }

    app := &cli.App{
        Name:     "ai",
        Usage:    "Run OpenAI prompts from config.yaml",
        Commands: append([]*cli.Command{cmd.InitCommand}, cmd.Load(cfg, client)...),
    }

    if err := app.Run(os.Args); err != nil {
        logger.Fatalf("Application error: %v", err)
    }
}
