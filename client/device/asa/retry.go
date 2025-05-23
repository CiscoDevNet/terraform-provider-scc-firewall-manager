package asa

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/statemachine"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/statemachine/state"
	"strings"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/retry"
)

func UntilStateDoneAndConnectivityOk(ctx context.Context, client http.Client, uid string) retry.Func {

	return func() (bool, error) {
		readOutp, err := Read(ctx, client, *NewReadInput(uid))
		if err != nil {
			return false, err
		}

		client.Logger.Printf("device state=%s\n", readOutp.State)

		if readOutp.State == state.DONE && strings.EqualFold(readOutp.Status, "IDLE") {

			if readOutp.ConnectivityState <= 0 {
				return false, fmt.Errorf("connectivity error: %s", readOutp.ConnectivityError)
			}

			return true, nil
		}
		if readOutp.State == state.ERROR {
			return false, statemachine.NewWorkflowErrorFromDetails(readOutp.StateMachineDetails)
		}
		if readOutp.State == state.BAD_CREDENTIALS {
			return false, statemachine.NewWorkflowErrorf("Bad Credentials")
		}
		return false, nil
	}
}
