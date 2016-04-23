package rest

import (
	"fmt"
	"sort"
	"strings"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
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
}

// NewQuery 组装查询对象
func NewQuery(
	auth *Auth,
	className string,
	where types.M,
	options types.M,
) (*Query, error) {
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
	}

	if auth.IsMaster == false {
		// 当前权限为 Master 时，findOptions 中不存在 acl 这个 key
		if auth.User != nil {
			query.findOptions["acl"] = []string{auth.User["objectId"].(string)}
		} else {
			query.findOptions["acl"] = nil
		}
		if className == "_Session" {
			if query.findOptions["acl"] == nil {
				return nil, errs.E(errs.InvalidSessionToken, "This session token is invalid.")
			}
			user := types.M{"__type": "Pointer", "className": "_User", "objectId": auth.User["objectId"]}
			and := types.S{where, user}
			query.Where = types.M{"$and": and}
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

	fmt.Println("keys       ", q.keys)
	fmt.Println("doCount    ", q.doCount)
	fmt.Println("findOptions", q.findOptions)
	fmt.Println("include    ", q.include)

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

// BuildRestWhere ...
func (q *Query) BuildRestWhere() error {
	q.getUserAndRoleACL()
	q.redirectClassNameForKey()
	q.validateClientClassCreation()
	q.replaceSelect()
	q.replaceDontSelect()
	q.replaceInQuery()
	q.replaceNotInQuery()
	return nil
}

func (q *Query) getUserAndRoleACL() error {
	if q.auth.IsMaster || q.auth.User == nil {
		return nil
	}
	roles := q.auth.GetUserRoles()
	roles = append(roles, q.auth.User["objectId"].(string))
	q.findOptions["acl"] = roles
	return nil
}

func (q *Query) redirectClassNameForKey() error {
	if q.redirectKey == "" {
		return nil
	}

	newClassName := orm.RedirectClassNameForKey(q.className, q.redirectKey)
	q.className = newClassName
	q.redirectClassName = newClassName

	return nil
}

func (q *Query) validateClientClassCreation() error {
	sysClass := []string{"_User", "_Installation", "_Role", "_Session", "_Product"}
	if config.TConfig.AllowClientClassCreation {
		return nil
	}
	if q.auth.IsMaster {
		return nil
	}
	for _, v := range sysClass {
		if v == q.className {
			return nil
		}
	}
	if orm.CollectionExists(q.className) {
		return nil
	}
	// TODO 无法操作不存在的表
	return nil
}

func (q *Query) replaceSelect() error {
	selectObject := findObjectWithKey(q.Where, "$select")
	if selectObject == nil {
		return nil
	}
	// 必须包含两个 key ： query key
	selectValue := utils.MapInterface(selectObject["$select"])
	if selectValue == nil ||
		selectValue["query"] == nil ||
		selectValue["key"] == nil {
		// TODO $select 用法不正确
		return nil
	}
	queryValue := utils.MapInterface(selectValue["query"])
	if queryValue == nil ||
		queryValue["className"] == nil ||
		queryValue["where"] == nil {
		// TODO $select 用法不正确
		return nil
	}

	values := types.S{}

	query, err := NewQuery(
		q.auth,
		utils.String(queryValue["className"]),
		utils.MapInterface(queryValue["where"]),
		types.M{})
	if err != nil {
		return err
	}
	response, err := query.Execute()
	if err != nil {
		return err
	}
	// 组装查询到的对象
	if utils.HasResults(response) == true {
		for _, v := range utils.SliceInterface(response["results"]) {
			result := utils.MapInterface(v)
			key := result[utils.String(selectValue["key"])]
			if key != nil {
				values = append(values, key)
			}
		}
	}
	// 替换 $select 为 $in
	delete(selectObject, "$select")
	if selectObject["$in"] != nil &&
		utils.SliceInterface(selectObject["$in"]) != nil {
		in := utils.SliceInterface(selectObject["$in"])
		selectObject["$in"] = append(in, values...)
	} else {
		selectObject["$in"] = values
	}
	// 继续搜索替换
	return q.replaceSelect()
}

func (q *Query) replaceDontSelect() error {
	dontSelectObject := findObjectWithKey(q.Where, "$dontSelect")
	if dontSelectObject == nil {
		return nil
	}
	// 必须包含两个 key ： query key
	dontSelectValue := utils.MapInterface(dontSelectObject["$dontSelect"])
	if dontSelectValue == nil ||
		dontSelectValue["query"] == nil ||
		dontSelectValue["key"] == nil {
		// TODO $dontSelect 用法不正确
		return nil
	}
	queryValue := utils.MapInterface(dontSelectValue["query"])
	if queryValue == nil ||
		queryValue["className"] == nil ||
		queryValue["where"] == nil {
		// TODO $dontSelect 用法不正确
		return nil
	}

	values := types.S{}

	query, err := NewQuery(
		q.auth,
		utils.String(queryValue["className"]),
		utils.MapInterface(queryValue["where"]),
		types.M{})
	if err != nil {
		return err
	}
	response, err := query.Execute()
	if err != nil {
		return err
	}
	// 组装查询到的对象
	if utils.HasResults(response) == true {
		for _, v := range utils.SliceInterface(response["results"]) {
			result := utils.MapInterface(v)
			key := result[utils.String(dontSelectValue["key"])]
			if key != nil {
				values = append(values, key)
			}
		}
	}
	// 替换 $dontSelect 为 $nin
	delete(dontSelectObject, "$dontSelect")
	if dontSelectObject["$nin"] != nil &&
		utils.SliceInterface(dontSelectObject["$nin"]) != nil {
		nin := utils.SliceInterface(dontSelectObject["$nin"])
		dontSelectObject["$nin"] = append(nin, values...)
	} else {
		dontSelectObject["$nin"] = values
	}
	// 继续搜索替换
	return q.replaceDontSelect()
}

func (q *Query) replaceInQuery() error {
	inQueryObject := findObjectWithKey(q.Where, "$inQuery")
	if inQueryObject == nil {
		return nil
	}
	// 必须包含两个 key ： where className
	inQueryValue := utils.MapInterface(inQueryObject["$inQuery"])
	if inQueryValue == nil ||
		inQueryValue["where"] == nil ||
		inQueryValue["className"] == nil {
		// TODO $inQuery 用法不正确
		return nil
	}

	values := types.S{}

	query, err := NewQuery(
		q.auth,
		utils.String(inQueryValue["className"]),
		utils.MapInterface(inQueryValue["where"]),
		types.M{})
	if err != nil {
		return err
	}
	response, err := query.Execute()
	if err != nil {
		return err
	}
	// 组装查询到的对象
	if utils.HasResults(response) == true {
		for _, v := range utils.SliceInterface(response["results"]) {
			result := utils.MapInterface(v)
			pointer := types.M{
				"__type":    "Pointer",
				"className": inQueryValue["className"],
				"objectId":  result["objectId"],
			}
			values = append(values, pointer)
		}
	}
	// 替换 $inQuery 为 $in
	delete(inQueryObject, "$inQuery")
	if inQueryObject["$in"] != nil &&
		utils.SliceInterface(inQueryObject["$in"]) != nil {
		in := utils.SliceInterface(inQueryObject["$in"])
		inQueryObject["$in"] = append(in, values...)
	} else {
		inQueryObject["$in"] = values
	}
	// 继续搜索替换
	return q.replaceInQuery()
}

func (q *Query) replaceNotInQuery() error {
	notInQueryObject := findObjectWithKey(q.Where, "$notInQuery")
	if notInQueryObject == nil {
		return nil
	}
	// 必须包含两个 key ： where className
	notInQueryValue := utils.MapInterface(notInQueryObject["$notInQuery"])
	if notInQueryValue == nil ||
		notInQueryValue["where"] == nil ||
		notInQueryValue["className"] == nil {
		// TODO $notInQuery 用法不正确
		return nil
	}

	values := types.S{}

	query, err := NewQuery(
		q.auth,
		utils.String(notInQueryValue["className"]),
		utils.MapInterface(notInQueryValue["where"]),
		types.M{})
	if err != nil {
		return err
	}
	response, err := query.Execute()
	if err != nil {
		return err
	}
	// 组装查询到的对象
	if utils.HasResults(response) == true {
		for _, v := range utils.SliceInterface(response["results"]) {
			result := utils.MapInterface(v)
			pointer := types.M{
				"__type":    "Pointer",
				"className": notInQueryValue["className"],
				"objectId":  result["objectId"],
			}
			values = append(values, pointer)
		}
	}
	// 替换 $notInQuery 为 $nin
	delete(notInQueryObject, "$notInQuery")
	if notInQueryObject["$nin"] != nil &&
		utils.SliceInterface(notInQueryObject["$nin"]) != nil {
		nin := utils.SliceInterface(notInQueryObject["$nin"])
		notInQueryObject["$nin"] = append(nin, values...)
	} else {
		notInQueryObject["$nin"] = values
	}
	// 继续搜索替换
	return q.replaceNotInQuery()
}

func (q *Query) runFind() error {
	response := orm.Find(q.className, q.Where, q.findOptions)
	if q.className == "_User" {
		for _, v := range response {
			user := utils.MapInterface(v)
			if user != nil {
				delete(user, "password")
			}
		}
	}
	// TODO 取出需要的 key   （TODO：通过数据库直接取key）
	results := types.S{}
	if len(q.keys) > 0 && len(response) > 0 {
		for _, v := range response {
			obj := utils.MapInterface(v)
			newObj := types.M{}
			for _, s := range q.keys {
				if obj[s] != nil {
					newObj[s] = obj[s]
				}
			}
			results = append(results, newObj)
		}
	}

	// TODO 展开文件类型

	q.response["results"] = results
	return nil
}

func (q *Query) runCount() error {
	if q.doCount == false {
		return nil
	}
	q.findOptions["count"] = true
	delete(q.findOptions, "skip")
	delete(q.findOptions, "limit")
	q.response["count"] = orm.Find(q.className, q.Where, q.findOptions)[0]
	return nil
}

func (q *Query) handleInclude() error {
	if len(q.include) == 0 {
		return nil
	}
	includePath(q.auth, q.response, q.include[0])

	if len(q.include) > 0 {
		q.include = q.include[1:]
		return q.handleInclude()
	}

	return nil
}

func includePath(auth *Auth, response types.M, path []string) error {
	pointers := findPointers(response["results"], path)
	if len(pointers) == 0 {
		return nil
	}
	className := ""
	objectIDs := []string{}
	for _, v := range pointers {
		pointer := utils.MapInterface(v)
		if className == "" {
			className = utils.String(pointer["className"])
		} else {
			if className != utils.String(pointer["className"]) {
				// TODO 对象类型不一致
				return nil
			}
		}
		objectIDs = append(objectIDs, utils.String(pointer["objectId"]))
	}
	if className == "" {
		// TODO 无效对象
		return nil
	}

	objectID := types.M{
		"$in": objectIDs,
	}
	where := types.M{
		"objectId": objectID,
	}
	query, err := NewQuery(auth, className, where, types.M{})
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
	results := utils.SliceInterface(includeResponse["results"])
	replace := types.M{}
	for _, v := range results {
		obj := utils.MapInterface(v)
		if className == "_User" {
			delete(obj, "sessionToken")
		}
		replace[utils.String(obj["objectId"])] = obj
	}

	replacePointers(pointers, replace)

	return nil
}

// 查询路径对应的对象列表
func findPointers(object interface{}, path []string) types.S {
	// 如果是对象数组，则遍历每一个对象
	if utils.SliceInterface(object) != nil {
		answer := types.S{}
		for _, v := range utils.SliceInterface(object) {
			answer = append(answer, findPointers(v, path)...)
		}
		return answer
	}

	obj := utils.MapInterface(object)
	if obj == nil {
		return types.S{}
	}
	// 如果当前是路径最后一个节点，判断是否为 Pointer
	if len(path) == 0 {
		if obj["__type"] == "Pointer" {
			return types.S{obj}
		}
		return types.S{}
	}
	// 取出下一个路径对应的对象，进行查找
	subobject := obj[path[0]]
	if subobject == nil {
		return types.S{}
	}
	return findPointers(subobject, path[1:])
}

func replacePointers(pointers types.S, replace types.M) error {
	for _, v := range pointers {
		pointer := utils.MapInterface(v)
		objectID := utils.String(pointer["objectId"])
		if replace[objectID] == nil {
			continue
		}
		rpl := utils.MapInterface(replace[objectID])
		for k, v := range rpl {
			pointer[k] = v
		}
		pointer["__type"] = "Object"
	}
	return nil
}

// 查找带有指定 key 的对象
func findObjectWithKey(root interface{}, key string) types.M {
	if s := utils.SliceInterface(root); s != nil {
		for _, v := range s {
			answer := findObjectWithKey(v, key)
			if answer != nil {
				return answer
			}
		}
	}
	if m := utils.MapInterface(root); m != nil {
		if m[key] != nil {
			return m
		}
		for _, v := range m {
			answer := findObjectWithKey(v, key)
			if answer != nil {
				return answer
			}
		}
	}
	return nil
}
