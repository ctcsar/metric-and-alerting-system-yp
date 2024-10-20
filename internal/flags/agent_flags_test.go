package flags

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAgentFlags(t *testing.T) {
	t.Run("SetAgentFlags", func(t *testing.T) {
		f := &agentFlags{}
		f.SetAgentFlags()

		assert.Equal(t, "localhost:8080", f.URL)
		assert.Equal(t, 10, f.SendTime)
		assert.Equal(t, 2, f.GetMetricTime)
	})

	t.Run("GetURLForSend", func(t *testing.T) {
		f := &agentFlags{
			URL: "localhost:8080",
		}

		assert.Equal(t, "localhost:8080", f.GetURLForSend())
	})

	t.Run("GetSendDuration", func(t *testing.T) {
		f := &agentFlags{
			SendTime: 10,
		}

		assert.Equal(t, time.Duration(10), f.GetSendDuration())
	})

	t.Run("GetMetricsGetDuration", func(t *testing.T) {
		f := &agentFlags{
			GetMetricTime: 2,
		}

		assert.Equal(t, time.Duration(2), f.GetMetricsGetDuration())
	})
}
