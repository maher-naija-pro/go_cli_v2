package cmd

import (
    "fmt"
    "os"
    "path/filepath"
    "ai/config"
    "github.com/urfave/cli/v2"
)

var InitCommand = &cli.Command{
    Name:  "init",
    Usage: "Generate a default config.yaml in ~/.ai/",
    Flags: []cli.Flag{
        &cli.StringFlag{
            Name:  "output",
            Value: "", // If not set, use ~/.ai/config.yaml
            Usage: "Path to output file (default: ~/.ai/config.yaml)",
        },
    },
    Action: func(c *cli.Context) error {
        outputPath := c.String("output")
        if outputPath == "" {
            home, err := os.UserHomeDir()
            if err != nil {
                return cli.Exit(fmt.Sprintf("failed to get user home directory: %v", err), 1)
            }
            aiDir := filepath.Join(home, ".ai")
            if err := os.MkdirAll(aiDir, 0700); err != nil {
                return cli.Exit(fmt.Sprintf("failed to create directory %s: %v", aiDir, err), 1)
            }
            outputPath = filepath.Join(aiDir, "config.yaml")
        } else {
            dir := filepath.Dir(outputPath)
            if err := os.MkdirAll(dir, 0700); err != nil {
                return cli.Exit(fmt.Sprintf("failed to create directory %s: %v", dir, err), 1)
            }
        }

        if _, err := os.Stat(outputPath); err == nil {
            return cli.Exit(fmt.Sprintf("File already exists at %s", outputPath), 1)
        }
        if err := config.WriteDefault(outputPath); err != nil {
            return cli.Exit(fmt.Sprintf("failed to write config: %v", err), 1)
        }
        fmt.Printf("Created default config at %s\n", outputPath)
        return nil
    },
}
