# coupang_spider

台湾地区便利店 / 物流单号追踪爬虫服务，基于 go-zero 框架。

## 目录结构

```
coupang_spider/
├── api_desc/coupang_spider.api    # goctl api 描述
├── etc/coupang_spider.yaml        # 运行时配置
├── resource/                      # 原始需求 docx 文档
├── internal/
│   ├── config/                    # 配置结构
│   ├── handler/                   # HTTP handler 层
│   ├── logic/                     # 业务逻辑层
│   ├── svc/                       # ServiceContext
│   ├── types/                     # 请求/响应类型
│   └── pkg/spider/                # 具体爬虫 client 实现
└── coupang_spider.go              # main 入口
```

## 对外 API

所有接口均为 `POST`，前缀 `/api/v1/spider`，请求体统一为：

```json
{ "trackNo": "单号" }
```

| 路径 | 说明 |
|------|------|
| `/api/v1/spider/seven/track`      | 台湾 7-ELEVEN 交货便 |
| `/api/v1/spider/familymart/track` | 台湾全家 FamilyMart |
| `/api/v1/spider/hilife/track`     | 台湾莱尔富 Hi-Life |
| `/api/v1/spider/okmart/track`     | 台湾 OK 超商 |
| `/api/v1/spider/kerry/track`      | 台湾嘉里大荣 |
| `/api/v1/spider/spx/track`        | 台湾虾皮 SPX |

响应：

```json
{
  "trackNo": "...",
  "carrier": "KERRY",
  "statusList": [
    {"time":"2025/04/01 12:00","status":"配送中","location":"台北"}
  ],
  "raw": "原始 HTML/JSON（可选）"
}
```

## 运行

```bash
go mod tidy
go run coupang_spider.go -f ./etc/coupang_spider.yaml
```
