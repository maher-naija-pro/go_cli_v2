package openai

import (
    "context"
    "fmt"
    "os"
    "github.com/sashabaranov/go-openai"
)

type Client struct {
    api   *openai.Client
    model string
}

func New(apiKey, model, baseURL string) *Client {
    cfg := openai.DefaultConfig(apiKey)
    cfg.BaseURL = baseURL
    return &Client{
        api:   openai.NewClientWithConfig(cfg),
        model: model,
    }
}

func (c *Client) AskStream(prompt string) error {
    req := openai.ChatCompletionRequest{
        Model:  c.model,
        Stream: true,
        Messages: []openai.ChatCompletionMessage{
            {Role: openai.ChatMessageRoleUser, Content: prompt},
        },
    }
    stream, err := c.api.CreateChatCompletionStream(context.Background(), req)
    if err != nil {
        return err
    }
    defer stream.Close()
    fmt.Print("📥 ")
    for {
        resp, err := stream.Recv()
        if err != nil {
            break
        }
        fmt.Print(resp.Choices[0].Delta.Content)
        os.Stdout.Sync()
    }
    fmt.Println()
    return nil
}
