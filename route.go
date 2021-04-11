package main

import (
    "fmt"
	"strings"
    "html/template"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
    "github.com/julienschmidt/httprouter"
    "encoding/json"
    "time"
)
//生成页面
func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
    var files []string
    for _, file := range filenames {
        files = append(files, fmt.Sprintf("templates/%s.html", file))
    }
    templates := template.Must(template.ParseFiles(files...))
    templates.ExecuteTemplate(w, "layout", data)
}
//首页
func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    generateHTML(w, nil, "layout", "navbar", "index")
}
//生成最新记录页面
func latest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    attn, _ := getLatest()
    p:=r.URL.Path
    if strings.HasPrefix(p,"/api/") {
        output, _ := json.MarshalIndent(&attn, "", "\t")
        w.Header().Set("Content-Type", "application/json")
        w.Write(output)
    } else {
        generateHTML(w, &attn, "layout", "navbar", "rec/latest")
    } 
}
//最近一周记录页面
func lastWeek(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    attn, _ := getLastWeek()
    p:=r.URL.Path
    if strings.HasPrefix(p,"/api/") {
        output, _ := json.MarshalIndent(&attn, "", "\t")
        w.Header().Set("Content-Type", "application/json")
        w.Write(output)
    } else {
        generateHTML(w, &attn, "layout", "navbar", "rec/last-week")
    } 
}
//最近一月记录页面
func lastMonth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    attn, _ := getLastMonth()
    p:=r.URL.Path
    if strings.HasPrefix(p,"/api/") {
        output, _ := json.MarshalIndent(&attn, "", "\t")
        w.Header().Set("Content-Type", "application/json")
        w.Write(output)
    } else {
        generateHTML(w, &attn, "layout", "navbar", "rec/last-month")
    } 
}
//月统计页
func monthStats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    monthHours:=monthHour()
    p:=r.URL.Path
    if strings.HasPrefix(p,"/api/") {
        output, _ := json.MarshalIndent(&monthHours, "", "\t")
        w.Header().Set("Content-Type", "application/json")
        w.Write(output)
    } else {
        generateHTML(w, &monthHours, "layout", "navbar", "stat/month")
    }
}
//周统计页
func weekStats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    weekHours:=weekHour()
    p:=r.URL.Path
    if strings.HasPrefix(p,"/api/") {
        output, _ := json.MarshalIndent(&weekHours, "", "\t")
        w.Header().Set("Content-Type", "application/json")
        w.Write(output)
    } else {
        generateHTML(w, &weekHours, "layout", "navbar", "stat/week")
    }
}
//年统计页
func yearStats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    yearHours:=yearHour()
    p:=r.URL.Path
    if strings.HasPrefix(p,"/api/") {
        output, _ := json.MarshalIndent(&yearHours, "", "\t")
        w.Header().Set("Content-Type", "application/json")
        w.Write(output)
    } else {
        generateHTML(w, &yearHours, "layout", "navbar", "stat/year")
    }
}
//新增记录
func add(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    if r.Method=="GET" {
        generateHTML(w, nil, "layout", "navbar", "create/add")
    } else if r.Method=="POST" {
        p:=r.URL.Path
        tf := "2006-01-02 15:04"
        var attn Attendance
        if strings.HasPrefix(p,"/api/") {
            type TmpAttn struct {
                Id        int `json:"id"`
                CheckIn      string `json:"checkin"`
                CheckOut     string `json:"checkout"`
                Comments    string `json:"comments"`
            }
            var ta TmpAttn
            len := r.ContentLength
            body := make([]byte, len)
            r.Body.Read(body)
            err:=json.Unmarshal(body, &ta)
            if err!=nil {
                logger.Printf("无法解析json数据到结构中：%v\n",err)
                return
            }
            local, _ := time.LoadLocation("Local")
            attn.CheckIn,_ = time.ParseInLocation(tf,ta.CheckIn,local)
            attn.CheckOut,_ = time.ParseInLocation(tf,ta.CheckOut,local)
            attn.Comments=ta.Comments
            if errInfo:=attn.NewAttn(); errInfo!="" {
                output, _ := json.MarshalIndent(&errInfo, "", "\t")
                w.Header().Set("Content-Type", "application/json")
                w.Write(output)
            } else  {
                http.Redirect(w, r, "/api/rec/latest", 302)
            }
        } else {
            err := r.ParseForm()
            if err != nil {
                logger.Printf("解析表单数据出错：%v\n",err)
                return
            }
            local, _ := time.LoadLocation("Local")
            attn.CheckIn,_ = time.ParseInLocation(tf,r.PostFormValue("checkin"),local)
            attn.CheckOut,_ = time.ParseInLocation(tf,r.PostFormValue("checkout"),local)
            attn.Comments = r.PostFormValue("comments")
            if errInfo:=attn.NewAttn(); errInfo!="" {
                generateHTML(w, &errInfo, "layout", "navbar", "create/invalid")
            } else {
                http.Redirect(w, r, "/rec/latest", 302)
            }
        }
    }
}