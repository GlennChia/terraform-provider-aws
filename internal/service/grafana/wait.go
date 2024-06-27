// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package grafana

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/service/managedgrafana"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func waitLicenseAssociationCreated(ctx context.Context, conn *managedgrafana.ManagedGrafana, id string, timeout time.Duration) (*managedgrafana.WorkspaceDescription, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{managedgrafana.WorkspaceStatusUpgrading},
		Target:  []string{managedgrafana.WorkspaceStatusActive},
		Refresh: statusWorkspace(ctx, conn, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*managedgrafana.WorkspaceDescription); ok {
		return output, err
	}

	return nil, err
}

func waitWorkspaceSAMLConfigurationCreated(ctx context.Context, conn *managedgrafana.ManagedGrafana, id string, timeout time.Duration) (*managedgrafana.SamlAuthentication, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{managedgrafana.SamlConfigurationStatusNotConfigured},
		Target:  []string{managedgrafana.SamlConfigurationStatusConfigured},
		Refresh: statusWorkspaceSAMLConfiguration(ctx, conn, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*managedgrafana.SamlAuthentication); ok {
		return output, err
	}

	return nil, err
}
