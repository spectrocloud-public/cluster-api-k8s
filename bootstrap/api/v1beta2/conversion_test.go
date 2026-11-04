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

// Round-trip tests for the v1beta2 (spoke) <-> v1beta3 (hub) conversion of CK8sConfig and
// CK8sConfigTemplate. See controlplane/api/v1beta2/conversion_test.go for the rationale: spec
// must survive both directions exactly; status losses are asserted explicitly.

import (
	"encoding/json"
	"testing"
	"time"

	apiequality "k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	"sigs.k8s.io/randfill"

	bootstrapv1 "github.com/canonical/cluster-api-k8s/bootstrap/api/v1beta3"
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

func TestCK8sConfigSpecRoundTrip(t *testing.T) {
	t.Run("v1beta2 -> v1beta3 -> v1beta2", func(t *testing.T) {
		for seed := int64(0); seed < fuzzIterations; seed++ {
			src := &CK8sConfig{}
			newFiller(seed).Fill(&src.Spec)
			newFiller(seed + 1_000_000).Fill(&src.Status) // exercise every status branch; only spec is compared

			hub := &bootstrapv1.CK8sConfig{}
			if err := src.ConvertTo(hub); err != nil {
				t.Fatalf("seed %d: ConvertTo: %v", seed, err)
			}
			got := &CK8sConfig{}
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
			src := &bootstrapv1.CK8sConfig{}
			newFiller(seed).Fill(&src.Spec)
			newFiller(seed + 1_000_000).Fill(&src.Status) // exercise every status branch; only spec is compared

			spoke := &CK8sConfig{}
			if err := spoke.ConvertFrom(src); err != nil {
				t.Fatalf("seed %d: ConvertFrom: %v", seed, err)
			}
			got := &bootstrapv1.CK8sConfig{}
			if err := spoke.ConvertTo(got); err != nil {
				t.Fatalf("seed %d: ConvertTo: %v", seed, err)
			}
			if !apiequality.Semantic.DeepEqual(src.Spec, got.Spec) {
				t.Fatalf("seed %d: spec changed\nwant %s\ngot  %s", seed, mustJSON(t, src.Spec), mustJSON(t, got.Spec))
			}
		}
	})
}

func TestCK8sConfigTemplateSpecRoundTrip(t *testing.T) {
	t.Run("v1beta2 -> v1beta3 -> v1beta2", func(t *testing.T) {
		for seed := int64(0); seed < fuzzIterations; seed++ {
			src := &CK8sConfigTemplate{}
			newFiller(seed).Fill(&src.Spec)

			hub := &bootstrapv1.CK8sConfigTemplate{}
			if err := src.ConvertTo(hub); err != nil {
				t.Fatalf("seed %d: ConvertTo: %v", seed, err)
			}
			got := &CK8sConfigTemplate{}
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
			src := &bootstrapv1.CK8sConfigTemplate{}
			newFiller(seed).Fill(&src.Spec)

			spoke := &CK8sConfigTemplate{}
			if err := spoke.ConvertFrom(src); err != nil {
				t.Fatalf("seed %d: ConvertFrom: %v", seed, err)
			}
			got := &bootstrapv1.CK8sConfigTemplate{}
			if err := spoke.ConvertTo(got); err != nil {
				t.Fatalf("seed %d: ConvertTo: %v", seed, err)
			}
			if !apiequality.Semantic.DeepEqual(src.Spec, got.Spec) {
				t.Fatalf("seed %d: spec changed\nwant %s\ngot  %s", seed, mustJSON(t, src.Spec), mustJSON(t, got.Spec))
			}
		}
	})
}

func TestCK8sConfigObjectMetaPreserved(t *testing.T) {
	src := &CK8sConfig{ObjectMeta: metav1.ObjectMeta{
		Name: "cfg", Namespace: "ns", UID: "uid-1", ResourceVersion: "42", Generation: 3,
		Labels:          map[string]string{"cluster.x-k8s.io/cluster-name": "c1"},
		Annotations:     map[string]string{"a": "b"},
		OwnerReferences: []metav1.OwnerReference{{APIVersion: "cluster.x-k8s.io/v1beta2", Kind: "Machine", Name: "m1", UID: "uid-m1"}},
	}}
	hub := &bootstrapv1.CK8sConfig{}
	if err := src.ConvertTo(hub); err != nil {
		t.Fatal(err)
	}
	got := &CK8sConfig{}
	if err := got.ConvertFrom(hub); err != nil {
		t.Fatal(err)
	}
	if !apiequality.Semantic.DeepEqual(src.ObjectMeta, hub.ObjectMeta) || !apiequality.Semantic.DeepEqual(src.ObjectMeta, got.ObjectMeta) {
		t.Fatalf("ObjectMeta changed\nhub %s\ngot %s", mustJSON(t, hub.ObjectMeta), mustJSON(t, got.ObjectMeta))
	}
}

var (
	t0 = metav1.NewTime(time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC))
	t1 = metav1.NewTime(time.Date(2026, 11, 4, 11, 0, 0, 0, time.UTC))
)

