package authpage

import (
	"net/http"

	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/errno"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/log"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/wx"
	wxbase "github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/wx/base"
	"github.com/gin-gonic/gin"
)

type getPreAuthCodeReq struct {
	ComponentAppid string `wx:"component_appid"`
	ComponentAccessToken string `wx:"component_access_token"`
}

type getPreAuthCodeResp struct {
	PreAuthCode string `wx:"pre_auth_code"`
}

func getPreAuthCodeHandler(c *gin.Context) {
	accessToken, tknErr := wx.GetComponentAccessToken()
	if tknErr != nil {
		c.JSON(http.StatusOK, errno.ErrSystemError.WithData(tknErr.Error()))
		return
	}
	req := getPreAuthCodeReq{
		ComponentAppid: wxbase.GetAppid(),
		ComponentAccessToken: accessToken,
	}
	_, body, err := wx.PostWxJsonWithComponentToken("/cgi-bin/component/api_create_preauthcode", "", req)
	if err != nil {
		c.JSON(http.StatusOK, errno.ErrSystemError.WithData(err.Error()))
		return
	}
	var resp getPreAuthCodeResp
	if err := wx.WxJson.Unmarshal(body, &resp); err != nil {
		log.Errorf("Unmarshal err, %v", err)
		c.JSON(http.StatusOK, errno.ErrSystemError.WithData(err.Error()))
		return
	}
	c.JSON(http.StatusOK, errno.OK.WithData(gin.H{
		"preAuthCode": resp.PreAuthCode,
	}))
}
