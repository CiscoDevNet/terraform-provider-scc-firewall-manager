package device

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/device/tags"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/url"
)

type UpdateInput struct {
	Uid string `json:"-"`

	Name              string    `json:"name,omitempty"`
	Tags              tags.Type `json:"tags,omitempty"`
	IgnoreCertificate bool      `json:"ignoreCertificate,omitempty"`
}

type UpdateOutput = ReadOutput

func NewUpdateInput(uid string, name string, ignoreCertificate bool, tags tags.Type) *UpdateInput {
	return &UpdateInput{
		Uid:               uid,
		Name:              name,
		IgnoreCertificate: ignoreCertificate,
		Tags:              tags,
	}
}

func NewUpdateRequest(ctx context.Context, client http.Client, updateReq UpdateInput) *http.Request {
	url := url.UpdateDevice(client.BaseUrl(), updateReq.Uid)

	req := client.NewPut(ctx, url, updateReq)

	return req
}

func Update(ctx context.Context, client http.Client, updateReq UpdateInput) (*UpdateOutput, error) {

	client.Logger.Println("updating device")

	req := NewUpdateRequest(ctx, client, updateReq)

	var outp UpdateOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
