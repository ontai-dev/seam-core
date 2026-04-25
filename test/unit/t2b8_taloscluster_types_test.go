// Package unit contains T-2B-8 tests for InfrastructureTalosCluster type reconciliation.
//
// Covers:
//   - InfrastructureProvider enum constants and serialization
//   - InfrastructureTalosClusterOrigin typed enum
//   - ClusterEndpoint rename from Endpoint
//   - NodeAddresses field
//   - Full InfrastructureCAPIConfig six-field struct
//   - InfrastructureLocalObjectRef and status.CAPIClusterRef
//   - status.Origin typed enum round-trip
//   - New condition type and reason constants added to seam-core/pkg/conditions
package unit

import (
	"encoding/json"
	"testing"

	v1alpha1 "github.com/ontai-dev/seam-core/api/v1alpha1"
	"github.com/ontai-dev/seam-core/pkg/conditions"
)

// --- InfrastructureProvider enum ---

func TestInfrastructureProvider_Constants(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		value v1alpha1.InfrastructureProvider
	}{
		{"native", v1alpha1.InfrastructureProviderNative},
		{"capi", v1alpha1.InfrastructureProviderCAPI},
		{"screen", v1alpha1.InfrastructureProviderScreen},
	}
	for _, tc := range cases {
		if string(tc.value) != tc.name {
			t.Errorf("InfrastructureProvider %q: got %q", tc.name, tc.value)
		}
	}
}

func TestInfrastructureProvider_SerializationRoundTrip(t *testing.T) {
	t.Parallel()
	spec := v1alpha1.InfrastructureTalosClusterSpec{
		Mode:                   v1alpha1.InfrastructureTalosClusterModeBootstrap,
		InfrastructureProvider: v1alpha1.InfrastructureProviderCAPI,
	}
	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got v1alpha1.InfrastructureTalosClusterSpec
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.InfrastructureProvider != v1alpha1.InfrastructureProviderCAPI {
		t.Errorf("InfrastructureProvider: got %q want %q", got.InfrastructureProvider, v1alpha1.InfrastructureProviderCAPI)
	}
}

func TestInfrastructureProvider_AbsentWhenZero(t *testing.T) {
	t.Parallel()
	spec := v1alpha1.InfrastructureTalosClusterSpec{
		Mode: v1alpha1.InfrastructureTalosClusterModeBootstrap,
		// InfrastructureProvider not set -- should be absent in JSON (omitempty)
	}
	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}
	if _, ok := raw["infrastructureProvider"]; ok {
		t.Error("infrastructureProvider present in JSON but should be absent (omitempty)")
	}
}

// --- ClusterEndpoint (renamed from Endpoint) ---

func TestInfrastructureTalosCluster_ClusterEndpointRoundTrip(t *testing.T) {
	t.Parallel()
	spec := v1alpha1.InfrastructureTalosClusterSpec{
		Mode:            v1alpha1.InfrastructureTalosClusterModeImport,
		Role:            v1alpha1.InfrastructureTalosClusterRoleManagement,
		ClusterEndpoint: "https://10.20.0.10:6443",
	}
	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got v1alpha1.InfrastructureTalosClusterSpec
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.ClusterEndpoint != "https://10.20.0.10:6443" {
		t.Errorf("ClusterEndpoint: got %q want %q", got.ClusterEndpoint, "https://10.20.0.10:6443")
	}
}

func TestInfrastructureTalosCluster_ClusterEndpoint_JSONKey(t *testing.T) {
	t.Parallel()
	spec := v1alpha1.InfrastructureTalosClusterSpec{
		Mode:            v1alpha1.InfrastructureTalosClusterModeImport,
		ClusterEndpoint: "https://192.168.0.1:6443",
	}
	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}
	if _, ok := raw["clusterEndpoint"]; !ok {
		t.Error("expected clusterEndpoint key in JSON, not found")
	}
	if _, ok := raw["endpoint"]; ok {
		t.Error("old endpoint key must not be present in JSON")
	}
}

// --- NodeAddresses ---

