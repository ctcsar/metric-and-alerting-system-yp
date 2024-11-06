package agent

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAgentFlags(t *testing.T) {
	t.Run("SetAgentFlags", func(t *testing.T) {
		f := &agentFlags{}
		f.SetAgentFlags()

		assert.Equal(t, "localhost:8080", f.url)
		assert.Equal(t, 10, f.sendTime)
		assert.Equal(t, 2, f.getMetricTime)
	})

	t.Run("GetURLForSend", func(t *testing.T) {
		f := &agentFlags{
			url: "localhost:8080",
		}

		assert.Equal(t, "localhost:8080", f.GetURLForSend())
	})

	t.Run("GetSendDuration", func(t *testing.T) {
		f := &agentFlags{
			sendTime: 10,
		}

		assert.Equal(t, time.Duration(10), f.GetSendDuration())
	})

	t.Run("GetMetricsGetDuration", func(t *testing.T) {
		f := &agentFlags{
			getMetricTime: 2,
		}

		assert.Equal(t, time.Duration(2), f.GetMetricsGetDuration())
	})
}
