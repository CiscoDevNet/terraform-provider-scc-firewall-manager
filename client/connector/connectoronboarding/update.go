package connectoronboarding

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-scc-firewall-manager/go-client/internal/http"
)

type UpdateInput struct {
}

func NewUpdateInput() UpdateInput {
	return UpdateInput{}
}

type UpdateOutput struct {
}

func Update(ctx context.Context, client http.Client, updateInp UpdateInput) (*UpdateOutput, error) {

	// empty

	return nil, nil
}
