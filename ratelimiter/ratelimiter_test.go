package ratelimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func BasicTest(t *testing.T) {
	limiter := NewInMemoryRateLimiter(2, time.Duration(5*int(time.Second)))
	test_ip := `test_ip`
	t.Run("New Request should be added properly", func(t *testing.T) {
		isAllowed := limiter.IsAllowed(test_ip)
		require.Equal(t, true, isAllowed)
		ipInfo := limiter.getRateLimitData(test_ip)
		require.Equal(t, test_ip, ipInfo.IP)
		require.Equal(t, 1, ipInfo.RequestCount)
		require.Equal(t, false, ipInfo.Blocked)
	})

	t.Run("IP shoudn't be blocked", func(test *testing.T) {
		isAllowed := limiter.IsAllowed(test_ip)
		require.Equal(test, true, isAllowed)
	})
	t.Run("Ip should be blocked", func(test *testing.T) {
		isAllowed := limiter.IsAllowed(test_ip)
		require.Equal(test, false, isAllowed)
	}) // You can add more test cases here
}

func Test_RateLimiter(test *testing.T) {
	limit := 2
	interval := time.Duration(5 * time.Second)
	limiter := NewInMemoryRateLimiter(limit, interval)
	test_ip := `test-ip-2`

	test.Run("Checking block and reset flow", func(test *testing.T) {
		for i := 1; i <= limit; i++ {
			res := limiter.IsAllowed(test_ip)
			require.Equal(test, true, res)
		}
		require.Equal(test, false, limiter.IsAllowed(test_ip))
		time.Sleep(interval)
		for i := 1; i <= limit; i++ {
			res := limiter.IsAllowed(test_ip)
			require.Equal(test, true, res)
		}
		require.Equal(test, false, limiter.IsAllowed(test_ip))
		time.Sleep(interval)
		for i := 1; i <= limit; i++ {
			res := limiter.IsAllowed(test_ip)
			require.Equal(test, true, res)
		}
		require.Equal(test, false, limiter.IsAllowed(test_ip))

	})
}

// func TestAdd(t *testing.T) {
// 	cases := struct {
// 		name     string
// 		a, b     int
// 		expected int
// 	}{
// 		{"2 + 3", 3, 3, 5},
// 		{"-1 + 1", -1, 1, 0},
// 		{"0 + 0", 0, 0, 0},
// 	}

// 	for _, c := range cases {
// 		t.Run(c.name, func(t *testing.T) {
// 			result := Add(c.a, c.b)
// 			require.Equal(t, c.expected, result)
// 		})
// 	}
// }
