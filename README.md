# AI CLI Tool

A flexible CLI for running prompts against any OpenAI-compatible API endpoint, organized by domain and subcommand, using a YAML configuration file. Easily extendable for DevOps, Kubernetes, and more.

## Features
- Run prompts from the command line using any OpenAI-compatible API (OpenAI, Azure OpenAI, local LLMs, etc.)
- Organize prompts by domain and subcommand
- Stream responses directly in your terminal
- Easily configurable via YAML
- Supports environment variable overrides for sensitive data

## Installation

1. **Clone the repository:**
   ```sh
   git clone <repo-url>
   cd <repo-directory>
   ```
2. **Build the CLI:**
   ```sh
   go build -o ai main.go
   ```

## Configuration


The CLI uses a `config.yaml` file to define prompts and OpenAI settings. By default, it looks for the config at `~/.ai/config.yaml`. You can generate a template with:

```sh
./ai init
```

### Example `config.yaml`
```yaml
openai_api_key: "sk-..." # Your OpenAI API key
model: "gpt-4"
base_url: "https://api.openai.com/v1"

commands:
  devops:
    cicd_expert:
      system_prompt: "Explain CI/CD with practical examples."
    terraform_guru:
      system_prompt: "Explain how remote backends work in Terraform."
  kubernetes:
    architect:
      system_prompt: "How does the Kubernetes HPA controller work?"
    security:
      system_prompt: "How do Kubernetes NetworkPolicies isolate traffic?"
```

- `openai_api_key`: Your OpenAI API key (can also be set via the `OPENAI_API_KEY` environment variable)
- `model`: OpenAI model to use (e.g., `gpt-4`, `gpt-3.5-turbo`)
- `base_url`: API endpoint (default: `https://api.openai.com/v1`)
- `commands`: Hierarchical structure for organizing prompts

## Usage

### List available commands
```sh
./ai --help
```

### Run a prompt
```sh
./ai <domain> <subcommand>
# Example:
./ai devops cicd_expert
```

### Generate a default config
```sh
./ai init
```

## Environment Variables

| Variable           | Description                                         | Example Value                        |
|--------------------|-----------------------------------------------------|--------------------------------------|
| `AI_CONFIG_PATH`   | Path to the config file (defaults to `config.yaml`) | `/home/user/myconfig.yaml`           |
| `OPENAI_API_KEY`   | Overrides the `openai_api_key` in config.yaml       | `sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxx`    |
| `OPENAI_MODEL`     | Overrides the `model` in config.yaml                | `gpt-4`                              |
| `OPENAI_BASE_URL`  | Overrides the `base_url` in config.yaml             | `https://api.openai.com/v1`          |

## Dependencies
- Go 1.21+
- [urfave/cli](https://github.com/urfave/cli)
- [sashabaranov/go-openai](https://github.com/sashabaranov/go-openai)
- [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3)

## License
MIT 