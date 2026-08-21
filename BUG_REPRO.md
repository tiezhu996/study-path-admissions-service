# BUG_REPRO

- Bug: nil相关问题
- 现象: 材料进度一打开就 panic，新建材料也崩；先别改代码，帮我查清楚崩在哪。
- 根因线索: 文件: internal/service/material_service.go, internal/handler/material_handler.go 符号: MaterialService.Progress, MaterialHandler.Create 机制: progressCache 和 summaryCache 是未初始化的 nil map，首次写入触发 assignment to entry in nil map
- 触发方式: 按题面步骤复现；评分测试见 verify_cmds。
