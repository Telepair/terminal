package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/telepair/terminal/server"
)

var (
	serverAddr string
	serverCtx  context.Context
	cancel     context.CancelFunc
)

var _ = BeforeSuite(func() {
	serverAddr = "127.0.0.1:18089"
	serverCtx, cancel = context.WithCancel(context.Background())
	go func() {
		_ = server.StartServerWithContext(serverCtx, serverAddr)
	}()
	time.Sleep(500 * time.Millisecond) // Wait for server to start
})

var _ = AfterSuite(func() {
	cancel()
	time.Sleep(200 * time.Millisecond)
})

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server E2E Suite")
}

var _ = Describe("Server E2E", func() {
	It("should return health ok", func() {
		resp, err := http.Get(fmt.Sprintf("http://%s/api/health", serverAddr))
		Expect(err).To(BeNil())
		defer func() {
			_ = resp.Body.Close() // Ignore error for closing response body
		}()
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		// Parse JSON response and check status field
		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		Expect(err).To(BeNil())
		Expect(result).To(HaveKeyWithValue("status", "ok"))
	})

	It("should return version info", func() {
		resp, err := http.Get(fmt.Sprintf("http://%s/api/version", serverAddr))
		Expect(err).To(BeNil())
		defer func() {
			_ = resp.Body.Close() // Ignore error for closing response body
		}()
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		// You can add more JSON field checks here if needed
	})
})
