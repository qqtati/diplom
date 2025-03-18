package api_http_model

import "helprepet/internal/common"

type VersionResponse struct {
	Version string `json:"version"`
	common.Response
}
