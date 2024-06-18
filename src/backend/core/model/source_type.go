package model

const (
	SourceTypeIngestionApi   SourceType = "ingestion_api"
	SourceTypeSlack          SourceType = "slack"
	SourceTypeWEB            SourceType = "web"
	SourceTypeGoogleDrive    SourceType = "google_drive"
	SourceTypeGMAIL          SourceType = "gmail"
	SourceTypeRequesttracker SourceType = "requesttracker"
	SourceTypeGithub         SourceType = "github"
	SourceTypeGitlab         SourceType = "gitlab"
	SourceTypeGuru           SourceType = "guru"
	SourceTypeBookstack      SourceType = "bookstack"
	SourceTypeConfluence     SourceType = "confluence"
	SourceTypeSlab           SourceType = "slab"
	SourceTypeJira           SourceType = "jira"
	SourceTypeProductboard   SourceType = "productboard"
	SourceTypeFile           SourceType = "file"
	SourceTypeNotion         SourceType = "notion"
	SourceTypeZulip          SourceType = "zulip"
	SourceTypeLinear         SourceType = "linear"
	SourceTypeHubspot        SourceType = "hubspot"
	SourceTypeDocument360    SourceType = "document360"
	SourceTypeGong           SourceType = "gong"
	SourceTypeGoogleSites    SourceType = "google_sites"
	SourceTypeZendesk        SourceType = "zendesk"
	SourceTypeLoopio         SourceType = "loopio"
	SourceTypeSharepoint     SourceType = "sharepoint"
	SourceTypeOneDrive       SourceType = "one-drive"
	SourceTypeMsTeams        SourceType = "msteams"
	SourceTypeYoutube        SourceType = "youtube"
)

type (
	SourceType            string
	SourceTypeDescription struct {
		ID            SourceType `json:"id"`
		Name          string     `json:"name"`
		IsImplemented bool       `json:"isImplemented"`
	}
)

var (
	sourceTypeFileDescription        = SourceTypeDescription{SourceTypeFile, "File", true}
	sourceTypeWEBDescription         = SourceTypeDescription{SourceTypeWEB, "Web", true}
	sourceTypeSlackDescription       = SourceTypeDescription{SourceTypeSlack, "Slack", false}
	sourceTypeGoogleDriveDescription = SourceTypeDescription{SourceTypeGoogleDrive, "Google Drive", false}
	sourceTypeGmailDescription       = SourceTypeDescription{SourceTypeGMAIL, "Gmail", false}
	sourceTypeSharepointDescription  = SourceTypeDescription{SourceTypeSharepoint, "Sharepoint", false}
	sourceTypeOneDriveDescription    = SourceTypeDescription{SourceTypeOneDrive, "OneDrive", true}
	sourceTypeYouTubeDescription     = SourceTypeDescription{SourceTypeYoutube, "Youtube", true}
)
var AllSourceTypes = map[SourceType]*SourceTypeDescription{
	SourceTypeFile:        &sourceTypeFileDescription,
	SourceTypeWEB:         &sourceTypeWEBDescription,
	SourceTypeSlack:       &sourceTypeSlackDescription,
	SourceTypeGoogleDrive: &sourceTypeGoogleDriveDescription,
	SourceTypeGMAIL:       &sourceTypeGmailDescription,
	SourceTypeSharepoint:  &sourceTypeSharepointDescription,
	SourceTypeOneDrive:    &sourceTypeOneDriveDescription,
	SourceTypeYoutube:     &sourceTypeYouTubeDescription,
}

var SourceTypesList = []*SourceTypeDescription{
	&sourceTypeFileDescription,
	&sourceTypeWEBDescription,
	&sourceTypeSlackDescription,
	&sourceTypeGoogleDriveDescription,
	&sourceTypeGmailDescription,
	&sourceTypeSharepointDescription,
	&sourceTypeOneDriveDescription,
	&sourceTypeYouTubeDescription,
}
