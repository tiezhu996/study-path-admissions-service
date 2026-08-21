# BUG_REPRO

- Bug: nil相关问题
- 现象: 存储没配好时上传接口直接 panic；先别改代码，帮我看看为什么会崩。
- 根因线索: 文件: internal/handler/upload_handler.go, internal/util/minio_client.go 符号: UploadHandler.Upload, MinIOClient.Upload 机制: nil 的 *MinIOClient 装进 ObjectStore 接口，typed-nil 判空失效，上传时 nil 指针解引用
- 触发方式: 按题面步骤复现；评分测试见 verify_cmds。
