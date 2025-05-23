package ftd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/device/publicapilabels"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/device/tags"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/ftd/license"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/device/cloudftd"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/ftd/tier"
	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ReadDataSource(ctx context.Context, resource *DataSource, stateData *DataSourceModel) error {

	// do read
	inp := cloudftd.NewReadByNameInput(stateData.Name.ValueString())
	res, err := resource.client.ReadCloudFtdByName(ctx, inp)
	if err != nil {
		return err
	}

	// map return struct to model
	stateData.ID = types.StringValue(res.Uid)
	stateData.Name = types.StringValue(res.Name)
	stateData.AccessPolicyName = types.StringValue(res.Metadata.AccessPolicyName)
	stateData.AccessPolicyUid = types.StringValue(res.Metadata.AccessPolicyUid)
	stateData.Virtual = types.BoolValue(res.Metadata.PerformanceTier != nil)
	stateData.Licenses = util.GoStringSliceToTFStringList(strings.Split(res.Metadata.LicenseCaps, ","))
	if res.Metadata.PerformanceTier != nil { // nil means physical ftd
		stateData.PerformanceTier = types.StringValue(string(*res.Metadata.PerformanceTier))
	}
	stateData.Hostname = types.StringValue(res.Metadata.CloudManagerDomain)
	stateData.Labels = util.GoStringSliceToTFStringList(res.Tags.UngroupedTags())
	stateData.GroupedLabels = util.GoMapToStringSetTFMap(res.Tags.GroupedTags())

	return nil
}

func Read(ctx context.Context, resource *Resource, stateData *ResourceModel) error {

	// do read
	inp := cloudftd.NewReadByNameInput(stateData.Name.ValueString())
	res, err := resource.client.ReadCloudFtdByName(ctx, inp)
	if err != nil {
		return err
	}

	// handle licenses
	licenseStrings, err := license.StringToCdoStrings(res.Metadata.LicenseCaps)
	if err != nil {
		return err
	}

	// map return struct to model
	stateData.ID = types.StringValue(res.Uid)
	stateData.Name = types.StringValue(res.Name)
	stateData.AccessPolicyName = types.StringValue(res.Metadata.AccessPolicyName)
	stateData.AccessPolicyUid = types.StringValue(res.Metadata.AccessPolicyUid)
	stateData.Virtual = types.BoolValue(res.Metadata.PerformanceTier != nil)
	stateData.Licenses = util.GoStringSliceToTFStringSet(licenseStrings)
	if res.Metadata.PerformanceTier != nil { // nil means physical cloudftd
		stateData.PerformanceTier = types.StringValue(string(*res.Metadata.PerformanceTier))
	}
	stateData.GeneratedCommand = types.StringValue(res.Metadata.GeneratedCommand)
	stateData.Hostname = types.StringValue(res.Metadata.CloudManagerDomain)
	stateData.NatId = types.StringValue(res.Metadata.NatID)
	stateData.RegKey = types.StringValue(res.Metadata.RegKey)
	stateData.Labels = util.GoStringSliceToTFStringSet(res.Tags.UngroupedTags())
	stateData.GroupedLabels = util.GoMapToStringSetTFMap(res.Tags.GroupedTags())

	return nil
}

