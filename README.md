# 欢迎使用 durl 短链服务

## 体验地址
# [durl.fun](https://durl.fun)


## 介绍
durl 是一个分布式的高性能短链服务,逻辑简单,并提供了相关api接口,开发人员可以快速接入,也可以作为go初学者练手项目.


## 特征
1. [beego](https://github.com/beego/beego) 为项目web框架.
2. 使用了 [xorm](https://github.com/xormplus/xorm) 来实现持久数据存储, 项目已测试 mysql 与 mongo.
3. 使用了 [mcache](https://github.com/songangweb/mcache) 来实现内存缓存.
4. 因使用内存缓存作为缓存池,实际使用中,项目本身的性能瓶颈更多体现在数据库自身.
5. 项目内存消耗大多为缓存内存所用容量,可通过配置文件进行内存大小限制.


## 如何使用
1. 数据导入数据库表结构
2. 修改配置文件
3. run ~


## 项目流程解析   详细了解durl项目时,此模块内容必看

[项目详解](https://github.com/songangweb/durl/wiki/Explain)

[项目目录结构](https://github.com/songangweb/durl/wiki/Directory)

[配置文件详解](https://github.com/songangweb/durl/wiki/Explain)

## openApi

[接口文档](https://github.com/songangweb/durl/wiki/OpenApi)


## JetBrains操作系统许可证

durl 是根据JetBrains sro授予的免费JetBrains开源许可证与GoLand一起开发的，因此在此我要表示感谢。

[免费申请 jetbrains 全家桶](https://zhuanlan.zhihu.com/p/264139984?utm_source=wechat_session)


## 交流
#### 如果文档中未能覆盖的任何疑问,欢迎您发送邮件到<songangweb@foxmail.com>,我会尽快答复。
#### 您可以在提出使用中需要改进的地方,我会考虑合理性并尽快修改。
#### 如果您发现 bug 请及时提 issue,我会尽快确认并修改。
#### 有劳点一下 star，一个小小的 star 是作者回答问题的动力 🤝
