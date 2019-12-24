package zddgo

import(
	"net/http"
)
func New() (z *Zddgo){
	z = &Zddgo{}
	return
}
type Zddgo struct{}
func(z *Zddgo) InitConfig() (err error) {
	err = ConfigInit()
	return
}
func(z *Zddgo) InitOrm() (err error) {
	err = MysqlInit(&Conf.Database)
	return
}
func(z *Zddgo) InitCache() (err error) {
	err = RedisInit(&Conf.Redis)
	return
}
func(z *Zddgo) HttpStart(handler http.Handler) (err error) {
	svr := NewHttpSvr(&Conf.Http, handler)
	err = svr.ListenAndServe()
	return
}