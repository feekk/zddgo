package validator

import(
	//"fmt"
	"reflect"
	"strconv"
	"unicode/utf8"
)
func AddValidator(name string, f ValidFunction) bool{
	if _, exists := ruleFuncMap[name]; exists{
		return false
	}
	ruleFuncMap[name] = f
	return true
}

type ValidFunction func(o interface{}, topStruct reflect.Value, value reflect.Value, param string) bool

var ruleFuncMap = map[string]ValidFunction{
	"required"	:	valueRequierd,
	"lte"		:	numberIsLte,
	"gte"		:	numberIsGte,
	"eq"		:	valueIsEq,
	"ne"		:   valueIsNe,
	"min"		:	lengthMin,
	"max"		:	lengthMax,
	"len"		:   lengthEq,
	"eqfield"	:	fieldEq,
	"email"		:	isEmail,
	"numeric"	:   isNumeric,
	"number"	:	isNumber,
	"lat"		:   isLatitude,
	"lon"		:   isLongitude,
	"phone"		:	isPhone,
	"chinese"	:	isChinese,
}



func valueRequierd(o interface{}, topStruct reflect.Value, value reflect.Value, param string) bool {
	if o == nil {
		return false
	}
	return true
}

func numberIsLte(o interface{}, topStruct reflect.Value, value reflect.Value, param string) bool{ 

	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := strconv.ParseInt(param, 0, 64)
		if err != nil {
			return false
		}
		return value.Int() <= p
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := strconv.ParseUint(param, 0, 64)
		if err != nil {
			return false
		}
		return value.Uint() <= p

	case reflect.Float32, reflect.Float64:
		p, err := strconv.ParseFloat(param, 64)
		if err != nil {
			return false
		}
		return value.Float() <= p

	default:
		return false
	}
}

func numberIsGte(o interface{}, topStruct reflect.Value, value reflect.Value, param string) bool{ 

	switch value.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := strconv.ParseInt(param, 0, 64)
		if err != nil {
			return false
		}
		return value.Int() >= p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := strconv.ParseUint(param, 0, 64)
		if err != nil {
			return false
		}
		return value.Uint() >= p

	case reflect.Float32, reflect.Float64:
		p, err := strconv.ParseFloat(param, 64)
		if err != nil {
			return false
		}
		return value.Float() >= p

	default:
		return false
	}
}

func valueIsEq(o interface{}, topStruct reflect.Value, value reflect.Value, param string) bool{

	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := strconv.ParseInt(param, 0, 64)
		if err != nil {
			return false
		}
		return value.Int() == p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := strconv.ParseUint(param, 0, 64)
		if err != nil {
			return false
		}
		return value.Uint() == p

	case reflect.Float32, reflect.Float64:
		p, err := strconv.ParseFloat(param, 64)
		if err != nil {
			return false
		}
		return value.Float() == p

	case reflect.String:
		return value.String() == param

	default:
		return false
	}
}
func valueIsNe(o interface{}, topStruct reflect.Value, value reflect.Value, param string) bool{

	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := strconv.ParseInt(param, 0, 64)
		if err != nil {
			return false
		}
		return value.Int() != p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := strconv.ParseUint(param, 0, 64)
		if err != nil {
			return false
		}
		return value.Uint() != p

	case reflect.Float32, reflect.Float64:
		p, err := strconv.ParseFloat(param, 64)
		if err != nil {
			return false
		}
		return value.Float() != p

	case reflect.String:
		return value.String() != param

	default:
		return false
	}
}

func fieldEq(o interface{}, topValue reflect.Value, value reflect.Value, param string) bool{
	sValue := topValue.FieldByName(param)
	if !sValue.IsValid() || (sValue.Kind() != value.Kind()){
		return false
	}
	return sValue.Interface() == value.Interface()
}