// v1beta2 -> v1beta3 -> v1beta2.
func TestCK8sConfigStatusSpokeHubSpoke(t *testing.T) {
	tests := []struct {
		name                  string
		in                    CK8sConfigStatus
		wantDataSecretCreated *bool
		wantConditions        []metav1.Condition
		want                  CK8sConfigStatus
	}{
		{
			name:                  "bootstrap data ready",
			in:                    CK8sConfigStatus{Ready: true, DataSecretName: ptr.To("m1-bootstrap"), ObservedGeneration: 3},
			wantDataSecretCreated: ptr.To(true),
			want:                  CK8sConfigStatus{Ready: true, DataSecretName: ptr.To("m1-bootstrap"), ObservedGeneration: 3},
		},
		{
			name:                  "not ready yet",
			in:                    CK8sConfigStatus{ObservedGeneration: 1},
			wantDataSecretCreated: nil,
			want:                  CK8sConfigStatus{ObservedGeneration: 1},
		},
		{
			// KNOWN: ready is re-derived from dataSecretCreated, which ConvertTo sets whenever
			// dataSecretName is present. The controller always sets both together.
			name:                  "dataSecretName set but ready=false",
			in:                    CK8sConfigStatus{Ready: false, DataSecretName: ptr.To("m1-bootstrap")},
			wantDataSecretCreated: ptr.To(true),
			want:                  CK8sConfigStatus{Ready: true, DataSecretName: ptr.To("m1-bootstrap")},
		},
		{
			name: "failure fields and conditions",
			in: CK8sConfigStatus{
				FailureReason: "InvalidConfig", FailureMessage: "bad", ObservedGeneration: 2,
				Conditions: clusterv1.Conditions{
					{Type: "Ready", Status: "False", Severity: clusterv1.ConditionSeverityError, Reason: "DataSecretGenerationFailed", Message: "x", LastTransitionTime: t0},
					{Type: "CertificatesAvailable", Status: "True", LastTransitionTime: t1},
				},
			},
			wantConditions: []metav1.Condition{
				{Type: "Ready", Status: "False", Reason: "DataSecretGenerationFailed", Message: "x", LastTransitionTime: t0},
				{Type: "CertificatesAvailable", Status: "True", Reason: "NoReasonReported", LastTransitionTime: t1},
			},
			want: CK8sConfigStatus{
				FailureReason: "InvalidConfig", FailureMessage: "bad", ObservedGeneration: 2,
				Conditions: clusterv1.Conditions{
					{Type: "Ready", Status: "False", Reason: "DataSecretGenerationFailed", Message: "x", LastTransitionTime: t0}, // KNOWN LOSS: severity
					{Type: "CertificatesAvailable", Status: "True", Reason: "NoReasonReported", LastTransitionTime: t1},          // KNOWN: reason defaulted
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := &CK8sConfig{Status: tt.in}
			hub := &bootstrapv1.CK8sConfig{}
			if err := src.ConvertTo(hub); err != nil {
				t.Fatal(err)
			}
			if !apiequality.Semantic.DeepEqual(tt.wantDataSecretCreated, hub.Status.Initialization.DataSecretCreated) {
				t.Errorf("dataSecretCreated: want %v got %v", mustJSON(t, tt.wantDataSecretCreated), mustJSON(t, hub.Status.Initialization.DataSecretCreated))
			}
			if !apiequality.Semantic.DeepEqual(tt.wantConditions, hub.Status.Conditions) {
				t.Errorf("v1beta3 conditions\nwant %s\ngot  %s", mustJSON(t, tt.wantConditions), mustJSON(t, hub.Status.Conditions))
			}

			got := &CK8sConfig{}
			if err := got.ConvertFrom(hub); err != nil {
				t.Fatal(err)
			}
			if !apiequality.Semantic.DeepEqual(tt.want, got.Status) {
				t.Errorf("v1beta2 status after round trip\nwant %s\ngot  %s", mustJSON(t, tt.want), mustJSON(t, got.Status))
			}
		})
	}
}

// v1beta3 -> v1beta2 -> v1beta3.
func TestCK8sConfigStatusHubSpokeHub(t *testing.T) {
	in := bootstrapv1.CK8sConfigStatus{
		Ready:              true,
		DataSecretName:     ptr.To("m1-bootstrap"),
		ObservedGeneration: 3,
		Initialization:     bootstrapv1.Ck8sConfigInitializationStatus{DataSecretCreated: ptr.To(true)},
		Conditions: []metav1.Condition{
			{Type: "Ready", Status: "True", Reason: "Ready", ObservedGeneration: 3, LastTransitionTime: t0},
		},
	}
	want := *in.DeepCopy()
	want.Conditions[0].ObservedGeneration = 0 // KNOWN LOSS: v1beta1 conditions have no observedGeneration

	hub := &bootstrapv1.CK8sConfig{Status: in}
	spoke := &CK8sConfig{}
	if err := spoke.ConvertFrom(hub); err != nil {
		t.Fatal(err)
	}
	if !spoke.Status.Ready {
		t.Errorf("v1beta2 view: ready should follow dataSecretCreated: %s", mustJSON(t, spoke.Status))
	}
	got := &bootstrapv1.CK8sConfig{}
	if err := spoke.ConvertTo(got); err != nil {
		t.Fatal(err)
	}
	if !apiequality.Semantic.DeepEqual(want, got.Status) {
		t.Errorf("v1beta3 status after round trip\nwant %s\ngot  %s", mustJSON(t, want), mustJSON(t, got.Status))
	}

	t.Run("dataSecretCreated without dataSecretName", func(t *testing.T) {
		// KNOWN LOSS: ConvertTo derives dataSecretCreated from dataSecretName only.
		hub := &bootstrapv1.CK8sConfig{Status: bootstrapv1.CK8sConfigStatus{
			Initialization: bootstrapv1.Ck8sConfigInitializationStatus{DataSecretCreated: ptr.To(true)},
		}}
		spoke := &CK8sConfig{}
		if err := spoke.ConvertFrom(hub); err != nil {
			t.Fatal(err)
		}
		got := &bootstrapv1.CK8sConfig{}
		if err := spoke.ConvertTo(got); err != nil {
			t.Fatal(err)
		}
		if got.Status.Initialization.DataSecretCreated != nil {
			t.Errorf("expected dataSecretCreated to be dropped, got %v", *got.Status.Initialization.DataSecretCreated)
		}
	})
}
