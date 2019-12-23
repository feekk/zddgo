package validator

import (
	//"fmt"
	"reflect"
	"strconv"
	"time"
	"encoding/json"

	"github.com/feekk/zddgo/errors"
)
const(
	NIL= "*"
	IGNORE = "-"
	EMPTYSTR = ""
	DEFAULT_TAG = "default" 
)
var (
	NeedPtrErr       = errors.New("Input ptr must be reflect.Ptr.")
	PtrNeedCanSetErr = errors.New("Input ptr Can't Set.")
	UnKnowTypeErr    = errors.New("Unknown type")
)

func mapNil(ptr interface{}) (map[string]interface{}, error){
	nMap := make(map[string][]string, 0)
	return mapSetterByTag(ptr, nilSource(nMap), NIL)
}

func mapFormByTag(ptr interface{}, form map[string][]string, tag string) (map[string]interface{}, error) {
	return mapSetterByTag(ptr, formSource(form), tag)
}

func mapSetterByTag(ptr interface{}, setter setter, tag string) (map[string]interface{}, error) {
	value := reflect.ValueOf(ptr)
	if value.Kind() != reflect.Ptr {
		return nil, NeedPtrErr
	}
	if !value.Elem().CanSet() {
		return nil, PtrNeedCanSetErr
	}

	return mapping(value.Elem(), setter, tag)
}

//setter tries to set value on a walking by fields of a struct
type setter interface {
	TrySet(value reflect.Value, field reflect.StructField, key string, opt setOptions) (interface{}, error)
}

//for default value
type setOptions struct {
	isDefaultExists bool
	defaultValue    string
}

var _ setter = formSource(nil)

type formSource map[string][]string

func (f formSource) TrySet(value reflect.Value, field reflect.StructField, key string, opt setOptions) (interface{}, error) {
	return trySet(value, field, f, key, opt)
}

var _ setter = nilSource(nil)

type nilSource map[string][]string
func (n nilSource) TrySet(value reflect.Value, field reflect.StructField, key string, opt setOptions) (interface{}, error){
	return nil, nil
}
//value = reflect.ValueOf(ptr)
func mapping(value reflect.Value, setter setter, tag string) (map[string]interface{}, error) {
	var err error
	var fieldStruct reflect.StructField
	var tagVal string
	var opt setOptions
	var errs VaildatorErrors = make(map[string]error, 0)
	var rq map[string]interface{} = make(map[string]interface{}, 0)

	tValue := value.Type()

	for i := 0; i < value.NumField(); i++ {
		fieldStruct = tValue.Field(i)
		opt = setOptions{}

		if tag != NIL{

			tagVal = fieldStruct.Tag.Get(tag)
			if tagVal == IGNORE {
				continue
			}
			if tagVal == EMPTYSTR {
				tagVal = fieldStruct.Name
			}
			if tagVal == EMPTYSTR {
				continue
			}
			if opt.defaultValue = fieldStruct.Tag.Get(DEFAULT_TAG); opt.defaultValue != EMPTYSTR {
				opt.isDefaultExists = true
			}
			//set value
			if rq[tagVal], err = setter.TrySet(value.Field(i), fieldStruct, tagVal, opt); err != nil {
				errs[tagVal] = err
				continue
			}
		}
	}
	if len(errs) > 0 {
		return rq, errs
	}
	return rq, nil
}

func trySet(value reflect.Value, field reflect.StructField, form map[string][]string, key string, opt setOptions) (interface{}, error)  {
	val, ok := form[key]
	if !ok && !opt.isDefaultExists {
		return nil, nil
	}
	
	if value.Kind() == reflect.Ptr {
		vPtr := reflect.New(value.Type().Elem())
		origin, err := trySet(vPtr.Elem(), field, form, key, opt)
		if err != nil {
			return origin, err
		}
		value.Set(vPtr)
		return origin, nil
	}

	err := setToField(ok, val, value, field, opt)

	return value.Interface(), err
}

