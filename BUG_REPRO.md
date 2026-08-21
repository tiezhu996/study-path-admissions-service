# BUG_REPRO

- Bug: error异常错误
- 现象: 不存在的账号登录返回 500，查用户也 500；这不该是服务器问题，帮我修好。
- 根因线索: 文件: internal/repository/user_repository.go, internal/service/user_service.go 符号: UserRepository.FindByUsername, UserService.Login, UserService.GetByID 机制: %v 断链使 errors.Is(ErrNotFound) 失效，登录/查询不存在用户误判为 500
- 触发方式: 按题面步骤复现；评分测试见 verify_cmds。
