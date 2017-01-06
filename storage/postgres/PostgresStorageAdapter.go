package postgres

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

const postgresSchemaCollectionName = "_SCHEMA"

const postgresRelationDoesNotExistError = "42P01"
const postgresDuplicateRelationError = "42P07"
const postgresDuplicateColumnError = "42701"
const postgresUniqueIndexViolationError = "23505"
const postgresTransactionAbortedError = "25P02"

// PostgresAdapter postgres 数据库适配器
type PostgresAdapter struct {
	collectionPrefix string
	collectionList   []string
}

// NewPostgresAdapter ...
func NewPostgresAdapter(collectionPrefix string) *PostgresAdapter {
	return &PostgresAdapter{
		collectionPrefix: collectionPrefix,
		collectionList:   []string{},
	}
}

// ClassExists ...
func (p *PostgresAdapter) ClassExists(name string) bool {
	return false
}

// SetClassLevelPermissions ...
func (p *PostgresAdapter) SetClassLevelPermissions(className string, CLPs types.M) error {
	return nil
}

// CreateClass ...
func (p *PostgresAdapter) CreateClass(className string, schema types.M) (types.M, error) {
	return nil, nil
}

// AddFieldIfNotExists ...
func (p *PostgresAdapter) AddFieldIfNotExists(className, fieldName string, fieldType types.M) error {
	return nil
}

// DeleteClass ...
func (p *PostgresAdapter) DeleteClass(className string) (types.M, error) {
	return nil, nil
}

// DeleteAllClasses ...
func (p *PostgresAdapter) DeleteAllClasses() error {
	return nil
}

// DeleteFields ...
func (p *PostgresAdapter) DeleteFields(className string, schema types.M, fieldNames []string) error {
	return nil
}

// CreateObject ...
func (p *PostgresAdapter) CreateObject(className string, schema, object types.M) error {
	return nil
}

// GetAllClasses ...
func (p *PostgresAdapter) GetAllClasses() ([]types.M, error) {
	return nil, nil
}

// GetClass ...
func (p *PostgresAdapter) GetClass(className string) (types.M, error) {
	return nil, nil
}

// DeleteObjectsByQuery ...
func (p *PostgresAdapter) DeleteObjectsByQuery(className string, schema, query types.M) error {
	return nil
}

// Find ...
func (p *PostgresAdapter) Find(className string, schema, query, options types.M) ([]types.M, error) {
	return nil, nil
}

// Count ...
func (p *PostgresAdapter) Count(className string, schema, query types.M) (int, error) {
	return 0, nil
}

// UpdateObjectsByQuery ...
func (p *PostgresAdapter) UpdateObjectsByQuery(className string, schema, query, update types.M) error {
	return nil
}

// FindOneAndUpdate ...
func (p *PostgresAdapter) FindOneAndUpdate(className string, schema, query, update types.M) (types.M, error) {
	return nil, nil
}

// UpsertOneObject ...
func (p *PostgresAdapter) UpsertOneObject(className string, schema, query, update types.M) error {
	return nil
}

// EnsureUniqueness ...
func (p *PostgresAdapter) EnsureUniqueness(className string, schema types.M, fieldNames []string) error {
	return nil
}

// PerformInitialization ...
func (p *PostgresAdapter) PerformInitialization(options types.M) error {
	return nil
}

var parseToPosgresComparator = map[string]string{
	"$gt":  ">",
	"$lt":  "<",
	"$gte": ">=",
	"$lte": "<=",
}

