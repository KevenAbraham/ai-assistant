package entity

// Intent categorises what the user wants the assistant to do.
type Intent string

const (
	IntentChat        Intent = "chat"        // general conversation
	IntentOpenApp     Intent = "open_app"    // open an application
	IntentSetAlarm    Intent = "set_alarm"   // create an alarm / reminder
	IntentSaveMemory  Intent = "save_memory" // persist information for later
	IntentQueryMemory Intent = "query_memory"
	IntentUnknown     Intent = "unknown"
)

// Action describes a concrete operation to perform locally on the device.
type Action struct {
	Type    string            // e.g. "open_app", "set_alarm"
	Payload map[string]string // arbitrary key-value parameters
}

// Command is the result of parsing raw user text into a structured request.
type Command struct {
	RawText string
	Intent  Intent
	Action  *Action // nil when the intent is pure chat
}
