# BUG_REPRO

- Bug: concurrency并发问题
- 现象: 申请列表和状态按钮一起点，服务报 map 并发读写，帮我修好。
- 根因线索: 文件: internal/service/application_service.go, internal/repository/application_project_repository.go 符号: ApplicationService.List, ApplicationService.UpdateStatus, ApplicationProjectRepository.ListByStudent 机制: 进程内缓存 map 无锁读写，仓库 listBuffer 原地复用，并发查看列表和修改状态触发 data race，旧快照被后续查询污染
- 触发方式: 按题面步骤复现；评分测试见 verify_cmds。
