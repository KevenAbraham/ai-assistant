package entity

import "errors"

// Domain-level sentinel errors. Callers compare with errors.Is.
var (
	ErrConversationNotFound = errors.New("conversation not found")
	ErrMemoryNotFound       = errors.New("memory not found")
	ErrEmptyInput           = errors.New("input text must not be empty")
	ErrAIClientFailure      = errors.New("AI client returned an error")
	ErrTranscriptionFailure = errors.New("audio transcription failed")
	ErrTTSFailure           = errors.New("text-to-speech synthesis failed")
)
