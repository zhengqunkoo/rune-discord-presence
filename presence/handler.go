package presence

import (
	"context"
	"log"

	"github.com/unstablebuild/rune-go-sdk/api/textapi"
)

// EventHandler implements textapi.EventHandler and forwards editor events
// into the Presence state machine. Keep it minimal: it runs on the
// editor's event path.
type EventHandler struct {
	Pres *Presence
}

// Handle implements textapi.EventHandler. Returning false keeps the
// subscription alive for the lifetime of the extension.
func (h *EventHandler) Handle(ctx context.Context, ev textapi.Event) bool {
	switch ev.Type {
	case textapi.EventTypeOpen, textapi.EventTypeFocus:
		log.Printf("focused on: %s", ev.URI.Name())
		h.Pres.SetActive(ev.URI)
	case textapi.EventTypeUnfocus:
		// The editor lost focus entirely (user in file manager, console,
		// or another app) — show idle until something regains focus.
		// this also triggers before a focus event when switching files
		// but the idle state should quickly be overwritten with the new state
		// before a presence update happens.
		log.Printf("Set idle")
		h.Pres.SetIdle()
	}

	return false // keep receiving events
}

// Compile-time check that we satisfy the event handler interface.
var _ textapi.EventHandler = (*EventHandler)(nil)
