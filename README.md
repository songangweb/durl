# 欢迎使用 durl 短链服务

## durl介绍:
durl 是一个分布式的高性能短链服务,逻辑简单,部署方便.经过多次迭代,现发布正式版本.

## 背景
发现在github中已有的短链服务中,非分布式服务无法做到快速扩容,并且有些项目是用的redis作为数据缓存,性能上不够优秀.

## 短链服务介绍
短链接，通俗来说，就是将长的URL网址，通过程序计算等方式，转换为简短的网址字符串。

## 使用场景
微博和Twitter都有140字数的限制，如果分享一个长网址，很容易就超出限制。
营销短信，字数的限制,当字数过长: 1.不美观 2.超出字符额外收费。
生成二维码的原始链接,当原始链接过长时,生成的二维码过于复杂,导致一些像素较低的手机无法扫描.

## 特征:
1. [beego](https://github.com/beego/beego) 为项目web框架.
2. 使用了 [xorm](https://github.com/xormplus/xorm) 来实现持久数据存储, 项目已测试 mysql 与 mongo.
3. 使用了 [mcache](https://github.com/songangweb/mcache) 来实现内存缓存.
4. 因使用内存缓存作为缓存池,实际使用中,项目本身的性能瓶颈更多体现在数据库自身.(单机qps轻松上w)
5. 项目内存消耗大多为缓存内存所用容量,可通过配置文件进行内存大小限制.

## durl的四个模块:

portal: 首页可以通过页面进行短链生成.公司内部或者公司外部可以通过页面生成短链接.
openApi: 对内开放api,增删改查. 一般来说可以部署为只内网访问.
jump: 只服务短链跳转.作为专门的跳转服务,当需要单机性能不够时,可直接横向扩容.
backend: 为后台管理页面,可管理短链接与黑名单. 可作为公司内部系统增加模块嵌入页面.
这样分为四个模块的原因,是因为根据需要进行部署,需要那个就部署那个.
因为这个项目的结构原因,整个项目四个模块之间没有耦合,可以随意增加pod数量,来提高系统性能.

## 系统架构

![avatar](https://github.com/songangweb/durl/wiki/durl.jpg)


## 如何使用

### 体验: docker-compose
1. 导入数据库文件 文件地址: durl/doc
2. 修改 durl/build/durl 目录下各个模块的配置信息
3. 在 durl/build/durl 目录下 执行 docker-compose up -d

### 系统部署:
在 durl/build 目录下提供有全模块的dockerfile demo. 可以根据需要进行修改后部署.

## 项目流程解析   详细了解durl项目时,此模块内容必看

[项目详解](https://github.com/songangweb/durl/wiki/Explain)

[comment]: <> ([配置文件详解]&#40;https://github.com/songangweb/durl/wiki/Explain&#41;)

## 接口文档
[接口文档](https://github.com/songangweb/durl/wiki/Api)

## JetBrains操作系统许可证

durl 是根据JetBrains sro授予的免费JetBrains开源许可证与GoLand一起开发的，因此在此我要表示感谢。

[免费申请 jetbrains 全家桶](https://zhuanlan.zhihu.com/p/264139984?utm_source=wechat_session)


## 交流
#### 如果文档中未能覆盖的任何疑问,欢迎您发送邮件到<songangweb@foxmail.com>,我会尽快答复。
#### 您可以在提出使用中需要改进的地方,我会考虑合理性并尽快修改。
#### 如果您发现 bug 请及时提 issue,我会尽快确认并修改。
#### 有劳点一下 star，一个小小的 star 是作者回答问题的动力 🤝
