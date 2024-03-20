package common

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	armNetwork "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v5"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/nexient-llc/lcaf-component-terratest-common/lib/azure/login"
	"github.com/nexient-llc/lcaf-component-terratest-common/types"
	"github.com/stretchr/testify/assert"
)

func TestNetworkWatcher(t *testing.T, ctx types.TestContext) {

	envVarMap := login.GetEnvironmentVariables()
	subscriptionID := envVarMap["subscriptionID"]

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatalf("Unable to get credentials: %e\n", err)
	}

	options := arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: cloud.AzurePublic,
		},
	}

	networkWatcherClient, err := armNetwork.NewWatchersClient(subscriptionID, credential, &options)
	if err != nil {
		t.Fatalf("Error getting Network Water client: %v", err)
	}

	t.Run("doesNetworkWatcherExist", func(t *testing.T) {
		resourceGroupName := terraform.Output(t, ctx.TerratestTerraformOptions(), "resource_group_name")
		networkWatcherName := terraform.Output(t, ctx.TerratestTerraformOptions(), "name")
		networkWatcherId := terraform.Output(t, ctx.TerratestTerraformOptions(), "id")

		networkWatcher, err := networkWatcherClient.Get(context.Background(), resourceGroupName, networkWatcherName, nil)
		if err != nil {
			t.Fatalf("Error getting Network Water: %v", err)
		}
		if networkWatcher.Name == nil {
			t.Fatalf("Network Water does not exist")
		}

		assert.Equal(t, *networkWatcher.ID, networkWatcherId)
	})
}
