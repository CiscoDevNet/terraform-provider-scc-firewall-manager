package msp_tenant_users

import "github.com/hashicorp/terraform-plugin-framework/types"

type MspManagedTenantUsersResourceModel struct {
	TenantUid types.String `tfsdk:"tenant_uid"`
	Users     []User       `tfsdk:"users"`
}

type User struct {
	Id                           types.String `tfsdk:"id"`
	Username                     types.String `tfsdk:"username"`
	UsernameInSccFirewallManager types.String `tfsdk:"username_in_scc_firewall_manager"`
	Roles                        types.List   `tfsdk:"roles"`
	ApiOnlyUser                  types.Bool   `tfsdk:"api_only_user"`
}
