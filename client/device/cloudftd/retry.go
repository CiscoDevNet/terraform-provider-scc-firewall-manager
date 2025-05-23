package cloudftd

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/retry"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/statemachine/state"
)

func UntilGeneratedCommandAvailable(ctx context.Context, client http.Client, uid string, metadata *Metadata) retry.Func {

	return func() (bool, error) {
		client.Logger.Println("checking if FTD generated command is available")

		readOutp, err := ReadByUid(ctx, client, NewReadByUidInput(uid))
		if err != nil {
			return false, err
		}

		client.Logger.Printf("device metadata=%v\n", readOutp.Metadata)

		if readOutp.Metadata.GeneratedCommand != "" {
			*metadata = readOutp.Metadata
			return true, nil
		} else {
			return false, fmt.Errorf("generated command not found in metadata: %+v", readOutp.Metadata)
		}
	}
}

func UntilSpecificStateDone(ctx context.Context, client http.Client, inp ReadSpecificInput) retry.Func {
	return func() (bool, error) {
		client.Logger.Println("check FTD specific device state")

		readOutp, err := ReadSpecific(ctx, client, inp)
		if err != nil {
			return false, err
		}

		if readOutp.State == state.DONE {
			return true, nil
		} else if readOutp.State == state.ERROR {
			return false, fmt.Errorf("workflow ended in error")
		} else {
			return false, nil
		}
	}
}
