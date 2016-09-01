package postgres

import "github.com/lfq7413/tomato/types"

const postgresSchemaCollectionName = "_SCHEMA"

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
