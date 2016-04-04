package main

import (
	"errors"
	"fmt"
	"gopkg/mgo"
	"gopkg/mgo/bson"
	"net/http"
	"reflect"
	"sync/atomic"
	"time"
	"sync"
"os"
"io"
)

var cnt int32 = 0
var maxQueue int32 = 100

//测试界面
func index(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><title>API文档</title><style type="text/css">table{border:1px solid #dbdbdb}th{background:-webkit-gradient(linear,0 0,0 100%,from(#616161),to(#555))}th b{color:#fff}</style></head><body><h1>http请求目录</h1><table class="t2"><tbody><tr><th width="110"><b>方法</b></th><th width="200"><b>地址</b></th><th width="60"><b>操作</b></th><th width="300"><b>post body示例</b></th><th width="180"><b>备注</b></th></tr><tr><td><b>帮助</b></td><td>http://222.185.195.206:8010/</td><td>GET</td><td></td><td></td></tr><tr><td><b>查询</b></td><td>http://222.185.195.206:8010/get</td><td>POST</td><td>{"Url":"","DBName":"spider","Business":"insocial","Data":{"zans":{"$gt":1}},"Limit":1,"Skip":1,"Select":{"id":1,"inid":1,"name":1}}</td><td></td></tr><tr><td><b>统计</b></td><td>http://222.185.195.206:8010/count</td><td>POST</td><td>{"Url":"","Business":"insocial","Data":{"zans":{"$gt":1}}</td><td></td></tr><tr><td><b>判断是否存在</b></td><td>http://222.185.195.206:8010/exists</td><td>POST</td><td>{"Business":"insocial","Data":{"id":"123456"},"Key":"id"}</td><td>需指定主键key，否则默认为key为id</td></tr><tr><td><b>新增</b></td><td>http://222.185.195.206:8010/add</td><td>POST</td><td>{"Business":"insocial","Data":{"id":"123456","name":"123123"},"Key":"id"}</td><td>Data值可以为map或者[]map</td></tr><tr><td><b>更新</b></td><td>http://222.185.195.206:8010/update</td><td>POST</td><td>{"Business":"insocial","Data":{"id":"123456","name":"123123"},"Key":"id"}</td><td>需指定主键key，否则默认为key为id</td></tr></tbody></table><h1>http post 参数详情</h1><table class="t2"><tbody><tr><th width="110"><b>参数</b></th><th width="195"><b>名称</b></th><th width="195"><b>备注</b></th><th width="195"><b>必选</b></th></tr><tr><td><b>Url</b></td><td>服务器IP</td><td>缺省为空，则为Server所在服务器，或者192.168.2.192:27027指定mongodb连接字符串</td><td><img src="http://qzonestyle.gtimg.cn/qzone/vas/opensns/res/img/FusionAPI_2.JPG"></td></tr><tr><td><b>DBName</b></td><td>数据库名称</td><td>缺省为spider</td><td><img src="http://qzonestyle.gtimg.cn/qzone/vas/opensns/res/img/FusionAPI_2.JPG"></td></tr><tr><td><b>Business</b></td><td>表名</td><td></td><td><img src="http://qzonestyle.gtimg.cn/qzone/vas/opensns/res/img/API_legend_6.png"></td></tr><tr><td><b>Data</b></td><td>数据</td><td>查询条件或者插入和更新的值</td><td><img src="http://qzonestyle.gtimg.cn/qzone/vas/opensns/res/img/API_legend_6.png"></td></tr><tr><td><b>OrderBy</b></td><td>排序条件</td><td>time,+time,-time</td><td><img src="http://qzonestyle.gtimg.cn/qzone/vas/opensns/res/img/FusionAPI_2.JPG"></td></tr><tr><td><b>Limit</b></td><td>选取条数</td><td>int型</td><td><img src="http://qzonestyle.gtimg.cn/qzone/vas/opensns/res/img/FusionAPI_2.JPG"></td></tr><tr><td><b>Skip</b></td><td>忽略条数</td><td>int型</td><td><img src="http://qzonestyle.gtimg.cn/qzone/vas/opensns/res/img/FusionAPI_2.JPG"></td></tr><tr><td><b>Select</b></td><td>指定返回列</td><td>{"time":1}</td><td><img src="http://qzonestyle.gtimg.cn/qzone/vas/opensns/res/img/FusionAPI_2.JPG"></td></tr><tr><td><b>Key</b></td><td>表主键</td><td>缺省为id,在做插入或者更新的时候需要指定</td><td><img src="http://qzonestyle.gtimg.cn/qzone/vas/opensns/res/img/FusionAPI_2.JPG"></td></tr><tr><td><b>Update</b></td><td>更新标识</td><td>bool类型，插入时如果已存在是否更新</td><td><img src="http://qzonestyle.gtimg.cn/qzone/vas/opensns/res/img/FusionAPI_2.JPG"></td></tr><tr><td><b>Distinct</b></td><td>去重查询标识</td><td>string类型，根据string执行去重查询</td><td><img src="http://qzonestyle.gtimg.cn/qzone/vas/opensns/res/img/FusionAPI_2.JPG"></td></tr><tr><td><b>LockKey</b></td><td>占用的键</td><td>用于并发获取数据，获取数据时先标识为-1,需同时指定key</td><td><img src="http://qzonestyle.gtimg.cn/qzone/vas/opensns/res/img/FusionAPI_2.JPG"></td></tr></tbody></table></body></html>`
	w.Write([]byte(html))
}

/*统计数据*/
func count(w http.ResponseWriter, r *http.Request) {
	defer Recovery("count")
	if body := loadRequest(w, r); body.Business != "-1" {
		if m, o := body.Data.(map[string]interface{}); o {
			session := InitMongo(body.Host)//初始化mongodb连接
			if session == nil {
				ResponseJson(w, 0, errors.New("mongodb connected failed"))
				return
			}
			defer session.Close()
			collection := session.DB(body.DBName).C(body.Business)
			var count = 0
			var er error
			//获取关键字对应的主键集合
			query := collection.Find(m)
			query.SetMaxTime(5 * time.Minute) //5分钟
			count, er = query.Count()
			if er == nil {
				ResponseJson(w, count, nil)
			} else {
				ResponseJson(w, nil, er)
			}
		} else {
			ResponseJson(w, 0, errors.New("data type is not map[string]interface{} but a "+reflect.TypeOf(body.Data).Name()))
		}
	}
}

//判断用户是否存在 所有数据格式  均为{业务名,数据}
func exists(w http.ResponseWriter, r *http.Request) {
	defer Recovery("exists")
	if body := loadRequest(w, r); body.Business != "-1" {
		if m, o := body.Data.(map[string]interface{}); o {
			session := InitMongo(body.Host)
			if session == nil {
				ResponseJson(w, 0, errors.New("mongodb connected failed"))
				return
			}
			defer session.Close()
			collection := session.DB(body.DBName).C(body.Business)
			c, e := collection.Find(m).Count()
			if e == nil && c > 0 {
				ResponseJson(w, 1, nil)
			}else {
				ResponseJson(w, 0, errors.New("not exists"))
			}
		} else {
			ResponseJson(w, 0, errors.New("data type is not map[string]interface{} but a " + reflect.TypeOf(body.Data).Name()))
		}
	}
}

//查询
func get(w http.ResponseWriter, r *http.Request) {
	if body := loadRequest(w, r); body.Business != "-1" {
		if m, o := body.Data.(map[string]interface{}); o {
			session := InitMongo(body.Host)//初始化mongodb连接
			if session == nil {
				ResponseJson(w, 0, errors.New("mongodb connected failed"))
				return
			}
			defer session.Close()
			collection := session.DB(body.DBName).C(body.Business)
			var users []interface{}
			var query *mgo.Query
			query = collection.Find(m)
			var err error
			if body.OrderBy != "" {
				query = query.Sort(body.OrderBy)
			}
			if body.Skip != 0 {
				query = query.Skip(body.Skip)
			}
			if body.Limit != 0 {
				query = query.Limit(body.Limit)
			}else {
				query = query.Limit(50)//默认返回50条数据
			}
			if body.Select != nil {
				query = query.Select(body.Select)
			}
			query.SetMaxTime(5 * time.Minute) //5分钟
			//如果是Distinct命令，则
			if body.Distinct != "" {
				var res []string
				err = query.Distinct(body.Distinct, &res)
				if err == nil {
					ResponseJson(w, res, nil)
				}else {
					ResponseJson(w, nil, err)
				}
				return
			}
			if body.LockKey == "" {
				err = query.All(&users)
				if err != nil {
					ResponseJson(w, nil, err)
				} else {
					ResponseJson(w, users, nil)
				}
			} else {
				//////////////锁定表//////////////////
				var lock *sync.Mutex
				if _lock, o := LockMaps[body.Business]; o {
					lock = _lock
				} else {
					lock = &sync.Mutex{}
					LockMaps[body.Business] = lock
				}
				lock.Lock()
				defer lock.Unlock()
				////////////////////////////////
				err = query.All(&users)
				if err != nil {
					ResponseJson(w, nil, err)
				} else {
					var arr []interface{}
					for i := 0; i < len(users); i++ {
						if v, oo := users[i].(bson.M); oo {
							if key, o1 := v["_id"]; o1 {
								arr = append(arr, key)
							} else {
								ResponseJson(w, nil, errors.New("lack of _id column"))
								return
							}
						}
					}
					collection.UpdateAll(bson.M{"_id": bson.M{"$in": arr}}, bson.M{"$set": bson.M{body.LockKey: -1}})
					ResponseJson(w, users, nil)
				}
			}
		} else {
			ResponseJson(w, 0, errors.New("data type is not map[string]interface{} but a " + reflect.TypeOf(body.Data).Name()))
		}
	}
}

//新增
func add(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt32(&cnt, 1)
	defer atomic.AddInt32(&cnt, -1)
	if body := loadRequest(w, r); body.Business != "-1" {
//		if cnt > maxQueue { //如果任务队列超过了阀值
//			AddToTaskQueue(w, body) //放入任务队列中
//			return
//		}
		var er error
		if m, o := body.Data.(map[string]interface{}); o { //如果是新增一条数据
			er = addFunc(m, body)
		} else if m, o := body.Data.([]interface{}); o { //如果是批量新增
			for i, _ := range m {
				if t, ok := m[i].(map[string]interface{}); ok {
					e := addFunc(t, body)
					if er == nil {
						er = e
					}
				}
			}
		} else { //否则提示错误
			ResponseJson(w, 0, errors.New("data type is not map[string]interface{} or []interface{} but a " + reflect.TypeOf(body.Data).Name()))
			return
		}
		ResponseJson(w, nil, er)
	}
}

//插入到数据库
func addFunc( m map[string]interface{}, body PostData) error {
	session := InitMongo(body.Host)
	defer session.Close()
	collection := session.DB(body.DBName).C(body.Business)
	m["timestamp"] = time.Now().Unix() //设置时间戳
	if _, ok := m["_id"]; ok { //去除OnjectID列
		 delete(m, "_id")
	}
	if body.Update {
		if id, ok := m[body.Key]; ok {
			c, e := collection.Find(bson.M{body.Key: id}).Count()
			if e == nil {
				if c == 0 {
					err := collection.Insert(m) //插入数据
					if err != nil {
						return err
					}
				} else {
					err := collection.Update(bson.M{body.Key: m[body.Key]}, bson.M{"$set": m}) //更新数据
					if err != nil {
						return err
					}
				}
			} else {
				return e
			}
		}else {
			return errors.New("can't find the key from the json")
		}
	} else {
		err := collection.Insert(m) //插入数据
		if err != nil {
			return err
		}
	}
	return nil
}

//更新
func update(w http.ResponseWriter, r *http.Request) {
	defer Recovery("update")
	if body := loadRequest(w, r); body.Business != "-1" {
		if m, o := body.Data.(map[string]interface{}); o {
			if k, ok := m[body.Key]; ok {
				session := InitMongo(body.Host)
				if session == nil {
					ResponseJson(w, 0, errors.New("mongodb connected failed"))
					return
				}
				defer session.Close()
				collection := session.DB(body.DBName).C(body.Business)
				m["timestamp"] = time.Now().Unix() //设置时间戳
				err := collection.Update(bson.M{body.Key: k}, bson.M{"$set": m})
				if err != nil {
					ResponseJson(w, nil, errors.New("update failed:" + err.Error()))
					return
				}
				ResponseJson(w, nil, nil)
			} else {
				ResponseJson(w, nil, errors.New("can't find the key from the json"))
			}
		} else {
			ResponseJson(w, 0, errors.New("data type is not map[string]interface{}   but a " + reflect.TypeOf(body.Data).Name()))
		}
	}
}


//上传文件
func upload(w http.ResponseWriter, r *http.Request) {
	//	ResponseJson(w, "小米克隆101", nil)
	//	return ;
	r.ParseForm()
	if r.Method == "GET" {
		ResponseJson(w, nil, errors.New("it should be Post Method"))
	} else {
		vps := "common"
		if len(r.Form["vps"]) > 0 {
			vps = r.Form["vps"][0]
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			ResponseJson(w, nil, err)
			return
		}
		defer file.Close()
		dic := time.Now().Format("20060102")
		dir := "d:/image/" + vps + "/" + dic
		if _, err := os.Stat(dir); err != nil {
			err = os.MkdirAll(dir, 0777)
			if err != nil {
				ResponseJson(w, nil, err)
			}
		}
		fileName := dir + "/" + fmt.Sprint(time.Now().Unix()) + ".png"
		f, err := os.OpenFile(fileName, os.O_CREATE, 0666)
		if err != nil {
			ResponseJson(w, fileName, err)
			return
		}
		io.Copy(f, file)
		defer f.Close()
		ResponseJson(w, fileName, nil)
	}
}

//下载文件
func getFile(w http.ResponseWriter, r *http.Request){

}
