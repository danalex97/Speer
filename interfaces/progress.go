package interfaces

// For more details on Progress properties see events.ProgressProperty
type Progress interface {
    // The function that should be called after progress has been made.
    Progress(id string)

    // The function that should block until we have Progressed enough.
    Advance()
}

type GroupProgress interface {
  Progress

  // The function that should be used to Add an entity to a Progress.
  Add()
}
