package interfaces

type Progress interface {
    // The function that should be called after progress has been made.
    Progress(id string)

    // The function that should block until we have Progressed enough.
    Advance()
}
