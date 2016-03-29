package orm

// CollectionExists ...
func CollectionExists(className string) bool {
	return TomatoDB.collectionExists(className)
}

// Find ...
func Find(className string, where map[string]interface{}, options map[string]interface{}) []interface{} {
	// TODO
	return []interface{}{}
}

// Count ...
func Count(className string, where map[string]interface{}, options map[string]interface{}) int {
	// TODO
	return 0
}

// Destroy ...
func Destroy(className string, where map[string]interface{}, options map[string]interface{}) {
	// TODO
}

// Update ...
func Update(className string, where map[string]interface{}, data map[string]interface{}, options map[string]interface{}) error {
	// TODO
	return nil
}

// UpdateAll ...
func UpdateAll(className string, where map[string]interface{}, data map[string]interface{}, options map[string]interface{}) error {
	// TODO
	return nil
}

// Create ...
func Create(className string, data map[string]interface{}, options map[string]interface{}) error {
	// TODO
	return nil
}

// ValidateObject ...
func ValidateObject(className string, object map[string]interface{}, where map[string]interface{}, options map[string]interface{}) error {
	// TODO
	return nil
}

// LoadSchema ...
func LoadSchema() {
	// TODO
}

// CanAddField ...
func CanAddField() {
	// TODO
}