func TestInfrastructureTalosCluster_NodeAddressesRoundTrip(t *testing.T) {
	t.Parallel()
	spec := v1alpha1.InfrastructureTalosClusterSpec{
		Mode:          v1alpha1.InfrastructureTalosClusterModeImport,
		NodeAddresses: []string{"10.20.0.11", "10.20.0.12", "10.20.0.13"},
	}
	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got v1alpha1.InfrastructureTalosClusterSpec
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got.NodeAddresses) != 3 {
		t.Fatalf("NodeAddresses: got %d entries want 3", len(got.NodeAddresses))
	}
	if got.NodeAddresses[0] != "10.20.0.11" {
		t.Errorf("NodeAddresses[0]: got %q want %q", got.NodeAddresses[0], "10.20.0.11")
	}
}

func TestInfrastructureTalosCluster_NodeAddresses_AbsentWhenEmpty(t *testing.T) {
	t.Parallel()
	spec := v1alpha1.InfrastructureTalosClusterSpec{
		Mode: v1alpha1.InfrastructureTalosClusterModeBootstrap,
	}
	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}
	if _, ok := raw["nodeAddresses"]; ok {
		t.Error("nodeAddresses present in JSON but should be absent (omitempty)")
	}
}

// --- Full InfrastructureCAPIConfig (6-field struct) ---

func TestInfrastructureCAPIConfig_FullFieldsRoundTrip(t *testing.T) {
	t.Parallel()
	spec := v1alpha1.InfrastructureTalosClusterSpec{
		Mode: v1alpha1.InfrastructureTalosClusterModeBootstrap,
		CAPI: &v1alpha1.InfrastructureCAPIConfig{
			Enabled:           true,
			TalosVersion:      "v1.9.3",
			KubernetesVersion: "v1.32.0",
			ControlPlane: &v1alpha1.InfrastructureCAPIControlPlaneConfig{
				Replicas: 3,
			},
			Workers: []v1alpha1.InfrastructureCAPIWorkerPool{
				{
					Name:                         "workers",
					Replicas:                     2,
					SeamInfrastructureMachineNames: []string{"sim-worker-01", "sim-worker-02"},
				},
			},
			CiliumPackRef: &v1alpha1.InfrastructureCAPICiliumPackRef{
				Name:    "cilium-ccs-test-v1.16.6-r1",
				Version: "v1.16.6-r1",
			},
		},
	}
	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got v1alpha1.InfrastructureTalosClusterSpec
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.CAPI == nil {
		t.Fatal("CAPI must survive round trip")
	}
	if !got.CAPI.Enabled {
		t.Error("CAPI.Enabled: got false want true")
	}
	if got.CAPI.TalosVersion != "v1.9.3" {
		t.Errorf("CAPI.TalosVersion: got %q want v1.9.3", got.CAPI.TalosVersion)
	}
	if got.CAPI.KubernetesVersion != "v1.32.0" {
		t.Errorf("CAPI.KubernetesVersion: got %q want v1.32.0", got.CAPI.KubernetesVersion)
	}
	if got.CAPI.ControlPlane == nil || got.CAPI.ControlPlane.Replicas != 3 {
		t.Errorf("CAPI.ControlPlane.Replicas: got %v", got.CAPI.ControlPlane)
	}
	if len(got.CAPI.Workers) != 1 || got.CAPI.Workers[0].Replicas != 2 {
		t.Errorf("CAPI.Workers: got %+v", got.CAPI.Workers)
	}
	if len(got.CAPI.Workers[0].SeamInfrastructureMachineNames) != 2 {
		t.Errorf("CAPI.Workers[0].SeamInfrastructureMachineNames: got %v", got.CAPI.Workers[0].SeamInfrastructureMachineNames)
	}
	if got.CAPI.CiliumPackRef == nil || got.CAPI.CiliumPackRef.Name != "cilium-ccs-test-v1.16.6-r1" {
		t.Errorf("CAPI.CiliumPackRef: got %+v", got.CAPI.CiliumPackRef)
	}
}

func TestInfrastructureCAPIConfig_EnabledFalse_OptionalFieldsAbsent(t *testing.T) {
	t.Parallel()
	spec := v1alpha1.InfrastructureTalosClusterSpec{
		Mode: v1alpha1.InfrastructureTalosClusterModeBootstrap,
		CAPI: &v1alpha1.InfrastructureCAPIConfig{
			Enabled: false,
		},
	}
	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}
	capi, ok := raw["capi"].(map[string]interface{})
	if !ok {
		t.Fatal("capi field missing from JSON")
	}
	for _, field := range []string{"talosVersion", "kubernetesVersion", "controlPlane", "workers", "ciliumPackRef"} {
		if _, ok := capi[field]; ok {
			t.Errorf("capi.%s present in JSON but should be absent (omitempty)", field)
		}
	}
}

