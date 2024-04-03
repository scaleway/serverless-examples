package app

import (
	"subscriber-server/internal/types"
	"sync"
)

// SharedState /* State holding the confirmation state and the notifications */
var SharedState = AppState{
	ReceivedConfirmation:   types.Confirmation{},
	FormattedNotifications: make([]string, 0),
}

type AppState struct {
	ReceivedConfirmation   types.Confirmation
	ConfirmationMutex      sync.Mutex
	FormattedNotifications []string
	NotificationMutex      sync.Mutex
}
