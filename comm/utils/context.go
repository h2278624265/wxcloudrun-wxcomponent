package utils

import (
	"fmt"
	"encoding/base64"

	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/encrypt"
	"github.com/gin-gonic/gin"
)

func VerifyReqContext(c *gin.Context) {

}

func DecryptReqContext(c *gin.Context) {
	EncodingAESKey := "1un5Ns81alsa192xx1mE80jiFfk49bd56o81a3ob941"
	AESKey, err := base64.StdEncoding.DecodeString(EncodingAESKey + "=")
	msgEncrypt := "nsv0Mgy35Y9++FUakGwEn4WG01FCbX0apzaO0cL8h+KjcNVTgmFcUlA+Iizqc4aQ08dD60a+jtHCjqV8pTjaTzLb6ILTZjC4X5ZblmyYP8EDQtG4QeyGOJd/vepdU/i2+1sEip1IacZbJjQN2NJD6tFZRtKudW+3BqkiP9H9zyMMp/xCt5POGNODXllvStYeLI2zwXEXF6Da/OSWU7RX9dMNYPRFtW8lpRppeIXlz0IArWFNUunbrCAsI6+67jkSyNG65wyldwM2WSyQQUyRMNNBbi99tbrtnfoF5Zl6ZJLDhZyXMSh42Mdhw4remY0BDNzW0+FwpXT5uJMEp7b0VzWQvsPgItXKljaAtJMJjE0Kph75H6wIYjUHfcUUSM78H3dwIUSHLNkFKtyIfGIMQZ+oDMbDPS5i/vZCrOYnT1Ah9TMvTVvE1PNmmRafaiBaCE5VqYdNj23aE6Su6yBwmQ=="
	tmpMsg, err := base64.StdEncoding.DecodeString(msgEncrypt)

	fullMsg, err := encrypt.AesDecrypt(tmpMsg, AESKey)
	if err != nil {
		fmt.Println("fullMsg err: ", err)
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
	}

}