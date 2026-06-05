package docs

import (
	_ "embed"
	"encoding/json"

	"github.com/swaggo/swag"
)

//go:embed swagger.json
var swaggerJSON []byte

type swaggerSpec struct {
	Version          string `json:"swagger"`
	SwaggerTemplate  string
	InfoInstanceName string
}

func init() {
	var spec map[string]interface{}
	json.Unmarshal(swaggerJSON, &spec)

	swag.Register("swagger", &swag.Spec{
		InfoInstanceName: "swagger",
		SwaggerTemplate:  string(swaggerJSON),
	})

	_ = spec
}
