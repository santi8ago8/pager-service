package entities

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/santi8ago8/pager-service/app/src/constants"
)

func TestNewAlert(t *testing.T) {
	type args struct {
		serviceID    string
		alertMessage string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Create ok",
			args: args{
				serviceID:    "FAKE_UUID_MONITORED_SERVICE",
				alertMessage: "Error rate > 10% in the las 15 minutes.",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAlert(tt.args.serviceID, tt.args.alertMessage); got == nil {
				t.Errorf("NewAlert() = %v", got)
			}
		})
	}
}

func TestAlert_Acknowledge(t *testing.T) {

	t.Run("Ok", func(t *testing.T) {
		alert := NewAlert("FAKE_UUID_MONITORED_SERVICE", "Error rate > 10% in the las 15 minutes.")
		alert.Acknowledge()
		assert.Equal(t, alert.Status, constants.AlertStatusAcknowledge)
	})
}

func TestAlert_GetNextLevelToNotify(t *testing.T) {

	t.Run("Zero", func(t *testing.T) {
		alert := NewAlert("FAKE_UUID_MONITORED_SERVICE", "Error rate > 10% in the las 15 minutes.")
		level := alert.GetNextLevelToNotify()
		assert.Equal(t, level, 0)
	})

	t.Run("One", func(t *testing.T) {
		alert := NewAlert("FAKE_UUID_MONITORED_SERVICE", "Error rate > 10% in the las 15 minutes.")
		tempLevel := 0
		alert.NotifiedLevelID = &tempLevel
		level := alert.GetNextLevelToNotify()
		assert.Equal(t, level, 1)
	})
}
