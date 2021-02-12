package handler

import (
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/m3o/services/v1api/proto/fakes"
	mstore "github.com/micro/micro/v3/service/store"
	"github.com/micro/micro/v3/service/store/memory"

	. "github.com/onsi/gomega"
)

func TestIsTimeForReset(t *testing.T) {
	g := NewWithT(t)

	tcs := []struct {
		name     string
		freq     resetFrequency
		timeStr  string
		expected bool
	}{
		{
			name:     "never",
			freq:     Never,
			timeStr:  "20210201 15:04:05",
			expected: false,
		},
		{
			name:     "neverAgain",
			freq:     Never,
			timeStr:  "20210205 15:04:05",
			expected: false,
		},
		{
			name:     "daily",
			freq:     Daily,
			timeStr:  "20210201 15:04:05",
			expected: true,
		},
		{
			name:     "dailyAgain",
			freq:     Daily,
			timeStr:  "20210205 15:04:05",
			expected: true,
		},
		{
			name:     "monthlyYes",
			freq:     Monthly,
			timeStr:  "20210201 15:04:05",
			expected: true,
		},
		{
			name:     "monthlyNo",
			freq:     Monthly,
			timeStr:  "20210205 15:04:05",
			expected: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tObj, _ := time.Parse("20060102 15:04:05", tc.timeStr)
			g.Expect(isTimeForReset(tc.freq, tObj)).To(Equal(tc.expected))
		})
	}
}

func TestResetQuotas(t *testing.T) {
	g := NewWithT(t)

	v1serv := fakes.FakeV1Service{}
	db, redisMock := redismock.NewClientMock()
	mstore.DefaultStore = memory.NewStore()

	qsvc := &Quota{
		v1Svc: &v1serv,
		c: counter{
			redisClient: db,
		},
	}
	g.Expect(qsvc.writeQuota(&quota{
		ID:             "hello_free",
		Limit:          10,
		ResetFrequency: Daily,
		Path:           "/hello/",
	})).To(BeNil())
	g.Expect(qsvc.writeMapping(&mapping{
		UserID:    "user_1234",
		Namespace: "my-user-ns",
		QuotaID:   "hello_free",
	})).To(BeNil())

	redisMock.ExpectSet(prefixCounter+":my-user-ns:user_1234:/hello/", 0, 0).SetVal("0")
	qsvc.ResetQuotasCron()
	g.Expect(redisMock.ExpectationsWereMet()).To(BeNil())
	g.Expect(v1serv.UpdateAllowedPathsCallCount()).To(Equal(1))
}
