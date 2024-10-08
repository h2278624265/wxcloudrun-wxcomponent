package wxcallback

import (
	// "fmt"
	// "io/ioutil"
	"net/http"
	"time"

	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/errno"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/log"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/wx"

	wxbase "github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/wx/base"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/db/dao"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/db/model"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type wxCallbackComponentRecord struct {
	CreateTime int64  `xml:"CreateTime" json:"CreateTime"`
	InfoType   string `xml:"InfoType" json:"InfoType"`
}

func componentHandler(c *gin.Context) {
	// 记录到数据库
	// body, _ := ioutil.ReadAll(c.Request.Body)
	// fmt.Println("body: ", body)
	decryptCtx, ok := c.Get("DecryptContext")
	body := decryptCtx.([]byte)
	if !ok {
		c.JSON(http.StatusOK, gin.H{ "ok": ok })
		return
	}
	
	var record wxCallbackComponentRecord
	if err := binding.XML.BindBody(body, &record); err != nil {
		c.JSON(http.StatusOK, errno.ErrInvalidParam.WithData(err.Error()))
		return
	}

	// if err := c.ShouldBind(body, &json); err != nil {
	// 	fmt.Println("1 error:", err.Error())
	// 	c.JSON(http.StatusOK, errno.ErrInvalidParam.WithData(err.Error()))
	// 	return
	// }

	// fmt.Println("record:", record)
	r := model.WxCallbackComponentRecord{
		CreateTime:  time.Unix(record.CreateTime, 0),
		ReceiveTime: time.Now(),
		InfoType:    record.InfoType,
		PostBody:    string(body),
	}
	if record.CreateTime == 0 {
		r.CreateTime = time.Unix(1, 0)
	}
	if err := dao.AddComponentCallBackRecord(&r); err != nil {
		c.JSON(http.StatusOK, errno.ErrSystemError.WithData(err.Error()))
		return
	}

	// 处理授权相关的消息
	var err error
	switch record.InfoType {
	case "component_verify_ticket":
		err = ticketHandler(&body)
	case "authorized":
		fallthrough
	case "updateauthorized":
		err = newAuthHander(&body)
	case "unauthorized":
		err = unAuthHander(&body)
	}
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, errno.ErrSystemError.WithData(err.Error()))
		return
	}

	// 转发到用户配置的地址
	var proxyOpen bool
	proxyOpen, err = proxyCallbackMsg(record.InfoType, "", "", string(body), c)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, errno.ErrSystemError.WithData(err.Error()))
		return
	}
	if !proxyOpen {
		c.String(http.StatusOK, "success")
	}
}

type ticketRecord struct {
	ComponentVerifyTicket string `xml:"ComponentVerifyTicket" json:"ComponentVerifyTicket"`
}

func ticketHandler(body *[]byte) error {
	var record ticketRecord
	if err := binding.XML.BindBody(*body, &record); err != nil {
		return err
	}

	log.Info("[new ticket]" + record.ComponentVerifyTicket)
	if err := wxbase.SetTicket(record.ComponentVerifyTicket); err != nil {
		return err
	}
	return nil
}

type newAuthRecord struct {
	CreateTime                   int64  `xml:"CreateTime" json:"CreateTime"`
	AuthorizerAppid              string `xml:"AuthorizerAppid" json:"AuthorizerAppid"`
	AuthorizationCode            string `xml:"AuthorizationCode" json:"AuthorizationCode"`
	AuthorizationCodeExpiredTime int64  `xml:"AuthorizationCodeExpiredTime" json:"AuthorizationCodeExpiredTime"`
}

func newAuthHander(body *[]byte) error {
	var record newAuthRecord
	var err error
	var refreshtoken string
	var appinfo wx.AuthorizerInfoResp
	if err = binding.XML.BindBody(*body, &record); err != nil {
		return err
	}

	if refreshtoken, err = queryAuth(record.AuthorizationCode); err != nil {
		return err
	}
	if err = wx.GetAuthorizerInfo(record.AuthorizerAppid, &appinfo); err != nil {
		return err
	}
	if err = dao.CreateOrUpdateAuthorizerRecord(&model.Authorizer{
		Appid:         record.AuthorizerAppid,
		AppType:       appinfo.AuthorizerInfo.AppType,
		ServiceType:   appinfo.AuthorizerInfo.ServiceType.Id,
		NickName:      appinfo.AuthorizerInfo.NickName,
		UserName:      appinfo.AuthorizerInfo.UserName,
		HeadImg:       appinfo.AuthorizerInfo.HeadImg,
		QrcodeUrl:     appinfo.AuthorizerInfo.QrcodeUrl,
		PrincipalName: appinfo.AuthorizerInfo.PrincipalName,
		RefreshToken:  refreshtoken,
		FuncInfo:      appinfo.AuthorizationInfo.StrFuncInfo,
		VerifyInfo:    appinfo.AuthorizerInfo.VerifyInfo.Id,
		AuthTime:      time.Unix(record.CreateTime, 0),
	}); err != nil {
		return err
	}
	return nil
}

type queryAuthReq struct {
	ComponentAppid    string `wx:"component_appid"`
	AuthorizationCode string `wx:"authorization_code"`
}

type authorizationInfo struct {
	AuthorizerRefreshToken string `wx:"authorizer_refresh_token"`
}
type queryAuthResp struct {
	AuthorizationInfo authorizationInfo `wx:"authorization_info"`
}

func queryAuth(authCode string) (string, error) {
	req := queryAuthReq{
		ComponentAppid:    wxbase.GetAppid(),
		AuthorizationCode: authCode,
	}
	var resp queryAuthResp
	_, body, err := wx.PostWxJsonWithComponentToken("/cgi-bin/component/api_query_auth", "", req)
	if err != nil {
		return "", err
	}
	if err := wx.WxJson.Unmarshal(body, &resp); err != nil {
		log.Errorf("Unmarshal err, %v", err)
		return "", err
	}
	return resp.AuthorizationInfo.AuthorizerRefreshToken, nil
}

type unAuthRecord struct {
	CreateTime      int64  `xml:"CreateTime json:"CreateTime"`
	AuthorizerAppid string `xml:"AuthorizerAppid" json:"AuthorizerAppid"`
}

func unAuthHander(body *[]byte) error {
	var record unAuthRecord
	var err error
	if err = binding.XML.BindBody(*body, &record); err != nil {
		log.Errorf("bind err %v", err)
		return err
	}
	if err := dao.DelAuthorizerRecord(record.AuthorizerAppid); err != nil {
		log.Errorf("DelAuthorizerRecord err %v", err)
		return err
	}
	return nil
}
