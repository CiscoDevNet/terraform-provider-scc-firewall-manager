package msp_tenant_users_test

import (
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"strings"
	"testing"
	"text/template"
)

type Users struct {
	Username    string
	Roles       []string
	ApiOnlyUser bool
}

var testMspManagedTenantUsersResource = struct {
	TenantUid string
	Users     []Users
}{
	Users: []Users{
		{Username: "user1@example.com", Roles: []string{"ROLE_SUPER_ADMIN"}, ApiOnlyUser: false},
		{Username: "example-api-user", Roles: []string{"ROLE_ADMIN"}, ApiOnlyUser: true},
	},
	TenantUid: acctest.Env.MspTenantId(),
}

// Join function to concatenate elements of a slice into a JSON array string.
func join(slice []string) string {
	quoted := make([]string, len(slice))
	for i, s := range slice {
		quoted[i] = fmt.Sprintf("%q", s) // Quotes each role to make it valid JSON
	}
	return strings.Join(quoted, ", ") // Joins with a comma
}

const testMspManagedTenantUsersTemplate = `
resource "sccfm_msp_managed_tenant_users" "test" {
	tenant_uid = "{{.TenantUid}}"
	users = [
		{
			"username": "{{(index .Users 0).Username}}"
			"roles": [{{ join (index .Users 0).Roles }}]
			"api_only_user": "{{(index .Users 0).ApiOnlyUser}}"
		},
		{
			"username": "{{(index .Users 1).Username}}"
			"roles": [{{ join (index .Users 1).Roles }}]
			"api_only_user": {{(index .Users 1).ApiOnlyUser}}
		}
	]
}`

var testMspManagedTenantUsersResourceConfig = acctest.MustParseTemplateWithFuncMap(testMspManagedTenantUsersTemplate, testMspManagedTenantUsersResource, template.FuncMap{
	"join": join,
})

func TestAccMspManagedTenantUsersResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.MspProviderConfig() + testMspManagedTenantUsersResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sccfm_msp_managed_tenant_users.test", "tenant_uid", testMspManagedTenantUsersResource.TenantUid),
					resource.TestCheckResourceAttr("sccfm_msp_managed_tenant_users.test", "users.0.username", testMspManagedTenantUsersResource.Users[0].Username),
					resource.TestCheckResourceAttr(
						"sccfm_msp_managed_tenant_users.test",
						"users.0.roles.0",
						testMspManagedTenantUsersResource.Users[0].Roles[0],
					),
					resource.TestCheckResourceAttr("sccfm_msp_managed_tenant_users.test", "users.0.api_only_user", fmt.Sprintf("%t", testMspManagedTenantUsersResource.Users[0].ApiOnlyUser)),
					resource.TestCheckResourceAttr("sccfm_msp_managed_tenant_users.test", "users.1.username", testMspManagedTenantUsersResource.Users[1].Username),
					resource.TestCheckResourceAttr(
						"sccfm_msp_managed_tenant_users.test",
						"users.1.roles.0",
						testMspManagedTenantUsersResource.Users[1].Roles[0],
					),
					resource.TestCheckResourceAttr("sccfm_msp_managed_tenant_users.test", "users.1.api_only_user", fmt.Sprintf("%t", testMspManagedTenantUsersResource.Users[1].ApiOnlyUser)),
				),
			},
		},
	})
}
