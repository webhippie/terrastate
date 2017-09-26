package model

// LockInfo gets sent by Terraform as locking payload.
type LockInfo struct {
	ID        string
	Operation string
	Info      string
	Who       string
	Version   string
	Created   string
	Path      string
}
