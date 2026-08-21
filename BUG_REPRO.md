# BUG_REPRO

- Bug: error异常错误
- 现象: 回滚不存在的文档版本返回 500，单独查不存在的文档也是 500，请帮我修好，404 才对。
- 根因线索: 文件: internal/repository/document_repository.go, internal/repository/document_version_repository.go, internal/service/document_service.go 符号: DocumentRepository.FindByID, DocumentVersionRepository.FindByDocumentAndVersion, DocumentService.Rollback 机制: %v 包装丢失 ErrNotFound，errors.Is 永不命中，文档/版本不存在被当成 500
- 触发方式: 按题面步骤复现；评分测试见 verify_cmds。
