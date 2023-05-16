### Micro_Tmpl 微服务架构模板

#### 包含模块
- gateway 网关模块 流量入口 负责组织下层逻辑 独立服务器 (BFF层)
  - gin
  - route
  - grpc
- player 玩家模块 对外提供处理玩家API 独立服务器
- common 通用模块 对外提供一些公共库和公共定义
- config 公共配置