//for length
func lengthMin(o interface{}, topValue reflect.Value, value reflect.Value, param string) bool{
	p, err := strconv.ParseInt(param, 0, 64)
	if err != nil {
		return false
	
	}
	switch value.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int64(utf8.RuneCountInString(strconv.FormatInt(value.Int(), 10))) <= p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return int64(utf8.RuneCountInString(strconv.FormatUint(value.Uint(), 10))) <= p

	case reflect.Float32, reflect.Float64:
		return int64(utf8.RuneCountInString(strconv.FormatFloat(value.Float(), 'f', -1, 64))) <= p

	case reflect.String:
		return int64(utf8.RuneCountInString(value.String())) <= p

	case reflect.Slice, reflect.Map, reflect.Array, reflect.Chan:
		return int64(value.Len()) <= p

	default:
		return false
	}
}
func lengthMax(o interface{}, topValue reflect.Value, value reflect.Value, param string) bool{
	p, err := strconv.ParseInt(param, 0, 64)
	if err != nil {
		return false
	
	}
	switch value.Kind() {
		
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int64(utf8.RuneCountInString(strconv.FormatInt(value.Int(), 10))) >= p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return int64(utf8.RuneCountInString(strconv.FormatUint(value.Uint(), 10))) >= p

	case reflect.Float32, reflect.Float64:
		return int64(utf8.RuneCountInString(strconv.FormatFloat(value.Float(), 'f', -1, 64))) >= p

	case reflect.String:
		return int64(utf8.RuneCountInString(value.String())) >= p

	case reflect.Slice, reflect.Map, reflect.Array, reflect.Chan:
		return int64(value.Len()) >= p

	default:
		return false
	}
}
func lengthEq(o interface{}, topValue reflect.Value, value reflect.Value, param string) bool{
	p, err := strconv.ParseInt(param, 0, 64)
	if err != nil {
		return false
	
	}
	switch value.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int64(utf8.RuneCountInString(strconv.FormatInt(value.Int(), 10))) == p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return int64(utf8.RuneCountInString(strconv.FormatUint(value.Uint(), 10))) == p

	case reflect.Float32, reflect.Float64:
		return int64(utf8.RuneCountInString(strconv.FormatFloat(value.Float(), 'f', -1, 64))) == p

	case reflect.String:
		return int64(utf8.RuneCountInString(value.String())) == p

	case reflect.Slice, reflect.Map, reflect.Array, reflect.Chan:
		return int64(value.Len()) == p

	default:
		return false
	}
}
func isEmail(o interface{}, topValue reflect.Value, value reflect.Value, param string) bool{
	switch value.Kind() {

		case reflect.String:
			return emailRegex.MatchString(value.String())

		default:
			return false
	}
}
func isNumeric(o interface{}, topValue reflect.Value, value reflect.Value, param string) bool{
	switch value.Kind() {

		case reflect.String:
			return numericRegex.MatchString(value.String())

		default:
			return false
	}
}
func isNumber(o interface{}, topValue reflect.Value, value reflect.Value, param string) bool{
	switch value.Kind() {

		case reflect.String:
			return numberRegex.MatchString(value.String())

		default:
			return false
	}
}

func isLatitude(o interface{}, topValue reflect.Value, value reflect.Value, param string) bool{
	switch value.Kind() {

		case reflect.String:
			return latitudeRegex.MatchString(value.String())

		default:
			return false
	}
}
func isLongitude(o interface{}, topValue reflect.Value, value reflect.Value, param string) bool{
	switch value.Kind() {

		case reflect.String:
			return longitudeRegex.MatchString(value.String())

		default:
			return false
	}
}
func isPhone(o interface{}, topStruct reflect.Value, value reflect.Value, param string) bool{
	switch value.Kind() {

	case reflect.String:
		return phoneRegex.MatchString(value.String())

	default:
		return false
	}

}

func isChinese(o interface{}, topStruct reflect.Value, value reflect.Value, param string) bool{
	switch value.Kind() {

	case reflect.String:
		return chineseRegex.MatchString(value.String())

	default:
		return false
	}
}
