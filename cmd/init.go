package cmd

import (
    "fmt"
    "os"
    "ai/config"
    "github.com/urfave/cli/v2"
)

var InitCommand = &cli.Command{
    Name:  "init",
    Usage: "Generate a default config.yaml",
    Flags: []cli.Flag{
        &cli.StringFlag{
            Name:  "output",
            Value: "config.yaml",
            Usage: "Path to output file",
        },
    },
    Action: func(c *cli.Context) error {
        path := c.String("output")
        if _, err := os.Stat(path); err == nil {
            return cli.Exit("File already exists", 1)
        }
        if err := config.WriteDefault(path); err != nil {
            return cli.Exit(fmt.Sprintf("failed to write config: %v", err), 1)
        }
        fmt.Printf("✅ created default config at %s\n", path)
        return nil
    },
}
