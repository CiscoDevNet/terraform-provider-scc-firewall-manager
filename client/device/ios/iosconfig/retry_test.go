package iosconfig_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/device/ios/iosconfig"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/statemachine/state"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/retry"
	internalTesting "github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/testing"
	"github.com/jarcoal/httpmock"
)

func TestIosConfigUntilState(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validIosConfig := iosconfig.ReadOutput{
		Uid:   iosConfigUid,
		State: state.DONE,
	}

	inProgressIosConfig := iosconfig.ReadOutput{
		Uid:   iosConfigUid,
		State: "SOME_INTERMEDIATE_STATE",
	}

	testCases := []struct {
		testName   string
		targetUid  string
		setupFunc  func()
		assertFunc func(err error, t *testing.T)
	}{
		{
			testName:  "successfully returns once state reaches done",
			targetUid: iosConfigUid,

			setupFunc: func() {
				configureIosConfigReadToSucceedInSubsequentCalls(iosConfigUid, []iosconfig.ReadOutput{
					inProgressIosConfig,
					inProgressIosConfig,
					validIosConfig,
				})
			},

			assertFunc: func(err error, t *testing.T) {
				assert.Nil(t, err)
				internalTesting.AssertEndpointCalledTimes("GET", buildIosConfigPath(iosConfigUid), 3, t)
			},
		},

		{
			testName:  "returns error if config state transitions to error",
			targetUid: iosConfigUid,

			setupFunc: func() {
				errorIosConfig := iosconfig.ReadOutput{
					Uid:   iosConfigUid,
					State: state.ERROR,
				}

				configureIosConfigReadToSucceedInSubsequentCalls(iosConfigUid, []iosconfig.ReadOutput{
					inProgressIosConfig,
					inProgressIosConfig,
					errorIosConfig,
				})
			},

			assertFunc: func(err error, t *testing.T) {
				assert.NotNil(t, err)
				internalTesting.AssertEndpointCalledTimes("GET", buildIosConfigPath(iosConfigUid), 3, t)
			},
		},

		{
			testName:  "return errors if config state transitions to bad credentials",
			targetUid: iosConfigUid,

			setupFunc: func() {
				badCredentialsIosConfig := iosconfig.ReadOutput{
					Uid:   iosConfigUid,
					State: state.BAD_CREDENTIALS,
				}

				configureIosConfigReadToSucceedInSubsequentCalls(iosConfigUid, []iosconfig.ReadOutput{
					inProgressIosConfig,
					inProgressIosConfig,
					badCredentialsIosConfig,
				})
			},

			assertFunc: func(err error, t *testing.T) {
				assert.NotNil(t, err)
				internalTesting.AssertEndpointCalledTimes("GET", buildIosConfigPath(iosConfigUid), 3, t)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			retryOptions := retry.DefaultOpts
			retryOptions.Delay = 1 * time.Millisecond

			err := retry.Do(context.Background(), iosconfig.UntilState(context.Background(), *http.MustNewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), testCase.targetUid, state.DONE), retryOptions)

			testCase.assertFunc(err, t)
		})
	}
}
