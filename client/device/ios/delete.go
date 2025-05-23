package ios

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
)

type DeleteInput = device.DeleteInput
type DeleteOutput = device.DeleteOutput

func NewDeleteInput(uid string) *DeleteInput {
	return &DeleteInput{
		Uid: uid,
	}
}

func NewDeleteRequest(ctx context.Context, client http.Client, deleteInp DeleteInput) *http.Request {
	return device.NewDeleteRequest(ctx, client, deleteInp)
}

func Delete(ctx context.Context, client http.Client, deleteInp DeleteInput) (*DeleteOutput, error) {

	client.Logger.Println("deleting ios device")

	return device.Delete(ctx, client, deleteInp)
}
