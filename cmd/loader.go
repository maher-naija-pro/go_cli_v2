package cmd

import (
	"ai/config"
	"ai/logger"
	"ai/openai"
	"fmt"

	"github.com/urfave/cli/v2"
)

func Load(cfg config.Config, client *openai.Client) []*cli.Command {
	var cmds []*cli.Command

	if client == nil {
		logger.Warnf("OpenAI client is not initialized")
		return cmds
	}

	if len(cfg.Commands) == 0 {
		logger.Warnf("No commands found in configuration")
	}

	for mainCmd, sub := range cfg.Commands {
		if len(sub) == 0 {
			logger.Warnf("No subcommands found for main command '%s'", mainCmd)
		}
		main := &cli.Command{
			Name:        mainCmd,
			Usage:       fmt.Sprintf("Run %s prompts", mainCmd),
			Subcommands: []*cli.Command{},
		}
		for subCmd, ctx := range sub {
			promptText := ctx.SystemPrompt
			if promptText == "" {
				logger.Warnf("System prompt is empty for %s/%s", mainCmd, subCmd)
			}
			// Capture variables for closure
			capturedMainCmd := mainCmd
			capturedSubCmd := subCmd
			capturedPromptText := promptText

			main.Subcommands = append(main.Subcommands, &cli.Command{
				Name:  capturedSubCmd,
				Usage: capturedPromptText,
				Action: func(c *cli.Context) error {
					logger.Infof("Running command: %s/%s", capturedMainCmd, capturedSubCmd)
					if capturedPromptText == "" {
						logger.Warnf("No prompt text provided for %s/%s", capturedMainCmd, capturedSubCmd)
						return fmt.Errorf("no prompt text provided for %s/%s", capturedMainCmd, capturedSubCmd)
					}
					userPrompt := c.Args().Get(0)
					if userPrompt == "" {
						logger.Warnf("No user prompt provided for %s/%s", capturedMainCmd, capturedSubCmd)
						fmt.Printf("Usage: %s %s <your prompt>\n", capturedMainCmd, capturedSubCmd)
						return fmt.Errorf("no user prompt provided for %s/%s", capturedMainCmd, capturedSubCmd)
					}
					fmt.Printf("📝 %s\n", capturedPromptText)
					err := client.AskStream(capturedPromptText, userPrompt)
					if err != nil {
						logger.Infof("Error from OpenAI client: %v", err)
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
