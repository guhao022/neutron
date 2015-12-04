package neutron

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"github.com/guhao022/neutron/log"
)

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
}

func (ctx *Context) Debug() {
	log.CLog("[INFO]  [ %s ] | < %s > | [ %s ] \n", ctx.Request.Method, ctx.Request.RequestURI, ctx.IP())
}

//==========================INTPUT--START=========================

func (ctx *Context) Uri() string {
	return ctx.Request.RequestURI
}

func (ctx *Context) Url() string {
	return ctx.Request.URL.Path
}

func (ctx *Context) Scheme() string {
	if ctx.Request.URL.Scheme != "" {
		return ctx.Request.URL.Scheme
	}
	if ctx.Request.TLS == nil {
		return "http"
	}
	return "https"
}

func (ctx *Context) Method() string {
	return ctx.Request.Method
}

func (ctx *Context) IsAjax() bool {
	return ctx.Header("X-Requested-With") == "XMLHttpRequest"
}

func (ctx *Context) IP() string {
	ips := ctx.Proxy()
	if len(ips) > 0 && ips[0] != "" {
		rip := strings.Split(ips[0], ":")
		return rip[0]
	}
	ip := strings.Split(ctx.Request.RemoteAddr, ":")
	if len(ip) > 0 {
		if ip[0] != "[" {
			return ip[0]
		}
	}
	return "127.0.0.1"
}

func (ctx *Context) Proxy() []string {
	if ips := ctx.Header("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}

func (ctx *Context) Port() int {
	parts := strings.Split(ctx.Request.Host, ":")
	if len(parts) == 2 {
		port, _ := strconv.Atoi(parts[1])
		return port
	}
	return 80
}

func (ctx *Context) UserAgent() string {
	return ctx.Header("User-Agent")
}

func (ctx *Context) Header(key string) string {
	return ctx.Request.Header.Get(key)
}

func (ctx *Context) GetStrings(key string) []string {
	return ctx.Request.Form[key]
}

func (ctx *Context) GetString(key string, def ...string) string {
	var defv string
	if len(def) > 0 {
		defv = def[0]
	}

	val := ctx.Request.FormValue(key)
	if val == "" {
		return defv
	}
	return val
}

func (ctx *Context) GetInt(key string, def ...int) (int, error) {
	if strv := ctx.GetString(key); strv != "" {
		return strconv.Atoi(strv)
	} else if len(def) > 0 {
		return def[0], nil
	} else {
		return strconv.Atoi(strv)
	}
}

func (ctx *Context) GetInt8(key string, def ...int8) (int8, error) {
	if strv := ctx.GetString(key); strv != "" {
		i64, err := strconv.ParseInt(strv, 10, 8)
		i8 := int8(i64)
		return i8, err
	} else if len(def) > 0 {
		return def[0], nil
	} else {
		i64, err := strconv.ParseInt(strv, 10, 8)
		i8 := int8(i64)
		return i8, err
	}
}

func (ctx *Context) GetInt16(key string, def ...int16) (int16, error) {
	if strv := ctx.GetString(key); strv != "" {
		i64, err := strconv.ParseInt(strv, 10, 16)
		i16 := int16(i64)
		return i16, err
	} else if len(def) > 0 {
		return def[0], nil
	} else {
		i64, err := strconv.ParseInt(strv, 10, 16)
		i16 := int16(i64)
		return i16, err
	}
}

func (ctx *Context) GetInt32(key string, def ...int32) (int32, error) {
	if strv := ctx.GetString(key); strv != "" {
		i64, err := strconv.ParseInt(strv, 10, 32)
		i32 := int32(i64)
		return i32, err
	} else if len(def) > 0 {
		return def[0], nil
	} else {
		i64, err := strconv.ParseInt(strv, 10, 32)
		i32 := int32(i64)
		return i32, err
	}
}

func (ctx *Context) GetInt64(key string, def ...int64) (int64, error) {
	if strv := ctx.GetString(key); strv != "" {
		return strconv.ParseInt(strv, 10, 64)
	} else if len(def) > 0 {
		return def[0], nil
	} else {
		return strconv.ParseInt(strv, 10, 64)
	}
}

func (ctx *Context) GetBool(key string, def ...bool) (bool, error) {
	if strv := ctx.GetString(key); strv != "" {
		return strconv.ParseBool(strv)
	} else if len(def) > 0 {
		return def[0], nil
	} else {
		return strconv.ParseBool(strv)
	}
}

func (ctx *Context) GetFloat(key string, def ...float64) (float64, error) {
	if strv := ctx.GetString(key); strv != "" {
		return strconv.ParseFloat(strv, 64)
	} else if len(def) > 0 {
		return def[0], nil
	} else {
		return strconv.ParseFloat(strv, 64)
	}
}

//=============================INTPUT--END=========================

func (ctx *Context) Cookie(key string, value ...string) string {
	if len(value) < 1 {
		c, e := ctx.Request.Cookie(key)
		if e != nil {
			return ""
		}
		return c.Value
	}
	if len(value) == 2 {
		t := time.Now()
		expire, _ := strconv.Atoi(value[1])
		t = t.Add(time.Duration(expire) * time.Second)
		cookie := &http.Cookie{
			Name:    key,
			Value:   value[0],
			Path:    "/",
			MaxAge:  expire,
			Expires: t,
		}
		http.SetCookie(ctx.Response, cookie)
		return ""
	}
	return ""
}

//============================OUTPUT--START============================

func (ctx *Context) Redirect(status int, localurl string) {
	ctx.Response.Header().Set("Location", localurl)
	ctx.Response.WriteHeader(status)
}

func (ctx *Context) SetHeader(key, val string) {
	ctx.Response.Header().Set(key, val)
}

func (ctx *Context) Abort(status int, body string) {
	ctx.Response.WriteHeader(status)
	ctx.WriteString(body)
	return
	//fmt.Printf("%s",body)
}

func (ctx *Context) WriteString(content string) {
	ctx.Response.Write([]byte(content))
}

func (ctx *Context) Json(i interface{}) {
	out, err := json.Marshal(i)
	if err != nil {
		return
	}
	ctx.Response.Header().Set("content-type", "application/json; charset=utf-8")
	ctx.Response.Write(out)
}

func (ctx *Context) Download(file string, filename ...string) {
	ctx.SetHeader("Content-Description", "File Transfer")
	ctx.SetHeader("Content-Type", "application/octet-stream")
	if len(filename) > 0 && filename[0] != "" {
		ctx.SetHeader("Content-Disposition", "attachment; filename="+filename[0])
	} else {
		ctx.SetHeader("Content-Disposition", "attachment; filename="+filepath.Base(file))
	}
	ctx.SetHeader("Content-Transfer-Encoding", "binary")
	ctx.SetHeader("Expires", "0")
	ctx.SetHeader("Cache-Control", "must-revalidate")
	ctx.SetHeader("Pragma", "public")
	http.ServeFile(ctx.Response, ctx.Request, file)
}

//============================OUTPUT--STOP============================
