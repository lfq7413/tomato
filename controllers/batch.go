package controllers

import (
	"errors"

	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// BatchController ...
type BatchController struct {
	ClassesController
}

// HandleBatch ...
// @router / [post]
func (b *BatchController) HandleBatch() {
	if b.JSONBody == nil {
		b.HandleError(errs.E(errs.InvalidJSON, "requests must be an array"), 0)
		return
	}
	requests := utils.A(b.JSONBody["requests"])
	if requests == nil {
		b.HandleError(errs.E(errs.InvalidJSON, "requests must be an array"), 0)
		return
	}

	headers := map[string]string{
		"X-Parse-Application-Id": b.Info.AppID,
	}
	if b.Info.MasterKey != "" {
		headers["X-Parse-Master-Key"] = b.Info.MasterKey
	}
	if b.Info.ClientKey != "" {
		headers["X-Parse-Client-Key"] = b.Info.ClientKey
	}
	if b.Info.JavascriptKey != "" {
		headers["X-Parse-Javascript-Key"] = b.Info.JavascriptKey
	}
	if b.Info.DotNetKey != "" {
		headers["X-Parse-Windows-Key"] = b.Info.DotNetKey
	}
	if b.Info.RestAPIKey != "" {
		headers["X-Parse-REST-API-Key"] = b.Info.RestAPIKey
	}
	if b.Info.SessionToken != "" {
		headers["X-Parse-Session-Token"] = b.Info.SessionToken
	}
	if b.Info.InstallationID != "" {
		headers["X-Parse-Installation-Id"] = b.Info.InstallationID
	}
	if b.Info.ClientVersion != "" {
		headers["X-Parse-Client-Version"] = b.Info.ClientVersion
	}
	if b.Ctx.Input.Header("Authorization") != "" {
		headers["Authorization"] = b.Ctx.Input.Header("Authorization")
	}

	b.Request(requests, headers, b.Ctx.Input.Site())
}

// Request ...
func (b *BatchController) Request(requests types.S, headers map[string]string, site string) {

}

// Get ...
// @router / [get]
func (b *BatchController) Get() {
	b.HandleError(errors.New("Method Not Allowed"), 405)
}

// Delete ...
// @router / [delete]
func (b *BatchController) Delete() {
	b.HandleError(errors.New("Method Not Allowed"), 405)
}

// Put ...
// @router / [put]
func (b *BatchController) Put() {
	b.HandleError(errors.New("Method Not Allowed"), 405)
}
