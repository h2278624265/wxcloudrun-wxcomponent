# wxcloudrun-wxcomponent
[![GitHub license](https://img.shields.io/github/license/WeixinCloud/wxcloudrun-wxcomponent)](https://github.com/WeixinCloud/wxcloudrun-wxcomponent)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/WeixinCloud/wxcloudrun-wxcomponent)

原项目官方文档：[微信第三方平台管理工具模版](https://github.com/WeixinCloud/wxcloudrun-wxcomponent)
微信服务商文档：[服务商微管家介绍](https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/product/management-tools.html)。

## 功能介绍

修改微信第三方平台管理工具模版，使其能够在非微信云的第三方服务器上部署（服务商文档中的“传统模式”）。

## 使用介绍

*开发模式*

./main.sh

*生产模式*

使用pm2管理常驻进程，启动：

npm start

也可以使用容器。

*配置文件*

1. 数据库环境变量模版在db.env.template

cp db.env.template db.env

根据实际修改db.env即可。

2. 修改系统配置文件comm/config/server.conf

AseKey=''  # 对应开发资料的 “消息加解密Key”
Token=''   # 对应开发资料的 “消息校验Token”

UseComponentAccessToken=true

## License

[MIT](./LICENSE)
