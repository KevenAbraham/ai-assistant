package config

import "os"

// LoadSystemPrompt reads the system prompt file specified in cfg.SystemPromptPath
// and returns its content as a string. It is registered as an Fx provider so the
// dependency graph can resolve: *Config → string → *service.ContextBuilder.
func LoadSystemPrompt(cfg *Config) (string, error) {
	data, err := os.ReadFile(cfg.SystemPromptPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
