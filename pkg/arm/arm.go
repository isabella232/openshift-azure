package arm

//go:generate go get github.com/golang/mock/mockgen
//go:generate mockgen -destination=../util/mocks/mock_$GOPACKAGE/arm.go -package=mock_$GOPACKAGE -source arm.go
//go:generate gofmt -s -l -w ../util/mocks/mock_$GOPACKAGE/arm.go
//go:generate go get golang.org/x/tools/cmd/goimports
//go:generate goimports -local=github.com/openshift/openshift-azure -e -w ../util/mocks/mock_$GOPACKAGE/arm.go

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-10-01/compute"
	"github.com/sirupsen/logrus"

	"github.com/openshift/openshift-azure/pkg/api"
	v15 "github.com/openshift/openshift-azure/pkg/arm/v15"
	v16 "github.com/openshift/openshift-azure/pkg/arm/v16"
	v17 "github.com/openshift/openshift-azure/pkg/arm/v17"
	v19 "github.com/openshift/openshift-azure/pkg/arm/v19"
	v20 "github.com/openshift/openshift-azure/pkg/arm/v20"
	v21 "github.com/openshift/openshift-azure/pkg/arm/v21"
)

type Interface interface {
	Generate(ctx context.Context, backupBlob string, isUpdate bool, suffix string) (map[string]interface{}, error)
	Vmss(app *api.AgentPoolProfile, backupBlob, suffix string) (*compute.VirtualMachineScaleSet, error)
	Hash(app *api.AgentPoolProfile) ([]byte, error)
}

func New(ctx context.Context, log *logrus.Entry, cs *api.OpenShiftManagedCluster, testConfig api.TestConfig) (Interface, error) {
	switch cs.Config.PluginVersion {
	case "v15.0":
		return v15.New(ctx, log, cs, testConfig), nil
	case "v16.0", "v16.1":
		return v16.New(ctx, log, cs, testConfig), nil
	case "v17.0":
		return v17.New(ctx, log, cs, testConfig), nil
	case "v19.0":
		return v19.New(ctx, log, cs, testConfig), nil
	case "v20.0":
		return v20.New(ctx, log, cs, testConfig), nil
	case "v21.0":
		return v21.New(ctx, log, cs, testConfig), nil
	}

	return nil, fmt.Errorf("version %q not found", cs.Config.PluginVersion)
}
