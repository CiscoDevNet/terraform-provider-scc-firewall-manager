package connectoronboarding

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/connector/connectoronboarding"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Read(ctx context.Context, resource *Resource, stateData *ResourceModel) error {

	// intentional empty, nothing to read

	return nil
}

func Create(ctx context.Context, resource *Resource, planData *ResourceModel) error {

	// do create

	readOutp, err := resource.client.CreateConnectorOnboarding(ctx, connectoronboarding.NewCreateInput(planData.Name.ValueString()))
	if err != nil {
		return err
	}

	// map return struct to sdc model
	planData.Id = types.StringValue(readOutp.Uid)
	planData.Name = types.StringValue(readOutp.Name)

	return nil
}

func Update(ctx context.Context, resource *Resource, planData *ResourceModel, stateData *ResourceModel) error {

	// map return struct to sdc model
	stateData.Name = planData.Name

	return nil
}

func Delete(ctx context.Context, resource *Resource, stateData *ResourceModel) error {

	// intentional empty, nothing to delete

	return nil
}