func parseTypeToPostgresType(t types.M) (string, error) {
	if t == nil {
		return "", nil
	}
	tp := utils.S(t["type"])
	switch tp {
	case "String":
		return "text", nil
	case "Date":
		return "timestamp with time zone", nil
	case "Object":
		return "jsonb", nil
	case "File":
		return "text", nil
	case "Boolean":
		return "boolean", nil
	case "Pointer":
		return "char(10)", nil
	case "Number":
		return "double precision", nil
	case "GeoPoint":
		return "point", nil
	case "Array":
		if contents := utils.M(t["contents"]); contents != nil {
			if utils.S(contents["type"]) == "String" {
				return "text[]", nil
			}
		}
		return "jsonb", nil
	default:
		return "", errs.E(errs.IncorrectType, "no type for "+tp+" yet")
	}
}

func toPostgresValue(value interface{}) interface{} {
	if v := utils.M(value); v != nil {
		if utils.S(v["__type"]) == "Date" {
			return v["iso"]
		}
		if utils.S(v["__type"]) == "File" {
			return v["name"]
		}
	}
	return value
}

func transformValue(value interface{}) interface{} {
	if v := utils.M(value); v != nil {
		if utils.S(v["__type"]) == "Pointer" {
			return v["objectId"]
		}
	}
	return value
}

var emptyCLPS = types.M{
	"find":     types.M{},
	"get":      types.M{},
	"create":   types.M{},
	"update":   types.M{},
	"delete":   types.M{},
	"addField": types.M{},
}

var defaultCLPS = types.M{
	"find":     types.M{"*": true},
	"get":      types.M{"*": true},
	"create":   types.M{"*": true},
	"update":   types.M{"*": true},
	"delete":   types.M{"*": true},
	"addField": types.M{"*": true},
}

func toParseSchema(schema types.M) types.M {
	if schema == nil {
		return nil
	}

	var fields types.M
	if fields = utils.M(schema["fields"]); fields == nil {
		fields = types.M{}
	}

	if utils.S(schema["className"]) == "_User" {
		if _, ok := fields["_hashed_password"]; ok {
			delete(fields, "_hashed_password")
		}
	}

	if _, ok := fields["_wperm"]; ok {
		delete(fields, "_wperm")
	}
	if _, ok := fields["_rperm"]; ok {
		delete(fields, "_rperm")
	}

	var clps types.M
	clps = utils.CopyMap(defaultCLPS)
	if classLevelPermissions := utils.M(schema["classLevelPermissions"]); classLevelPermissions != nil {
		// clps = utils.CopyMap(emptyCLPS)
		// 不存在的 action 默认为公共权限
		for k, v := range classLevelPermissions {
			clps[k] = v
		}
	}

	return types.M{
		"className":             schema["className"],
		"fields":                fields,
		"classLevelPermissions": clps,
	}
}

func toPostgresSchema(schema types.M) types.M {
	if schema == nil {
		return nil
	}

	var fields types.M
	if fields = utils.M(schema["fields"]); fields == nil {
		fields = types.M{}
	}

	fields["_wperm"] = types.M{
		"type":     "Array",
		"contents": types.M{"type": "String"},
	}
	fields["_rperm"] = types.M{
		"type":     "Array",
		"contents": types.M{"type": "String"},
	}

	if utils.S(schema["className"]) == "_User" {
		fields["_hashed_password"] = types.M{"type": "String"}
		fields["_password_history"] = types.M{"type": "Array"}
	}

	schema["fields"] = fields

	return schema
}

func handleDotFields(object types.M) types.M {
	// TODO
	return nil
}

func validateKeys(object interface{}) error {
	// TODO
	return nil
}

func joinTablesForSchema(schema types.M) []string {
	// TODO
	return nil
}

func buildWhereClause(schema, query types.M, index int) (types.M, error) {
	// TODO
	// removeWhiteSpace
	// processRegexPattern
	return nil, nil
}

func removeWhiteSpace(s string) string {
	// TODO
	return ""
}

func processRegexPattern(s string) string {
	// TODO
	// literalizeRegexPart
	return ""
}

func createLiteralRegex(s string) string {
	// TODO
	return ""
}

func literalizeRegexPart(s string) string {
	// TODO
	// createLiteralRegex
	return ""
}
