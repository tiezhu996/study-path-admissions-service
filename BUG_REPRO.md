# BUG_REPRO

- Bug: context相关问题
- 现象: 提醒扫描任务取消不掉，cancel 后还一直跑；先别改代码，帮我定位取消为什么没生效。
- 根因线索: 文件: internal/service/notification_service.go, internal/repository/timeline_node_repository.go 符号: NotificationService.RunReminders, NotificationService.ScanDeadlinesContext, TimelineNodeRepository.ListDue, ReminderCancelError 机制: RunReminders 传入 context.Background，循环与 ScanDeadlinesContext/ListDue 都忽略 ctx.Done，取消信号不传播
- 触发方式: 按题面步骤复现；评分测试见 verify_cmds。
