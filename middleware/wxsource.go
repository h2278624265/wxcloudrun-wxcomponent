package middleware

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/xml"

	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/errno"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/utils"

	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin/binding"
)

type xmlCallbackComponentRecord struct {
	AppId string `xml:"AppId"`
	Encrypt string `xml:"Encrypt"`
}

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
	fmt.Println("body: ", body)

	xmlBody := xmlCallbackComponentRecord{}
	err := xml.Unmarshal(body, &xmlBody)

	if err != nil {
		c.JSON(http.StatusOK, errno.ErrInvalidParam.WithData(err.Error()))
		return
	}
	// fmt.Println("XML body: ", xmlBody)
	fmt.Println("XML body: ", xmlBody)
	utils.DecryptReqContext(xmlBody.Encrypt)
	// if errXml := c.ShouldBindBodyWith(&xmlBody, binding.XML); errXml == nil {
	// 	fmt.Println("XML body: ", xmlBody)
	// 	utils.DecryptReqContext(xmlBody)
	// }
	

	// body, _ := ioutil.ReadAll(c.Request.Body)
}
