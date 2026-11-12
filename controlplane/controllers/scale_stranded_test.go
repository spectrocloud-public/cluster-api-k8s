package controllers

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util/collections"

	controlplanev1 "github.com/canonical/cluster-api-k8s/controlplane/api/v1beta3"
	"github.com/canonical/cluster-api-k8s/pkg/ck8s"
)

var strandedTestEpoch = time.Date(2026, 11, 6, 6, 0, 0, 0, time.UTC)

// strandedTestMachine returns a control plane machine created `age` minutes after the epoch, in failure
// domain fd, that is Ready unless nodeDeleted is set, in which case it reports what the Cluster API
// Machine controller reports once the Node is gone: Ready=False and NodeHealthy=False/NodeDeleted.
func strandedTestMachine(name, fd string, age int, nodeDeleted bool) *clusterv1.Machine {
	m := &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			Namespace:         "default",
			CreationTimestamp: metav1.NewTime(strandedTestEpoch.Add(time.Duration(age) * time.Minute)),
		},
		Spec: clusterv1.MachineSpec{FailureDomain: fd},
		Status: clusterv1.MachineStatus{
			NodeRef: clusterv1.MachineNodeReference{Name: name + "-node"},
			Conditions: []metav1.Condition{
				{Type: clusterv1.MachineReadyCondition, Status: metav1.ConditionTrue, Reason: clusterv1.MachineReadyReason},
				{Type: clusterv1.MachineNodeHealthyCondition, Status: metav1.ConditionTrue, Reason: clusterv1.MachineNodeHealthyReason},
			},
		},
	}
	if nodeDeleted {
		m.Status.Conditions = []metav1.Condition{
			{Type: clusterv1.MachineReadyCondition, Status: metav1.ConditionFalse, Reason: clusterv1.MachineNotReadyReason},
			{Type: clusterv1.MachineNodeHealthyCondition, Status: metav1.ConditionFalse, Reason: clusterv1.MachineNodeDeletedReason},
		}
	}
	return m
}

func strandedTestControlPlane(failureDomains []string, machines ...*clusterv1.Machine) *ck8s.ControlPlane {
	cluster := &clusterv1.Cluster{}
	for _, fd := range failureDomains {
		cluster.Status.FailureDomains = append(cluster.Status.FailureDomains, clusterv1.FailureDomain{Name: fd, ControlPlane: ptr.To(true)})
	}
	return &ck8s.ControlPlane{
		KCP:      &controlplanev1.CK8sControlPlane{},
		Cluster:  cluster,
		Machines: collections.FromMachines(machines...),
	}
}

func TestSelectMachineForScaleDownPrefersMachineWithDeletedNode(t *testing.T) {
	t.Run("no failure domains: picks the stranded machine, not the oldest", func(t *testing.T) {
		g := NewWithT(t)
		cp := strandedTestControlPlane(nil,
			strandedTestMachine("oldest", "", 0, false),
			strandedTestMachine("stranded", "", 10, true),
			strandedTestMachine("newest", "", 20, false),
		)
		m, err := selectMachineForScaleDown(context.Background(), cp, collections.Machines{})
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(m.Name).To(Equal("stranded"))
	})

	t.Run("failure domains: picks the stranded machine outside the failure domain with most machines", func(t *testing.T) {
		g := NewWithT(t)
		cp := strandedTestControlPlane([]string{"az1", "az2"},
			strandedTestMachine("az1-old", "az1", 0, false),
			strandedTestMachine("az1-new", "az1", 5, false),
			strandedTestMachine("az2-stranded", "az2", 10, true),
		)
		m, err := selectMachineForScaleDown(context.Background(), cp, collections.Machines{})
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(m.Name).To(Equal("az2-stranded"))
	})

	t.Run("rollout: picks the stranded machine over other outdated machines", func(t *testing.T) {
		g := NewWithT(t)
		outdatedOld := strandedTestMachine("outdated-old", "", 0, false)
		outdatedStranded := strandedTestMachine("outdated-stranded", "", 5, true)
		cp := strandedTestControlPlane(nil, outdatedOld, outdatedStranded, strandedTestMachine("new", "", 20, false))
		m, err := selectMachineForScaleDown(context.Background(), cp, collections.FromMachines(outdatedOld, outdatedStranded))
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(m.Name).To(Equal("outdated-stranded"))
	})

	t.Run("a machine already being deleted is not preferred", func(t *testing.T) {
		g := NewWithT(t)
		deleting := strandedTestMachine("deleting", "", 10, true)
		deleting.DeletionTimestamp = ptr.To(metav1.NewTime(strandedTestEpoch.Add(time.Hour)))
		deleting.Finalizers = []string{"test"}
		cp := strandedTestControlPlane(nil,
			strandedTestMachine("oldest", "", 0, false),
			deleting,
		)
		m, err := selectMachineForScaleDown(context.Background(), cp, collections.Machines{})
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(m.Name).To(Equal("oldest"))
	})

	t.Run("no stranded machine: unchanged, picks the oldest", func(t *testing.T) {
		g := NewWithT(t)
		cp := strandedTestControlPlane(nil,
			strandedTestMachine("oldest", "", 0, false),
			strandedTestMachine("newest", "", 20, false),
		)
		m, err := selectMachineForScaleDown(context.Background(), cp, collections.Machines{})
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(m.Name).To(Equal("oldest"))
	})
}

func TestNotReadyMachinesHaveOnlyDeletedNodes(t *testing.T) {
	notReadyNodeUp := strandedTestMachine("not-ready", "", 5, false)
	notReadyNodeUp.Status.Conditions = []metav1.Condition{
		{Type: clusterv1.MachineReadyCondition, Status: metav1.ConditionFalse, Reason: clusterv1.MachineNotReadyReason},
		{Type: clusterv1.MachineNodeHealthyCondition, Status: metav1.ConditionFalse, Reason: clusterv1.MachineNodeNotReadyReason},
	}

	tests := []struct {
		name     string
		machines []*clusterv1.Machine
		wantOK   bool
		wantLen  int
	}{
		{
			name:     "only a machine with a deleted node is not ready: continue",
			machines: []*clusterv1.Machine{strandedTestMachine("ready", "", 0, false), strandedTestMachine("stranded", "", 10, true)},
			wantOK:   true,
			wantLen:  1,
		},
		{
			name:     "another machine is not ready for a different reason: requeue",
			machines: []*clusterv1.Machine{strandedTestMachine("stranded", "", 10, true), notReadyNodeUp},
			wantOK:   false,
			wantLen:  2,
		},
		{
			name:     "a machine is not ready with its node still present: requeue",
			machines: []*clusterv1.Machine{strandedTestMachine("ready", "", 0, false), notReadyNodeUp},
			wantOK:   false,
			wantLen:  1,
		},
		{
			name:     "all machines ready: nothing to excuse, requeue as before",
			machines: []*clusterv1.Machine{strandedTestMachine("ready", "", 0, false)},
			wantOK:   false,
			wantLen:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			notReady, ok := notReadyMachinesHaveOnlyDeletedNodes(collections.FromMachines(tt.machines...))
			g.Expect(ok).To(Equal(tt.wantOK))
			g.Expect(notReady.Len()).To(Equal(tt.wantLen))
		})
	}
}
