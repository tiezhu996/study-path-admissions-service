# gb-12-1 留学申请管理系统 执行记录

- 项目编号/名称：gb-12-1 / 留学申请管理系统（gbstudyapply）
- 日期：2026-08-16
- 短名：gbstudyapply
- 端口：前端 8012 / 后端 3012 / PostgreSQL 5505 / MinIO 9009+9010
- 技术栈：React 18 + TypeScript + Ant Design + ECharts + Vite（前端）；Go 1.22 + Gin + GORM（后端）；PostgreSQL + MinIO（minio-go/v7）；JWT + RBAC（student/counselor/admin）

## Docker Compose 结果

| 容器 | 状态 | 端口 |
| --- | --- | --- |
| gbstudyapply_db | healthy | 0.0.0.0:5505->5432 |
| gbstudyapply_minio | healthy | 0.0.0.0:9009->9000, 0.0.0.0:9010->9001 |
| gbstudyapply_backend | healthy | 0.0.0.0:3012->8080 |
| gbstudyapply_frontend | up | 0.0.0.0:8012->80 |

`docker compose config --quiet` 通过；验证后 `docker compose down -v --remove-orphans` 无残留。

## 关键 API 冒烟结果（29 项，全部通过）

| 接口 | 方法 | 状态码 | 结果摘要 |
| --- | --- | --- | --- |
| /api/v1/users/register | POST | 201 | 学生注册返回 JWT |
| /api/v1/users/login | POST | 200 | 学生/顾问/管理员登录 |
| /api/v1/users/me | GET | 200 | 当前用户资料 |
| /api/v1/universities?country=美国 | GET | 200 | 院校国家筛选 |
| /api/v1/universities/:id | GET | 200 | 院校详情 |
| /api/v1/universities（管理员） | POST | 201 | 管理员创建院校 |
| /api/v1/universities（学生） | POST | 403 | RBAC 拒绝非管理员 |
| /api/v1/applications（学生） | POST | 201 | 创建申请项目 |
| /api/v1/applications/:id | GET | 200 | 项目详情（学生/顾问均可读） |
| /api/v1/applications/:id/status | PUT | 200 | 状态流转 planning->preparing |
| /api/v1/applications/:id/status（跳级） | PUT | 409 | 非法流转拒绝 |
| /api/v1/applications/:id/documents | POST | 201 | 创建文书 |
| /api/v1/documents/:id | PUT | 200 | 保存新版本 v2 |
| /api/v1/documents/:id/versions | GET | 200 | 版本历史 |
| /api/v1/documents/:id/rollback | POST | 200 | 回滚到 v1 |
| /api/v1/documents/:id/annotations（顾问） | POST | 201 | 顾问批注 |
| /api/v1/documents/:id/annotations（学生） | POST | 403 | 非顾问拒绝 |
| /api/v1/applications/:id/materials | POST | 201 | 创建材料项 |
| /api/v1/materials/:id/status | PUT | 200 | 材料状态 pending->uploaded |
| /api/v1/applications/:id/timeline | POST | 201 | 创建时间线节点 |
| /api/v1/timeline/:id/done | PUT | 200 | 标记节点完成 |
| /api/v1/recommendations（顾问） | POST | 201 | 生成选校方案 |
| /api/v1/recommendations（学生） | POST | 403 | 非顾问拒绝 |
| /api/v1/messages | POST | 201 | 发送站内消息 |
| /api/v1/messages/:id/read | PUT | 200 | 标记已读 |
| /api/v1/dashboard/stats | GET | 200 | 看板统计（状态分布/录取率/材料进度） |
| /api/v1/users/me（无 token） | GET | 401 | 未认证拒绝 |
| /api/v1/universities/999999 | GET | 404 | 不存在返回 404 |
| /api/v1/users/register（非法） | POST | 400 | 参数校验失败 |

## 浏览器验证结论（内置 playwright，无外部 Chrome）

- /login：fill+click 登录成功，跳转 /dashboard 并显示"李同学（学生）" ✅
- /dashboard：申请项目总数 4 / 已提交 3 / 已录取 1 / 材料平均进度 16% + 项目列表（已录取等状态标签）+ ECharts 状态分布饼图（真实数据）✅
- /universities：国家筛选 + 院校卡片（哥伦比亚大学#3、伦敦大学学院#8、新加坡国立大学#11 等真实数据）✅
- 截图：output/dashboard.png

## README 检查项

- Docker Compose 一键启动命令置顶 ✅；本地开发命令 ✅；技术栈表格（后端 Go 1.22 + Gin + GORM）✅
- 目录结构、环境变量表、部署说明、MinIO 控制台说明 ✅
- 枚举出现位置清单：ApplicationStatus / DocumentType / MaterialStatus / UserRole（前后端全部位置）✅
- 横切关注点清单（含截止提醒）✅、License ✅

## 其他质量项

- 后端：`go build ./...`、`go vet ./...` 通过；`go test ./...` 通过（config/util/service 表驱动测试）
- 前端：`npm run build` 零错误（tsc + vite）
- 屎山设计：log_templates.go 30+ 条模板、错误信息含实体/字段/角色拼接、formatters 多职责耦合、申请状态机跨多处定义、枚举前后端重复定义 ✅

- 执行状态：完成（API 冒烟 29/29 通过，浏览器验证通过，docker down -v 无残留）

- 执行状态：完成（API 冒烟 29/29 通过，浏览器验证通过，docker down -v 无残留）
