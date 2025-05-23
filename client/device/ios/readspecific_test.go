package ios_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/device/ios"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/statemachine/state"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/jarcoal/httpmock"
)

func TestIosReadSpecific(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	deviceUid := "00000000-0000-0000-0000-000000000000"

	specificDevice := ios.ReadSpecificOutput{
		SpecificUid: "11111111-1111-1111-1111-111111111111",
		State:       state.DONE,
		Namespace:   "targets",
		Type:        "device",
	}

	testCases := []struct {
		testName   string
		input      ios.ReadSpecificInput
		setupFunc  func()
		assertFunc func(output *ios.ReadSpecificOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully reads iOS specific device",
			input: ios.ReadSpecificInput{
				Uid: deviceUid,
			},

			setupFunc: func() {
				configureDeviceReadSpecificToRespondSuccessfully(deviceUid, device.ReadSpecificOutput(specificDevice))
			},

			assertFunc: func(output *ios.ReadSpecificOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, specificDevice, *output)
			},
		},

		{
			testName: "returns error when the remote service reading the iOS specific device encounters an issue",
			input: ios.ReadSpecificInput{
				Uid: deviceUid,
			},

			setupFunc: func() {
				configureDeviceReadSpecificToRespondWithError(deviceUid)
			},

			assertFunc: func(output *ios.ReadSpecificOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := ios.ReadSpecific(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
