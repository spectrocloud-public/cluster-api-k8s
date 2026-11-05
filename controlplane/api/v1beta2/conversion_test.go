/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta2

// Round-trip tests for the v1beta2 (spoke) <-> v1beta3 (hub) conversion.
//
// Objects written by releases before v1beta3 existed stay stored as v1beta2 in etcd until
// something rewrites them, so every read by the v1beta3 controllers goes through ConvertTo, and
// every v1beta2 client read goes through ConvertFrom. Spec must survive both directions exactly
// (a lost spec field changes desired state and can trigger a rollout). Status is lossy by design;
// every known loss is asserted explicitly below so a change in behavior is visible.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	apiequality "k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
	"sigs.k8s.io/randfill"

	controlplanev1 "github.com/canonical/cluster-api-k8s/controlplane/api/v1beta3"
)

const fuzzIterations = 500

// newFiller returns a deterministic filler. The explicit funcs matter: types such as metav1.Time
// and intstr.IntOrString have their own fill methods that leave a nil pointer field nil, so without
// these funcs fields like rolloutAfter and maxSurge would never be exercised. They also keep values
// to what the API server can store (second-precision timestamps, valid managedFields JSON).
func newFiller(seed int64) *randfill.Filler {
	return randfill.NewWithSeed(seed).NilChance(0.3).NumElements(0, 3).Funcs(
		func(t *metav1.Time, c randfill.Continue) {
			*t = metav1.Unix(c.Int63n(4_000_000_000), 0).Rfc3339Copy()
		},
		func(f *metav1.FieldsV1, c randfill.Continue) {
			f.Raw = []byte(`{}`)
		},
		func(i *intstr.IntOrString, c randfill.Continue) {
			if c.Bool() {
				*i = intstr.FromInt32(c.Int31())
			} else {
				*i = intstr.FromString(c.String(0))
			}
		},
	)
}

func mustJSON(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(b)
}

// strictDecode fails on any JSON field the target type does not know, i.e. a field that the
// API server would prune.
func strictDecode(data []byte, into any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode(into)
}

func TestCK8sControlPlaneSpecRoundTrip(t *testing.T) {
	t.Run("v1beta2 -> v1beta3 -> v1beta2", func(t *testing.T) {
		for seed := int64(0); seed < fuzzIterations; seed++ {
			src := &CK8sControlPlane{}
			newFiller(seed).Fill(&src.Spec)
			newFiller(seed + 1_000_000).Fill(&src.Status) // exercise every status branch; only spec is compared

			hub := &controlplanev1.CK8sControlPlane{}
			if err := src.ConvertTo(hub); err != nil {
				t.Fatalf("seed %d: ConvertTo: %v", seed, err)
			}
			got := &CK8sControlPlane{}
			if err := got.ConvertFrom(hub); err != nil {
				t.Fatalf("seed %d: ConvertFrom: %v", seed, err)
			}
			if !apiequality.Semantic.DeepEqual(src.Spec, got.Spec) {
				t.Fatalf("seed %d: spec changed\nwant %s\ngot  %s", seed, mustJSON(t, src.Spec), mustJSON(t, got.Spec))
			}
		}
	})

	t.Run("v1beta3 -> v1beta2 -> v1beta3", func(t *testing.T) {
		for seed := int64(0); seed < fuzzIterations; seed++ {
			src := &controlplanev1.CK8sControlPlane{}
			newFiller(seed).Fill(&src.Spec)
			newFiller(seed + 1_000_000).Fill(&src.Status) // exercise every status branch; only spec is compared

			spoke := &CK8sControlPlane{}
			if err := spoke.ConvertFrom(src); err != nil {
				t.Fatalf("seed %d: ConvertFrom: %v", seed, err)
			}
			got := &controlplanev1.CK8sControlPlane{}
			if err := spoke.ConvertTo(got); err != nil {
				t.Fatalf("seed %d: ConvertTo: %v", seed, err)
			}
			if !apiequality.Semantic.DeepEqual(src.Spec, got.Spec) {
				t.Fatalf("seed %d: spec changed\nwant %s\ngot  %s", seed, mustJSON(t, src.Spec), mustJSON(t, got.Spec))
			}
		}
	})
}

