package e2e_test

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("管理员 API", func() {
	var (
		sysAdminToken    string
		tenantID         string
		adminEmail       string
		tenantAdminToken string
		tenantName       string
		systemTenantName string
	)

	BeforeEach(func() {
		// 先创建一个系统租户
		sysID := uuid.New()
		systemTenantName = fmt.Sprintf("system-%d", GinkgoRandomSeed())
		systemTenant, err := client.Tenant.Create().
			SetID(sysID).
			SetName(systemTenantName).
			SetEnabled(true).
			SetAdminEmail("system@example.com"). // 添加管理员邮箱
			Save(ctx)
		Expect(err).NotTo(HaveOccurred())

		// 使用 Ent 客户端直接创建系统管理员
		sysAdmin, err := client.User.Create().
			SetUsername(fmt.Sprintf("sysadmin-%d", GinkgoRandomSeed())).
			SetEmail(fmt.Sprintf("sysadmin-%d@example.com", GinkgoRandomSeed())).
			SetPasswordBcrypt("$2a$10$DwPN5Y3aTk5.He/lJh/XG.KaHSXviCuB/u0/YSjSImkzKhVMQJ/4C"). // "Password123!"
			SetIsSuperuser(true).
			SetEmailVerified(true).
			SetEnabled(true).
			SetTenantID(systemTenant.ID). // 设置系统租户 ID
			Save(ctx)
		Expect(err).NotTo(HaveOccurred())

		// 模拟系统管理员登录
		username := sysAdmin.Username
		sysAdminToken = getAuthToken(username, "Password123!", systemTenantName)

		// 创建测试租户
		tenantID, adminEmail, tenantName = setupTestTenant()

		// 获取租户管理员令牌
		tenantAdminToken = getAuthToken(adminEmail, "Password123!", tenantName)
	})

	Context("系统管理员权限", func() {
		It("系统管理员应该能列出所有租户", func() {
			headers := map[string]string{
				"Authorization": "Bearer " + sysAdminToken,
			}

			w := makeRequest(http.MethodGet, "/admin/tenants", nil, headers)
			Expect(w.Code).To(Equal(http.StatusOK))

			var resp struct {
				Tenants []map[string]interface{} `json:"tenants"`
			}
			parseResponse(w, &resp)

			// 验证至少有我们创建的租户
			Expect(resp.Tenants).NotTo(BeEmpty())

			// 查找我们创建的租户
			var found bool
			for _, tenant := range resp.Tenants {
				if tenant["id"] == tenantID {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		It("系统管理员应该能禁用租户", func() {
			headers := map[string]string{
				"Authorization": "Bearer " + sysAdminToken,
			}

			path := fmt.Sprintf("/admin/tenants/%s/disable", tenantID)
			w := makeRequest(http.MethodPut, path, nil, headers)
			Expect(w.Code).To(Equal(http.StatusOK))

			// 验证租户管理员无法再使用其令牌进行操作
			tenantAdminHeaders := map[string]string{
				"Authorization": "Bearer " + tenantAdminToken,
			}

			tenantPath := fmt.Sprintf("/tenants/%s/admin", tenantID)
			w = makeRequest(http.MethodPut, tenantPath, map[string]string{"name": "new-name"}, tenantAdminHeaders)
			Expect(w.Code).To(Equal(http.StatusUnauthorized))
		})

		It("系统管理员应该能启用租户", func() {
			headers := map[string]string{
				"Authorization": "Bearer " + sysAdminToken,
			}

			// 先禁用租户
			disablePath := fmt.Sprintf("/admin/tenants/%s/disable", tenantID)
			w := makeRequest(http.MethodPut, disablePath, nil, headers)
			Expect(w.Code).To(Equal(http.StatusOK))

			// 再启用租户
			enablePath := fmt.Sprintf("/admin/tenants/%s/enable", tenantID)
			w = makeRequest(http.MethodPut, enablePath, nil, headers)
			Expect(w.Code).To(Equal(http.StatusOK))

			// 验证租户管理员现在可以使用其令牌进行操作
			tenantAdminHeaders := map[string]string{
				"Authorization": "Bearer " + tenantAdminToken,
			}

			tenantPath := fmt.Sprintf("/tenants/%s/admin", tenantID)
			w = makeRequest(http.MethodPut, tenantPath, map[string]string{"name": "reactivated-name"}, tenantAdminHeaders)
			Expect(w.Code).To(Equal(http.StatusOK))
		})
	})

	Context("租户管理员权限", func() {
		It("租户管理员应该能看到其租户下的所有用户", func() {
			// 创建几个测试用户
			for i := 0; i < 3; i++ {
				setupTestUser(tenantID, tenantAdminToken)
			}

			headers := map[string]string{
				"Authorization": "Bearer " + tenantAdminToken,
			}

			path := fmt.Sprintf("/tenants/%s/admin/users", tenantID)
			w := makeRequest(http.MethodGet, path, nil, headers)
			Expect(w.Code).To(Equal(http.StatusOK))

			var resp struct {
				Users []map[string]interface{} `json:"users"`
			}
			parseResponse(w, &resp)

			// 应该能看到至少 4 个用户（租户管理员 + 3 个测试用户）
			Expect(len(resp.Users)).To(BeNumerically(">=", 4))
		})

		It("租户管理员应该能将用户提升为管理员", func() {
			// 创建测试用户
			userID := setupTestUser(tenantID, tenantAdminToken)

			headers := map[string]string{
				"Authorization": "Bearer " + tenantAdminToken,
			}

			// 将用户设置为管理员
			path := fmt.Sprintf("/tenants/%s/admin/users/%s/admin", tenantID, userID)
			w := makeRequest(http.MethodPut, path, nil, headers)
			Expect(w.Code).To(Equal(http.StatusOK))

			// 获取用户列表并验证用户现在是管理员
			usersPath := fmt.Sprintf("/tenants/%s/admin/users", tenantID)
			w = makeRequest(http.MethodGet, usersPath, nil, headers)
			Expect(w.Code).To(Equal(http.StatusOK))

			var resp struct {
				Users []map[string]interface{} `json:"users"`
			}
			parseResponse(w, &resp)

			// 查找我们的用户
			var found bool
			for _, user := range resp.Users {
				if user["id"] == userID {
					found = true
					// 验证用户是管理员
					Expect(user["is_tenant_admin"]).To(BeTrue())
					break
				}
			}
			Expect(found).To(BeTrue())
		})
	})
})
