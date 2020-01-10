package controllers

import (
	"compress/gzip"
	"encoding/json"
	"github.com/astaxie/beego"
	"io"
	"strings"
)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) JsonResult(errCode int, errMsg string, data ...interface{}) {
	jsonData := make(map[string]interface{})
	jsonData["errcode"] = errCode
	jsonData["message"] = errMsg

	if len(data) > 0 {
		jsonData["data"] = data[0]
	}
	resJson, err := json.Marshal(jsonData)
	if err != nil {
		beego.Error(err)
	}
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf8")
	if strings.Contains(strings.ToLower(c.Ctx.Request.Header.Get("Accept-Encoding")), "gzip") {
		c.Ctx.ResponseWriter.Header().Set("Content-Encoding", "gzip")
		w := gzip.NewWriter(c.Ctx.ResponseWriter)
		defer w.Close()
		w.Write(resJson)
		w.Flush()
	} else {
		io.WriteString(c.Ctx.ResponseWriter, string(resJson))
	}
	c.StopRun()
}
