package sdk

// A set of useful internal functions or implementations exposed in
// the SDK.

import (
  "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/Speer/events"
)

func NewWGProgress(start int) interfaces.GroupProgress {
  return events.NewWGProgress(start)
}