func TestCK8sControlPlaneObjectMetaPreserved(t *testing.T) {
	src := &CK8sControlPlane{ObjectMeta: metav1.ObjectMeta{
		Name: "cp", Namespace: "ns", UID: "uid-1", ResourceVersion: "42", Generation: 7,
		Labels:      map[string]string{"cluster.x-k8s.io/cluster-name": "c1"},
		Annotations: map[string]string{"a": "b"},
		Finalizers:  []string{"ck8s.controlplane.cluster.x-k8s.io"},
	}}
	hub := &controlplanev1.CK8sControlPlane{}
	if err := src.ConvertTo(hub); err != nil {
		t.Fatal(err)
	}
	if !apiequality.Semantic.DeepEqual(src.ObjectMeta, hub.ObjectMeta) {
		t.Fatalf("ConvertTo changed ObjectMeta: %s", mustJSON(t, hub.ObjectMeta))
	}
	got := &CK8sControlPlane{}
	if err := got.ConvertFrom(hub); err != nil {
		t.Fatal(err)
	}
	if !apiequality.Semantic.DeepEqual(src.ObjectMeta, got.ObjectMeta) {
		t.Fatalf("ConvertFrom changed ObjectMeta: %s", mustJSON(t, got.ObjectMeta))
	}
}

// PCP-7634 (fork): an empty machineTemplate.metadata must not be emitted as `metadata: {}`, which
// fails the v1beta3 CRD's minProperties:1 validation.
func TestCK8sControlPlaneEmptyMachineTemplateMetadataOmitted(t *testing.T) {
	src := &CK8sControlPlane{Spec: CK8sControlPlaneSpec{Version: "v1.32.0", Replicas: ptr.To[int32](1)}}
	hub := &controlplanev1.CK8sControlPlane{}
	if err := src.ConvertTo(hub); err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(mustJSON(t, hub.Spec)), &out); err != nil {
		t.Fatal(err)
	}
	if mt, ok := out["machineTemplate"].(map[string]any); ok {
		if _, has := mt["metadata"]; has {
			t.Fatalf("converted v1beta3 spec emits machineTemplate.metadata: %s", mustJSON(t, hub.Spec))
		}
	}
}

var (
	t0 = metav1.NewTime(time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC))
	t1 = metav1.NewTime(time.Date(2026, 11, 4, 11, 0, 0, 0, time.UTC))
)

