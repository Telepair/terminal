package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// makeRequest 创建并执行一个 HTTP 请求
func makeRequest(method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	// 创建请求体
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		Expect(err).NotTo(HaveOccurred())
		reqBody = bytes.NewBuffer(jsonData)
	}

	// 创建请求
	req, err := http.NewRequest(method, path, reqBody)
	Expect(err).NotTo(HaveOccurred())

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 执行请求
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	return w
}

// parseResponse 解析响应数据
func parseResponse(w *httptest.ResponseRecorder, resp interface{}) {
	Expect(w.Body).NotTo(BeNil())
	err := json.Unmarshal(w.Body.Bytes(), resp)
	Expect(err).NotTo(HaveOccurred())
}

// getAuthToken 登录并获取认证令牌
func getAuthToken(username, password string, tenantName string) string {
	// 创建登录请求，包含租户名称
	loginReq := map[string]string{
		"username":    username,
		"password":    password,
		"tenant_name": tenantName,
	}

	// 打印登录请求信息
	fmt.Printf("Login request for user: %s in tenant: %s\n", username, tenantName)

	w := makeRequest(http.MethodPost, "/users/login", loginReq, nil)

	// 打印响应结果
	fmt.Printf("Login response status: %d\n", w.Code)
	fmt.Printf("Login response body: %s\n", w.Body.String())

	Expect(w.Code).To(Equal(http.StatusOK))

	var resp struct {
		Token string `json:"token"`
	}
	parseResponse(w, &resp)

	Expect(resp.Token).NotTo(BeEmpty())

	return resp.Token
}

// setupTestTenant 创建测试租户
func setupTestTenant() (string, string, string) {
	// 生成随机租户名称以避免冲突
	tenantName := fmt.Sprintf("test-tenant-%d", GinkgoRandomSeed())
	adminEmail := fmt.Sprintf("admin-%d@example.com", GinkgoRandomSeed())

	// 注册租户
	tenantReq := map[string]interface{}{
		"name":           tenantName,
		"admin_email":    adminEmail,
		"admin_password": "Password123!",
		"given_name":     "Test Admin",
	}

	// 打印出请求体以便调试
	fmt.Printf("Tenant request body: %+v\n", tenantReq)

	w := makeRequest(http.MethodPost, "/tenants/register", tenantReq, nil)

	// 打印响应状态码和响应体以便调试
	fmt.Printf("Tenant register status: %d\n", w.Code)
	fmt.Printf("Tenant register response: %s\n", w.Body.String())

	Expect(w.Code).To(Equal(http.StatusCreated))

	var resp struct {
		ID string `json:"id"`
	}
	parseResponse(w, &resp)

	return resp.ID, adminEmail, tenantName
}

// setupTestUser 创建测试用户
func setupTestUser(tenantID, token string) string {
	// 生成随机用户名以避免冲突
	username := fmt.Sprintf("test-user-%d", GinkgoRandomSeed())
	email := fmt.Sprintf("user-%d@example.com", GinkgoRandomSeed())

	// 创建用户
	userReq := map[string]string{
		"username":   username,
		"email":      email,
		"password":   "Password123!",
		"given_name": "Test User",
	}

	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	path := fmt.Sprintf("/tenants/%s/admin/users", tenantID)
	w := makeRequest(http.MethodPost, path, userReq, headers)
	Expect(w.Code).To(Equal(http.StatusCreated))

	var resp struct {
		ID string `json:"id"`
	}
	parseResponse(w, &resp)

	return resp.ID
}
