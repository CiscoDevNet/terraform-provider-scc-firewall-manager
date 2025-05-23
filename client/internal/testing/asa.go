package testing

import (
	"fmt"
	"time"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/device/asa"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/device/publicapilabels"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/device/tags"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/devicetype"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/statemachine"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/statemachine/state"
)

func (m Model) AsaReadOutput() asa.ReadOutput {
	return asa.ReadOutput{
		Uid:                 m.AsaUid.String(),
		Name:                m.AsaName,
		CreatedDate:         m.AsaCreatedDate.Unix(),
		LastUpdatedDate:     time.Now().Unix(),
		DeviceType:          devicetype.Asa,
		ConnectorUid:        m.CdgUid.String(),
		ConnectorType:       "CDG",
		SocketAddress:       fmt.Sprintf("%s:%s", m.AsaHost, m.AsaPort),
		Port:                m.AsaPort,
		Host:                m.AsaHost,
		Tags:                tags.Type{},
		IgnoreCertificate:   false,
		ConnectivityState:   0,
		ConnectivityError:   "",
		State:               state.DONE,
		Status:              "",
		StateMachineDetails: statemachine.Details{},
	}
}

func (m Model) AsaReadSpecificDeviceOutput() asa.ReadSpecificOutput {
	return asa.ReadSpecificOutput{
		SpecificUid: m.AsaUid.String(),
		Metadata: asa.SpecificDeviceMetadata{
			AsdmVersion: "7.6(2)",
		},
	}
}

func (m Model) AsaCreateInput() asa.CreateInput {
	return asa.CreateInput{
		Name:              m.AsaName,
		ConnectorUid:      m.CdgUid.String(),
		ConnectorType:     "CDG",
		SocketAddress:     fmt.Sprintf("%s:%s", m.AsaHost, m.AsaPort),
		Labels:            publicapilabels.Empty(),
		Username:          m.AsaUsername,
		Password:          m.AsaPassword,
		IgnoreCertificate: false,
	}
}
