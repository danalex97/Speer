package interfaces

type Progress interface {
    // The function that should be called after progress has been made.
    Progress()

    // The function that should block until we have Progressed enough.
    Advance()
}
