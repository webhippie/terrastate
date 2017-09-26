package model

type LockInfo struct {
	ID        string
	Operation string
	Info      string
	Who       string
	Version   string
	Created   string
	Path      string
}
