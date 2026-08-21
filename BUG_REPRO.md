# BUG_REPRO

- Bug: slice相关问题
- 现象: 院校列表先查美国再查英国，之前美国的记录被改成英国的了，帮我修好。
- 根因线索: 文件: internal/repository/university_repository.go, internal/service/university_service.go 符号: UniversityRepository.List, UniversityRepository.ListByIDs, UniversityService.List 机制: s[:0] 复用底层数组，后续查询原地覆写之前返回的院校切片
- 触发方式: 按题面步骤复现；评分测试见 verify_cmds。
