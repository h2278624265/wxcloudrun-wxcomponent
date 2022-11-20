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

type DataContext struct {
	Random string
	MsgLen int
	Data string
	AppId string
}

type xmlEncryptCallbackComponentRecord struct {
	AppId string `xml:"AppId"`
	Encrypt string `xml:"Encrypt"`
}

type xmlCallbackComponentRecord struct {
	CreateTime int64  `xml:"CreateTime"`
	InfoType   string `xml:"InfoType"`
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
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	encryptType := c.Query("encrypt_type")
	msgSignature := c.Query("msg_signature")

	fmt.Println("signature: ", signature)
	fmt.Println("timestamp: ", timestamp)
	fmt.Println("nonce: ", nonce)
	fmt.Println("encrypt_type: ", encryptType)
	fmt.Println("msg_signature: ", msgSignature)

	body, _ := ioutil.ReadAll(c.Request.Body)
	xmlBody := xmlEncryptCallbackComponentRecord{}
	err := xml.Unmarshal(body, &xmlBody)

	if err != nil {
		c.JSON(http.StatusOK, errno.ErrInvalidParam.WithData(err.Error()))
		return
	}
	toSign := []string{xmlBody.Encrypt, timestamp, nonce}
	if ok := utils.VerifyReqContext(toSign, msgSignature); ok == false {
		c.JSON(http.StatusUnauthorized, errno.ErrNotAuthorized)
		return
	}

	ctx, err := utils.DecryptReqContext(xmlBody.Encrypt)
	if err != nil {
		c.JSON(http.StatusOK, errno.ErrInvalidParam.WithData(err.Error()))
		return
	}
	// fmt.Println("ctx: ", ctx)
	c.Set("DecryptContext", []byte(ctx.Data))
	c.Next()
}
