package user

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model"
)

func ReadByUid(ctx context.Context, client http.Client, readInp ReadByUidInput) (*ReadUserOutput, error) {

	readReq := NewReadByUidRequest(ctx, client, readInp.Uid)
	var userDetails model.UserDetails
	if readErr := readReq.Send(&userDetails); readErr != nil {
		return nil, readErr
	}

	return &userDetails, nil
}