// v1beta2 -> v1beta3 -> v1beta2: the status a v1beta2 client (e.g. Palette) reads back for an
// object that was written as v1beta2.
func TestCK8sControlPlaneStatusSpokeHubSpoke(t *testing.T) {
	baseStatus := func() CK8sControlPlaneStatus {
		return CK8sControlPlaneStatus{
			Selector:           "cluster.x-k8s.io/cluster-name=c1",
			Replicas:           3,
			Version:            ptr.To("v1.32.0"),
			UpdatedReplicas:    3,
			ReadyReplicas:      3,
			Initialized:        true,
			Ready:              true,
			FailureMessage:     ptr.To("boom"),
			ObservedGeneration: 7,
			LastRemediation:    &LastRemediationStatus{Machine: "m1", Timestamp: t0, RetryCount: 2},
		}
	}

	tests := []struct {
		name    string
		in      CK8sControlPlaneStatus
		wantHub controlplanev1.CK8sControlPlaneStatus // only the mapped fields are checked
		want    CK8sControlPlaneStatus
	}{
		{
			name: "healthy 3/3",
			in:   baseStatus(),
			wantHub: controlplanev1.CK8sControlPlaneStatus{
				Initialization:    controlplanev1.CK8sControlPlaneInitializationStatus{ControlPlaneInitialized: ptr.To(true)},
				ReadyReplicas:     ptr.To[int32](3),
				AvailableReplicas: ptr.To[int32](3),
				UpToDateReplicas:  ptr.To[int32](3),
			},
			want: baseStatus(),
		},
		{
			name: "rolling: 4 replicas, 1 updated, 3 ready, 1 unavailable",
			in: func() CK8sControlPlaneStatus {
				s := baseStatus()
				s.Replicas, s.UpdatedReplicas, s.ReadyReplicas, s.UnavailableReplicas = 4, 1, 3, 1
				return s
			}(),
			wantHub: controlplanev1.CK8sControlPlaneStatus{
				Initialization:    controlplanev1.CK8sControlPlaneInitializationStatus{ControlPlaneInitialized: ptr.To(true)},
				ReadyReplicas:     ptr.To[int32](3),
				AvailableReplicas: ptr.To[int32](3), // KNOWN: v1beta2 has no availableReplicas; mirrors readyReplicas
				UpToDateReplicas:  ptr.To[int32](1),
			},
			want: func() CK8sControlPlaneStatus {
				s := baseStatus()
				s.Replicas, s.UpdatedReplicas, s.ReadyReplicas = 4, 1, 3
				s.UnavailableReplicas = 0 // KNOWN LOSS: v1beta3 has no unavailableReplicas
				return s
			}(),
		},
		{
			name: "not initialized, nothing ready",
			in: func() CK8sControlPlaneStatus {
				s := baseStatus()
				s.UpdatedReplicas, s.ReadyReplicas, s.UnavailableReplicas = 0, 0, 3
				s.Initialized, s.Ready = false, false
				return s
			}(),
			wantHub: controlplanev1.CK8sControlPlaneStatus{
				Initialization:    controlplanev1.CK8sControlPlaneInitializationStatus{ControlPlaneInitialized: nil}, // false -> unset
				ReadyReplicas:     ptr.To[int32](0),
				AvailableReplicas: ptr.To[int32](0),
				UpToDateReplicas:  ptr.To[int32](0),
			},
			want: func() CK8sControlPlaneStatus {
				s := baseStatus()
				s.UpdatedReplicas, s.ReadyReplicas = 0, 0
				s.UnavailableReplicas = 0 // KNOWN LOSS
				s.Initialized, s.Ready = false, false
				return s
			}(),
		},
		{
			// User-managed CNI (enableDefaultNetwork=false): the controller reports the control plane
			// as ready once initialized, before any node is Ready. v1beta3 has no `ready` field (the
			// controller sets the Available condition instead), and ConvertFrom derives `ready` from
			// readyReplicas only, so a v1beta2 reader sees ready=false here.
			// KNOWN LOSS, open item on PCP-7766 (Palette MAAS/Hypershift wait on status.ready). If the
			// conversion is changed to also consult the Available condition, update this case.
			name: "initialized, ready with 0 ready replicas (user-managed CNI)",
			in: func() CK8sControlPlaneStatus {
				s := baseStatus()
				s.UpdatedReplicas, s.ReadyReplicas, s.UnavailableReplicas = 0, 0, 3
				s.Initialized, s.Ready = true, true
				return s
			}(),
			wantHub: controlplanev1.CK8sControlPlaneStatus{
				Initialization:    controlplanev1.CK8sControlPlaneInitializationStatus{ControlPlaneInitialized: ptr.To(true)},
				ReadyReplicas:     ptr.To[int32](0),
				AvailableReplicas: ptr.To[int32](0),
				UpToDateReplicas:  ptr.To[int32](0),
			},
			want: func() CK8sControlPlaneStatus {
				s := baseStatus()
				s.UpdatedReplicas, s.ReadyReplicas = 0, 0
				s.UnavailableReplicas = 0 // KNOWN LOSS
				s.Initialized = true
				s.Ready = false // KNOWN LOSS: recomputed as readyReplicas > 0
				return s
			}(),
		},
		{
			name: "conditions: severity dropped, empty reason defaulted",
			in: func() CK8sControlPlaneStatus {
				s := baseStatus()
				s.Conditions = clusterv1.Conditions{
					{Type: "Available", Status: "True", LastTransitionTime: t0},
					{Type: "MachinesReady", Status: "False", Severity: clusterv1.ConditionSeverityWarning, Reason: "Scaling", Message: "1 of 3", LastTransitionTime: t1},
				}
				return s
			}(),
			wantHub: controlplanev1.CK8sControlPlaneStatus{
				Initialization:    controlplanev1.CK8sControlPlaneInitializationStatus{ControlPlaneInitialized: ptr.To(true)},
				ReadyReplicas:     ptr.To[int32](3),
				AvailableReplicas: ptr.To[int32](3),
				UpToDateReplicas:  ptr.To[int32](3),
				Conditions: []metav1.Condition{
					{Type: "Available", Status: "True", Reason: "NoReasonReported", LastTransitionTime: t0},
					{Type: "MachinesReady", Status: "False", Reason: "Scaling", Message: "1 of 3", LastTransitionTime: t1},
				},
			},
			want: func() CK8sControlPlaneStatus {
				s := baseStatus()
				s.Conditions = clusterv1.Conditions{
					{Type: "Available", Status: "True", Reason: "NoReasonReported", LastTransitionTime: t0},                // KNOWN: reason defaulted
					{Type: "MachinesReady", Status: "False", Reason: "Scaling", Message: "1 of 3", LastTransitionTime: t1}, // KNOWN LOSS: severity
				}
				return s
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := &CK8sControlPlane{Status: tt.in}
			hub := &controlplanev1.CK8sControlPlane{}
			if err := src.ConvertTo(hub); err != nil {
				t.Fatal(err)
			}
			gotHub := controlplanev1.CK8sControlPlaneStatus{
				Initialization:    hub.Status.Initialization,
				ReadyReplicas:     hub.Status.ReadyReplicas,
				AvailableReplicas: hub.Status.AvailableReplicas,
				UpToDateReplicas:  hub.Status.UpToDateReplicas,
				Conditions:        hub.Status.Conditions,
			}
			if !apiequality.Semantic.DeepEqual(tt.wantHub, gotHub) {
				t.Errorf("v1beta3 status\nwant %s\ngot  %s", mustJSON(t, tt.wantHub), mustJSON(t, gotHub))
			}
			// Fields with the same name and type in both versions are copied as-is.
			if hub.Status.Selector != tt.in.Selector || hub.Status.Replicas != tt.in.Replicas ||
				!apiequality.Semantic.DeepEqual(hub.Status.Version, tt.in.Version) ||
				hub.Status.ObservedGeneration != tt.in.ObservedGeneration ||
				!apiequality.Semantic.DeepEqual(hub.Status.FailureMessage, tt.in.FailureMessage) ||
				mustJSON(t, hub.Status.LastRemediation) != mustJSON(t, tt.in.LastRemediation) {
				t.Errorf("v1beta3 status: copied fields differ: %s", mustJSON(t, hub.Status))
			}

			got := &CK8sControlPlane{}
			if err := got.ConvertFrom(hub); err != nil {
				t.Fatal(err)
			}
			if !apiequality.Semantic.DeepEqual(tt.want, got.Status) {
				t.Errorf("v1beta2 status after round trip\nwant %s\ngot  %s", mustJSON(t, tt.want), mustJSON(t, got.Status))
			}
		})
	}
}

