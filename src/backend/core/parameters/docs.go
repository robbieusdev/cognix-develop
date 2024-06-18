package parameters

import "github.com/shopspring/decimal"

type DocumentSetParam struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DocumentSetConnectorPairsParam struct {
	DocumentSetID decimal.Decimal   `json:"document_set_id"`
	ConnectorIDs  []decimal.Decimal `json:"connector_ids"`
}

type DocumentUploadResponse struct {
	FileName string      `json:"file_name"`
	Error    string      `json:"error"`
	Document interface{} `json:"document"`
}
