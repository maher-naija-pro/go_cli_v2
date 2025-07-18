package main

import (
    "log"
    "os"
    "ai/cmd"
    "ai/config"
    "ai/openai"
    "github.com/urfave/cli/v2"
)

func main() {
    cfg := config.Load("config.yaml")
    client := openai.New(cfg.OpenAIAPIKey, cfg.Model, cfg.BaseURL)
    app := &cli.App{
        Name:     "ai",
        Usage:    "Run OpenAI prompts from config.yaml",
        Commands: append([]*cli.Command{cmd.InitCommand}, cmd.Load(cfg, client)...),
    }
    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}
