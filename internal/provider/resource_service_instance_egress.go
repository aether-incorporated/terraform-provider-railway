package provider

import (
	"context"
	"slices"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type staticEgressAction int

const (
	staticEgressActionNone staticEgressAction = iota
	staticEgressActionCreate
	staticEgressActionClear
)

// getStaticEgressAction compares the configured switch with the number of gateways Railway has assigned.
// Gateways that already exist are kept as they are, because recreating them would change the addresses
// that downstream firewall rules allow.
func getStaticEgressAction(isEnabled bool, currentCount int) staticEgressAction {
	if isEnabled && currentCount == 0 {
		return staticEgressActionCreate
	}
	if !isEnabled && currentCount > 0 {
		return staticEgressActionClear
	}
	return staticEgressActionNone
}

// buildStaticEgressIpsValue returns the addresses sorted and deduplicated, so the order Railway lists
// gateways in never shows up as a plan diff.
func buildStaticEgressIpsValue(ips []string) types.List {
	sorted := slices.Clone(ips)
	slices.Sort(sorted)
	sorted = slices.Compact(sorted)

	values := make([]attr.Value, 0, len(sorted))
	for _, ip := range sorted {
		values = append(values, types.StringValue(ip))
	}
	return types.ListValueMust(types.StringType, values)
}

func getStaticEgressIps(ctx context.Context, client graphql.Client, environmentId string, serviceId string) ([]string, error) {
	response, err := getServiceInstanceEgressGateways(ctx, client, environmentId, serviceId)
	if err != nil {
		return nil, err
	}

	ips := make([]string, 0, len(response.EgressGateways))
	for _, gateway := range response.EgressGateways {
		ips = append(ips, gateway.Ipv4)
	}
	return ips, nil
}

func applyStaticEgress(ctx context.Context, client graphql.Client, environmentId string, serviceId string, isEnabled bool) error {
	current, err := getStaticEgressIps(ctx, client, environmentId, serviceId)
	if err != nil {
		return err
	}

	switch getStaticEgressAction(isEnabled, len(current)) {
	case staticEgressActionCreate:
		_, err = createServiceInstanceEgressGateway(ctx, client, EgressGatewayCreateInput{
			EnvironmentId: environmentId,
			ServiceId:     serviceId,
		})
	case staticEgressActionClear:
		_, err = clearServiceInstanceEgressGateways(ctx, client, EgressGatewayServiceTargetInput{
			EnvironmentId: environmentId,
			ServiceId:     serviceId,
		})
	}
	return err
}

// staticEgressIpsPlanModifier keeps the known addresses while the switch is unchanged and marks them
// unknown when it flips, since only Railway can say which addresses an enable assigns.
type staticEgressIpsPlanModifier struct{}

func (staticEgressIpsPlanModifier) Description(context.Context) string {
	return "Keeps the assigned addresses unless `static_egress_enabled` changes."
}

func (m staticEgressIpsPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (staticEgressIpsPlanModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.State.Raw.IsNull() || !req.PlanValue.IsUnknown() {
		return
	}

	var planEnabled, stateEnabled types.Bool
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("static_egress_enabled"), &planEnabled)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("static_egress_enabled"), &stateEnabled)...)
	if resp.Diagnostics.HasError() || planEnabled.IsUnknown() {
		return
	}

	if planEnabled.Equal(stateEnabled) {
		resp.PlanValue = req.StateValue
	}
}
