package cdfmc

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sccFwMgrClient "github.com/CiscoDevnet/terraform-provider-scc-firewall-manager/go-client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &Resource{}

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	client *sccFwMgrClient.Client
}

type ResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Hostname   types.String `tfsdk:"hostname"`
	DomainUuid types.String `tfsdk:"domain_uuid"`
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdfmc"
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Use this resource to create a Cloud-delivered FMC (cdFMC) in your tenant. You can have only one cdFMC in your tenant. In addition, you cannot delete a created cdFMC using this resource; to delete your cdFMC, please contact Cisco support.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the cdFMC. This is automatically generated by SCC Firewall Manager.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the cdFMC. This is automatically generated from the tenant name by SCC Firewall Manager.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "The hostname of the cdFMC.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_uuid": schema.StringAttribute{
				MarkdownDescription: "The domain UUID of the cdFMC.",
				Computed:            true,
			},
		},
	}
}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sccFwMgrClient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *sccFwMgrClient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	tflog.Trace(ctx, "read cdfmc resource")

	// 1. read state data
	var stateData ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 2. do read
	if err := Read(ctx, r, &stateData); err != nil {
		resp.Diagnostics.AddError("failed to read cdfmc resource", err.Error())
	}

	// 3. save data into terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {

	tflog.Trace(ctx, "create cdfmc resource")

	// 1. read plan data into planData
	var planData ResourceModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if res.Diagnostics.HasError() {
		return
	}

	// 2. use plan data to create device and fill up rest of the model
	if err := Create(ctx, r, &planData); err != nil {
		res.Diagnostics.AddError("failed to create cdfmc resource", err.Error())
		return
	}

	// 3. set state using filled model
	res.Diagnostics.Append(res.State.Set(ctx, &planData)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	// update is noop
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	// delete is noop
	res.Diagnostics.AddWarning("Delete cdFMC is an noop", "Please reach out to SCC Firewall Manager TAC if you really want to delete a cdFMC.")
}
