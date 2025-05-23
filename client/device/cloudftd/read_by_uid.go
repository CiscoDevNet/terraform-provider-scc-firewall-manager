package cloudftd

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/device/tags"
)

type ReadByUidInput struct {
	Uid string `json:"uid"`
}

func NewReadByUidInput(uid string) ReadByUidInput {
	return ReadByUidInput{
		Uid: uid,
	}
}

type ReadOutput struct {
	Uid               string    `json:"uid"`
	DeviceType        string    `json:"deviceType"`
	Name              string    `json:"name"`
	Metadata          Metadata  `json:"metadata,omitempty"`
	State             string    `json:"state"`
	ConnectivityState int       `json:"connectivityState"`
	Tags              tags.Type `json:"tags"`
	SoftwareVersion   string    `json:"softwareVersion"`
}

type FtdDevice = ReadOutput

func ReadByUid(ctx context.Context, client http.Client, readInp ReadByUidInput) (*ReadOutput, error) {

	readUrl := url.ReadDevice(client.BaseUrl(), readInp.Uid)
	req := client.NewGet(ctx, readUrl)

	var readOutp ReadOutput
	if err := req.Send(&readOutp); err != nil {
		return nil, err
	}

	return &readOutp, nil
}
