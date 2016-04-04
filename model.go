package main

//请求数据实体类
type PostData struct {
	Host     string         //数据库服务器地址
	DBName   string         //数据库名称
	Business string         //数据表名称
	Data     interface{}
	OrderBy  string         //排序
	Limit    int            //取值数量
	Skip     int            //skip数量
	Select   map[string]int //返回值列
	LockKey  string         //防重复锁定键
	Key  string         	//执行更新所制定的主键字段
	Update   bool           //新增时是否判断需要更新
	Distinct string         //是否去重，返回一个[]string
}

//返回值实体类
type ResponseModel struct {
	State   int         //状态值
	Message string      //状态信息
	Data    interface{} //返回值数据
}




