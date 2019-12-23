package validator

import (
	"testing"
	"time"
)

func TestMapNil(t *testing.T){
	var st testConvert_1

	if _, err := mapNil(&st); err != nil {
		t.Errorf("TestMapNil found err: %+v\n", err)
	}
}


type testConvert_1 struct{
	Id int64 `form:"id"`
	Uid int64 `form:"uid" default:"9"`
	Pid *int64 `form:"pid"`
	Did *int64 `form:"did" default:"99"`
	Cid *int64 `form:"cid"`
}
func TestMapByTag_1(t *testing.T) {
	data := make(map[string][]string, 0)
	data["id"] = []string{"4"}
	data["pid"] = []string{"44"}
	
	var st testConvert_1

	rq, err := mapFormByTag(&st, data, "form")
	if err != nil {
		t.Errorf("mapFormByTag found err: %+v\n", err)
	}
	if st.Id != 4 {
		t.Errorf("Id not match: %+v\n", st.Id)
	}
	if st.Uid != 9 {
		t.Errorf("Uid not match: %+v\n", st.Uid)
	}
	if *st.Pid != 44 {
		t.Errorf("Pid not match: %+v\n", st.Pid)
	}
	if *st.Did != 99 {
		t.Errorf("Did not match: %+v\n", st.Did)
	}
	if st.Cid != nil {
		t.Errorf("Cid not match: %+v\n", st.Cid)
	}
	if rq["id"] == nil {
		t.Errorf("rq id not match: %+v\n", rq["id"])
	}
	if rq["uid"] == nil {
		t.Errorf("rq uid not match: %+v\n", rq["uid"])
	}
	if rq["pid"] == nil {
		t.Errorf("rq pid not match: %+v\n", rq["pid"])
	}
	if rq["did"] == nil {
		t.Errorf("rq did not match: %+v\n", rq["did"])
	}
	if rq["cid"] != nil {
		t.Errorf("rq cid not match: %+v\n", rq["cid"])
	}
}

type testConvert_2 struct{
	Time time.Time `form:"time"  time_format:"2006-01-02 15:04:05"`
	Utime time.Time `form:"utime"  time_format:"2006-01-02 15:04:05" default:"2019-09-09 09:09:09"`
	Ptime *time.Time `form:"ptime" time_format:"2006-01-02 15:04:05"`
	Dtime *time.Time `form:"dtime"  time_format:"2006-01-02 15:04:05" default:"2019-10-10 19:19:19"`
	Ctime *time.Time `form:"ctime"  time_format:"2006-01-02 15:04:05"`
}
func TestMapByTag_2(t *testing.T) {
	data := make(map[string][]string, 0)
	data["time"] = []string{"2018-08-08 08:08:08"}
	data["ptime"] = []string{"2018-08-18 18:18:18"}
	
	var st testConvert_2

	rq, err := mapFormByTag(&st, data, "form")
	if err != nil {
		t.Errorf("mapFormByTag found err: %+v\n", err)
	}
	if st.Time.Format("2006-01-02 15:04:05") != "2018-08-08 08:08:08" {
		t.Errorf("Time not match: %+v\n", st.Time.Format("2006-01-02 15:04:05"))
	}
	if st.Utime.Format("2006-01-02 15:04:05") != "2019-09-09 09:09:09" {
		t.Errorf("Utime not match: %+v\n", st.Utime.Format("2006-01-02 15:04:05"))
	}
	if st.Ptime.Format("2006-01-02 15:04:05") != "2018-08-18 18:18:18" {
		t.Errorf("Ptime not match: %+v\n", st.Ptime.Format("2006-01-02 15:04:05"))
	}
	if st.Dtime.Format("2006-01-02 15:04:05") != "2019-10-10 19:19:19" {
		t.Errorf("Dtime not match: %+v\n", st.Dtime.Format("2006-01-02 15:04:05"))
	}
	if rq["time"] == nil {
		t.Errorf("rq time not match: %+v\n", rq["time"])
	}
	if rq["utime"] == nil {
		t.Errorf("rq utime not match: %+v\n", rq["utime"])
	}
	if rq["ptime"] == nil {
		t.Errorf("rq ptime not match: %+v\n", rq["ptime"])
	}
	if rq["dtime"] == nil {
		t.Errorf("rq dtime not match: %+v\n", rq["dtime"])
	}
	if rq["ctime"] != nil {
		t.Errorf("rq ctime not match: %+v\n", rq["ctime"])
	}

}
