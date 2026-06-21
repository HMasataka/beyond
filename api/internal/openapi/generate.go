package openapi

// OpenAPI 仕様 (リポジトリ直下 openapi/openapi.yaml) から型と chi サーバを生成する。
//go:generate go tool oapi-codegen -config config.yaml ../../../openapi/openapi.yaml
