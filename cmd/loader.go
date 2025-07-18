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

    if client == nil {
        log.Println("OpenAI client is not initialized")
        return cmds
    }

    if len(cfg.Commands) == 0 {
        log.Println("No commands found in configuration")
    }

    for mainCmd, sub := range cfg.Commands {
        if len(sub) == 0 {
            log.Printf("No subcommands found for main command '%s'", mainCmd)
        }
        main := &cli.Command{
            Name:        mainCmd,
            Usage:       fmt.Sprintf("Run %s prompts", mainCmd),
            Subcommands: []*cli.Command{},
        }
        for subCmd, ctx := range sub {
            promptText := ctx.SystemPrompt
            if promptText == "" {
                log.Printf("System prompt is empty for %s/%s", mainCmd, subCmd)
            }
            // Capture variables for closure
            capturedMainCmd := mainCmd
            capturedSubCmd := subCmd
            capturedPromptText := promptText

            main.Subcommands = append(main.Subcommands, &cli.Command{
                Name:  capturedSubCmd,
                Usage: capturedPromptText,
                Action: func(c *cli.Context) error {
                    log.Printf("Running command: %s/%s", capturedMainCmd, capturedSubCmd)
                    if capturedPromptText == "" {
                        log.Printf("No prompt text provided for %s/%s", capturedMainCmd, capturedSubCmd)
                        return fmt.Errorf("no prompt text provided for %s/%s", capturedMainCmd, capturedSubCmd)
                    }
                    fmt.Printf("📝 %s\n", capturedPromptText)
                    err := client.AskStream(capturedPromptText)
                    if err != nil {
                        log.Printf("Error from OpenAI client: %v", err)
                        fmt.Println("------")
                        return err
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
