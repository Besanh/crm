package validator

import (
	"contactcenter-api/common/response"
	"fmt"
	"os"

	"github.com/xeipuuv/gojsonschema"
)

func getSchema(schema string) ([]byte, error) {
	return os.ReadFile("./common/schema/" + schema)
}

func GetSchema(schema string) ([]byte, error) {
	return getSchema(schema)
}

func CheckSchema(schema string, value any) (int, any) {
	schemaChecker, err := getSchema(schema)
	if err != nil {
		return response.NotFoundMsg(schema + " is not existed")
	}
	schemaLoader := gojsonschema.NewStringLoader(fmt.Sprintf("%v", string(schemaChecker)))
	documentLoader := gojsonschema.NewGoLoader(value)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return response.ServiceUnavailableMsg(err.Error())
	}
	if result.Valid() {
		return response.NewOKResponse("")
	} else {
		var errArr []map[string]any
		for _, msg := range result.Errors() {
			if msg.Description() == gojsonschema.Locale.ConditionThen() || msg.Description() == gojsonschema.Locale.NumberAllOf() {
				continue
			}
			errMsg := map[string]any{
				msg.Field(): msg.Description(),
			}
			errArr = append(errArr, errMsg)
		}
		return response.NewErrorResponse(400, errArr)
	}
}