// v1beta3 -> v1beta2 -> v1beta3: what the v1beta3 controllers see after a v1beta2 client
// (e.g. Palette) writes the object back.
func TestCK8sControlPlaneStatusHubSpokeHub(t *testing.T) {
	in := controlplanev1.CK8sControlPlaneStatus{
		Selector:           "cluster.x-k8s.io/cluster-name=c1",
		Replicas:           4,
		Version:            ptr.To("v1.32.0"),
		ObservedGeneration: 7,
		Initialization:     controlplanev1.CK8sControlPlaneInitializationStatus{ControlPlaneInitialized: ptr.To(true)},
		ReadyReplicas:      ptr.To[int32](3),
		AvailableReplicas:  ptr.To[int32](2),
		UpToDateReplicas:   ptr.To[int32](1),
		Conditions: []metav1.Condition{
			{Type: "Available", Status: "True", Reason: "Available", Message: "ok", ObservedGeneration: 7, LastTransitionTime: t0},
		},
	}
	want := *in.DeepCopy()
	want.AvailableReplicas = ptr.To[int32](3) // KNOWN LOSS: v1beta2 has no availableReplicas; mirrors readyReplicas
	want.Conditions[0].ObservedGeneration = 0 // KNOWN LOSS: v1beta1 conditions have no observedGeneration

	hub := &controlplanev1.CK8sControlPlane{Status: in}
	spoke := &CK8sControlPlane{}
	if err := spoke.ConvertFrom(hub); err != nil {
		t.Fatal(err)
	}
	if spoke.Status.UpdatedReplicas != 1 || spoke.Status.ReadyReplicas != 3 || !spoke.Status.Initialized || !spoke.Status.Ready {
		t.Errorf("v1beta2 view: %s", mustJSON(t, spoke.Status))
	}
	got := &controlplanev1.CK8sControlPlane{}
	if err := spoke.ConvertTo(got); err != nil {
		t.Fatal(err)
	}
	if !apiequality.Semantic.DeepEqual(want, got.Status) {
		t.Errorf("v1beta3 status after round trip\nwant %s\ngot  %s", mustJSON(t, want), mustJSON(t, got.Status))
	}

	t.Run("unset counters and initialized=false", func(t *testing.T) {
		hub := &controlplanev1.CK8sControlPlane{Status: controlplanev1.CK8sControlPlaneStatus{
			Initialization: controlplanev1.CK8sControlPlaneInitializationStatus{ControlPlaneInitialized: ptr.To(false)},
		}}
		spoke := &CK8sControlPlane{}
		if err := spoke.ConvertFrom(hub); err != nil {
			t.Fatal(err)
		}
		got := &controlplanev1.CK8sControlPlane{}
		if err := spoke.ConvertTo(got); err != nil {
			t.Fatal(err)
		}
		want := controlplanev1.CK8sControlPlaneStatus{
			// KNOWN: nil counters come back as explicit 0; initialized=false comes back unset.
			ReadyReplicas:     ptr.To[int32](0),
			AvailableReplicas: ptr.To[int32](0),
			UpToDateReplicas:  ptr.To[int32](0),
		}
		if !apiequality.Semantic.DeepEqual(want, got.Status) {
			t.Errorf("want %s\ngot  %s", mustJSON(t, want), mustJSON(t, got.Status))
		}
	})
}