// --- InfrastructureLocalObjectRef and status.CAPIClusterRef ---

func TestInfrastructureLocalObjectRef_RoundTrip(t *testing.T) {
	t.Parallel()
	ref := v1alpha1.InfrastructureLocalObjectRef{
		Name:      "ccs-test",
		Namespace: "seam-tenant-ccs-test",
	}
	data, err := json.Marshal(ref)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got v1alpha1.InfrastructureLocalObjectRef
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Name != "ccs-test" {
		t.Errorf("Name: got %q want ccs-test", got.Name)
	}
	if got.Namespace != "seam-tenant-ccs-test" {
		t.Errorf("Namespace: got %q want seam-tenant-ccs-test", got.Namespace)
	}
}

func TestInfrastructureTalosClusterStatus_CAPIClusterRef_RoundTrip(t *testing.T) {
	t.Parallel()
	status := v1alpha1.InfrastructureTalosClusterStatus{
		CAPIClusterRef: &v1alpha1.InfrastructureLocalObjectRef{
			Name:      "ccs-test",
			Namespace: "seam-tenant-ccs-test",
		},
	}
	data, err := json.Marshal(status)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got v1alpha1.InfrastructureTalosClusterStatus
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.CAPIClusterRef == nil {
		t.Fatal("CAPIClusterRef must survive round trip")
	}
	if got.CAPIClusterRef.Name != "ccs-test" {
		t.Errorf("CAPIClusterRef.Name: got %q want ccs-test", got.CAPIClusterRef.Name)
	}
}

func TestInfrastructureTalosClusterStatus_CAPIClusterRef_AbsentWhenNil(t *testing.T) {
	t.Parallel()
	status := v1alpha1.InfrastructureTalosClusterStatus{
		ObservedGeneration: 1,
	}
	data, err := json.Marshal(status)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}
	if _, ok := raw["capiClusterRef"]; ok {
		t.Error("capiClusterRef present in JSON but should be absent when nil (omitempty)")
	}
}

// --- InfrastructureTalosClusterOrigin typed enum ---

func TestInfrastructureTalosClusterOrigin_Constants(t *testing.T) {
	t.Parallel()
	if string(v1alpha1.InfrastructureTalosClusterOriginBootstrapped) != "bootstrapped" {
		t.Errorf("OriginBootstrapped: got %q", v1alpha1.InfrastructureTalosClusterOriginBootstrapped)
	}
	if string(v1alpha1.InfrastructureTalosClusterOriginImported) != "imported" {
		t.Errorf("OriginImported: got %q", v1alpha1.InfrastructureTalosClusterOriginImported)
	}
}

func TestInfrastructureTalosClusterStatus_Origin_TypedEnum_RoundTrip(t *testing.T) {
	t.Parallel()
	status := v1alpha1.InfrastructureTalosClusterStatus{
		Origin: v1alpha1.InfrastructureTalosClusterOriginBootstrapped,
	}
	data, err := json.Marshal(status)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got v1alpha1.InfrastructureTalosClusterStatus
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Origin != v1alpha1.InfrastructureTalosClusterOriginBootstrapped {
		t.Errorf("Origin: got %q want bootstrapped", got.Origin)
	}
}

// --- New conditions constants ---

func TestConditionType_Bootstrapped(t *testing.T) {
	t.Parallel()
	if conditions.ConditionTypeBootstrapped != "Bootstrapped" {
		t.Errorf("ConditionTypeBootstrapped: got %q", conditions.ConditionTypeBootstrapped)
	}
}

func TestConditionType_ScreenProviderNotImplemented(t *testing.T) {
	t.Parallel()
	if conditions.ConditionTypeScreenProviderNotImplemented != "ScreenProviderNotImplemented" {
		t.Errorf("ConditionTypeScreenProviderNotImplemented: got %q", conditions.ConditionTypeScreenProviderNotImplemented)
	}
}

func TestConditionType_PhaseFailed(t *testing.T) {
	t.Parallel()
	if conditions.ConditionTypePhaseFailed != "PhaseFailed" {
		t.Errorf("ConditionTypePhaseFailed: got %q", conditions.ConditionTypePhaseFailed)
	}
}

