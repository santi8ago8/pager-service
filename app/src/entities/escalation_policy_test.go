package entities

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/santi8ago8/pager-service/app/src/constants"
)

func TestEscalationPolicy_GetLevelToNotify(t *testing.T) {
	level0 := &Level{
		ID: 0,
		Targets: []*Target{
			{
				ID:          uuid.New().String(),
				Type:        constants.TargetTypeSms,
				PhoneNumber: "+54923959292",
			},
		},
	}
	level1 := &Level{
		ID: 1,
		Targets: []*Target{
			{
				ID:          uuid.New().String(),
				Type:        constants.TargetTypeSms,
				PhoneNumber: "+3493223463",
			},
			{
				ID:    uuid.New().String(),
				Type:  constants.TargetTypeEmail,
				Email: "test@test.com",
			},
		},
	}
	escPolicy := EscalationPolicy{
		ID:                 uuid.New().String(),
		MonitoredServiceID: "FAKE_MON_SERVICE_UUID",
		Levels: []*Level{
			level0,
			level1,
		},
	}

	tests := []struct {
		name             string
		escalationPolicy EscalationPolicy
		levelNumber      int
		want             *Level
	}{
		{
			name:             "Zero",
			escalationPolicy: escPolicy,
			levelNumber:      0,
			want:             level0,
		},
		{
			name:             "One",
			escalationPolicy: escPolicy,
			levelNumber:      1,
			want:             level1,
		},
		{
			name:             "Nil",
			escalationPolicy: escPolicy,
			levelNumber:      2,
			want:             nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.escalationPolicy.GetLevelToNotify(tt.levelNumber); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EscalationPolicy.GetLevelToNotify() = %v, want %v", got, tt.want)
			}
		})
	}
}
