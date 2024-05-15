package status

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	stack "github.com/rhobs/observability-operator/pkg/apis/monitoring/v1alpha1"
)

const (
	AvailableReason                = "MonitoringStackAvailable"
	ReconciledReason               = "MonitoringStackReconciled"
	FailedToReconcileReason        = "FailedToReconcile"
	PrometheusNotAvailable         = "PrometheusNotAvailable"
	PrometheusNotReconciled        = "PrometheusNotReconciled"
	PrometheusDegraded             = "PrometheusDegraded"
	ResourceSelectorIsNil          = "ResourceSelectorNil"
	CannotReadPrometheusConditions = "Cannot read Prometheus status conditions"
	AvailableMessage               = "Monitoring Stack is available"
	SuccessfullyReconciledMessage  = "Monitoring Stack is successfully reconciled"
	ResourceSelectorIsNilMessage   = "No resources will be discovered, ResourceSelector is nil"
	ResourceDiscoveryOnMessage     = "Resource discovery is operational"
	NoReason                       = "None"
)

type Condition interface {
	GetStackCondition() stack.Condition
}

// UpdateAvailable gets existing "Available" condition and updates its parameters
// based on the operand "Available" condition
func UpdateAvailable[C Condition](conditions []C) stack.Condition {
	ac := stack.Condition{
		Type:               stack.AvailableCondition,
		Status:             stack.ConditionUnknown,
		Reason:             NoReason,
		LastTransitionTime: metav1.Now(),
	}

	for _, condition := range conditions {
		availableCondition := condition.GetStackCondition()
		ac = updateStackCondition(ac, availableCondition)
	}
	return ac
}

func updateStackCondition(oldCondition, newCondition stack.Condition) stack.Condition {
	// update oldCondition with new, e.g. a new available condition doesn't change
	// an old degraded condtition
	return newCondition
}
