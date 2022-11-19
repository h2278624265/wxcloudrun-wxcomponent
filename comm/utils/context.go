package utils

import (
	"fmt"
	"encoding/base64"

	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/encrypt"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/config"
	"github.com/gin-gonic/gin"
)

type DataContext struct {
	random string
	msgLen int
	data string
	appId string
}

func VerifyReqContext(c *gin.Context) {

}

func DecryptReqContext(msgEncrypt string) (context DataContext, err error) {
	AesKeyDecode, err := base64.StdEncoding.DecodeString(config.ServerConf.AesKey + "=")
	tmpMsg, err := base64.StdEncoding.DecodeString(msgEncrypt)
	var ctx DataContext
	if fullMsg, err := encrypt.AesDecrypt(tmpMsg, AesKeyDecode); err != nil {
		fmt.Println("fullMsg err: ", err)
		return ctx, err
	} else {
		fmt.Println("fullMsg", fullMsg)
		var random string = string(fullMsg[:16])
		fmt.Println("random:", random)

		msgLen := BytesToInt(fullMsg[16:20])
		fmt.Println("msgLen:", msgLen)

		var data string = string(fullMsg[20:])
		fmt.Println("data:", data)

		var msg string = data[:msgLen]
		fmt.Println("msg:", msg);

		var appId string = data[msgLen:]
		fmt.Println("appId:", appId)

		ctx := DataContext{random, msgLen, data, appId};
		return ctx, nil
	}
}