package e2e_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("租户 API", func() {
	Context("租户注册", func() {
		It("应该能成功注册新租户", func() {
			// 准备注册数据
			tenantName := fmt.Sprintf("test-tenant-%d", GinkgoRandomSeed())
			adminEmail := fmt.Sprintf("admin-%d@example.com", GinkgoRandomSeed())

			registerReq := map[string]string{
				"name":           tenantName,
				"admin_email":    adminEmail,
				"admin_password": "Password123!",
				"given_name":     "Test Admin",
			}

			// 发送注册请求
			w := makeRequest(http.MethodPost, "/tenants/register", registerReq, nil)
			Expect(w.Code).To(Equal(http.StatusCreated))

			// 验证响应
			var resp struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}
			parseResponse(w, &resp)

			Expect(resp.Name).To(Equal(tenantName))
			Expect(resp.ID).NotTo(BeEmpty())
		})

		It("重复的租户名称应该注册失败", func() {
			// 准备第一个租户的注册数据
			tenantName := fmt.Sprintf("duplicate-tenant-%d", GinkgoRandomSeed())
			adminEmail1 := fmt.Sprintf("admin1-%d@example.com", GinkgoRandomSeed())

			// 注册第一个租户
			registerReq1 := map[string]string{
				"name":           tenantName,
				"admin_email":    adminEmail1,
				"admin_password": "Password123!",
				"given_name":     "Test Admin 1",
			}

			w := makeRequest(http.MethodPost, "/tenants/register", registerReq1, nil)
			Expect(w.Code).To(Equal(http.StatusCreated))

			// 尝试注册同名租户
			adminEmail2 := fmt.Sprintf("admin2-%d@example.com", GinkgoRandomSeed())
			registerReq2 := map[string]string{
				"name":           tenantName,
				"admin_email":    adminEmail2,
				"admin_password": "Password123!",
				"given_name":     "Test Admin 2",
			}

			w = makeRequest(http.MethodPost, "/tenants/register", registerReq2, nil)
			// 应该返回冲突错误
			Expect(w.Code).To(Equal(http.StatusConflict))
		})
	})

	Context("租户管理", func() {
		var (
			tenantID   string
			adminEmail string
			adminToken string
			tenantName string
		)

		BeforeEach(func() {
			// 创建测试租户
			tenantID, adminEmail, tenantName = setupTestTenant()

			// 获取管理员令牌
			adminToken = getAuthToken(adminEmail, "Password123!", tenantName)
		})

		It("租户管理员应该能更新租户信息", func() {
			// 准备更新数据
			updateReq := map[string]string{
				"name": fmt.Sprintf("updated-tenant-%d", GinkgoRandomSeed()),
			}

			headers := map[string]string{
				"Authorization": "Bearer " + adminToken,
			}

			// 发送更新请求
			path := fmt.Sprintf("/tenants/%s/admin", tenantID)
			w := makeRequest(http.MethodPut, path, updateReq, headers)
			Expect(w.Code).To(Equal(http.StatusOK))

			// 验证响应
			var resp struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}
			parseResponse(w, &resp)

			Expect(resp.ID).To(Equal(tenantID))
			Expect(resp.Name).To(Equal(updateReq["name"]))
		})

		It("租户管理员应该能更新管理员邮箱", func() {
			// 准备更新数据
			newEmail := fmt.Sprintf("new-admin-%d@example.com", GinkgoRandomSeed())
			updateReq := map[string]string{
				"admin_email": newEmail,
			}

			headers := map[string]string{
				"Authorization": "Bearer " + adminToken,
			}

			// 发送更新请求
			path := fmt.Sprintf("/tenants/%s/admin/admin-email", tenantID)
			w := makeRequest(http.MethodPut, path, updateReq, headers)
			Expect(w.Code).To(Equal(http.StatusOK))

			// 验证响应
			var resp map[string]interface{}
			parseResponse(w, &resp)

			Expect(resp).To(HaveKey("admin_email"))
			Expect(resp["admin_email"]).To(Equal(newEmail))
		})
	})
})
