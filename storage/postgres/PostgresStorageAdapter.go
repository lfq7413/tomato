package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"regexp"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
	"github.com/lib/pq"
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
	db               *sql.DB
}

// NewPostgresAdapter ...
func NewPostgresAdapter(collectionPrefix string, db *sql.DB) *PostgresAdapter {
	return &PostgresAdapter{
		collectionPrefix: collectionPrefix,
		collectionList:   []string{},
		db:               db,
	}
}

// ensureSchemaCollectionExists 确保 _SCHEMA 表存在，不存在则创建表
func (p *PostgresAdapter) ensureSchemaCollectionExists() error {
	_, err := p.db.Exec(`CREATE TABLE IF NOT EXISTS "_SCHEMA" ( "className" varChar(120), "schema" jsonb, "isParseClass" bool, PRIMARY KEY ("className") )`)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			if e.Code == postgresDuplicateRelationError || e.Code == postgresUniqueIndexViolationError {
				// _SCHEMA 表已经存在，已经由其他请求创建，忽略错误
				return nil
			}
		} else {
			return err
		}
	}
	return nil
}

// ClassExists 检测数据库中是否存在指定类
func (p *PostgresAdapter) ClassExists(name string) bool {
	var result bool
	err := p.db.QueryRow(`SELECT EXISTS (SELECT 1 FROM   information_schema.tables WHERE table_name = $1)`, name).Scan(&result)
	if err != nil {
		return false
	}
	return result
}

// SetClassLevelPermissions 设置类级别权限
func (p *PostgresAdapter) SetClassLevelPermissions(className string, CLPs types.M) error {
	err := p.ensureSchemaCollectionExists()
	if err != nil {
		return err
	}
	if CLPs == nil {
		CLPs = types.M{}
	}
	b, err := json.Marshal(CLPs)
	if err != nil {
		return err
	}

	qs := `UPDATE "_SCHEMA" SET "schema" = json_object_set_key("schema", $1::text, $2::jsonb) WHERE "className"=$3 `
	_, err = p.db.Exec(qs, "classLevelPermissions", string(b), className)
	if err != nil {
		return err
	}

	return nil
}

// CreateClass 创建类
func (p *PostgresAdapter) CreateClass(className string, schema types.M) (types.M, error) {
	b, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}

	err = p.createTable(className, schema)
	if err != nil {
		return nil, err
	}

	_, err = p.db.Exec(`INSERT INTO "_SCHEMA" ("className", "schema", "isParseClass") VALUES ($1, $2, $3)`, className, string(b), true)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			if e.Code == postgresUniqueIndexViolationError {
				return nil, errs.E(errs.DuplicateValue, "Class "+className+" already exists.")
			}
		}
		return nil, err
	}

	return toParseSchema(schema), nil
}

// createTable 仅创建表，不加入 schema 中
func (p *PostgresAdapter) createTable(className string, schema types.M) error {
	if schema == nil {
		schema = types.M{}
	}
	valuesArray := types.S{}
	patternsArray := []string{}
	var fields types.M
	if f := utils.M(schema["fields"]); f != nil {
		fields = f
	}

	if className == "_User" {
		fields["_email_verify_token_expires_at"] = types.M{"type": "Date"}
		fields["_email_verify_token"] = types.M{"type": "String"}
		fields["_account_lockout_expires_at"] = types.M{"type": "Date"}
		fields["_failed_login_count"] = types.M{"type": "Number"}
		fields["_perishable_token"] = types.M{"type": "String"}
		fields["_perishable_token_expires_at"] = types.M{"type": "Date"}
		fields["_password_changed_at"] = types.M{"type": "Date"}
		fields["_password_history"] = types.M{"type": "Array"}
	}

	relations := []string{}

	for fieldName, t := range fields {
		parseType := utils.M(t)
		if parseType == nil {
			parseType = types.M{}
		}

		if utils.S(parseType["type"]) == "Relation" {
			relations = append(relations, fieldName)
			continue
		}

		if fieldName == "_rperm" || fieldName == "_wperm" {
			parseType["contents"] = types.M{"type": "String"}
		}

		valuesArray = append(valuesArray, fieldName)
		postgresType, err := parseTypeToPostgresType(parseType)
		if err != nil {
			return err
		}
		valuesArray = append(valuesArray, postgresType)

		patternsArray = append(patternsArray, `"%s" %s`)
		if fieldName == "objectId" {
			valuesArray = append(valuesArray, fieldName)
			patternsArray = append(patternsArray, `PRIMARY KEY ("%s")`)
		}
	}

	qs := `CREATE TABLE IF NOT EXISTS "%s" (` + strings.Join(patternsArray, ",") + `)`
	values := append(types.S{className}, valuesArray...)
	qs = fmt.Sprintf(qs, values...)

	err := p.ensureSchemaCollectionExists()
	if err != nil {
		return err
	}

	_, err = p.db.Exec(qs)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			if e.Code == postgresDuplicateRelationError {
				// 表已经存在，已经由其他请求创建，忽略错误
			} else {
				return err
			}
		} else {
			return err
		}
	}

	// 创建 relation 表
	for _, fieldName := range relations {
		name := fmt.Sprintf(`_Join:%s:%s`, fieldName, className)
		_, err = p.db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" ("relatedId" varChar(120), "owningId" varChar(120), PRIMARY KEY("relatedId", "owningId") )`, name))
		if err != nil {
			return err
		}
	}

	return nil
}

