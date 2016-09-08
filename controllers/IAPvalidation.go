package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/astaxie/beego"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

const iapSandboxURL = "https://sandbox.itunes.apple.com/verifyReceipt"
const iapProductionURL = "https://buy.itunes.apple.com/verifyReceipt"

var appStoreErrors = map[int]string{
	21000: "The App Store could not read the JSON object you provided.",
	21002: "The data in the receipt-data property was malformed or missing.",
	21003: "The receipt could not be authenticated.",
	21004: "The shared secret you provided does not match the shared secret on file for your account.",
	21005: "The receipt server is not currently available.",
	21006: "This receipt is valid but the subscription has expired.",
	21007: "This receipt is from the test environment, but it was sent to the production environment for verification. Send it to the test environment instead.",
	21008: "This receipt is from the production environment, but it was sent to the test environment for verification. Send it to the production environment instead.",
}

// IAPValidationController ...
type IAPValidationController struct {
	ClassesController
}

// HandlePost ...
// @router / [post]
func (i *IAPValidationController) HandlePost() {
	receipt := i.JSONBody["receipt"]
	productIdentifier := i.JSONBody["productIdentifier"]
	if receipt == nil || productIdentifier == nil {
		i.HandleError(errs.E(errs.InvalidJSON, "missing receipt or productIdentifier"), 0)
		return
	}

	if o := utils.M(receipt); o != nil {
		if utils.S(o["__type"]) == "Bytes" {
			receipt = o["base64"]
		}
	}

	if beego.AppConfig.String("runmode") == "dev" && i.JSONBody["bypassAppStoreValidation"] != nil {
		i.getFileForProductIdentifier(utils.S(productIdentifier))
		return
	}

	result := validateWithAppStore(iapProductionURL, utils.S(receipt))
	if result == nil {
		i.getFileForProductIdentifier(utils.S(productIdentifier))
		return
	}
	if v, ok := result["status"].(float64); ok {
		if v == 21007 {
			r := validateWithAppStore(iapSandboxURL, utils.S(receipt))
			if r == nil {
				i.getFileForProductIdentifier(utils.S(productIdentifier))
				return
			}
			i.Data["json"] = appStoreError(r)
			i.ServeJSON()
			return
		}
	}
	i.Data["json"] = appStoreError(result)
	i.ServeJSON()
}

func (i *IAPValidationController) getFileForProductIdentifier(productIdentifier string) {
	r, err := rest.Find(i.Auth, "_Product", types.M{"productIdentifier": productIdentifier}, types.M{}, i.Info.ClientSDK)
	if err != nil {
		i.HandleError(err, 0)
		return
	}
	products := utils.A(r["results"])
	if products == nil || len(products) != 1 {
		i.HandleError(errs.E(errs.ObjectNotFound, "Object not found."), 0)
		return
	}
	product := utils.M(products[0])
	if product == nil {
		i.HandleError(errs.E(errs.ObjectNotFound, "Object not found."), 0)
		return
	}
	i.Data["json"] = product["download"]
	i.ServeJSON()
}

func validateWithAppStore(URL, receipt string) types.M {
	jsonParams := `{ "receipt-data": ` + receipt + ` }`
	encodeParams := url.QueryEscape(jsonParams)
	request, err := http.NewRequest("POST", URL, strings.NewReader(encodeParams))
	if err != nil {
		return types.M{}
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml,application/json;q=0.9,*/*;q=0.8")
	request.Header.Add("Accept-Language", "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return types.M{}
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return types.M{}
	}

	var result types.M
	err = json.Unmarshal(body, &result)
	if err != nil {
		return types.M{}
	}
	if status, ok := result["result"].(float64); ok && status == 0 {
		return nil
	}

	return result
}

func appStoreError(r types.M) types.M {
	var status float64
	var errorString = "unknown error."
	if v, ok := r["status"].(float64); ok {
		status = v
		errorString = appStoreErrors[int(status)]
		if errorString == "" {
			errorString = "unknown error."
		}
	}
	return types.M{"status": status, "error": errorString}
}

// Get ...
// @router / [get]
func (i *IAPValidationController) Get() {
	i.ClassesController.Get()
}

// Delete ...
// @router / [delete]
func (i *IAPValidationController) Delete() {
	i.ClassesController.Delete()
}

// Put ...
// @router / [put]
func (i *IAPValidationController) Put() {
	i.ClassesController.Put()
}
