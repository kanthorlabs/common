package entities

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestEvaluationValidateOnGrant(t *testing.T) {
	evaluation := &Evaluation{
		Tenant:   uuid.NewString(),
		Username: uuid.NewString(),
	}
	require.ErrorContains(t, EvaluationValidateOnGrant(evaluation), "GATEKEEPER.EVALUATION.")
}

func TestEvaluationValidateOnRevoke(t *testing.T) {
	evaluation := &Evaluation{
		Tenant: uuid.NewString(),
	}
	require.ErrorContains(t, EvaluationValidateOnRevoke(evaluation), "GATEKEEPER.EVALUATION.")
}

func TestEvaluationValidateOnEnforce(t *testing.T) {
	evaluation := &Evaluation{
		Tenant: uuid.NewString(),
	}
	require.ErrorContains(t, EvaluationValidateOnEnforce(evaluation), "GATEKEEPER.EVALUATION.")
}
