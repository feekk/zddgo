package validator

import(
	"testing"
	"net/url"
)

type testQuery struct{
	//eqfield=Name
	Id int64 `query:"id" validator:"required,lte=2|gte=10,eq=10"`
	Name string `query:"name" validator:"required,chinese"`
	Fee float64 `query:"fee" validator:"required,eq=100"`
	Fee1 float64 `query:"fee1" validator:"eqfield=Fee"`
}
func(t testQuery) Message() map[string]map[string]string {
	return map[string]map[string]string{
		"id":map[string]string{
			"required":"id 不能为空",
			"lte":"id 必须小于等于1",
			"gte":"id 必须大于等于10",
			"eq" :"id 不等于10",
		},
		"name":map[string]string{
			"required":"name 不能为空",
			"eq" : "name 值不相等",
			"chinese": "name 必须为中文",
		},
		"fee":map[string]string{
			"required":"fee 不能为空",
			"eq":"fee 值不相等",
		},
		"fee1":map[string]string{
			"eqfield":"fee1 值与fee不相等",
		},
	}
}

func TestQuery(t *testing.T){
	var obj testQuery

	u, _ := url.Parse("https://127.0.0.1?id=10&name=钟&fee=100")
	values := u.Query()

	rq, err := mapFormByTag(&obj, values, "query")
	t.Logf("mapFormByTag resp:%+v\n", obj)
	_, ok := err.(VaildatorErrors)
	if err != nil && !ok {
		t.Errorf("TestQuery mapFormByTag found err: %+v\n", err)
	}
	if err = vaildator(rq, &obj, "query", err); err != nil {
		t.Logf("vaildator resp:%+v\n", err.(VaildatorErrors).ToResponse())
	}
}