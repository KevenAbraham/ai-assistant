package entity

type Intent string

const (
	IntentChat        Intent = "chat"
	IntentOpenApp     Intent = "open_app"
	IntentSetAlarm    Intent = "set_alarm"
	IntentSaveMemory  Intent = "save_memory"
	IntentQueryMemory Intent = "query_memory"
	IntentUnknown     Intent = "unknown"
)

type Action struct {
	Type    string
	Payload map[string]string
}

type Command struct {
	RawText string
	Intent  Intent
	Action  *Action
}
