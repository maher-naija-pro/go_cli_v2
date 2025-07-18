package openai

import (
    "context"
    "fmt"
    "os"
    "ai/logger"
    "github.com/sashabaranov/go-openai"
)

type Client struct {
    api   *openai.Client
    model string
}

func New(apiKey, model, baseURL string) *Client {
    cfg := openai.DefaultConfig(apiKey)
    cfg.BaseURL = baseURL
    logger.Infof("Initializing OpenAI client with model: %s, baseURL: %s", model, baseURL)
    return &Client{
        api:   openai.NewClientWithConfig(cfg),
        model: model,
    }
}

func (c *Client) AskStream(prompt string) error {
    logger.Debugf("Sending prompt to OpenAI: %s", prompt)
    req := openai.ChatCompletionRequest{
        Model:  c.model,
        Stream: true,
        Messages: []openai.ChatCompletionMessage{
            {Role: openai.ChatMessageRoleUser, Content: prompt},
        },
    }
    stream, err := c.api.CreateChatCompletionStream(context.Background(), req)
    if err != nil {
        logger.Errorf("Failed to create chat completion stream: %v", err)
        return err
    }
    defer stream.Close()
    fmt.Print(" ")
    for {
        resp, err := stream.Recv()
        if err != nil {
            if err.Error() != "EOF" {
                logger.Warnf("Error receiving from stream: %v", err)
            }
            break
        }
        fmt.Print(resp.Choices[0].Delta.Content)
        os.Stdout.Sync()
    }
    fmt.Println()
    logger.Infof("Completed streaming response from OpenAI")
    return nil
}
