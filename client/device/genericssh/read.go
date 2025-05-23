package genericssh

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
)

type ReadInput struct {
	Uid string
}

type ReadOutput = device.ReadOutput

func NewReadInput(uid string) *ReadInput {
	return &ReadInput{
		Uid: uid,
	}
}

func Read(ctx context.Context, client http.Client, readInp ReadInput) (*ReadOutput, error) {

	client.Logger.Println("reading generic ssh")

	readOutp, err := device.ReadByUid(ctx, client, *device.NewReadByUidInput(readInp.Uid))
	if err != nil {
		return nil, err
	}

	return readOutp, nil
}
