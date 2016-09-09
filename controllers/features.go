package controllers

import "github.com/lfq7413/tomato/types"

// FeaturesController ...
type FeaturesController struct {
	ClassesController
}

// HandleGet ...
// @router / [get]
func (f *FeaturesController) HandleGet() {
	if f.EnforceMasterKeyAccess() == false {
		return
	}

	features := types.M{
		"globalConfig": types.M{
			"create": true,
			"read":   true,
			"update": true,
			"delete": true,
		},
		"hooks": types.M{
			"create": false,
			"read":   false,
			"update": false,
			"delete": false,
		},
		"cloudCode": types.M{
			"jobs": true,
		},
		"logs": types.M{
			"level": true,
			"size":  true,
			"order": true,
			"until": true,
			"from":  true,
		},
		"push": types.M{
			"immediatePush":  true,
			"scheduledPush":  false,
			"storedPushData": true,
			"pushAudiences":  false,
		},
		"schemas": types.M{
			"addField":                  true,
			"removeField":               true,
			"addClass":                  true,
			"removeClass":               true,
			"clearAllDataFromClass":     true,
			"exportClass":               false,
			"editClassLevelPermissions": true,
			"editPointerPermissions":    true,
		},
	}
	f.Data["json"] = types.M{
		"features":           features,
		"parseServerVersion": "1.0",
	}
	f.ServeJSON()
}

// Post ...
// @router / [post]
func (f *FeaturesController) Post() {
	f.ClassesController.Post()
}

// Delete ...
// @router / [delete]
func (f *FeaturesController) Delete() {
	f.ClassesController.Delete()
}

// Put ...
// @router / [put]
func (f *FeaturesController) Put() {
	f.ClassesController.Put()
}
