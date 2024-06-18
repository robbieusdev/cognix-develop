package handler

const (
	ActionDelete  = "delete"
	ActionRestore = "restore"
)

type TokenResponse struct {
	Token string `json:"token"`
}
