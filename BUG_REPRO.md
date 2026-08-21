# BUG_REPRO

- Bug: 其他问题
- 现象: 把选校方案从已发出改成已接受总是不成功，接口还返回成功；先别改代码，帮我查为什么。
- 根因线索: 文件: internal/constants/recommendation.go, internal/service/recommendation_service.go, internal/repository/recommendation_repository.go, internal/handler/recommendation_handler.go 符号: RecommendationService.UpdateStatus, RecommendationRepository.UpdateStatus, RecommendationHandler.UpdateStatus 机制: 状态转换表漏 sent→accepted 边，仓库写回旧状态，handler 把失败当成功返回
- 触发方式: 按题面步骤复现；评分测试见 verify_cmds。
