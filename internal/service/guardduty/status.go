// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package guardduty

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/guardduty"
	awstypes "github.com/aws/aws-sdk-go-v2/service/guardduty/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	// AdminStatus NotFound
	adminStatusNotFound = "NotFound"

	// AdminStatus Unknown
	adminStatusUnknown = "Unknown"

	// Constants not currently provided by the AWS Go SDK
	publishingStatusFailed  = "Failed"
	publishingStatusUnknown = "Unknown"
)

// statusAdminAccountAdmin fetches the AdminAccount and its AdminStatus
func statusAdminAccountAdmin(ctx context.Context, conn *guardduty.Client, adminAccountID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		adminAccount, err := getOrganizationAdminAccount(ctx, conn, adminAccountID)

		if err != nil {
			return nil, adminStatusUnknown, err
		}

		if adminAccount == nil {
			return adminAccount, adminStatusNotFound, nil
		}

		return adminAccount, string(adminAccount.AdminStatus), nil
	}
}

// statusPublishingDestination fetches the PublishingDestination and its Status
func statusPublishingDestination(ctx context.Context, conn *guardduty.Client, destinationID, detectorID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &guardduty.DescribePublishingDestinationInput{
			DetectorId:    aws.String(detectorID),
			DestinationId: aws.String(destinationID),
		}

		output, err := conn.DescribePublishingDestination(ctx, input)

		if err != nil {
			return output, publishingStatusFailed, err
		}

		if output == nil {
			return output, publishingStatusUnknown, nil
		}

		return output, aws.ToString(output.Status), nil
	}
}

// TODO: Migrate to shared internal package guardduty
func getOrganizationAdminAccount(ctx context.Context, conn *guardduty.Client, adminAccountID string) (*awstypes.AdminAccount, error) {
	input := &guardduty.ListOrganizationAdminAccountsInput{}
	var result *awstypes.AdminAccount

	err := conn.ListOrganizationAdminAccountsPagesWithContext(ctx, input, func(page *guardduty.ListOrganizationAdminAccountsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, adminAccount := range page.AdminAccounts {
			if adminAccount == nil {
				continue
			}

			if aws.ToString(adminAccount.AdminAccountId) == adminAccountID {
				result = adminAccount
				return false
			}
		}

		return !lastPage
	})

	return result, err
}
