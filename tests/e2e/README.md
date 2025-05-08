# API E2E 测试

本目录包含使用 Ginkgo 和 Gomega 框架编写的 API 端到端测试，用于验证 API 的功能和集成是否正常工作。

## 测试结构

测试按照 API 功能进行组织：

- `user_api_test.go`: 用户 API 测试（注册、登录、个人资料等）
- `tenant_api_test.go`: 租户 API 测试（租户注册、更新等）
- `admin_api_test.go`: 管理员 API 测试（系统管理员和租户管理员）

## 运行测试

运行所有 E2E 测试：

```bash
make test-e2e
```

或者直接运行：

```bash
cd tests/e2e
go test -v
```

使用特定标签运行测试：

```bash
cd tests/e2e
go test -v -ginkgo.focus="用户 API"
```

## 测试实现细节

这些测试用例使用内存数据库（SQLite）运行，以确保测试速度和隔离性。它们：

1. 初始化测试环境并设置内存数据库
2. 创建测试路由器，类似于实际应用程序
3. 为每个测试用例创建必要的测试数据
4. 验证 API 的行为和响应

这些测试不使用外部依赖项，因此可以在任何环境中运行，无需设置外部服务。

## 编写新测试

要添加新的测试用例，请按照以下模式：

1. 为 API 分组创建一个新的测试文件（如果尚不存在）
2. 使用 `Describe` 和 `Context` 组织测试结构
3. 使用 `BeforeEach` 设置必要的测试数据
4. 使用 `It` 函数编写测试用例
5. 使用 `helpers_test.go` 中提供的辅助函数执行请求和验证响应

例如：

```go
var _ = Describe("新的 API 功能", func() {
    BeforeEach(func() {
        // 设置测试数据
    })

    It("应该按预期工作", func() {
        // 调用 API
        w := makeRequest(...)

        // 验证结果
        Expect(w.Code).To(Equal(http.StatusOK))
    })
})
```
