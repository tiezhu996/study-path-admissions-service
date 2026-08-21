# BUG_REPRO

- Bug: defer相关问题
- 现象: 给几个学生群发消息，发完数据库一条都没有，帮我修好。
- 根因线索: 文件: internal/repository/message_repository.go, internal/service/message_service.go 符号: MessageRepository.CreateMany, MessageService.SendBatch 机制: 循环内 defer tx.Rollback 堆积到函数尾，已提交消息被回滚，SendBatch 命名返回值被 defer 清零
- 触发方式: 按题面步骤复现；评分测试见 verify_cmds。
