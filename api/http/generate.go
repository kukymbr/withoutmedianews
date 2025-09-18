package apihttpspec

//go:generate go tool oapi-codegen -config server.conf.yaml openapi.yaml
//go:generate go tool oapi-codegen -config testclient.conf.yaml openapi.yaml
