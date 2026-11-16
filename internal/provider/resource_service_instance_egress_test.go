package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestEgressGatewayCreateInputOmitsRegion(t *testing.T) {
	input := EgressGatewayCreateInput{
		EnvironmentId: "environment-id",
		ServiceId:     "service-id",
	}

	encoded, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("marshalling input: %v", err)
	}

	const expected = `{"environmentId":"environment-id","serviceId":"service-id"}`
	if string(encoded) != expected {
		t.Fatalf("unexpected wire format:\n got: %s\nwant: %s", encoded, expected)
	}
}

func TestGetStaticEgressAction(t *testing.T) {
	cases := []struct {
		name      string
		isEnabled bool
		current   int
		expected  staticEgressAction
	}{
		{name: "enable when none assigned", isEnabled: true, current: 0, expected: staticEgressActionCreate},
		{name: "keep when already assigned", isEnabled: true, current: 2, expected: staticEgressActionNone},
		{name: "clear when disabled", isEnabled: false, current: 1, expected: staticEgressActionClear},
		{name: "nothing when disabled and none assigned", isEnabled: false, current: 0, expected: staticEgressActionNone},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := getStaticEgressAction(tc.isEnabled, tc.current)
			if actual != tc.expected {
				t.Fatalf("got %v, want %v", actual, tc.expected)
			}
		})
	}
}

func TestBuildStaticEgressIpsValueSortsAndDeduplicates(t *testing.T) {
	gateways := []string{"203.0.113.9", "198.51.100.4", "203.0.113.9"}

	actual := buildStaticEgressIpsValue(gateways)

	expected := types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("198.51.100.4"),
		types.StringValue("203.0.113.9"),
	})
	if !actual.Equal(expected) {
		t.Fatalf("got %s, want %s", actual, expected)
	}
}

func TestBuildStaticEgressIpsValueEmpty(t *testing.T) {
	actual := buildStaticEgressIpsValue(nil)

	expected := types.ListValueMust(types.StringType, []attr.Value{})
	if !actual.Equal(expected) {
		t.Fatalf("got %s, want %s", actual, expected)
	}
}

func TestStaticEgressIpsPlanModifier(t *testing.T) {
	stateIps := types.ListValueMust(types.StringType, []attr.Value{types.StringValue("198.51.100.4")})

	cases := []struct {
		name          string
		hasState      bool
		stateEnabled  bool
		planEnabled   bool
		isUnknownWant bool
	}{
		{name: "unchanged keeps state", hasState: true, stateEnabled: true, planEnabled: true, isUnknownWant: false},
		{name: "enabling is unknown", hasState: true, stateEnabled: false, planEnabled: true, isUnknownWant: true},
		{name: "disabling is unknown", hasState: true, stateEnabled: true, planEnabled: false, isUnknownWant: true},
		{name: "create is unknown", hasState: false, planEnabled: true, isUnknownWant: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			plan := buildEgressTestState(t, tc.planEnabled, types.ListUnknown(types.StringType))
			req := planmodifier.ListRequest{
				Path:        path.Root("static_egress_ips"),
				Plan:        tfsdk.Plan{Raw: plan.Raw, Schema: plan.Schema},
				PlanValue:   types.ListUnknown(types.StringType),
				ConfigValue: types.ListNull(types.StringType),
			}
			if tc.hasState {
				req.State = buildEgressTestState(t, tc.stateEnabled, stateIps)
				req.StateValue = stateIps
			} else {
				req.State = tfsdk.State{Raw: tftypes.NewValue(egressTestObjectType(), nil), Schema: egressTestSchema()}
				req.StateValue = types.ListNull(types.StringType)
			}
			resp := &planmodifier.ListResponse{PlanValue: req.PlanValue}

			staticEgressIpsPlanModifier{}.PlanModifyList(context.Background(), req, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
			if resp.PlanValue.IsUnknown() != tc.isUnknownWant {
				t.Fatalf("unknown = %v, want %v (value %s)", resp.PlanValue.IsUnknown(), tc.isUnknownWant, resp.PlanValue)
			}
			if !tc.isUnknownWant && !resp.PlanValue.Equal(stateIps) {
				t.Fatalf("got %s, want state value %s", resp.PlanValue, stateIps)
			}
		})
	}
}

func egressTestSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"static_egress_enabled": schema.BoolAttribute{Optional: true},
			"static_egress_ips":     schema.ListAttribute{Computed: true, ElementType: types.StringType},
		},
	}
}

func egressTestObjectType() tftypes.Object {
	return tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"static_egress_enabled": tftypes.Bool,
		"static_egress_ips":     tftypes.List{ElementType: tftypes.String},
	}}
}

func buildEgressTestState(t *testing.T, isEnabled bool, ips types.List) tfsdk.State {
	t.Helper()
	ipsValue, err := ips.ToTerraformValue(context.Background())
	if err != nil {
		t.Fatalf("converting ips: %v", err)
	}
	raw := tftypes.NewValue(egressTestObjectType(), map[string]tftypes.Value{
		"static_egress_enabled": tftypes.NewValue(tftypes.Bool, isEnabled),
		"static_egress_ips":     ipsValue,
	})
	return tfsdk.State{Raw: raw, Schema: egressTestSchema()}
}
