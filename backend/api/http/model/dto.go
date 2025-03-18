package api_http_model

import "backend/internal/common"

type VersionResponse struct {
	Version string `json:"version"`
	common.Response
}