func setToField(ok bool, val []string, value reflect.Value, field reflect.StructField, opt setOptions) error {
	switch value.Kind() {
	case reflect.Slice:
		if !ok {
			val = []string{opt.defaultValue}
		}
		return setSlice(val, value, field)
	case reflect.Array:
		if !ok {
			val = []string{opt.defaultValue}
		}
		if len(val) != value.Len() {
			return errors.Errorf("%q's len was not match %s", val, value.Type().String())
		}
		return setArray(val, value, field)
	default:
		var dVal string
		if !ok {
			dVal = opt.defaultValue
		}

		if len(val) > 0 {
			dVal = val[0]
		}
		return setWithProperType(dVal, value, field)
	}
	return nil
}

func setWithProperType(val string, value reflect.Value, field reflect.StructField) error {
	switch value.Kind() {
	case reflect.Int:
		return setIntField(val, 0, value)
	case reflect.Int8:
		return setIntField(val, 8, value)
	case reflect.Int16:
		return setIntField(val, 16, value)
	case reflect.Int32:
		return setIntField(val, 32, value)
	case reflect.Int64:
		switch value.Interface().(type) {
		case time.Duration:
			return setTimeDuration(val, value, field)
		}
		return setIntField(val, 64, value)
	case reflect.Uint:
		return setUintField(val, 0, value)
	case reflect.Uint8:
		return setUintField(val, 8, value)
	case reflect.Uint16:
		return setUintField(val, 16, value)
	case reflect.Uint32:
		return setUintField(val, 32, value)
	case reflect.Uint64:
		return setUintField(val, 64, value)
	case reflect.Bool:
		return setBoolField(val, value)
	case reflect.Float32:
		return setFloatField(val, 32, value)
	case reflect.Float64:
		return setFloatField(val, 64, value)
	case reflect.String:
		value.SetString(val)
	case reflect.Struct:
		switch value.Interface().(type) {
		case time.Time:
			return setTimeField(val, field, value)
		}
		return json.Unmarshal([]byte(val), value.Addr().Interface())
	case reflect.Map:
		return json.Unmarshal([]byte(val), value.Addr().Interface())
	default:
		return UnKnowTypeErr
	}
	return nil
}
func setArray(vals []string, value reflect.Value, field reflect.StructField) error {
	for i, s := range vals {
		err := setWithProperType(s, value.Index(i), field)
		if err != nil {
			return err
		}
	}
	return nil
}
func setSlice(vals []string, value reflect.Value, field reflect.StructField) error {
	slice := reflect.MakeSlice(value.Type(), len(vals), len(vals))
	err := setArray(vals, slice, field)
	if err != nil {
		return err
	}
	value.Set(slice)
	return nil
}
func setIntField(val string, bitSize int, field reflect.Value) error {
	if val == EMPTYSTR {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == EMPTYSTR {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(val string, field reflect.Value) error {
	if val == EMPTYSTR {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == EMPTYSTR {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func setTimeField(val string, structField reflect.StructField, value reflect.Value) error {
	timeFormat := structField.Tag.Get("time_format")
	if timeFormat == EMPTYSTR {
		timeFormat = time.RFC3339
	}

	if val == EMPTYSTR {
		value.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	l := time.Local
	if isUTC, _ := strconv.ParseBool(structField.Tag.Get("time_utc")); isUTC {
		l = time.UTC
	}

	if locTag := structField.Tag.Get("time_location"); locTag != EMPTYSTR {
		loc, err := time.LoadLocation(locTag)
		if err != nil {
			return err
		}
		l = loc
	}

	t, err := time.ParseInLocation(timeFormat, val, l)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(t))
	return nil
}
func setTimeDuration(val string, value reflect.Value, field reflect.StructField) error {
	d, err := time.ParseDuration(val)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(d))
	return nil
}
