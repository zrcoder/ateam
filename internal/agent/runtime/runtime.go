package runtime

import (
	"bufio"
	"context"
	"os/exec"
	"strings"
	"time"
)

// Provider represents a discovered agent CLI
type Provider struct {
	Name string // CLI name: "claude", "codex", "opencode"
	Path string // Absolute path to the executable
}

// Model represents a supported model
type Model struct {
	ID       string
	Label    string
	Provider string
	Default  bool
}

// All supported provider CLI names
var providerNames = []string{
	"claude", "codex", "opencode", "openclaw", "hermes",
	"gemini", "pi", "cursor-agent", "copilot", "kimi", "kiro-cli",
}

// DiscoverProviders finds all available agent CLIs on PATH
func DiscoverProviders() []Provider {
	var found []Provider
	for _, name := range providerNames {
		if path, err := exec.LookPath(name); err == nil {
			found = append(found, Provider{Name: name, Path: path})
		}
	}
	return found
}

// GetModels returns the supported models for a provider
func GetModels(ctx context.Context, provider string) []Model {
	switch provider {
	case "claude":
		return claudeModels()
	case "codex":
		return codexModels()
	case "opencode":
		return opencodeModels()
	case "gemini":
		return geminiModels()
	case "pi":
		return discoverPiModels(ctx)
	case "cursor-agent":
		return discoverCursorModels(ctx)
	case "copilot":
		return copilotModels()
	case "hermes":
		return hermesModels()
	case "kimi":
		return kimiModels()
	case "kiro-cli":
		return kiroModels()
	case "openclaw":
		return openclawModels()
	default:
		return nil
	}
}

// ── Model lists ─────────────────────────────────────────────────────────────

func claudeModels() []Model {
	return []Model{
		{ID: "claude-sonnet-4-7", Label: "Claude Sonnet 4.7", Provider: "anthropic", Default: true},
		{ID: "claude-opus-4-7", Label: "Claude Opus 4.7", Provider: "anthropic"},
		{ID: "claude-haiku-4-5", Label: "Claude Haiku 4.5", Provider: "anthropic"},
		{ID: "claude-opus-4-6", Label: "Claude Opus 4.6", Provider: "anthropic"},
		{ID: "claude-sonnet-4-6", Label: "Claude Sonnet 4.6", Provider: "anthropic"},
	}
}

func codexModels() []Model {
	return []Model{
		{ID: "gpt-4o", Label: "GPT-4o", Provider: "openai", Default: true},
		{ID: "gpt-4o-mini", Label: "GPT-4o mini", Provider: "openai"},
		{ID: "o3", Label: "o3", Provider: "openai"},
		{ID: "o3-mini", Label: "o3-mini", Provider: "openai"},
	}
}

func geminiModels() []Model {
	return []Model{
		{ID: "auto", Label: "Auto (Gemini 3)", Provider: "google", Default: true},
		{ID: "gemini-2-5-flash", Label: "Gemini 2.5 Flash", Provider: "google"},
		{ID: "gemini-2-5-pro", Label: "Gemini 2.5 Pro", Provider: "google"},
		{ID: "gemini-2-0-flash", Label: "Gemini 2.0 Flash", Provider: "google"},
	}
}

func opencodeModels() []Model {
	return []Model{
		{ID: "anthropic/claude-sonnet-4-7", Label: "Claude Sonnet 4.7", Provider: "anthropic", Default: true},
		{ID: "anthropic/claude-opus-4-7", Label: "Claude Opus 4.7", Provider: "anthropic"},
		{ID: "openai/gpt-4o", Label: "GPT-4o", Provider: "openai"},
		{ID: "openai/gpt-4o-mini", Label: "GPT-4o mini", Provider: "openai"},
		{ID: "google/gemini-2-5-flash", Label: "Gemini 2.5 Flash", Provider: "google"},
	}
}

func copilotModels() []Model {
	return []Model{
		{ID: "gpt-4o", Label: "GPT-4o", Provider: "openai", Default: true},
		{ID: "gpt-4o-mini", Label: "GPT-4o mini", Provider: "openai"},
		{ID: "claude-opus-4-7", Label: "Claude Opus 4.7", Provider: "anthropic"},
		{ID: "claude-sonnet-4-6", Label: "Claude Sonnet 4.6", Provider: "anthropic"},
	}
}

func openclawModels() []Model {
	return []Model{
		{ID: "auto", Label: "Auto", Provider: "openclaw", Default: true},
	}
}

func hermesModels() []Model {
	return []Model{
		{ID: "auto", Label: "Auto", Provider: "hermes", Default: true},
		{ID: "claude-opus-4-7", Label: "Claude Opus 4.7", Provider: "anthropic"},
	}
}

func kimiModels() []Model {
	return []Model{
		{ID: "auto", Label: "Auto", Provider: "kimi", Default: true},
		{ID: "moonshot-v1-8k", Label: "Moonshot V1 8K", Provider: "kimi"},
	}
}

func kiroModels() []Model {
	return []Model{
		{ID: "auto", Label: "Auto", Provider: "kiro", Default: true},
		{ID: "kiro-1-5", Label: "Kiro 1.5", Provider: "kiro"},
	}
}

// ── Dynamic discovery ──────────────────────────────────────────────────────

func discoverPiModels(ctx context.Context) []Model {
	return listModels(ctx, "pi", func(out []byte) []Model {
		var models []Model
		seen := make(map[string]bool)
		scanner := bufio.NewScanner(strings.NewReader(string(out)))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "provider") {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}
			id := fields[0] + "/" + fields[1]
			if seen[id] {
				continue
			}
			seen[id] = true
			models = append(models, Model{ID: id, Label: id, Provider: fields[0]})
		}
		return models
	})
}

func discoverCursorModels(ctx context.Context) []Model {
	return listModels(ctx, "cursor-agent", func(out []byte) []Model {
		var models []Model
		seen := make(map[string]bool)
		scanner := bufio.NewScanner(strings.NewReader(string(out)))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || seen[line] {
				continue
			}
			seen[line] = true
			models = append(models, Model{ID: line, Label: line, Provider: "cursor"})
		}
		return models
	})
}

// listModels runs cliName --list-models and parses output with the given parser.
// Returns nil if the CLI is not found or fails.
func listModels(ctx context.Context, cliName string, parse func([]byte) []Model) []Model {
	path, err := exec.LookPath(cliName)
	if err != nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, path, "--list-models")
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	models := parse(out)
	if len(models) == 0 {
		return nil
	}
	models[0].Default = true
	return models
}
