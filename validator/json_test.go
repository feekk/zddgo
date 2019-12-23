package validator

import(
	"testing"
	//"time"
	"github.com/feekk/zddgo/ztime"
)

type testJsonStruct struct {
	Id   *int64     `json:"id" binding:"required" default:"9"`
	Time *ztime.JsonTime `json:"time" binding:"required" time_format:"2006-01-02 15:04:05"`
}

func TestJson(t *testing.T){

	var data = `{"id":2,"time":"2019-11-02 18:22:44"}`

	var j testJsonStruct

	if err := Json.BindBody([]byte(data), &j); err != nil {
		t.Errorf("err:%+v\n", err)
	}
	tStr, err := j.Time.MarshalJSON()
	if err != nil{
		t.Errorf("time err:%+v\n", err)
	}
	if string(tStr) != "2019-11-02 18:22:44" {
		t.Errorf("time value err: %T %+v\n", string(tStr), string(tStr))
	}
}
