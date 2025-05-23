package asaconfig

import (
	"context"
	"encoding/json"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/crypto"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/statemachine/state"
)

type UpdateInput struct {
	SpecificUid string
	Username    string
	Password    string
	PublicKey   *model.PublicKey
	State       state.Type
}

type UpdateOutput struct {
	Uid string `json:"uid"`
}

func NewUpdateInput(specificUid string, username string, password string, publicKey *model.PublicKey, state state.Type) *UpdateInput {
	return &UpdateInput{
		SpecificUid: specificUid,
		Username:    username,
		Password:    password,
		PublicKey:   publicKey,
		State:       state,
	}
}

func Update(ctx context.Context, client http.Client, updateInp UpdateInput) (*UpdateOutput, error) {

	client.Logger.Println("updating asaconfig")

	url := url.UpdateAsaConfig(client.BaseUrl(), updateInp.SpecificUid)

	creds, err := makeCredentials(updateInp)
	if err != nil {
		return nil, err
	}

	req := client.NewPut(ctx, url, makeReqBody(creds))

	var outp UpdateOutput
	err = req.Send(&outp)
	if err != nil {
		return nil, err
	}

	return &outp, nil
}

func UpdateCredentials(ctx context.Context, client http.Client, updateInput UpdateInput) (*UpdateOutput, error) {

	client.Logger.Println("updating asaconfig credentials")

	url := url.UpdateAsaConfig(client.BaseUrl(), updateInput.SpecificUid)

	creds, err := makeCredentials(updateInput)
	if err != nil {
		return nil, err
	}

	isWaitForUserToUpdateCreds := updateInput.State == state.WAIT_FOR_USER_TO_UPDATE_CREDS || updateInput.State == state.PRE_WAIT_FOR_USER_TO_UPDATE_CREDS
	req := client.NewPut(ctx, url, makeUpdateCredentialsReqBody(isWaitForUserToUpdateCreds, creds))

	var outp UpdateOutput
	err = req.Send(&outp)
	if err != nil {
		return nil, err
	}

	return &outp, nil
}

type UpdateLocationOptions struct {
	SpecificUid string
	Location    string
}

type updateLocationRequestBody struct {
	QueueTriggerState string                         `json:"queueTriggerState"`
	SmContext         pendingLocationUpdateSmContext `json:"stateMachineContext"`
}

type pendingLocationUpdateSmContext struct {
	Ipv4                string `json:"ipv4"`
	CertificateAccepted bool   `json:"certificateAccepted"`
}

func UpdateLocation(ctx context.Context, client http.Client, options UpdateLocationOptions) (*UpdateOutput, error) {
	url := url.UpdateAsaConfig(client.BaseUrl(), options.SpecificUid)

	req := client.NewPut(ctx, url, updateLocationRequestBody{
		QueueTriggerState: "PENDING_LOCATION_UPDATE",
		SmContext: pendingLocationUpdateSmContext{
			options.Location,
			true,
		},
	})

	var outp UpdateOutput
	err := req.Send(&outp)
	if err != nil {
		return nil, err
	}

	return &outp, nil
}

func makeReqBody(creds []byte) *UpdateBody {
	return &UpdateBody{
		State:       "CERT_VALIDATED", // question: should this be hardcoded?
		Credentials: string(creds),
	}
}

func makeUpdateCredentialsReqBody(isWaitForUserToUpdateCreds bool, creds []byte) interface{} {
	if isWaitForUserToUpdateCreds {
		return &UpdateCredentialsBody{
			SmContext: SmContext{
				Credentials: string(creds),
			},
		}
	} else {
		return &UpdateCredentialsBodyWithState{
			QueueTriggerState: "WAIT_FOR_USER_TO_UPDATE_CREDS",
			SmContext: SmContext{
				Credentials: string(creds),
			},
		}
	}
}

type UpdateBody struct {
	State       string `json:"state"`
	Credentials string `json:"credentials"`
}

type UpdateCredentialsBodyWithState struct {
	QueueTriggerState string    `json:"queueTriggerState"`
	SmContext         SmContext `json:"stateMachineContext"`
}

type UpdateCredentialsBody struct {
	SmContext SmContext `json:"stateMachineContext"`
}

type SmContext struct {
	Credentials string `json:"credentials"`
}

func makeCredentials(updateInp UpdateInput) ([]byte, error) {
	if updateInp.PublicKey != nil {

		encryptedCredentials, err := crypto.EncryptCredentials(*updateInp.PublicKey, updateInp.Username, updateInp.Password)
		if err != nil {
			return nil, err
		}
		return json.Marshal(encryptedCredentials)
	}

	return json.Marshal(model.NewCredentials(updateInp.Username, updateInp.Password))
}
