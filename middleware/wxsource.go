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
	fmt.Println("body: ", body)

	xmlBody := xmlEncryptCallbackComponentRecord{}
	err := xml.Unmarshal(body, &xmlBody)

	if err != nil {
		c.JSON(http.StatusOK, errno.ErrInvalidParam.WithData(err.Error()))
		return
	}
	// fmt.Println("XML body: ", xmlBody)
	ctx, err := utils.DecryptReqContext(xmlBody.Encrypt)

	if err != nil {
		c.JSON(http.StatusOK, errno.ErrInvalidParam.WithData(err.Error()))
		return
	}
	fmt.Println("ctx: ", ctx)

	c.Set("Decrypt", []byte(ctx.Data))
	c.Next()

	// ctxData := xmlCallbackComponentRecord{}

	// err2 := xml.Unmarshal([]byte(ctx.Data), &ctxData)
	// fmt.Println("ctx body: ", ctxData)

	// if err := binding.XML.BindBody([]byte(ctx.Data), &ctxData); err != nil {
	// 	fmt.Println("err: ", err.Error())
	// 	c.JSON(http.StatusOK, errno.ErrInvalidParam.WithData(err.Error()))
	// 	return
	// } else {
	// 	fmt.Println("ctx body infoType: ", ctxData.InfoType)
	// 	fmt.Println("ctx body createTime: ", ctxData.CreateTime)
	// 	c.Set("Body", ctxData)
	// 	c.Next()
	// }


	// if errXml := c.ShouldBindBodyWith(&xmlBody, binding.XML); errXml == nil {
	// 	fmt.Println("XML body: ", xmlBody)
	// 	utils.DecryptReqContext(xmlBody)
	// }
	

	// body, _ := ioutil.ReadAll(c.Request.Body)
}
