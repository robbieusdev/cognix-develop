package model

const (
	ProviderCustom    OAuthProvider = "custom"
	ProviderMicrosoft OAuthProvider = "microsoft"
)

var ConnectorAuthProvider = map[SourceType]OAuthProvider{
	SourceTypeOneDrive: ProviderMicrosoft,
	SourceTypeMsTeams:  ProviderMicrosoft,
}

// OAuthProvider represents enum for oauth providers
type OAuthProvider string
