package validator
//目前只支持一维struct

import(
		//"fmt"
	"reflect"
	"strings"

	"github.com/feekk/zddgo/errors"
)

var (
	Query    = queryBinding{}
	Form     = formBinding{}
	FormPost = formPostBinding{}
	FormMultipart = formMultipartBinding{}
	Uri = uriBinding{}
	Json = jsonBinding{}
)



const (
	bindTag = "validator"
	bindSeparator = ","
	orSeparator = "|"
	tagKeySeparator = "="

	validErr = "field:%s validator:%s is failed."
)
var (
	
)

//
func vaildator(rq map[string]interface{}, ptr interface{}, tag string, mapZerrs error) error {
	var zerrs VaildatorErrors = make(VaildatorErrors)
	if mapZerrs != nil {
		zerrs = mapZerrs.(VaildatorErrors)
	}

	var message map[string]map[string]string = make(map[string]map[string]string)
	value := reflect.ValueOf(ptr)
	value = elem(value)
	//获取自定义错误信息
	method := value.MethodByName("Message")
	if method.IsValid(){
		r := method.Call(make([]reflect.Value, 0))
		message = r[0].Interface().(map[string]map[string]string)
	}

	tValue := value.Type()
	var fieldStruct reflect.StructField
	var ruleStr, msg string
	var rules, rs, vals []string
	var exists, isPass bool
	var vfunc ValidFunction 
	var fMsg map[string]string
	//field
	for i := 0; i < value.NumField(); i++ {
		fieldStruct = tValue.Field(i)
		
		ruleStr = fieldStruct.Tag.Get(bindTag)
		//拆分&&条件
		rules = strings.Split(ruleStr, bindSeparator)

		//validator
		for j := 0; j < len(rules); j++ {
			//&&条件下只要有一个error就直接跳过本字段
			if _, exists = zerrs[fieldStruct.Tag.Get(tag)]; exists{
				break
			}

			//拆分||条件
			rs = strings.Split(rules[j], orSeparator)

			//vfunc
			for k := 0; k < len(rs); k++ {
				//&&条件下，只要有一个nil，则本字段通过

				//拆分值
				vals = strings.Split(rs[k], tagKeySeparator)
				
				vfunc, exists = ruleFuncMap[vals[0]]
				if !exists || vfunc == nil{
					zerrs[fieldStruct.Tag.Get(tag)] = errors.Errorf("field:%s validator:%s not found.", fieldStruct.Tag.Get(tag), vals[0])
					break
				}
				if len(vals) == 1 {
					vals = append(vals, EMPTYSTR)
				}

				if isPass = vfunc(rq[fieldStruct.Tag.Get(tag)], elem(value), elem(value.Field(i)), vals[1]); isPass{
					//如果有，删除错误信息
					delete(zerrs, fieldStruct.Tag.Get(tag))
					break
				}else{
					fMsg, exists = message[fieldStruct.Tag.Get(tag)]
					if !exists{
						zerrs[fieldStruct.Tag.Get(tag)] = errors.Errorf(validErr, fieldStruct.Tag.Get(tag), vals[0])
						continue
					}
					msg, exists = fMsg[vals[0]]
					if !exists{
						zerrs[fieldStruct.Tag.Get(tag)] = errors.Errorf(validErr, fieldStruct.Tag.Get(tag), vals[0])
						continue
					}
					zerrs[fieldStruct.Tag.Get(tag)] = errors.New(msg)
				}
			}
		}
	}
	if len(zerrs) > 0 {
		return zerrs
	}
	return nil
}

func elem(value reflect.Value) reflect.Value{
	if value.Kind() == reflect.Ptr{
		return value.Elem()
	}
	return value
}