// CK8sControlPlaneTemplate has no conversion functions: its CRD uses the default `None`
// strategy, where the API server only rewrites apiVersion. That is only safe while the two
// versions' JSON is interchangeable, which this test enforces in both directions.
func TestCK8sControlPlaneTemplateSchemaCompatible(t *testing.T) {
	if _, ok := any(&CK8sControlPlaneTemplate{}).(conversion.Convertible); ok {
		t.Fatal("CK8sControlPlaneTemplate now implements conversion; replace this test with a ConvertTo/ConvertFrom round trip")
	}

	roundTrip := func(src, mid, dst any) error {
		b, err := json.Marshal(src)
		if err != nil {
			return err
		}
		if err := strictDecode(b, mid); err != nil {
			return fmt.Errorf("decode into other version: %w (json %s)", err, b)
		}
		b2, err := json.Marshal(mid)
		if err != nil {
			return err
		}
		if err := strictDecode(b2, dst); err != nil {
			return fmt.Errorf("decode back: %w", err)
		}
		return nil
	}

	for seed := int64(0); seed < fuzzIterations; seed++ {
		src := &CK8sControlPlaneTemplate{}
		newFiller(seed).Fill(&src.Spec)
		mid := &controlplanev1.CK8sControlPlaneTemplate{}
		got := &CK8sControlPlaneTemplate{}
		if err := roundTrip(&src.Spec, &mid.Spec, &got.Spec); err != nil {
			t.Fatalf("seed %d: v1beta2 -> v1beta3: %v", seed, err)
		}
		if !apiequality.Semantic.DeepEqual(src.Spec, got.Spec) {
			t.Fatalf("seed %d: v1beta2 spec changed\nwant %s\ngot  %s", seed, mustJSON(t, src.Spec), mustJSON(t, got.Spec))
		}

		hsrc := &controlplanev1.CK8sControlPlaneTemplate{}
		newFiller(seed).Fill(&hsrc.Spec)
		hmid := &CK8sControlPlaneTemplate{}
		hgot := &controlplanev1.CK8sControlPlaneTemplate{}
		if err := roundTrip(&hsrc.Spec, &hmid.Spec, &hgot.Spec); err != nil {
			t.Fatalf("seed %d: v1beta3 -> v1beta2: %v", seed, err)
		}
		if !apiequality.Semantic.DeepEqual(hsrc.Spec, hgot.Spec) {
			t.Fatalf("seed %d: v1beta3 spec changed\nwant %s\ngot  %s", seed, mustJSON(t, hsrc.Spec), mustJSON(t, hgot.Spec))
		}
	}
}
