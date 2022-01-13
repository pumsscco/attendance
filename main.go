package main

import (
    "fmt"
    "database/sql"
    "log"
    "net/http"
    yaml "gopkg.in/yaml.v2"
	"github.com/go-redis/redis"
	"io/ioutil"
    "github.com/julienschmidt/httprouter"
    _ "github.com/go-sql-driver/mysql"
    "os"
    "path/filepath"
)
type Conf struct {
    Listen struct {
        Port int `yaml:"port"`
    }
    MySQL struct {
        Db string `yaml:"db"`
        Host string `yaml:"host"`
        Port int `yaml:"port"`
        User string `yaml:"user"`
        Pass string `yaml:"pass"`
    }
    Redis struct {
        Host string `yaml:"host"`
        Port int `yaml:"port"`
        Db int `yaml:"db"`
        Pass string `yaml:"pass"`
    }
    Logfile string `yaml:"logfile"`
}
var (
    Db *sql.DB
    logger *log.Logger
    client *redis.Client
    cnf Conf
)
//所有初始化操作
func init() {
    //获得本可执行文件的当前路径
    dir,_:=filepath.Abs(filepath.Dir(os.Args[0]))
    //抓全部的配置信息
    yamlBytes, err := ioutil.ReadFile(dir+"/config.yml")
    if err!=nil {
        log.Fatalf("无法打开环境配置文件: %v",err)
    }
    yaml.Unmarshal(yamlBytes,&cnf)
    //基于相对路径打开日志
    file, err := os.OpenFile(dir+"/"+cnf.Logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("无法打开日志文件：%v\n", err)
    }
    logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
    client=redis.NewClient(&redis.Options{
        Addr:       fmt.Sprintf("%s:%d",cnf.Redis.Host,cnf.Redis.Port),
        Password:   cnf.Redis.Pass,
        DB:         cnf.Redis.Db,
    })
    _, err = client.Ping().Result()
    if err!=nil {
        logger.Fatalf("redis连接异常：%v\n",err)
    } 
    dsn:=fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?loc=Local&parseTime=true", cnf.MySQL.User, cnf.MySQL.Pass, cnf.MySQL.Host, cnf.MySQL.Port, cnf.MySQL.Db)
    //fmt.Println("Data Source Name: ",dsn)
    Db,err=sql.Open("mysql",dsn)
    if err!=nil {
        logger.Fatalf("mysql连接异常：%v\n",err)
    }
}
//
func main() {
    // handle static assets
    router := httprouter.New()
    router.ServeFiles("/static/*filepath", http.Dir("static"))
    //首页
    router.GET("/", index)
    //路由，先抓一个最新记录
    router.GET("/rec/latest", latest)
    router.GET("/rec/last-week", lastWeek)
    router.GET("/rec/last-month", lastMonth)
    router.GET("/api/rec/latest", latest)
    router.GET("/api/rec/last-week", lastWeek)
    router.GET("/api/rec/last-month", lastMonth)
    //统计
    router.GET("/stat/week", weekStats)
    router.GET("/stat/month", monthStats)
    router.GET("/stat/year", yearStats)
    router.GET("/api/stat/week", weekStats)
    router.GET("/api/stat/month", monthStats)
    router.GET("/api/stat/year", yearStats)
    //增加新记录
    router.GET("/add", add)
    router.POST("/add", add)
    router.POST("/api/add", add)
    //以最简方式启动服务
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d",cnf.Listen.Port),router))
}