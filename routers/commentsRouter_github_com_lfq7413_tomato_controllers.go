package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:AnalyticsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:AnalyticsController"],
		beego.ControllerComments{
			"HandleEvent",
			`/:eventName`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:AnalyticsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:AnalyticsController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:AnalyticsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:AnalyticsController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:AnalyticsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:AnalyticsController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:AnalyticsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:AnalyticsController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ClassesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ClassesController"],
		beego.ControllerComments{
			"HandleCreate",
			`/:className`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ClassesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ClassesController"],
		beego.ControllerComments{
			"HandleGet",
			`/:className/:objectId`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ClassesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ClassesController"],
		beego.ControllerComments{
			"HandleUpdate",
			`/:className/:objectId`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ClassesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ClassesController"],
		beego.ControllerComments{
			"HandleFind",
			`/:className`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ClassesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ClassesController"],
		beego.ControllerComments{
			"HandleDelete",
			`/:className/:objectId`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ClassesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ClassesController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ClassesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ClassesController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ClassesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ClassesController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ClassesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ClassesController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:CloudCodeController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:CloudCodeController"],
		beego.ControllerComments{
			"HandleGet",
			`/jobs`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:CloudCodeController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:CloudCodeController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:CloudCodeController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:CloudCodeController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:CloudCodeController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:CloudCodeController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:CloudCodeController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:CloudCodeController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FeaturesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FeaturesController"],
		beego.ControllerComments{
			"HandleGet",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FeaturesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FeaturesController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FeaturesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FeaturesController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FeaturesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FeaturesController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FilesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FilesController"],
		beego.ControllerComments{
			"HandleGet",
			`/:appId/:filename`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FilesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FilesController"],
		beego.ControllerComments{
			"HandleCreate",
			`/:filename`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FilesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FilesController"],
		beego.ControllerComments{
			"HandleDelete",
			`/:filename`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FilesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FilesController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FilesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FilesController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FilesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FilesController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FilesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FilesController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FunctionsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FunctionsController"],
		beego.ControllerComments{
			"HandleCloudFunction",
			`/:functionName`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FunctionsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FunctionsController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FunctionsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FunctionsController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FunctionsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FunctionsController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FunctionsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:FunctionsController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:GlobalConfigController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:GlobalConfigController"],
		beego.ControllerComments{
			"HandleGet",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:GlobalConfigController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:GlobalConfigController"],
		beego.ControllerComments{
			"HandlePut",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:GlobalConfigController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:GlobalConfigController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:GlobalConfigController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:GlobalConfigController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"],
		beego.ControllerComments{
			"HandleGetAllFunctions",
			`/functions`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"],
		beego.ControllerComments{
			"HandleGetFunction",
			`/functions/:functionName`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"],
		beego.ControllerComments{
			"HandleCreateFunction",
			`/functions`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"],
		beego.ControllerComments{
			"HandleUpdateFunction",
			`/functions/:functionName`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"],
		beego.ControllerComments{
			"HandleGetAllTriggers",
			`/triggers`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"],
		beego.ControllerComments{
			"HandleGetTrigger",
			`/triggers/:className/:triggerName`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"],
		beego.ControllerComments{
			"HandleCreateTrigger",
			`/triggers`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"],
		beego.ControllerComments{
			"HandleUpdateTrigger",
			`/triggers/:className/:triggerName`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:HooksController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:IAPValidationController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:IAPValidationController"],
		beego.ControllerComments{
			"HandlePost",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:IAPValidationController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:IAPValidationController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:IAPValidationController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:IAPValidationController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:IAPValidationController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:IAPValidationController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:InstallationsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:InstallationsController"],
		beego.ControllerComments{
			"HandleFind",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:InstallationsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:InstallationsController"],
		beego.ControllerComments{
			"HandleGet",
			`/:objectId`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:InstallationsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:InstallationsController"],
		beego.ControllerComments{
			"HandleCreate",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:InstallationsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:InstallationsController"],
		beego.ControllerComments{
			"HandleUpdate",
			`/:objectId`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:InstallationsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:InstallationsController"],
		beego.ControllerComments{
			"HandleDelete",
			`/:objectId`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:InstallationsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:InstallationsController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:InstallationsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:InstallationsController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:JobsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:JobsController"],
		beego.ControllerComments{
			"HandleCloudJob",
			`/:jobName`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:JobsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:JobsController"],
		beego.ControllerComments{
			"HandlePost",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:JobsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:JobsController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:JobsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:JobsController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:JobsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:JobsController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LoginController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LoginController"],
		beego.ControllerComments{
			"HandleLogIn",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LoginController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LoginController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LoginController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LoginController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LoginController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LoginController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogoutController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogoutController"],
		beego.ControllerComments{
			"HandleLogOut",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogoutController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogoutController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogoutController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogoutController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogoutController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogoutController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogsController"],
		beego.ControllerComments{
			"HandleGet",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogsController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogsController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:LogsController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PublicController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PublicController"],
		beego.ControllerComments{
			"VerifyEmail",
			`/verify_email`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PublicController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PublicController"],
		beego.ControllerComments{
			"ChangePassword",
			`/choose_password`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PublicController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PublicController"],
		beego.ControllerComments{
			"ResetPassword",
			`/request_password_reset`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PublicController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PublicController"],
		beego.ControllerComments{
			"RequestResetPassword",
			`/request_password_reset`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PublicController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PublicController"],
		beego.ControllerComments{
			"InvalidLink",
			`/invalid_link`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PublicController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PublicController"],
		beego.ControllerComments{
			"PasswordResetSuccess",
			`/password_reset_success`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PublicController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PublicController"],
		beego.ControllerComments{
			"VerifyEmailSuccess",
			`/verify_email_success`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PurgeController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PurgeController"],
		beego.ControllerComments{
			"HandleDelete",
			`/:className`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PurgeController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PurgeController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PurgeController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PurgeController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PurgeController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PurgeController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PushController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PushController"],
		beego.ControllerComments{
			"HandlePost",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PushController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PushController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PushController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PushController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PushController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:PushController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ResetController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ResetController"],
		beego.ControllerComments{
			"HandleResetRequest",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ResetController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ResetController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ResetController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ResetController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ResetController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:ResetController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:RolesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:RolesController"],
		beego.ControllerComments{
			"HandleFind",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:RolesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:RolesController"],
		beego.ControllerComments{
			"HandleGet",
			`/:objectId`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:RolesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:RolesController"],
		beego.ControllerComments{
			"HandleCreate",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:RolesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:RolesController"],
		beego.ControllerComments{
			"HandleUpdate",
			`/:objectId`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:RolesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:RolesController"],
		beego.ControllerComments{
			"HandleDelete",
			`/:objectId`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:RolesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:RolesController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:RolesController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:RolesController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SchemasController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SchemasController"],
		beego.ControllerComments{
			"HandleFind",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SchemasController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SchemasController"],
		beego.ControllerComments{
			"HandleGet",
			`/:className`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SchemasController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SchemasController"],
		beego.ControllerComments{
			"HandleCreate",
			`/:className`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SchemasController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SchemasController"],
		beego.ControllerComments{
			"HandleUpdate",
			`/:className`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SchemasController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SchemasController"],
		beego.ControllerComments{
			"HandleDelete",
			`/:className`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SchemasController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SchemasController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SchemasController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SchemasController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SessionsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SessionsController"],
		beego.ControllerComments{
			"HandleFind",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SessionsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SessionsController"],
		beego.ControllerComments{
			"HandleGet",
			`/:objectId`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SessionsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SessionsController"],
		beego.ControllerComments{
			"HandleCreate",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SessionsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SessionsController"],
		beego.ControllerComments{
			"HandleUpdate",
			`/:objectId`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SessionsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SessionsController"],
		beego.ControllerComments{
			"HandleDelete",
			`/:objectId`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SessionsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SessionsController"],
		beego.ControllerComments{
			"HandleGetMe",
			`/me`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SessionsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SessionsController"],
		beego.ControllerComments{
			"HandleUpdateMe",
			`/me`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SessionsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SessionsController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SessionsController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:SessionsController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"],
		beego.ControllerComments{
			"HandleFind",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"],
		beego.ControllerComments{
			"HandleGet",
			`/:objectId`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"],
		beego.ControllerComments{
			"HandleCreate",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"],
		beego.ControllerComments{
			"HandleUpdate",
			`/:objectId`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"],
		beego.ControllerComments{
			"HandleDelete",
			`/:objectId`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"],
		beego.ControllerComments{
			"HandleMe",
			`/me`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"] = append(beego.GlobalControllerRouter["github.com/lfq7413/tomato/controllers:UsersController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

}