func TestConditionType_KubeconfigUnavailable(t *testing.T) {
	t.Parallel()
	if conditions.ConditionTypeKubeconfigUnavailable != "KubeconfigUnavailable" {
		t.Errorf("ConditionTypeKubeconfigUnavailable: got %q", conditions.ConditionTypeKubeconfigUnavailable)
	}
}

func TestReason_ScreenNotImplemented(t *testing.T) {
	t.Parallel()
	if conditions.ReasonScreenNotImplemented != "ScreenNotImplemented" {
		t.Errorf("ReasonScreenNotImplemented: got %q", conditions.ReasonScreenNotImplemented)
	}
}

func TestReason_TalosVersionRequired(t *testing.T) {
	t.Parallel()
	if conditions.ReasonTalosVersionRequired != "TalosVersionRequired" {
		t.Errorf("ReasonTalosVersionRequired: got %q", conditions.ReasonTalosVersionRequired)
	}
}

func TestReason_TalosConfigSecretAbsent(t *testing.T) {
	t.Parallel()
	if conditions.ReasonTalosConfigSecretAbsent != "TalosConfigSecretAbsent" {
		t.Errorf("ReasonTalosConfigSecretAbsent: got %q", conditions.ReasonTalosConfigSecretAbsent)
	}
}

// --- Full InfrastructureTalosCluster round-trip ---

func TestInfrastructureTalosCluster_FullSpecRoundTrip(t *testing.T) {
	t.Parallel()
	tc := v1alpha1.InfrastructureTalosCluster{
		Spec: v1alpha1.InfrastructureTalosClusterSpec{
			Mode:                   v1alpha1.InfrastructureTalosClusterModeBootstrap,
			TalosVersion:           "v1.9.3",
			ClusterEndpoint:        "https://10.20.0.10:6443",
			NodeAddresses:          []string{"10.20.0.11", "10.20.0.12"},
			InfrastructureProvider: v1alpha1.InfrastructureProviderNative,
			CAPI: &v1alpha1.InfrastructureCAPIConfig{
				Enabled:           true,
				TalosVersion:      "v1.9.3",
				KubernetesVersion: "v1.32.0",
				ControlPlane:      &v1alpha1.InfrastructureCAPIControlPlaneConfig{Replicas: 1},
				Workers:           []v1alpha1.InfrastructureCAPIWorkerPool{{Name: "workers", Replicas: 2}},
				CiliumPackRef:     &v1alpha1.InfrastructureCAPICiliumPackRef{Name: "cilium-ccs-test", Version: "v1.16.6-r1"},
			},
		},
		Status: v1alpha1.InfrastructureTalosClusterStatus{
			Origin:             v1alpha1.InfrastructureTalosClusterOriginBootstrapped,
			ObservedGeneration: 2,
			CAPIClusterRef: &v1alpha1.InfrastructureLocalObjectRef{
				Name:      "ccs-test",
				Namespace: "seam-tenant-ccs-test",
			},
		},
	}
	data, err := json.Marshal(tc)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got v1alpha1.InfrastructureTalosCluster
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Spec.ClusterEndpoint != "https://10.20.0.10:6443" {
		t.Errorf("Spec.ClusterEndpoint: got %q", got.Spec.ClusterEndpoint)
	}
	if len(got.Spec.NodeAddresses) != 2 {
		t.Errorf("Spec.NodeAddresses: got %d want 2", len(got.Spec.NodeAddresses))
	}
	if got.Spec.InfrastructureProvider != v1alpha1.InfrastructureProviderNative {
		t.Errorf("Spec.InfrastructureProvider: got %q", got.Spec.InfrastructureProvider)
	}
	if got.Spec.CAPI == nil || !got.Spec.CAPI.Enabled {
		t.Error("Spec.CAPI.Enabled must be true")
	}
	if got.Status.Origin != v1alpha1.InfrastructureTalosClusterOriginBootstrapped {
		t.Errorf("Status.Origin: got %q", got.Status.Origin)
	}
	if got.Status.CAPIClusterRef == nil || got.Status.CAPIClusterRef.Name != "ccs-test" {
		t.Errorf("Status.CAPIClusterRef: got %v", got.Status.CAPIClusterRef)
	}
}
