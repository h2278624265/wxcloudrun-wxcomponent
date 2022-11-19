package utils

import (
	"fmt"
	"encoding/base64"

	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/encrypt"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/config"
	"github.com/gin-gonic/gin"
)

type DataContext struct {
	Random string
	MsgLen int
	Data string
	AppId string
}

func VerifyReqContext(c *gin.Context) {

}

func DecryptReqContext(msgEncrypt string) (context *DataContext, err error) {
	AesKeyDecode, err := base64.StdEncoding.DecodeString(config.ServerConf.AesKey + "=")
	tmpMsg, err := base64.StdEncoding.DecodeString(msgEncrypt)
	var ctx DataContext
	if fullMsg, err := encrypt.AesDecrypt(tmpMsg, AesKeyDecode); err != nil {
		fmt.Println("fullMsg err: ", err)
		return &ctx, err
	} else {
		fmt.Println("fullMsg", fullMsg)
		var random string = string(fullMsg[:16])
		fmt.Println("random:", random)

		msgLen := BytesToInt(fullMsg[16:20])
		fmt.Println("msgLen:", msgLen)

		var remain string = string(fullMsg[20:])
		// fmt.Println("data:", remain)

		var data string = remain[:msgLen]
		fmt.Println("msg:", remain);

		var appId string = remain[msgLen:]
		fmt.Println("appId:", appId)

		ctx := DataContext{
			Random: random, 
			MsgLen: msgLen, 
			Data: data,
			AppId: appId,
		}
		return &ctx, nil
	}
}