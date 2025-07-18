package cmd

import (
    "fmt"
    "log"
    "ai/config"
    "ai/openai"
    "github.com/urfave/cli/v2"
)

func Load(cfg config.Config, client *openai.Client) []*cli.Command {
    var cmds []*cli.Command
    for mainCmd, sub := range cfg.Commands {
        main := &cli.Command{
            Name:        mainCmd,
            Usage:       fmt.Sprintf("Run %s prompts", mainCmd),
            Subcommands: []*cli.Command{},
        }
        for subCmd, ctx := range sub {
            promptText := ctx.SystemPrompt
            main.Subcommands = append(main.Subcommands, &cli.Command{
                Name:  subCmd,
                Usage: promptText,
                Action: func(c *cli.Context) error {
                    log.Printf("▶️ %s/%s", mainCmd, subCmd)
                    fmt.Printf("📝 %s\n", promptText)
                    err := client.AskStream(promptText)
                    if err != nil {
                        log.Printf("❌ Error: %v", err)
                    }
                    fmt.Println("------")
                    return nil
                },
            })
        }
        cmds = append(cmds, main)
    }
    return cmds
}