// AddFieldIfNotExists 添加字段定义
func (p *PostgresAdapter) AddFieldIfNotExists(className, fieldName string, fieldType types.M) error {
	if fieldType == nil {
		fieldType = types.M{}
	}

	if utils.S(fieldType["type"]) != "Relation" {
		tp, err := parseTypeToPostgresType(fieldType)
		if err != nil {
			return err
		}
		qs := fmt.Sprintf(`ALTER TABLE "%s" ADD COLUMN "%s" %s`, className, fieldName, tp)
		_, err = p.db.Exec(qs)
		if err != nil {
			if e, ok := err.(*pq.Error); ok {
				if e.Code == postgresRelationDoesNotExistError {
					// TODO 添加默认字段
					_, ce := p.CreateClass(className, types.M{"fields": types.M{fieldName: fieldType}})
					if ce != nil {
						return ce
					}
				} else if e.Code == postgresDuplicateColumnError {
					// Column 已经存在，由其他请求创建
				} else {
					return err
				}
			} else {
				return err
			}
		}
	} else {
		name := fmt.Sprintf(`_Join:%s:%s`, fieldName, className)
		qs := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" ("relatedId" varChar(120), "owningId" varChar(120), PRIMARY KEY("relatedId", "owningId") )`, name)
		_, err := p.db.Exec(qs)
		if err != nil {
			return err
		}
	}

	qs := `SELECT "schema" FROM "_SCHEMA" WHERE "className" = $1`
	rows, err := p.db.Query(qs, className)
	if err != nil {
		return err
	}
	if rows.Next() {
		var sch types.M
		var v []byte
		err := rows.Scan(&v)
		if err != nil {
			return err
		}
		err = json.Unmarshal(v, &sch)
		if err != nil {
			return err
		}
		if sch == nil {
			sch = types.M{}
		}
		var fields types.M
		if v := utils.M(sch["fields"]); v != nil {
			fields = v
		} else {
			fields = types.M{}
		}
		if _, ok := fields[fieldName]; ok {
			// 当表不存在时，会进行新建表，所以也会走到这里，不再处理错误
			// Attempted to add a field that already exists
			return nil
		}
		fields[fieldName] = fieldType
		sch["fields"] = fields
		b, err := json.Marshal(sch)
		qs := `UPDATE "_SCHEMA" SET "schema"=$1 WHERE "className"=$2`
		_, err = p.db.Exec(qs, b, className)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteClass 删除指定表
func (p *PostgresAdapter) DeleteClass(className string) (types.M, error) {
	qs := fmt.Sprintf(`DROP TABLE IF EXISTS "%s"`, className)
	_, err := p.db.Exec(qs)
	if err != nil {
		return nil, err
	}

	qs = `DELETE FROM "_SCHEMA" WHERE "className"=$1`
	_, err = p.db.Exec(qs, className)
	if err != nil {
		return nil, err
	}

	return types.M{}, nil
}

// DeleteAllClasses 删除所有表，仅用于测试
func (p *PostgresAdapter) DeleteAllClasses() error {
	qs := `SELECT "className","schema" FROM "_SCHEMA"`
	rows, err := p.db.Query(qs)
	if err != nil {
		if e, ok := err.(*pq.Error); ok && e.Code == postgresRelationDoesNotExistError {
			// _SCHEMA 不存在，则不删除
			return nil
		}
		return err
	}

	classNames := []string{}
	schemas := []types.M{}

	for rows.Next() {
		var clsName string
		var sch types.M
		var v []byte
		err := rows.Scan(&clsName, &v)
		if err != nil {
			return err
		}
		err = json.Unmarshal(v, &sch)
		if err != nil {
			return err
		}
		classNames = append(classNames, clsName)
		schemas = append(schemas, sch)
	}

	joins := []string{}
	for _, sch := range schemas {
		joins = append(joins, joinTablesForSchema(sch)...)
	}

	classes := []string{"_SCHEMA", "_PushStatus", "_JobStatus", "_Hooks", "_GlobalConfig"}
	classes = append(classes, classNames...)
	classes = append(classes, joins...)

	for _, name := range classes {
		qs = fmt.Sprintf(`DROP TABLE IF EXISTS "%s"`, name)
		p.db.Exec(qs)
	}

	return nil
}

// DeleteFields 删除字段
func (p *PostgresAdapter) DeleteFields(className string, schema types.M, fieldNames []string) error {
	if schema == nil {
		schema = types.M{}
	}

	fields := utils.M(schema["fields"])
	if fields == nil {
		fields = types.M{}
	}
	fldNames := types.S{}
	for _, fieldName := range fieldNames {
		field := utils.M(fields[fieldName])
		if field != nil && utils.S(field["type"]) == "Relation" {
			// 不处理 Relation 类型字段
		} else {
			fldNames = append(fldNames, fieldName)
		}
		delete(fields, fieldName)
	}
	schema["fields"] = fields

	values := append(types.S{className}, fldNames...)
	columnArray := []string{}
	for _ = range fldNames {
		columnArray = append(columnArray, `"%s"`)
	}
	columns := strings.Join(columnArray, ",")

	b, err := json.Marshal(schema)
	if err != nil {
		return err
	}
	qs := `UPDATE "_SCHEMA" SET "schema"=$1 WHERE "className"=$2`
	_, err = p.db.Exec(qs, b, className)
	if err != nil {
		return err
	}

	if len(values) > 1 {
		qs = fmt.Sprintf(`ALTER TABLE "%%s" DROP COLUMN %s`, columns)
		qs = fmt.Sprintf(qs, values...)
		_, err = p.db.Exec(qs)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateObject 创建对象
func (p *PostgresAdapter) CreateObject(className string, schema, object types.M) error {
	columnsArray := []string{}
	valuesArray := types.S{}
	geoPoints := types.M{}
	if schema == nil {
		schema = types.M{}
	}
	if object == nil {
		object = types.M{}
	}
	schema = toPostgresSchema(schema)
	object = handleDotFields(object)

	err := validateKeys(object)
	if err != nil {
		return err
	}

	for fieldName := range object {
		re := regexp.MustCompile(`^_auth_data_([a-zA-Z0-9_]+)$`)
		authDataMatch := re.FindStringSubmatch(fieldName)
		if authDataMatch != nil && len(authDataMatch) == 2 {
			provider := authDataMatch[1]
			authData := utils.M(object["authData"])
			if authData == nil {
				authData = types.M{}
			}
			authData[provider] = object[fieldName]
			delete(object, fieldName)
			fieldName = "authData"
			object["authData"] = authData
		}
		columnsArray = append(columnsArray, fieldName)

		fields := utils.M(schema["fields"])
		if fields == nil {
			fields = types.M{}
		}
		if fields[fieldName] == nil && className == "_User" {
			if fieldName == "_email_verify_token" ||
				fieldName == "_failed_login_count" ||
				fieldName == "_perishable_token" {
				valuesArray = append(valuesArray, object[fieldName])
			}

			if fieldName == "_password_history" {
				b, err := json.Marshal(object[fieldName])
				if err != nil {
					return err
				}
				valuesArray = append(valuesArray, b)
			}

			if fieldName == "_email_verify_token_expires_at" ||
				fieldName == "_account_lockout_expires_at" ||
				fieldName == "_perishable_token_expires_at" ||
				fieldName == "_password_changed_at" {
				if v := utils.M(object[fieldName]); v != nil && utils.S(v["iso"]) != "" {
					valuesArray = append(valuesArray, v["iso"])
				} else {
					valuesArray = append(valuesArray, nil)
				}
			}

			continue
		}

		tp := utils.M(fields[fieldName])
		if tp == nil {
			tp = types.M{}
		}
		switch utils.S(tp["type"]) {
		case "Date":
			if v := utils.M(object[fieldName]); v != nil && utils.S(v["iso"]) != "" {
				valuesArray = append(valuesArray, v["iso"])
			} else {
				valuesArray = append(valuesArray, nil)
			}
		case "Pointer":
			if v := utils.M(object[fieldName]); v != nil && utils.S(v["objectId"]) != "" {
				valuesArray = append(valuesArray, v["objectId"])
			} else {
				valuesArray = append(valuesArray, "")
			}
		case "Array":
			b, err := json.Marshal(object[fieldName])
			if err != nil {
				return err
			}
			if fieldName == "_rperm" || fieldName == "_wperm" {
				// '[' => '{'
				if b[0] == 91 {
					b[0] = 123
				}
				// ']' => '}'
				if len(b) > 0 && b[len(b)-1] == 93 {
					b[len(b)-1] = 125
				}
			}
			valuesArray = append(valuesArray, b)
		case "Object":
			b, err := json.Marshal(object[fieldName])
			if err != nil {
				return err
			}
			valuesArray = append(valuesArray, b)
		case "String", "Number", "Boolean":
			valuesArray = append(valuesArray, object[fieldName])
		case "File":
			if v := utils.M(object[fieldName]); v != nil && utils.S(v["name"]) != "" {
				valuesArray = append(valuesArray, v["name"])
			} else {
				valuesArray = append(valuesArray, "")
			}
		case "GeoPoint":
			geoPoints[fieldName] = object[fieldName]
			columnsArray = columnsArray[:len(columnsArray)-1]
		default:
			return errs.E(errs.OtherCause, "Type "+utils.S(tp["type"])+" not supported yet")
		}

	}

	for key := range geoPoints {
		columnsArray = append(columnsArray, key)
	}

	initialValues := []string{}
	for index := range valuesArray {
		termination := ""
		fieldName := columnsArray[index]
		if fieldName == "_rperm" || fieldName == "_wperm" {
			termination = "::text[]"
		} else {
			fields := utils.M(schema["fields"])
			if fields == nil {
				fields = types.M{}
			}
			tp := utils.M(fields[fieldName])
			if tp == nil {
				tp = types.M{}
			}
			if utils.S(tp["type"]) == "Array" {
				termination = "::jsonb"
			}
		}
		initialValues = append(initialValues, fmt.Sprintf(`$%d%s`, index+1, termination))
	}

	geoPointsInjects := []string{}
	for _, v := range geoPoints {
		value := utils.M(v)
		if value == nil {
			value = types.M{}
		}
		valuesArray = append(valuesArray, value["longitude"], value["latitude"])
		l := len(valuesArray)
		geoPointsInjects = append(geoPointsInjects, fmt.Sprintf(`POINT($%d, $%d)`, l-1, l))
	}

	columnsPatternArray := []string{}
	for _, key := range columnsArray {
		columnsPatternArray = append(columnsPatternArray, fmt.Sprintf(`"%s"`, key))
	}
	columnsPattern := strings.Join(columnsPatternArray, ",")

	initialValues = append(initialValues, geoPointsInjects...)
	valuesPattern := strings.Join(initialValues, ",")

	qs := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s)`, className, columnsPattern, valuesPattern)
	_, err = p.db.Exec(qs, valuesArray...)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			if e.Code == postgresUniqueIndexViolationError {
				return errs.E(errs.DuplicateValue, "A duplicate value for a field with unique values was provided")
			}
		}
		return err
	}

	return nil
}

// GetAllClasses ...
func (p *PostgresAdapter) GetAllClasses() ([]types.M, error) {
	err := p.ensureSchemaCollectionExists()
	if err != nil {
		return nil, err
	}
	qs := `SELECT "className","schema" FROM "_SCHEMA"`
	rows, err := p.db.Query(qs)
	if err != nil {
		return nil, err
	}

	schemas := []types.M{}

	for rows.Next() {
		var clsName string
		var sch types.M
		var v []byte
		err := rows.Scan(&clsName, &v)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(v, &sch)
		if err != nil {
			return nil, err
		}
		sch["className"] = clsName
		schemas = append(schemas, toParseSchema(sch))
	}

	return schemas, nil
}

// GetClass ...
func (p *PostgresAdapter) GetClass(className string) (types.M, error) {
	qs := `SELECT "schema" FROM "_SCHEMA" WHERE "className"=$1`
	rows, err := p.db.Query(qs, className)
	if err != nil {
		return nil, err
	}

	schema := types.M{}
	if rows.Next() {
		var v []byte
		err = rows.Scan(&v)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(v, &schema)
		if err != nil {
			return nil, err
		}
	} else {
		return schema, nil
	}

	return toParseSchema(schema), nil
}

// DeleteObjectsByQuery 删除符合条件的所有对象
func (p *PostgresAdapter) DeleteObjectsByQuery(className string, schema, query types.M) error {
	where, err := buildWhereClause(schema, query, 1)
	if err != nil {
		return err
	}

	if len(query) == 0 {
		where.pattern = "TRUE"
	}

	qs := fmt.Sprintf(`WITH deleted AS (DELETE FROM "%s" WHERE %s RETURNING *) SELECT count(*) FROM deleted`, className, where.pattern)
	row := p.db.QueryRow(qs, where.values...)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return errs.E(errs.ObjectNotFound, "Object not found.")
	}

	return nil
}

// Find ...
func (p *PostgresAdapter) Find(className string, schema, query, options types.M) ([]types.M, error) {
	if schema == nil {
		schema = types.M{}
	}
	if options == nil {
		options = types.M{}
	}

	var hasLimit bool
	var hasSkip bool
	if _, ok := options["limit"]; ok {
		hasLimit = true
	}
	if _, ok := options["skip"]; ok {
		hasSkip = true
	}

	values := types.S{}
	where, err := buildWhereClause(schema, query, 1)
	if err != nil {
		return nil, err
	}
	values = append(values, where.values...)

	var wherePattern string
	var limitPattern string
	var skipPattern string
	if where.pattern != "" {
		wherePattern = `WHERE ` + where.pattern
	}
	if hasLimit {
		limitPattern = fmt.Sprintf(`LIMIT $%d`, len(values)+1)
		values = append(values, options["limit"])
	}
	if hasSkip {
		skipPattern = fmt.Sprintf(`OFFSET $%d`, len(values)+1)
		values = append(values, options["skip"])
	}

	var sortPattern string
	if _, ok := options["sort"]; ok {
		if keys, ok := options["sort"].([]string); ok {
			postgresSort := []string{}
			for _, key := range keys {
				var postgresKey string
				if strings.HasPrefix(key, "-") {
					key = key[1:]
					postgresKey = fmt.Sprintf(`"%s" DESC`, key)
				} else {
					postgresKey = fmt.Sprintf(`"%s" ASC`, key)
				}
				postgresSort = append(postgresSort, postgresKey)
			}
			sorting := strings.Join(postgresSort, ",")
			if len(postgresSort) > 0 {
				sortPattern = fmt.Sprintf(`ORDER BY %s`, sorting)
			}
		}
	}
	if len(where.sorts) > 0 {
		sortPattern = fmt.Sprintf(`ORDER BY %s`, strings.Join(where.sorts, ","))
	}

	columns := "*"
	if _, ok := options["keys"]; ok {
		if keys, ok := options["keys"].([]string); ok {
			postgresKeys := []string{}
			for _, key := range keys {
				if key != "" {
					postgresKeys = append(postgresKeys, fmt.Sprintf(`"%s"`, key))
				}
			}
			if len(postgresKeys) > 0 {
				columns = strings.Join(postgresKeys, ",")
			}
		}
	}

	qs := fmt.Sprintf(`SELECT %s FROM "%s" %s %s %s %s`, columns, className, wherePattern, sortPattern, limitPattern, skipPattern)
	rows, err := p.db.Query(qs, values...)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			// 表不存在返回空
			if e.Code == postgresRelationDoesNotExistError {
				return []types.M{}, nil
			}
		}
		return nil, err
	}

	fields := utils.M(schema["fields"])
	if fields == nil {
		fields = types.M{}
	}

	results := []types.M{}
	var resultColumns []string
	for rows.Next() {
		if resultColumns == nil {
			resultColumns, err = rows.Columns()
			if err != nil {
				return nil, err
			}
		}
		resultValues := []*interface{}{}
		values := types.S{}
		for i := 0; i < len(resultColumns); i++ {
			var v interface{}
			resultValues = append(resultValues, &v)
			values = append(values, &v)
		}
		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}
		object := types.M{}
		for i, field := range resultColumns {
			object[field] = *resultValues[i]
		}

		object, err = postgresObjectToParseObject(object, fields)
		if err != nil {
			return nil, err
		}

		results = append(results, object)
	}

	return results, nil
}

// Count ...
func (p *PostgresAdapter) Count(className string, schema, query types.M) (int, error) {
	where, err := buildWhereClause(schema, query, 1)
	if err != nil {
		return 0, err
	}

	wherePattern := ""
	if len(where.pattern) > 0 {
		wherePattern = `WHERE ` + where.pattern
	}

	qs := fmt.Sprintf(`SELECT count(*) FROM "%s" %s`, className, wherePattern)
	rows, err := p.db.Query(qs, where.values...)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			if e.Code == postgresRelationDoesNotExistError {
				return 0, nil
			}
		}
		return 0, err
	}
	var count int
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, nil
		}
	}

	return count, nil
}

// UpdateObjectsByQuery ...
func (p *PostgresAdapter) UpdateObjectsByQuery(className string, schema, query, update types.M) error {
	_, err := p.FindOneAndUpdate(className, schema, query, update)
	return err
}

// FindOneAndUpdate ...
func (p *PostgresAdapter) FindOneAndUpdate(className string, schema, query, update types.M) (types.M, error) {
	updatePatterns := []string{}
	values := types.S{}
	index := 1

	if schema == nil {
		schema = types.M{}
	}
	schema = toPostgresSchema(schema)

	fields := utils.M(schema["fields"])
	if fields == nil {
		fields = types.M{}
	}

	originalUpdate := utils.CopyMapM(update)
	update = handleDotFields(update)

	for fieldName, v := range update {
		re := regexp.MustCompile(`^_auth_data_([a-zA-Z0-9_]+)$`)
		authDataMatch := re.FindStringSubmatch(fieldName)
		if authDataMatch != nil && len(authDataMatch) == 2 {
			provider := authDataMatch[1]
			delete(update, fieldName)
			authData := utils.M(update["authData"])
			if authData == nil {
				authData = types.M{}
			}
			authData[provider] = v
			update["authData"] = authData
		}
	}

	for fieldName, fieldValue := range update {
		if fieldValue == nil {
			updatePatterns = append(updatePatterns, fmt.Sprintf(`"%s" = NULL`, fieldName))
			continue
		}

		if fieldName == "authData" {
			generate := func(jsonb, key, value string) string {
				return fmt.Sprintf(`json_object_set_key(COALESCE(%s, '{}'::jsonb), %s, %s)::jsonb`, jsonb, key, value)
			}
			lastKey := fmt.Sprintf(`"%s"`, fieldName)
			authData := utils.M(fieldValue)
			if authData == nil {
				continue
			}
			for key, value := range authData {
				lastKey = generate(lastKey, fmt.Sprintf(`$%d::text`, index), fmt.Sprintf(`$%d::jsonb`, index+1))
				index = index + 2
				if value != nil {
					if v := utils.M(value); v != nil && utils.S(v["__op"]) == "Delete" {
						value = nil
					} else {
						b, err := json.Marshal(v)
						if err != nil {
							return nil, err
						}
						value = string(b)
					}
				}
				values = append(values, key, value)
			}
			updatePatterns = append(updatePatterns, fmt.Sprintf(`"%s" = %s`, fieldName, lastKey))
			continue
		}

		if fieldName == "updatedAt" {
			updatePatterns = append(updatePatterns, fmt.Sprintf(`"%s" = $%d`, fieldName, index))
			values = append(values, fieldValue)
			index = index + 1
			continue
		}

		switch fieldValue.(type) {
		case string, bool, float64, int:
			updatePatterns = append(updatePatterns, fmt.Sprintf(`"%s" = $%d`, fieldName, index))
			values = append(values, fieldValue)
			index = index + 1
			continue
		case time.Time:
			updatePatterns = append(updatePatterns, fmt.Sprintf(`"%s" = $%d`, fieldName, index))
			values = append(values, utils.TimetoString(fieldValue.(time.Time)))
			index = index + 1
			continue
		}

		if object := utils.M(fieldValue); object != nil {
			switch utils.S(object["__op"]) {
			case "Increment":
				updatePatterns = append(updatePatterns, fmt.Sprintf(`"%s" = COALESCE("%s", 0) + $%d`, fieldName, fieldName, index))
				values = append(values, object["amount"])
				index = index + 1
				continue
			case "Add":
				updatePatterns = append(updatePatterns, fmt.Sprintf(`"%s" = array_add(COALESCE("%s", '[]'::jsonb), $%d::jsonb)`, fieldName, fieldName, index))
				b, err := json.Marshal(object["objects"])
				if err != nil {
					return nil, err
				}
				values = append(values, string(b))
				index = index + 1
				continue
			case "Delete":
				updatePatterns = append(updatePatterns, fmt.Sprintf(`"%s" = $%d`, fieldName, index))
				values = append(values, nil)
				index = index + 1
				continue
			case "Remove":
				updatePatterns = append(updatePatterns, fmt.Sprintf(`"%s" = array_remove(COALESCE("%s", '[]'::jsonb), $%d::jsonb)`, fieldName, fieldName, index))
				b, err := json.Marshal(object["objects"])
				if err != nil {
					return nil, err
				}
				values = append(values, string(b))
				index = index + 1
				continue
			case "AddUnique":
				updatePatterns = append(updatePatterns, fmt.Sprintf(`"%s" = array_add_unique(COALESCE("%s", '[]'::jsonb), $%d::jsonb)`, fieldName, fieldName, index))
				b, err := json.Marshal(object["objects"])
				if err != nil {
					return nil, err
				}
				values = append(values, string(b))
				index = index + 1
				continue
			}

			switch utils.S(object["__type"]) {
			case "Pointer":
				updatePatterns = append(updatePatterns, fmt.Sprintf(`"%s" = $%d`, fieldName, index))
				values = append(values, object["objectId"])
				index = index + 1
				continue
			case "Date", "File":
				updatePatterns = append(updatePatterns, fmt.Sprintf(`"%s" = $%d`, fieldName, index))
				values = append(values, toPostgresValue(object))
				index = index + 1
				continue
			case "GeoPoint":
				updatePatterns = append(updatePatterns, fmt.Sprintf(`"%s" = POINT($%d, $%d)`, fieldName, index, index+1))
				values = append(values, object["longitude"], object["latitude"])
				index = index + 2
				continue
			case "Relation":
				continue
			}

			if tp := utils.M(fields[fieldName]); tp != nil && utils.S(tp["type"]) == "Object" {
				keysToDelete := []string{}
				for k, v := range originalUpdate {
					if o := utils.M(v); o != nil && utils.S(o["__op"]) == "Delete" {
						if keys := strings.Split(k, "."); len(keys) == 2 && keys[0] == fieldName {
							keysToDelete = append(keysToDelete, keys[1])
						}
					}
				}

				deletePatterns := ""
				for _, k := range keysToDelete {
					deletePatterns = deletePatterns + fmt.Sprintf(` - '%s'`, k)
				}

				updatePatterns = append(updatePatterns, fmt.Sprintf(`"%s" = ( COALESCE("%s", '{}'::jsonb) %s || $%d::jsonb )`, fieldName, fieldName, deletePatterns, index))
				b, err := json.Marshal(object)
				if err != nil {
					return nil, err
				}
				values = append(values, string(b))
				index = index + 1
				continue
			}
		}

		if array := utils.A(fieldValue); array != nil {
			if tp := utils.M(fields[fieldName]); tp != nil && utils.S(tp["type"]) == "Array" {
				expectedType, err := parseTypeToPostgresType(tp)
				if err != nil {
					return nil, err
				}

				b, err := json.Marshal(fieldValue)
				if err != nil {
					return nil, err
				}

				if expectedType == "text[]" {
					// '[' => '{'
					if b[0] == 91 {
						b[0] = 123
					}
					// ']' => '}'
					if len(b) > 0 && b[len(b)-1] == 93 {
						b[len(b)-1] = 125
					}
					fieldValue = string(b)
					updatePatterns = append(updatePatterns, fmt.Sprintf(`"%s" = $%d::text[]`, fieldName, index))
				} else {
					updatePatterns = append(updatePatterns, fmt.Sprintf(`"%s" = $%d::jsonb`, fieldName, index))
				}

				values = append(values, fieldValue)
				index = index + 1
				continue
			}
		}

		b, _ := json.Marshal(fieldValue)
		return nil, errs.E(errs.OperationForbidden, "Postgres doesn't support update "+string(b)+" yet")
	}

	where, err := buildWhereClause(schema, query, index)
	if err != nil {
		return nil, err
	}
	values = append(values, where.values...)

	qs := fmt.Sprintf(`UPDATE "%s" SET %s WHERE %s RETURNING *`, className, strings.Join(updatePatterns, ","), where.pattern)
	rows, err := p.db.Query(qs, values...)
	if err != nil {
		return nil, err
	}

	object := types.M{}
	if rows.Next() {
		resultColumns, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		resultValues := []*interface{}{}
		values := types.S{}
		for i := 0; i < len(resultColumns); i++ {
			var v interface{}
			resultValues = append(resultValues, &v)
			values = append(values, &v)
		}
		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}
		for i, field := range resultColumns {
			object[field] = *resultValues[i]
		}

		object, err = postgresObjectToParseObject(object, fields)
		if err != nil {
			return nil, err
		}
	}

	return object, nil
}

// UpsertOneObject 仅用于 config 和 hooks
func (p *PostgresAdapter) UpsertOneObject(className string, schema, query, update types.M) error {
	object, err := p.FindOneAndUpdate(className, schema, query, update)
	if err != nil {
		return err
	}
	if len(object) == 0 {
		createValue := types.M{}
		for k, v := range query {
			createValue[k] = v
		}
		for k, v := range update {
			createValue[k] = v
		}

		err = p.CreateObject(className, schema, createValue)
		if err != nil {
			return err
		}
	}
	return nil
}

// EnsureUniqueness ...
func (p *PostgresAdapter) EnsureUniqueness(className string, schema types.M, fieldNames []string) error {
	// TODO
	return nil
}

// PerformInitialization ...
func (p *PostgresAdapter) PerformInitialization(options types.M) error {
	if options == nil {
		options = types.M{}
	}

	if volatileClassesSchemas, ok := options["VolatileClassesSchemas"].([]types.M); ok {
		for _, schema := range volatileClassesSchemas {
			err := p.createTable(utils.S(schema["className"]), schema)
			if err != nil {
				if e, ok := err.(*pq.Error); ok {
					if e.Code != postgresDuplicateRelationError {
						return err
					}
				} else if e, ok := err.(*errs.TomatoError); ok {
					if e.Code != errs.InvalidClassName {
						return err
					}
				} else {
					return err
				}
			}
		}
	}

	_, err := p.db.Exec(jsonObjectSetKey)
	if err != nil {
		return err
	}

	_, err = p.db.Exec(arrayAdd)
	if err != nil {
		return err
	}

	_, err = p.db.Exec(arrayAddUnique)
	if err != nil {
		return err
	}

	_, err = p.db.Exec(arrayRemove)
	if err != nil {
		return err
	}

	_, err = p.db.Exec(arrayContainsAll)
	if err != nil {
		return err
	}

	_, err = p.db.Exec(arrayContains)
	if err != nil {
		return err
	}

	return nil
}

func postgresObjectToParseObject(object, fields types.M) (types.M, error) {
	if len(object) == 0 {
		return object, nil
	}
	for fieldName, v := range fields {
		tp := utils.M(v)
		if tp == nil {
			continue
		}
		objectType := utils.S(tp["type"])

		if objectType == "Pointer" && object[fieldName] != nil {
			if v, ok := object[fieldName].([]byte); ok {
				object[fieldName] = types.M{
					"objectId":  string(v),
					"__type":    "Pointer",
					"className": tp["targetClass"],
				}
			} else {
				object[fieldName] = nil
			}
		} else if objectType == "Relation" {
			object[fieldName] = types.M{
				"__type":    "Relation",
				"className": tp["targetClass"],
			}
		} else if objectType == "GeoPoint" && object[fieldName] != nil {
			// object[fieldName] = (10,20) (longitude, latitude)
			resString := ""
			if v, ok := object[fieldName].([]byte); ok {
				resString = string(v)
			}
			if len(resString) < 5 {
				object[fieldName] = nil
				continue
			}
			pointString := strings.Split(resString[1:len(resString)-1], ",")
			if len(pointString) != 2 {
				object[fieldName] = nil
				continue
			}
			longitude, err := strconv.ParseFloat(pointString[0], 64)
			if err != nil {
				return nil, err
			}
			latitude, err := strconv.ParseFloat(pointString[1], 64)
			if err != nil {
				return nil, err
			}
			object[fieldName] = types.M{
				"__type":    "GeoPoint",
				"longitude": longitude,
				"latitude":  latitude,
			}
		} else if objectType == "File" && object[fieldName] != nil {
			if v, ok := object[fieldName].([]byte); ok {
				object[fieldName] = types.M{
					"__type": "File",
					"name":   string(v),
				}
			} else {
				object[fieldName] = nil
			}
		} else if objectType == "String" && object[fieldName] != nil {
			if v, ok := object[fieldName].([]byte); ok {
				object[fieldName] = string(v)
			} else {
				object[fieldName] = nil
			}
		} else if objectType == "Object" && object[fieldName] != nil {
			if v, ok := object[fieldName].([]byte); ok {
				var r types.M
				err := json.Unmarshal(v, &r)
				if err != nil {
					return nil, err
				}
				object[fieldName] = r
			} else {
				object[fieldName] = nil
			}
		} else if objectType == "Array" && object[fieldName] != nil {
			if fieldName == "_rperm" || fieldName == "_wperm" {
				continue
			}
			if v, ok := object[fieldName].([]byte); ok {
				var r types.S
				err := json.Unmarshal(v, &r)
				if err != nil {
					return nil, err
				}
				object[fieldName] = r
			} else {
				object[fieldName] = nil
			}
		}
	}

	if object["_rperm"] != nil {
		// object["_rperm"] = {hello,world}
		// 在添加 _rperm 时已保证值里不含 ','
		resString := ""
		if v, ok := object["_rperm"].([]byte); ok {
			resString = string(v)
		}
		if len(resString) < 2 {
			object["_rperm"] = nil
		} else {
			keys := strings.Split(resString[1:len(resString)-1], ",")
			rperm := make(types.S, len(keys))
			for i, k := range keys {
				rperm[i] = k
			}
			object["_rperm"] = rperm
		}
	}

	if object["_wperm"] != nil {
		// object["_wperm"] = {hello,world}
		// 在添加 _wperm 时已保证值里不含 ','
		resString := ""
		if v, ok := object["_wperm"].([]byte); ok {
			resString = string(v)
		}
		if len(resString) < 2 {
			object["_wperm"] = nil
		} else {
			keys := strings.Split(resString[1:len(resString)-1], ",")
			wperm := make(types.S, len(keys))
			for i, k := range keys {
				wperm[i] = k
			}
			object["_wperm"] = wperm
		}
	}

	if object["createdAt"] != nil {
		if v, ok := object["createdAt"].(time.Time); ok {
			object["createdAt"] = utils.TimetoString(v)
		} else {
			object["createdAt"] = nil
		}
	}

	if object["updatedAt"] != nil {
		if v, ok := object["updatedAt"].(time.Time); ok {
			object["updatedAt"] = utils.TimetoString(v)
		} else {
			object["updatedAt"] = nil
		}
	}

	if object["expiresAt"] != nil {
		object["expiresAt"] = valueToDate(object["expiresAt"])
	}

	if object["_email_verify_token_expires_at"] != nil {
		object["_email_verify_token_expires_at"] = valueToDate(object["_email_verify_token_expires_at"])
	}

	if object["_account_lockout_expires_at"] != nil {
		object["_account_lockout_expires_at"] = valueToDate(object["_account_lockout_expires_at"])
	}

	if object["_perishable_token_expires_at"] != nil {
		object["_perishable_token_expires_at"] = valueToDate(object["_perishable_token_expires_at"])
	}

	if object["_password_changed_at"] != nil {
		object["_password_changed_at"] = valueToDate(object["_password_changed_at"])
	}

	for fieldName := range object {
		if object[fieldName] == nil {
			delete(object, fieldName)
		}
		if v, ok := object[fieldName].(time.Time); ok {
			object[fieldName] = types.M{
				"__type": "Date",
				"iso":    utils.TimetoString(v),
			}
		}
	}
	return object, nil
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
		return "char(24)", nil
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

type whereClause struct {
	pattern string
	values  types.S
	sorts   []string
}

func buildWhereClause(schema, query types.M, index int) (*whereClause, error) {
	patterns := []string{}
	values := types.S{}
	sorts := []string{}

	schema = toPostgresSchema(schema)
	if schema == nil {
		schema = types.M{}
	}
	fields := utils.M(schema["fields"])
	if fields == nil {
		fields = types.M{}
	}
	for fieldName, fieldValue := range query {
		isArrayField := false
		if fields != nil {
			if tp := utils.M(fields[fieldName]); tp != nil {
				if utils.S(tp["type"]) == "Array" {
					isArrayField = true
				}
			}
		}
		initialPatternsLength := len(patterns)

		if fields[fieldName] == nil {
			if v := utils.M(fieldValue); v != nil {
				if b, ok := v["$exists"].(bool); ok && b == false {
					continue
				}
			}
		}

		if strings.Contains(fieldName, ".") {
			components := strings.Split(fieldName, ".")
			for index, cmpt := range components {
				if index == 0 {
					components[index] = `"` + cmpt + `"`
				} else {
					components[index] = `'` + cmpt + `'`
				}
			}
			name := strings.Join(components, "->")
			b, err := json.Marshal(fieldValue)
			if err != nil {
				return nil, err
			}
			patterns = append(patterns, fmt.Sprintf(`%s = '%v'`, name, string(b)))
		} else if _, ok := fieldValue.(string); ok {
			patterns = append(patterns, fmt.Sprintf(`"%s" = $%d`, fieldName, index))
			values = append(values, fieldValue)
			index = index + 1
		} else if _, ok := fieldValue.(bool); ok {
			patterns = append(patterns, fmt.Sprintf(`"%s" = $%d`, fieldName, index))
			values = append(values, fieldValue)
			index = index + 1
		} else if _, ok := fieldValue.(float64); ok {
			patterns = append(patterns, fmt.Sprintf(`"%s" = $%d`, fieldName, index))
			values = append(values, fieldValue)
			index = index + 1
		} else if _, ok := fieldValue.(int); ok {
			patterns = append(patterns, fmt.Sprintf(`"%s" = $%d`, fieldName, index))
			values = append(values, fieldValue)
			index = index + 1
		} else if fieldName == "$or" || fieldName == "$and" {
			clauses := []string{}
			clauseValues := types.S{}
			if array := utils.A(fieldValue); array != nil {
				for _, v := range array {
					if subQuery := utils.M(v); subQuery != nil {
						clause, err := buildWhereClause(schema, subQuery, index)
						if err != nil {
							return nil, err
						}
						if len(clause.pattern) > 0 {
							clauses = append(clauses, clause.pattern)
							clauseValues = append(clauseValues, clause.values...)
							index = index + len(clause.values)
						}
					}
				}
			}
			var orOrAnd string
			if fieldName == "$or" {
				orOrAnd = " OR "
			} else {
				orOrAnd = " AND "
			}
			patterns = append(patterns, fmt.Sprintf(`(%s)`, strings.Join(clauses, orOrAnd)))
			values = append(values, clauseValues...)
		}

		if value := utils.M(fieldValue); value != nil {

			if v, ok := value["$ne"]; ok {
				if isArrayField {
					j, _ := json.Marshal(types.S{v})
					value["$ne"] = string(j)
					patterns = append(patterns, fmt.Sprintf(`NOT array_contains("%s", $%d)`, fieldName, index))
					values = append(values, value["$ne"])
					index = index + 1
				} else {
					if v == nil {
						patterns = append(patterns, fmt.Sprintf(`"%s" IS NOT NULL`, fieldName))
					} else {
						patterns = append(patterns, fmt.Sprintf(`("%s" <> $%d OR "%s" IS NULL)`, fieldName, index, fieldName))
						values = append(values, value["$ne"])
						index = index + 1
					}
				}
			}

			if v, ok := value["$eq"]; ok {
				if v == nil {
					patterns = append(patterns, fmt.Sprintf(`"%s" IS NULL`, fieldName))
				} else {
					patterns = append(patterns, fmt.Sprintf(`"%s" = $%d`, fieldName, index))
					values = append(values, v)
					index = index + 1
				}
			}

			inArray := utils.A(value["$in"])
			ninArray := utils.A(value["$nin"])
			isInOrNin := (inArray != nil) || (ninArray != nil)
			isTypeString := false
			if tp := utils.M(fields[fieldName]); tp != nil {
				if contents := utils.M(tp["contents"]); contents != nil {
					if utils.S(contents["type"]) == "String" {
						isTypeString = true
					}
				}
			}
			if inArray != nil && isArrayField && isTypeString {
				inPatterns := []string{}
				allowNull := false

				for listIndex, listElem := range inArray {
					if listElem == nil {
						allowNull = true
					} else {
						values = append(values, listElem)
						i := 0
						if allowNull {
							i = index + listIndex - 1
						} else {
							i = index + listIndex
						}
						inPatterns = append(inPatterns, fmt.Sprintf("$%d", i))
					}
				}

				if allowNull {
					patterns = append(patterns, fmt.Sprintf(`("%s" IS NULL OR "%s" && ARRAY[%s])`, fieldName, fieldName, strings.Join(inPatterns, ",")))
				} else {
					patterns = append(patterns, fmt.Sprintf(`("%s" && ARRAY[%s])`, fieldName, strings.Join(inPatterns, ",")))
				}
				index = index + len(inPatterns)
			} else if isInOrNin {
				createConstraint := func(baseArray types.S, notIn bool) {
					if len(baseArray) > 0 {
						not := ""
						if notIn {
							not = " NOT "
						}
						if isArrayField {
							patterns = append(patterns, fmt.Sprintf(`%s array_contains("%s", $%d)`, not, fieldName, index))
							j, _ := json.Marshal(baseArray)
							values = append(values, string(j))
							index = index + 1
						} else {
							inPatterns := []string{}
							for listIndex, listElem := range baseArray {
								values = append(values, listElem)
								inPatterns = append(inPatterns, fmt.Sprintf("$%d", index+listIndex))
							}
							patterns = append(patterns, fmt.Sprintf(`"%s" %s IN (%s)`, fieldName, not, strings.Join(inPatterns, ",")))
							index = index + len(inPatterns)
						}
					} else if !notIn {
						patterns = append(patterns, fmt.Sprintf(`"%s" IS NULL`, fieldName))
					}
				}
				if inArray != nil {
					createConstraint(inArray, false)
				}
				if ninArray != nil {
					createConstraint(ninArray, true)
				}
			}

			allArray := utils.A(value["$all"])
			if allArray != nil && isArrayField {
				patterns = append(patterns, fmt.Sprintf(`array_contains_all("%s", $%d::jsonb)`, fieldName, index))
				j, _ := json.Marshal(allArray)
				values = append(values, string(j))
				index = index + 1
			}

			if b, ok := value["$exists"].(bool); ok {
				if b {
					patterns = append(patterns, fmt.Sprintf(`"%s" IS NOT NULL`, fieldName))
				} else {
					patterns = append(patterns, fmt.Sprintf(`"%s" IS NULL`, fieldName))
				}
			}

			if point := utils.M(value["$nearSphere"]); point != nil {
				var distance float64
				if v, ok := value["$maxDistance"].(float64); ok {
					distance = v
				}
				distanceInKM := distance * 6371 * 1000
				patterns = append(patterns, fmt.Sprintf(`ST_distance_sphere("%s"::geometry, POINT($%d, $%d)::geometry) <= $%d`, fieldName, index, index+1, index+2))
				sorts = append(sorts, fmt.Sprintf(`ST_distance_sphere("%s"::geometry, POINT($%d, $%d)::geometry) ASC`, fieldName, index, index+1))
				values = append(values, point["longitude"], point["latitude"], distanceInKM)
				index = index + 3
			}

			if within := utils.M(value["$within"]); within != nil {
				if box := utils.A(within["$box"]); len(box) == 2 {
					box1 := utils.M(box[0])
					box2 := utils.M(box[1])
					if box1 != nil && box2 != nil {
						left := box1["longitude"]
						bottom := box1["latitude"]
						right := box2["longitude"]
						top := box2["latitude"]

						patterns = append(patterns, fmt.Sprintf(`"%s"::point <@ $%d::box`, fieldName, index))
						values = append(values, fmt.Sprintf("((%v, %v), (%v, %v))", left, bottom, right, top))
						index = index + 1
					}
				}
			}

			if regex := utils.S(value["$regex"]); regex != "" {
				operator := "~"
				opts := utils.S(value["$options"])
				if opts != "" {
					if strings.Contains(opts, "i") {
						operator = "~*"
					}
					if strings.Contains(opts, "x") {
						regex = removeWhiteSpace(regex)
					}
				}

				regex = processRegexPattern(regex)

				patterns = append(patterns, fmt.Sprintf(`"%s" %s '%s'`, fieldName, operator, regex))
			}

			if utils.S(value["__type"]) == "Pointer" {
				if isArrayField {
					patterns = append(patterns, fmt.Sprintf(`array_contains("%s", $%d)`, fieldName, index))
					j, _ := json.Marshal(types.S{value})
					values = append(values, string(j))
					index = index + 1
				} else {
					patterns = append(patterns, fmt.Sprintf(`"%s" = $%d`, fieldName, index))
					values = append(values, value["objectId"])
					index = index + 1
				}
			}

			if utils.S(value["__type"]) == "Date" {
				patterns = append(patterns, fmt.Sprintf(`"%s" = $%d`, fieldName, index))
				values = append(values, value["iso"])
				index = index + 1
			}

			for cmp, pgComparator := range parseToPosgresComparator {
				if v, ok := value[cmp]; ok {
					patterns = append(patterns, fmt.Sprintf(`"%s" %s $%d`, fieldName, pgComparator, index))
					values = append(values, toPostgresValue(v))
					index = index + 1
				}
			}
		}

		if fieldValue == nil {
			patterns = append(patterns, fmt.Sprintf(`"%s" IS NULL`, fieldName))
		}

		if initialPatternsLength == len(patterns) {
			s, _ := json.Marshal(fieldValue)
			return nil, errs.E(errs.OperationForbidden, "Postgres doesn't support this query type yet "+string(s))
		}
	}
	for i, v := range values {
		values[i] = transformValue(v)
	}
	return &whereClause{strings.Join(patterns, " AND "), values, sorts}, nil
}

func removeWhiteSpace(s string) string {
	if strings.HasSuffix(s, "\n") == false {
		s = s + "\n"
	}

	re := regexp.MustCompile(`(?im)^#.*\n`)
	s = re.ReplaceAllString(s, "")
	re = regexp.MustCompile(`(?im)([^\\])#.*\n`)
	s = re.ReplaceAllString(s, "$1")
	re = regexp.MustCompile(`(?im)([^\\])\s+`)
	s = re.ReplaceAllString(s, "$1")
	re = regexp.MustCompile(`^\s+`)
	s = re.ReplaceAllString(s, "")
	s = strings.TrimSpace(s)

	return s
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

func valueToDate(v interface{}) types.M {
	if v, ok := v.(time.Time); ok {
		return types.M{
			"__type": "Date",
			"iso":    utils.TimetoString(v),
		}
	}
	return nil
}

// Function to set a key on a nested JSON document
const jsonObjectSetKey = `CREATE OR REPLACE FUNCTION "json_object_set_key"(
  "json"          jsonb,
  "key_to_set"    TEXT,
  "value_to_set"  anyelement
)
  RETURNS jsonb 
  LANGUAGE sql 
  IMMUTABLE 
  STRICT 
AS $function$
SELECT concat('{', string_agg(to_json("key") || ':' || "value", ','), '}')::jsonb
  FROM (SELECT *
          FROM jsonb_each("json")
         WHERE "key" <> "key_to_set"
         UNION ALL
        SELECT "key_to_set", to_json("value_to_set")::jsonb) AS "fields"
$function$`

const arrayAdd = `CREATE OR REPLACE FUNCTION "array_add"(
  "array"   jsonb,
  "values"  jsonb
)
  RETURNS jsonb 
  LANGUAGE sql 
  IMMUTABLE 
  STRICT 
AS $function$ 
  SELECT array_to_json(ARRAY(SELECT unnest(ARRAY(SELECT DISTINCT jsonb_array_elements("array")) ||  ARRAY(SELECT jsonb_array_elements("values")))))::jsonb;
$function$`

const arrayAddUnique = `CREATE OR REPLACE FUNCTION "array_add_unique"(
  "array"   jsonb,
  "values"  jsonb
)
  RETURNS jsonb 
  LANGUAGE sql 
  IMMUTABLE 
  STRICT 
AS $function$ 
  SELECT array_to_json(ARRAY(SELECT DISTINCT unnest(ARRAY(SELECT DISTINCT jsonb_array_elements("array")) ||  ARRAY(SELECT DISTINCT jsonb_array_elements("values")))))::jsonb;
$function$`

const arrayRemove = `CREATE OR REPLACE FUNCTION "array_remove"(
  "array"   jsonb,
  "values"  jsonb
)
  RETURNS jsonb 
  LANGUAGE sql 
  IMMUTABLE 
  STRICT 
AS $function$ 
  SELECT array_to_json(ARRAY(SELECT * FROM jsonb_array_elements("array") as elt WHERE elt NOT IN (SELECT * FROM (SELECT jsonb_array_elements("values")) AS sub)))::jsonb;
$function$`

const arrayContainsAll = `CREATE OR REPLACE FUNCTION "array_contains_all"(
  "array"   jsonb,
  "values"  jsonb
)
  RETURNS boolean 
  LANGUAGE sql 
  IMMUTABLE 
  STRICT 
AS $function$ 
  SELECT RES.CNT = jsonb_array_length("values") FROM (SELECT COUNT(*) as CNT FROM jsonb_array_elements("array") as elt WHERE elt IN (SELECT jsonb_array_elements("values"))) as RES ;
$function$`

const arrayContains = `CREATE OR REPLACE FUNCTION "array_contains"(
  "array"   jsonb,
  "values"  jsonb
)
  RETURNS boolean 
  LANGUAGE sql 
  IMMUTABLE 
  STRICT 
AS $function$ 
  SELECT RES.CNT >= 1 FROM (SELECT COUNT(*) as CNT FROM jsonb_array_elements("array") as elt WHERE elt IN (SELECT jsonb_array_elements("values"))) as RES ;
$function$`
