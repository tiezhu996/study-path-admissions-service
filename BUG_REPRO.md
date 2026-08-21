# BUG_REPRO

- Bug: concurrency并发问题
- 现象: 好几个人同时看工作台统计，日志里报 map 并发写，帮我修好。
- 根因线索: 文件: internal/service/stats_service.go, internal/handler/dashboard_handler.go 符号: FillAppStats, DashboardHandler.Stats 机制: 共享 lastStats 和 ByStatus map 无锁原地复用，并发统计请求触发 data race
- 触发方式: 按题面步骤复现；评分测试见 verify_cmds。
