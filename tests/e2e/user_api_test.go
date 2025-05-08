package e2e_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("用户 API", func() {
	var (
		tenantID   string
		adminEmail string
		adminToken string
		userID     string // 在用户配置文件测试中使用
		tenantName string // 添加租户名称变量
	)

	BeforeEach(func() {
		// 创建测试租户
		tenantID, adminEmail, tenantName = setupTestTenant()

		// 获取管理员令牌
		adminToken = getAuthToken(adminEmail, "Password123!", tenantName)

		// 使 linter 满意，声明我们会在后面的测试中使用 userID
		_ = userID
	})

	Context("用户注册和登录", func() {
		It("应该能成功注册新用户", func() {
			// 准备注册数据
			username := fmt.Sprintf("test-user-%d", GinkgoRandomSeed())
			email := fmt.Sprintf("user-%d@example.com", GinkgoRandomSeed())

			registerReq := map[string]string{
				"username":   username,
				"email":      email,
				"password":   "Password123!",
				"given_name": "Test User",
				"tenant_id":  tenantID, // 添加租户 ID
			}

			// 打印请求
			fmt.Printf("User register request: %+v\n", registerReq)

			// 发送注册请求
			w := makeRequest(http.MethodPost, "/users/register", registerReq, nil)

			// 打印响应
			fmt.Printf("User register response: %d\n%s\n", w.Code, w.Body.String())

			Expect(w.Code).To(Equal(http.StatusCreated))

			// 验证响应
			var resp struct {
				ID       string `json:"id"`
				Username string `json:"username"`
				Email    string `json:"email"`
			}
			parseResponse(w, &resp)

			Expect(resp.Username).To(Equal(username))
			Expect(resp.Email).To(Equal(email))
		})

		It("应该能成功登录", func() {
			// 准备用户数据
			username := fmt.Sprintf("login-test-user-%d", GinkgoRandomSeed())
			email := fmt.Sprintf("login-user-%d@example.com", GinkgoRandomSeed())
			password := "Password123!"

			// 先注册用户
			registerReq := map[string]string{
				"username":   username,
				"email":      email,
				"password":   password,
				"given_name": "Login Test User",
				"tenant_id":  tenantID, // 添加租户 ID
			}

			w := makeRequest(http.MethodPost, "/users/register", registerReq, nil)
			Expect(w.Code).To(Equal(http.StatusCreated))

			// 然后尝试登录
			loginReq := map[string]string{
				"username":    username,
				"password":    password,
				"tenant_name": tenantName, // 添加租户名称
			}

			w = makeRequest(http.MethodPost, "/users/login", loginReq, nil)
			Expect(w.Code).To(Equal(http.StatusOK))

			// 验证获取到了令牌
			var resp struct {
				Token string `json:"token"`
			}
			parseResponse(w, &resp)

			Expect(resp.Token).NotTo(BeEmpty())
		})

		It("不正确的密码应该登录失败", func() {
			// 准备用户数据
			username := fmt.Sprintf("failed-login-test-user-%d", GinkgoRandomSeed())
			email := fmt.Sprintf("failed-login-user-%d@example.com", GinkgoRandomSeed())
			password := "Password123!"

			// 先注册用户
			registerReq := map[string]string{
				"username":   username,
				"email":      email,
				"password":   password,
				"given_name": "Failed Login User",
				"tenant_id":  tenantID, // 添加租户 ID
			}

			w := makeRequest(http.MethodPost, "/users/register", registerReq, nil)
			Expect(w.Code).To(Equal(http.StatusCreated))

			// 然后尝试使用错误密码登录
			loginReq := map[string]string{
				"username":    username,
				"password":    "WrongPassword123!",
				"tenant_name": tenantName, // 添加租户名称
			}

			w = makeRequest(http.MethodPost, "/users/login", loginReq, nil)
			Expect(w.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Context("用户配置文件", func() {
		var userToken string

		BeforeEach(func() {
			// 创建一个测试用户
			userID = setupTestUser(tenantID, adminToken)

			// 使用管理员 API 创建的用户默认密码为 "Password123!"
			username := fmt.Sprintf("user-%d@example.com", GinkgoRandomSeed()-1)
			userToken = getAuthToken(username, "Password123!", tenantName)
		})

		It("应该能获取用户配置文件", func() {
			headers := map[string]string{
				"Authorization": "Bearer " + userToken,
			}

			w := makeRequest(http.MethodGet, "/users/me", nil, headers)
			Expect(w.Code).To(Equal(http.StatusOK))

			var resp map[string]interface{}
			parseResponse(w, &resp)

			// 验证响应包含关键字段
			Expect(resp).To(HaveKey("id"))
			Expect(resp).To(HaveKey("username"))
			Expect(resp).To(HaveKey("email"))
		})

		It("应该能更新用户配置文件", func() {
			updateReq := map[string]string{
				"given_name": "Updated Name",
			}

			headers := map[string]string{
				"Authorization": "Bearer " + userToken,
			}

			w := makeRequest(http.MethodPut, "/users/me", updateReq, headers)
			Expect(w.Code).To(Equal(http.StatusOK))

			var resp map[string]interface{}
			parseResponse(w, &resp)

			// 验证更新成功
			Expect(resp["given_name"]).To(Equal("Updated Name"))
		})

		It("应该能更新密码", func() {
			updateReq := map[string]string{
				"old_password": "Password123!",
				"new_password": "NewPassword123!",
			}

			headers := map[string]string{
				"Authorization": "Bearer " + userToken,
			}

			w := makeRequest(http.MethodPut, "/users/me/password", updateReq, headers)
			Expect(w.Code).To(Equal(http.StatusOK))

			// 使用新密码尝试登录
			username := fmt.Sprintf("user-%d@example.com", GinkgoRandomSeed()-1)
			newToken := getAuthToken(username, "NewPassword123!", tenantName)
			Expect(newToken).NotTo(BeEmpty())
		})
	})
})
