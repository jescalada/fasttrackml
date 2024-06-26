package repositories

import (
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"

	"github.com/G-Research/fasttrackml/pkg/api/mlflow/dao/models"
)

// makeSqlPlaceholders collects a string of "(?,?,?), (?,?,?)" and so on,
// for use as sql parameters
func makeSqlPlaceholders(numberInEachSet, numberOfSets int) string {
	set := fmt.Sprintf("(%s)", strings.Repeat("?,", numberInEachSet-1)+"?")
	return strings.Repeat(set+",", numberOfSets-1) + set
}

// makeParamConflictPlaceholdersAndValues provides sql placeholders and concatenates
// Key, Value, RunID from each input Param for use in sql values replacement
func makeParamConflictPlaceholdersAndValues(params []models.Param, dialector string) (string, []interface{}) {
	var placeholders string
	// make place holders of 5 fields for each param
	if (sqlite.Dialector{}.Name() == dialector) {
		placeholders = fmt.Sprintf("VALUES %s", makeSqlPlaceholders(5, len(params)))
	} else {
		set := "SELECT ?::text, ?::text, ?::int, ?::float, ?::text"
		placeholders = strings.Repeat(set+"\nUNION ALL\n", len(params)-1) + set
	}
	// values array is params * 5 in length since using 5 fields from each
	valuesArray := make([]interface{}, len(params)*5)
	index := 0
	for _, param := range params {
		valuesArray[index] = param.Key
		valuesArray[index+1] = param.RunID
		if param.ValueInt != nil {
			valuesArray[index+2] = *param.ValueInt
		} else if param.ValueFloat != nil {
			valuesArray[index+3] = *param.ValueFloat
		} else if param.ValueStr != nil {
			valuesArray[index+4] = *param.ValueStr
		}
		index = index + 5
	}
	return placeholders, valuesArray
}

// BuildJsonCondition creates sql and values for where condition to select items having the specified map of json paths
// and values in the given json column. Json path is expressed as "key" or "outerkey.nestedKey".
func BuildJsonCondition(
	dialector string,
	jsonColumnName string,
	jsonPathValueMap map[string]string,
) (sql string, args []any) {
	if len(jsonPathValueMap) == 0 {
		return sql, args
	}
	var conditionTemplate string
	args = make([]any, len(jsonPathValueMap)*2)
	switch dialector {
	case postgres.Dialector{}.Name():
		conditionTemplate = "%s#>>? = ?"
		idx := 0
		for k, v := range jsonPathValueMap {
			path := strings.ReplaceAll(k, ".", ",")
			args[idx] = fmt.Sprintf("{%s}", path)
			args[idx+1] = v
			idx = idx + 2
		}
	default:
		conditionTemplate = "%s->>? = ?"
		idx := 0
		for k, v := range jsonPathValueMap {
			args[idx] = fmt.Sprintf("$.%s", k)
			args[idx+1] = v
			idx = idx + 2
		}
	}
	conditionTemplate = fmt.Sprintf(conditionTemplate, jsonColumnName)
	sql = strings.Repeat(conditionTemplate+" AND ", len(jsonPathValueMap)-1) + conditionTemplate
	return sql, args
}
