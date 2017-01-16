package postgres

import (
	"strings"

	"regexp"

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
	for fieldName := range object {
		if strings.Index(fieldName, ".") == -1 {
			continue
		}
		components := strings.Split(fieldName, ".")

		value := object[fieldName]
		if v := utils.M(value); v != nil {
			if utils.S(v["__op"]) == "Delete" {
				value = nil
			}
		}

		currentObj := object
		for i, next := range components {
			if i == (len(components) - 1) {
				currentObj[next] = value
				break
			}
			obj := currentObj[next]
			if obj == nil {
				obj = types.M{}
				currentObj[next] = obj
			}
			currentObj = utils.M(currentObj[next])
		}

		delete(object, fieldName)
	}
	return object
}

func validateKeys(object interface{}) error {
	if obj := utils.M(object); obj != nil {
		for key, value := range obj {
			err := validateKeys(value)
			if err != nil {
				return err
			}

			if strings.Contains(key, "$") || strings.Contains(key, ".") {
				return errs.E(errs.InvalidNestedKey, "Nested keys should not contain the '$' or '.' characters")
			}
		}
	}
	return nil
}

func joinTablesForSchema(schema types.M) []string {
	list := []string{}
	if schema != nil {
		if fields := utils.M(schema["fields"]); fields != nil {
			className := utils.S(schema["className"])
			for field, v := range fields {
				if tp := utils.M(v); tp != nil {
					if utils.S(tp["type"]) == "Relation" {
						list = append(list, "_Join:"+field+":"+className)
					}
				}
			}
		}
	}
	return list
}

func buildWhereClause(schema, query types.M, index int) (types.M, error) {
	// TODO
	// removeWhiteSpace
	return nil, nil
}

func removeWhiteSpace(s string) string {
	// TODO
	return ""
}

func processRegexPattern(s string) string {
	if strings.HasPrefix(s, "^") {
		return "^" + literalizeRegexPart(s[1:])
	} else if strings.HasSuffix(s, "$") {
		return literalizeRegexPart(s[:len(s)-1]) + "$"
	}
	return literalizeRegexPart(s)
}

func createLiteralRegex(s string) string {
	chars := strings.Split(s, "")
	for i, c := range chars {
		if m, _ := regexp.MatchString(`[0-9a-zA-Z]`, c); m == false {
			if c == `'` {
				chars[i] = `''`
			} else {
				chars[i] = `\` + c
			}
		}
	}
	return strings.Join(chars, "")
}

func literalizeRegexPart(s string) string {
	// go 不支持 (?!) 语法，需要进行等价替换
	// /\\Q((?!\\E).*)\\E$/
	// /\\Q(\\[^E\n\r].*|[^\\\n\r].*|.??)\\E$/
	matcher1 := regexp.MustCompile(`\\Q(\\[^E\n\r].*|[^\\\n\r].*|.??)\\E$`)
	result1 := matcher1.FindStringSubmatch(s)
	if len(result1) > 1 {
		index := strings.Index(s, result1[0])
		prefix := s[:index]
		remaining := result1[1]
		return literalizeRegexPart(prefix) + createLiteralRegex(remaining)
	}

	// /\\Q((?!\\E).*)$/
	// /\\Q(\\[^E\n\r].*|[^\\\n\r].*|.??)$/
	matcher2 := regexp.MustCompile(`\\Q(\\[^E\n\r].*|[^\\\n\r].*|.??)$`)
	result2 := matcher2.FindStringSubmatch(s)
	if len(result2) > 1 {
		index := strings.Index(s, result2[0])
		prefix := s[:index]
		remaining := result2[1]
		return literalizeRegexPart(prefix) + createLiteralRegex(remaining)
	}

	re := regexp.MustCompile(`([^\\])(\\E)`)
	s = re.ReplaceAllString(s, "$1")
	re = regexp.MustCompile(`([^\\])(\\Q)`)
	s = re.ReplaceAllString(s, "$1")
	re = regexp.MustCompile(`^\\E`)
	s = re.ReplaceAllString(s, "")
	re = regexp.MustCompile(`^\\Q`)
	s = re.ReplaceAllString(s, "")
	re = regexp.MustCompile(`([^'])'`)
	s = re.ReplaceAllString(s, "$1''")
	re = regexp.MustCompile(`^'([^'])`)
	s = re.ReplaceAllString(s, "''$1")
	return s
}
