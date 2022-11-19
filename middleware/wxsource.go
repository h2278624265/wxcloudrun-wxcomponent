package middleware

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/errno"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/utils"

	"github.com/gin-gonic/gin"
)

// WXSourceMiddleWare 中间件 判断是否来源于微信
func WXSourceMiddleWare(c *gin.Context) {
	if _, ok := c.Request.Header[http.CanonicalHeaderKey("x-wx-source")]; ok {
		fmt.Println("[WXSourceMiddleWare]from wx")
		c.Next()
	} else {
		c.Abort()
		c.JSON(http.StatusUnauthorized, errno.ErrNotAuthorized)
	}
}

func DecryptContext(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	var xmlBody string
	if errXml := c.ShouldBindBodyWith(&xmlBody, binding.XML); errXml == nil {
		fmt.Println("XML body: ", xmlBody)
		utils.DecryptReqContext(xmlBody)
	}
	

	// body, _ := ioutil.ReadAll(c.Request.Body)
}
