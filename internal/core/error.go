package core

type WorkspaceNotFoundError struct {
}

func (err WorkspaceNotFoundError) Error() string {
	return "cannot find the workspace"
}
