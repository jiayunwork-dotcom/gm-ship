# gm-ship：船舶初稳性计算工具

本仓库是一个用 Go 编写的船舶初稳性（intact & damaged stability）计算库与服务。
给定排水量、浮心高度、重心高度、水线面惯性矩以及自由液面等装载参数，它能计算：

- **BM**（横向稳心半径）= IT / ∇
- **GM**（初稳性高）= KB + BM − KG
- **GM_free**（经自由液面修正后的稳性高）= GM − i/∇
- **GZ**（复原力臂）= GM_free · sin(φ)，以及对应的复原力矩
- 破损稳性（进水工况的稳性高与横倾）、波浪稳性、纵倾（trim）、船型系数、
  阻力估算与排水量/载重量换算等一整套派生量

计算遵循经典小角度稳性模型；当 heel 超过 ±10° 时结果仍是同一公式，但会在
`Result.Warning` 中提示近似误差被放大。

## 模块划分

- `internal/stability` — 稳性核心计算（metacentric / righting / damage / wave / trim / index / report）
- `internal/hull` — 船体几何与系数（水线面、棱形系数、排水量、阻力）
- `internal/api` — 同进程托管的 JSON API 与上述模块的编排

## 输入 / 输出

核心输入 `stability.Input`：

| 字段 | 含义 | 单位 |
| --- | --- | --- |
| Volume | 排水体积 ∇ | m³ |
| KB | 浮心距基线高 | m |
| KG | 重心距基线高 | m |
| IT | 水线面横向惯性矩 | m⁴ |
| HeelDeg | 横倾角 | 度 |
| Density | 水密度 ρ（0 = 默认海水 1025 kg/m³） | kg/m³ |
| FreeSurface | 自由液面惯性矩 i | m⁴ |

计算结果 `stability.Result` 给出 BM、GM、GM_free、GZ（m）、复原力矩（N）以及 Warning。

## 命令行 / 服务用法

```text
go build ./...
go test ./...
go run . -http :8080
```

启动后访问 http://localhost:8080/ 查看内置页面；API：

- `POST /api/gm`：提交 JSON 形式的 `Input`，返回 GM / GZ 等结果
- `POST /api/scan`：提交 `Input` 与一个 heel 范围，返回 GZ 曲线采样点

示例输入：`example/barge.json`。

## 关键约定

- 中国船舶约定：涨/跌不在本工具涉及；单位默认 SI（米 / 千克 / 秒）。
- 角度一律在边界以「度」传入，仅在 `righting.go` 内部转弧度。
- 自由液面只降低 GM，不改变 GZ 的几何关系。
