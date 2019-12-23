package zddgo

import(
	"net/http"
)
//const UPDATE_FAIL = 20001
//Add(200, UPDATE_FAIL, "update fail", "update fail response")


func AddResponseMessage(status, code int, message, desc string){
	rt := responseTemp{
		status : status,
		code : code,
		message : message,
		desc : desc,
	}
	_responseTempMap[rt.code] = rt	
}

const(
	// 0 ~ 1000
	OK = 0
	FAIL = 1
	UNKNOW = 99
	PARAMERR = 100
)

var _responseTempMap responseTempMap = map[int]responseTemp{
	OK : responseTemp{http.StatusOK, OK, "success", "success response"},
	FAIL: responseTemp{http.StatusOK, FAIL, "fail", "fail response"},
	UNKNOW: responseTemp{http.StatusOK, UNKNOW, "unkow", "unkow response"},
	PARAMERR: responseTemp{http.StatusOK, PARAMERR, "param check error", "param check error"},
}

type responseTempMap map[int]responseTemp

type responseTemp struct {
	status int
	code int
	message string
	desc string
}

func(r responseTempMap) get(code int) (val responseTemp){
	var ok bool
	val, _ = r[UNKNOW]
	if val ,ok = r[code]; ok{
		return 
	}
	return 
}