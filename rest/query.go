package rest

import (
	"sort"
	"strings"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/files"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// Query 处理查询请求的结构体
type Query struct {
	auth              *Auth
	className         string
	Where             types.M
	findOptions       types.M
	response          types.M
	doCount           bool
	include           [][]string
	keys              []string
	redirectKey       string
	redirectClassName string
	clientSDK         map[string]string
}

// NewQuery 组装查询对象
func NewQuery(
	auth *Auth,
	className string,
	where types.M,
	options types.M,
	clientSDK map[string]string,
) (*Query, error) {
	if auth == nil {
		auth = Nobody()
	}
	if where == nil {
		where = types.M{}
	}
	query := &Query{
		auth:              auth,
		className:         className,
		Where:             where,
		findOptions:       types.M{},
		response:          types.M{},
		doCount:           false,
		include:           [][]string{},
		keys:              []string{},
		redirectKey:       "",
		redirectClassName: "",
		clientSDK:         clientSDK,
	}

	if auth.IsMaster == false {
		// 当前权限为 Master 时，findOptions 中不存在 acl 这个 key
		if auth.User != nil {
			query.findOptions["acl"] = []string{utils.S(auth.User["objectId"])}
		} else {
			query.findOptions["acl"] = nil
		}
		if className == "_Session" {
			if query.findOptions["acl"] == nil {
				return nil, errs.E(errs.InvalidSessionToken, "This session token is invalid.")
			}
			user := types.M{
				"user": types.M{
					"__type":    "Pointer",
					"className": "_User",
					"objectId":  auth.User["objectId"],
				},
			}
			query.Where = types.M{
				"$and": types.S{where, user},
			}
		}
	}

	for k, v := range options {
		switch k {
		case "keys":
			if s, ok := v.(string); ok {
				query.keys = strings.Split(s, ",")
				query.keys = append(query.keys, "objectId", "createdAt", "updatedAt")
			}
		case "count":
			query.doCount = true
		case "skip":
			query.findOptions["skip"] = v
		case "limit":
			query.findOptions["limit"] = v
		case "order":
			if s, ok := v.(string); ok {
				fields := strings.Split(s, ",")
				// sortMap := map[string]int{}
				// for _, v := range fields {
				// 	if strings.HasPrefix(v, "-") {
				// 		sortMap[v[1:]] = -1
				// 	} else {
				// 		sortMap[v] = 1
				// 	}
				// }
				// query.findOptions["sort"] = sortMap
				query.findOptions["sort"] = fields
			}
		case "include":
			if s, ok := v.(string); ok { // v = "user.session,name.friend"
				paths := strings.Split(s, ",") // paths = ["user.session","name.friend"]
				pathSet := []string{}
				for _, path := range paths {
					parts := strings.Split(path, ".") // parts = ["user","session"]
					for lenght := 1; lenght <= len(parts); lenght++ {
						pathSet = append(pathSet, strings.Join(parts[0:lenght], "."))
					} // pathSet = ["user","user.session"]
				} // pathSet = ["user","user.session","name","name.friend"]
				sort.Strings(pathSet) // pathSet = ["name","name.friend","user","user.session"]
				for _, set := range pathSet {
					query.include = append(query.include, strings.Split(set, "."))
				} // query.include = [["name"],["name","friend"],["user"],["user","seeeion"]]
			}
		case "redirectClassNameForKey":
			if s, ok := v.(string); ok {
				query.redirectKey = s
				query.redirectClassName = ""
			}
		default:
			return nil, errs.E(errs.InvalidJSON, "bad option: "+k)
		}
	}

	return query, nil
}

// Execute 执行查询请求，返回的数据包含 results count 两个字段
func (q *Query) Execute() (types.M, error) {

	err := q.BuildRestWhere()
	if err != nil {
		return nil, err
	}
	err = q.runFind()
	if err != nil {
		return nil, err
	}
	err = q.runCount()
	if err != nil {
		return nil, err
	}
	err = q.handleInclude()
	if err != nil {
		return nil, err
	}
	return q.response, nil
}

// BuildRestWhere 展开查询参数，组装设置项
func (q *Query) BuildRestWhere() error {
	err := q.getUserAndRoleACL()
	if err != nil {
		return err
	}
	err = q.redirectClassNameForKey()
	if err != nil {
		return err
	}
	err = q.validateClientClassCreation()
	if err != nil {
		return err
	}
	err = q.replaceSelect()
	if err != nil {
		return err
	}
	err = q.replaceDontSelect()
	if err != nil {
		return err
	}
	err = q.replaceInQuery()
	if err != nil {
		return err
	}
	err = q.replaceNotInQuery()
	if err != nil {
		return err
	}
	return nil
}

// getUserAndRoleACL 获取当前用户角色信息，以及用户 id，添加到设置项 acl 中
func (q *Query) getUserAndRoleACL() error {
	if q.auth.IsMaster || q.auth.User == nil {
		return nil
	}
	acl := []string{utils.S(q.auth.User["objectId"])}
	roles := q.auth.GetUserRoles()
	acl = append(acl, roles...)
	q.findOptions["acl"] = acl
	return nil
}

// redirectClassNameForKey 修改 className 为 redirectKey 字段对应的相关类型
func (q *Query) redirectClassNameForKey() error {
	if q.redirectKey == "" {
		return nil
	}

	newClassName := orm.TomatoDBController.RedirectClassNameForKey(q.className, q.redirectKey)
	q.className = newClassName
	q.redirectClassName = newClassName

	return nil
}

// validateClientClassCreation 验证当前请求是否能创建类
func (q *Query) validateClientClassCreation() error {
	// 检测配置项是否允许
	if config.TConfig.AllowClientClassCreation {
		return nil
	}
	if q.auth.IsMaster {
		return nil
	}
	// 允许操作系统表
	for _, v := range orm.SystemClasses {
		if v == q.className {
			return nil
		}
	}
	// 允许操作已存在的表
	schema := orm.TomatoDBController.LoadSchema(nil)
	hasClass := schema.HasClass(q.className)
	if hasClass {
		return nil
	}

	// 无法操作不存在的表
	return errs.E(errs.OperationForbidden, "This user is not allowed to access non-existent class: "+q.className)
}

// replaceSelect 执行 $select 中的查询语句，把结果放入 $in 中，替换掉 $select
// 替换前的格式如下：
// {
//     "hometown":{
//         "$select":{
//             "query":{
//                 "className":"Team",
//                 "where":{
//                     "winPct":{
//                         "$gt":0.5
//                     }
//                 }
//             },
//             "key":"city"
//         }
//     }
// }
// 转换后格式如下
// {
//     "hometown":{
//         "$in":["abc","cba"]
//     }
// }
func (q *Query) replaceSelect() error {
	selectObject := findObjectWithKey(q.Where, "$select")
	if selectObject == nil {
		return nil
	}
	// 必须包含两个 key ： query key
	selectValue := utils.M(selectObject["$select"])
	if selectValue == nil ||
		selectValue["query"] == nil ||
		utils.S(selectValue["key"]) == "" ||
		len(selectValue) != 2 {
		return errs.E(errs.InvalidQuery, "improper usage of $select")
	}
	queryValue := utils.M(selectValue["query"])
	// iOS SDK 中不设置 where 时，没有 where 字段，所以此处不检测 where
	if queryValue == nil ||
		utils.S(queryValue["className"]) == "" {
		return errs.E(errs.InvalidQuery, "improper usage of $select")
	}

	additionalOptions := types.M{
		"redirectClassNameForKey": queryValue["redirectClassNameForKey"],
	}

	var where types.M
	if w := utils.M(queryValue["where"]); w == nil {
		where = types.M{}
	} else {
		where = w
	}
	query, err := NewQuery(
		q.auth,
		utils.S(queryValue["className"]),
		where,
		additionalOptions, q.clientSDK)
	if err != nil {
		return err
	}
	response, err := query.Execute()
	if err != nil {
		return err
	}
	// 组装查询到的对象
	values := []types.M{}
	if utils.HasResults(response) == true {
		for _, v := range utils.A(response["results"]) {
			result := utils.M(v)
			if result == nil {
				continue
			}
			values = append(values, result)
		}
	}
	// 替换 $select 为 $in
	transformSelect(selectObject, utils.S(selectValue["key"]), values)
	// 继续搜索替换
	return q.replaceSelect()
}

// replaceDontSelect 执行 $dontSelect 中的查询语句，把结果放入 $nin 中，替换掉 $select
// 数据结构与 replaceSelect 类似
func (q *Query) replaceDontSelect() error {
	dontSelectObject := findObjectWithKey(q.Where, "$dontSelect")
	if dontSelectObject == nil {
		return nil
	}
	// 必须包含两个 key ： query key
	dontSelectValue := utils.M(dontSelectObject["$dontSelect"])
	if dontSelectValue == nil ||
		dontSelectValue["query"] == nil ||
		utils.S(dontSelectValue["key"]) == "" ||
		len(dontSelectValue) != 2 {
		return errs.E(errs.InvalidQuery, "improper usage of $dontSelect")
	}
	queryValue := utils.M(dontSelectValue["query"])
	if queryValue == nil ||
		utils.S(queryValue["className"]) == "" {
		return errs.E(errs.InvalidQuery, "improper usage of $dontSelect")
	}

	additionalOptions := types.M{
		"redirectClassNameForKey": queryValue["redirectClassNameForKey"],
	}

	var where types.M
	if w := utils.M(queryValue["where"]); w == nil {
		where = types.M{}
	} else {
		where = w
	}
	query, err := NewQuery(
		q.auth,
		utils.S(queryValue["className"]),
		where,
		additionalOptions, q.clientSDK)
	if err != nil {
		return err
	}
	response, err := query.Execute()
	if err != nil {
		return err
	}
	// 组装查询到的对象
	values := []types.M{}
	if utils.HasResults(response) == true {
		for _, v := range utils.A(response["results"]) {
			result := utils.M(v)
			if result == nil {
				continue
			}
			values = append(values, result)
		}
	}
	// 替换 $dontSelect 为 $nin
	transformDontSelect(dontSelectObject, utils.S(dontSelectValue["key"]), values)
	// 继续搜索替换
	return q.replaceDontSelect()
}

// replaceInQuery 执行 $inQuery 中的查询语句，把结果放入 $in 中，替换掉 $inQuery
// 替换前的格式：
// {
//     "post":{
//         "$inQuery":{
//             "where":{
//                 "image":{
//                     "$exists":true
//                 }
//             },
//             "className":"Post"
//         }
//     }
// }
// 替换后的格式
// {
//     "post":{
//         "$in":[
// 			{
// 				"__type":    "Pointer",
// 				"className": "className",
// 				"objectId":  "objectId",
// 			},
// 			{...}
// 		]
//     }
// }
func (q *Query) replaceInQuery() error {
	inQueryObject := findObjectWithKey(q.Where, "$inQuery")
	if inQueryObject == nil {
		return nil
	}
	// 必须包含两个 key ： where className
	inQueryValue := utils.M(inQueryObject["$inQuery"])
	if inQueryValue == nil ||
		utils.M(inQueryValue["where"]) == nil ||
		utils.S(inQueryValue["className"]) == "" ||
		len(inQueryValue) != 2 {
		return errs.E(errs.InvalidQuery, "improper usage of $inQuery")
	}

	additionalOptions := types.M{
		"redirectClassNameForKey": inQueryValue["redirectClassNameForKey"],
	}

	query, err := NewQuery(
		q.auth,
		utils.S(inQueryValue["className"]),
		utils.M(inQueryValue["where"]),
		additionalOptions, q.clientSDK)
	if err != nil {
		return err
	}
	response, err := query.Execute()
	if err != nil {
		return err
	}
	// 组装查询到的对象
	values := []types.M{}
	if utils.HasResults(response) == true {
		for _, v := range utils.A(response["results"]) {
			result := utils.M(v)
			if result == nil {
				continue
			}
			values = append(values, result)
		}
	}
	// 替换 $inQuery 为 $in
	transformInQuery(inQueryObject, query.className, values)
	// 继续搜索替换
	return q.replaceInQuery()
}

// replaceNotInQuery 执行 $notInQuery 中的查询语句，把结果放入 $nin 中，替换掉 $notInQuery
// 数据格式与 replaceInQuery 类似
func (q *Query) replaceNotInQuery() error {
	notInQueryObject := findObjectWithKey(q.Where, "$notInQuery")
	if notInQueryObject == nil {
		return nil
	}
	// 必须包含两个 key ： where className
	notInQueryValue := utils.M(notInQueryObject["$notInQuery"])
	if notInQueryValue == nil ||
		utils.M(notInQueryValue["where"]) == nil ||
		utils.S(notInQueryValue["className"]) == "" ||
		len(notInQueryValue) != 2 {
		return errs.E(errs.InvalidQuery, "improper usage of $notInQuery")
	}

	additionalOptions := types.M{
		"redirectClassNameForKey": notInQueryValue["redirectClassNameForKey"],
	}

	query, err := NewQuery(
		q.auth,
		utils.S(notInQueryValue["className"]),
		utils.M(notInQueryValue["where"]),
		additionalOptions, q.clientSDK)
	if err != nil {
		return err
	}
	response, err := query.Execute()
	if err != nil {
		return err
	}
	// 组装查询到的对象
	values := []types.M{}
	if utils.HasResults(response) == true {
		for _, v := range utils.A(response["results"]) {
			result := utils.M(v)
			if result == nil {
				continue
			}
			values = append(values, result)
		}
	}
	// 替换 $notInQuery 为 $nin
	transformNotInQuery(notInQueryObject, query.className, values)
	// 继续搜索替换
	return q.replaceNotInQuery()
}

// runFind 从数据库查找数据，并处理返回结果
func (q *Query) runFind() error {
	if q.findOptions["limit"] != nil {
		if l, ok := q.findOptions["limit"].(float64); ok {
			if l == 0 {
				q.response["results"] = types.S{}
				return nil
			}
		} else if l, ok := q.findOptions["limit"].(int); ok {
			if l == 0 {
				q.response["results"] = types.S{}
				return nil
			}
		}
	}
	response, err := orm.TomatoDBController.Find(q.className, q.Where, q.findOptions)
	if err != nil {
		return err
	}
	// 从 _User 表中删除密码字段
	if q.className == "_User" {
		for _, v := range response {
			if user := utils.M(v); user != nil {
				delete(user, "password")
				if authData := utils.M(user["authData"]); authData != nil {
					for provider, v := range authData {
						if v == nil {
							delete(authData, provider)
						}
					}
					if len(authData) == 0 {
						delete(user, "authData")
					}
				}
			}
		}
	}

	// 展开文件类型
	files.ExpandFilesInObject(response)

	// 取出需要的 key   （TODO：通过数据库直接取key）
	results := types.S{}
	if len(q.keys) > 0 && len(response) > 0 {
		for _, v := range response {
			obj := utils.M(v)
			newObj := types.M{}
			for _, s := range q.keys {
				if obj[s] != nil {
					newObj[s] = obj[s]
				}
			}
			results = append(results, newObj)
		}
	} else {
		results = append(results, response...)
	}

	if q.redirectClassName != "" {
		for _, v := range results {
			if r := utils.M(v); r != nil {
				r["className"] = q.redirectClassName
			}
		}
	}

	q.response["results"] = results
	return nil
}

// runCount 查询符合条件的结果数量
func (q *Query) runCount() error {
	if q.doCount == false {
		return nil
	}
	q.findOptions["count"] = true
	delete(q.findOptions, "skip")
	delete(q.findOptions, "limit")
	// 当需要取 count 时，数据库返回结果的第一个即为 count
	result, err := orm.TomatoDBController.Find(q.className, q.Where, q.findOptions)
	if err != nil {
		return err
	}
	if result == nil || len(result) == 0 {
		q.response["count"] = 0
	} else {
		q.response["count"] = result[0]
	}
	return nil
}

// handleInclude 展开 include 对应的内容
func (q *Query) handleInclude() error {
	if len(q.include) == 0 {
		return nil
	}
	// includePath 中会直接更新 q.response
	err := includePath(q.auth, q.response, q.include[0])
	if err != nil {
		return err
	}

	if len(q.include) > 0 {
		q.include = q.include[1:]
		return q.handleInclude()
	}

	return nil
}

// includePath 在 response 中搜索 path 路径中对应的节点，
// 查询出该节点对应的对象，然后用对象替换该节点
func includePath(auth *Auth, response types.M, path []string) error {
	// 查找路径对应的所有节点
	pointers := findPointers(response["results"], path)
	if len(pointers) == 0 {
		return nil
	}
	pointersHash := map[string]types.S{}
	for _, pointer := range pointers {
		// 不再区分不同 className ，添加不为空的 className
		className := utils.S(pointer["className"])
		objectID := utils.S(pointer["objectId"])
		if className != "" && objectID != "" {
			if v, ok := pointersHash[className]; ok {
				v = append(v, objectID)
				pointersHash[className] = v
			} else {
				pointersHash[className] = types.S{objectID}
			}
		}

	}

	replace := types.M{}
	for clsName, ids := range pointersHash {
		// 获取所有 ids 对应的对象
		objectID := types.M{
			"$in": ids,
		}
		where := types.M{
			"objectId": objectID,
		}
		query, err := NewQuery(auth, clsName, where, types.M{}, nil)
		if err != nil {
			return err
		}
		includeResponse, err := query.Execute()
		if err != nil {
			return err
		}
		if utils.HasResults(includeResponse) == false {
			return nil
		}

		// 组装查询到的对象
		results := utils.A(includeResponse["results"])
		for _, v := range results {
			obj := utils.M(v)
			if obj == nil {
				continue
			}
			obj["__type"] = "Object"
			obj["className"] = clsName
			if clsName == "_User" && auth.IsMaster == false {
				delete(obj, "sessionToken")
				delete(obj, "authData")
			}
			replace[utils.S(obj["objectId"])] = obj
		}
	}
	// 使用查询到的对象替换对应的节点
	replacePointers(pointers, replace)

	return nil
}

// findPointers 查询路径对应的对象列表，对象必须为 Pointer 类型
func findPointers(object interface{}, path []string) []types.M {
	if object == nil {
		return []types.M{}
	}
	// 如果是对象数组，则遍历每一个对象
	if s := utils.A(object); s != nil {
		answer := []types.M{}
		for _, v := range s {
			p := findPointers(v, path)
			answer = append(answer, p...)
		}
		return answer
	}

	// 如果不能转成 map ，则返回错误
	obj := utils.M(object)
	if obj == nil {
		return []types.M{}
	}
	// 如果当前是路径最后一个节点，判断是否为 Pointer
	if len(path) == 0 {
		if utils.S(obj["__type"]) == "Pointer" {
			return []types.M{obj}
		}
		return []types.M{}
	}
	// 取出下一个路径对应的对象，进行查找
	subobject := obj[path[0]]
	if subobject == nil {
		// 对象不存在，则不进行处理
		return []types.M{}
	}
	return findPointers(subobject, path[1:])
}

// replacePointers 把 replace 保存的对象，添加到 pointers 对应的节点中
// pointers 中保存的是指向 response 的引用，修改 pointers 中的内容，即可同时修改 response 的内容
func replacePointers(pointers []types.M, replace types.M) {
	if replace == nil {
		return
	}
	for _, pointer := range pointers {
		if pointer == nil {
			continue
		}
		objectID := utils.S(pointer["objectId"])
		if objectID == "" {
			continue
		}
		if rpl := utils.M(replace[objectID]); rpl != nil {
			// 把对象中的所有字段写入节点
			for k, v := range rpl {
				pointer[k] = v
			}
		}
	}
}

// findObjectWithKey 查找带有指定 key 的对象，root 可以是 Slice 或者 map
// 查找到一个符合条件的对象之后立即返回
func findObjectWithKey(root interface{}, key string) types.M {
	if root == nil {
		return nil
	}
	// 如果是 Slice 则遍历查找
	if s := utils.A(root); s != nil {
		for _, v := range s {
			answer := findObjectWithKey(v, key)
			if answer != nil {
				return answer
			}
		}
	}

	if m := utils.M(root); m != nil {
		// 当前 map 中存在指定的 key，表示已经找到，立即返回
		if m[key] != nil {
			return m
		}
		// 不存在指定 key 时，则遍历 map 中各对象进行查找
		for _, v := range m {
			answer := findObjectWithKey(v, key)
			if answer != nil {
				return answer
			}
		}
	}
	return nil
}

// transformSelect 转换对象中的 $select
func transformSelect(selectObject types.M, key string, objects []types.M) {
	if selectObject == nil || selectObject["$select"] == nil {
		return
	}
	values := types.S{}
	for _, result := range objects {
		if result == nil || result[key] == nil {
			continue
		}
		values = append(values, result[key])
	}

	delete(selectObject, "$select")
	var in types.S
	if v := utils.A(selectObject["$in"]); v != nil {
		in = v
		in = append(in, values...)
	} else {
		in = values
	}
	selectObject["$in"] = in
}

// transformDontSelect 转换对象中的 $dontSelect
func transformDontSelect(dontSelectObject types.M, key string, objects []types.M) {
	if dontSelectObject == nil || dontSelectObject["$dontSelect"] == nil {
		return
	}
	values := types.S{}
	for _, result := range objects {
		if result == nil || result[key] == nil {
			continue
		}
		values = append(values, result[key])
	}

	delete(dontSelectObject, "$dontSelect")
	var nin types.S
	if v := utils.A(dontSelectObject["$nin"]); v != nil {
		nin = v
		nin = append(nin, values...)
	} else {
		nin = values
	}
	dontSelectObject["$nin"] = nin
}

// transformInQuery 转换对象中的 $inQuery
func transformInQuery(inQueryObject types.M, className string, results []types.M) {
	if inQueryObject == nil || inQueryObject["$inQuery"] == nil {
		return
	}
	values := types.S{}
	for _, result := range results {
		if result == nil || utils.S(result["objectId"]) == "" {
			continue
		}
		o := types.M{
			"__type":    "Pointer",
			"className": className,
			"objectId":  result["objectId"],
		}
		values = append(values, o)
	}

	delete(inQueryObject, "$inQuery")
	var in types.S
	if v := utils.A(inQueryObject["$in"]); v != nil {
		in = v
		in = append(in, values...)
	} else {
		in = values
	}
	inQueryObject["$in"] = in
}

// transformNotInQuery 转换对象中的 $notInQuery
func transformNotInQuery(notInQueryObject types.M, className string, results []types.M) {
	if notInQueryObject == nil || notInQueryObject["$notInQuery"] == nil {
		return
	}
	values := types.S{}
	for _, result := range results {
		if result == nil || utils.S(result["objectId"]) == "" {
			continue
		}
		o := types.M{
			"__type":    "Pointer",
			"className": className,
			"objectId":  result["objectId"],
		}
		values = append(values, o)
	}

	delete(notInQueryObject, "$notInQuery")
	var nin types.S
	if v := utils.A(notInQueryObject["$nin"]); v != nil {
		nin = v
		nin = append(nin, values...)
	} else {
		nin = values
	}
	notInQueryObject["$nin"] = nin
}