func Create(ctx context.Context, resource *Resource, planData *ResourceModel) error {

	// do create
	var performanceTier *tier.Type = nil
	if planData.PerformanceTier.ValueString() != "" {
		t, err := tier.Parse(planData.PerformanceTier.ValueString())
		performanceTier = &t
		if err != nil {
			return err
		}
	}

	// convert tf licenses to go license
	licenses, err := util.TFStringSetToLicenses(ctx, planData.Licenses)
	if err != nil {
		return err
	}

	// convert tf tags to go tags
	planTags, err := labelsFromResourceModel(ctx, planData)
	if err != nil {
		return err
	}

	createInp := cloudftd.NewCreateInput(
		planData.Name.ValueString(),
		planData.AccessPolicyName.ValueString(),
		performanceTier,
		planData.Virtual.ValueBool(),
		&licenses,
		planTags,
	)
	res, err := resource.client.CreateCloudFtd(ctx, createInp)
	if err != nil {
		return err
	}

	// convert licenses
	licenseStrings, err := license.StringToCdoStrings(res.Metadata.LicenseCaps)
	if err != nil {
		return err
	}

	// map return struct to model
	planData.ID = types.StringValue(res.Uid)
	planData.Name = types.StringValue(res.Name)
	planData.AccessPolicyName = types.StringValue(res.Metadata.AccessPolicyName)
	planData.AccessPolicyUid = types.StringValue(res.Metadata.AccessPolicyUid)
	planData.Licenses = util.GoStringSliceToTFStringSet(licenseStrings)
	planData.Labels = util.GoStringSliceToTFStringSet(res.Labels.UngroupedLabels)
	planData.GroupedLabels = util.GoMapToStringSetTFMap(res.Labels.GroupedLabels)
	if res.Metadata.PerformanceTier != nil { // nil means physical cloud ftd
		planData.PerformanceTier = types.StringValue(string(*res.Metadata.PerformanceTier))
	}
	planData.GeneratedCommand = types.StringValue(res.Metadata.GeneratedCommand)
	planData.Hostname = types.StringValue(res.Metadata.CloudManagerDomain)
	planData.NatId = types.StringValue(res.Metadata.NatID)
	planData.RegKey = types.StringValue(res.Metadata.RegKey)

	return nil
}

func Update(ctx context.Context, resource *Resource, planData *ResourceModel, stateData *ResourceModel) error {

	// do update

	// convert tf tags to go tags
	planTags, err := tagsFromResourceModel(ctx, planData)
	if err != nil {
		return err
	}

	// convert tf license to go license
	licenses, err := util.TFStringSetToLicenses(ctx, planData.Licenses)
	if err != nil {
		return err
	}

	inp := cloudftd.NewUpdateInput(
		planData.ID.ValueString(),
		planData.Name.ValueString(),
		planTags,
		licenses,
	)
	res, err := resource.client.UpdateCloudFtd(ctx, inp)
	if err != nil {
		return err
	}

	licensesStrings, err := license.StringToCdoStrings(res.Metadata.LicenseCaps)
	if err != nil {
		return err
	}
	// map return struct to model
	stateData.Name = types.StringValue(res.Name)
	stateData.Labels = planData.Labels
	stateData.GroupedLabels = planData.GroupedLabels
	stateData.Licenses = util.GoStringSliceToTFStringSet(licensesStrings)

	return nil
}

func Delete(ctx context.Context, resource *Resource, stateData *ResourceModel) error {

	// do delete
	inp := cloudftd.NewDeleteInput(stateData.ID.ValueString())
	_, err := resource.client.DeleteCloudFtd(ctx, inp)

	return err
}

func ungroupedAndGroupedLabelsFromResourceModel(ctx context.Context, resourceModel *ResourceModel) ([]string, map[string][]string, error) {
	if resourceModel == nil {
		return nil, nil, errors.New("resource model cannot be nil")
	}

	ungroupedLabels, err := util.TFStringSetToGoStringList(ctx, resourceModel.Labels)
	if err != nil {
		return nil, nil, fmt.Errorf("error while converting terraform labels to go slice, %s", resourceModel.Labels)
	}

	groupedLabels, err := util.TFMapToGoMapOfStringSlices(ctx, resourceModel.GroupedLabels)
	if err != nil {
		return nil, nil, fmt.Errorf("error while converting terraform grouped labels to go map, %v", resourceModel.GroupedLabels)
	}

	return ungroupedLabels, groupedLabels, nil
}

func tagsFromResourceModel(ctx context.Context, resourceModel *ResourceModel) (tags.Type, error) {
	ungroupedLabels, groupedLabels, err := ungroupedAndGroupedLabelsFromResourceModel(ctx, resourceModel)
	if err != nil {
		return nil, err
	}

	return tags.New(ungroupedLabels, groupedLabels), nil
}

func labelsFromResourceModel(ctx context.Context, resourceModel *ResourceModel) (publicapilabels.Type, error) {
	ungroupedLabels, groupedLabels, err := ungroupedAndGroupedLabelsFromResourceModel(ctx, resourceModel)
	if err != nil {
		return publicapilabels.Empty(), err
	}

	return publicapilabels.New(ungroupedLabels, groupedLabels), nil
}
