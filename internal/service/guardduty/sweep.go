// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package guardduty

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/guardduty"
	"github.com/hashicorp/aws-sdk-go-base/v2/tfawserr"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep/awsv1"
)

func RegisterSweepers() {
	resource.AddTestSweepers("aws_guardduty_detector", &resource.Sweeper{
		Name:         "aws_guardduty_detector",
		F:            sweepDetectors,
		Dependencies: []string{"aws_guardduty_publishing_destination"},
	})

	resource.AddTestSweepers("aws_guardduty_publishing_destination", &resource.Sweeper{
		Name: "aws_guardduty_publishing_destination",
		F:    sweepPublishingDestinations,
	})
}

func sweepDetectors(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(ctx, region)

	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}

	conn := client.GuardDutyClient(ctx)
	input := &guardduty.ListDetectorsInput{}
	var sweeperErrs *multierror.Error

	err = conn.ListDetectorsPagesWithContext(ctx, input, func(page *guardduty.ListDetectorsOutput, lastPage bool) bool {
		for _, detectorID := range page.DetectorIds {
			input := &guardduty.DeleteDetectorInput{
				DetectorId: aws.String(detectorID),
			}

			log.Printf("[INFO] Deleting GuardDuty Detector: %s", detectorID)
			_, err := conn.DeleteDetector(ctx, input)
			if tfawserr.ErrCodeContains(err, "AccessDenied") {
				log.Printf("[WARN] Skipping GuardDuty Detector (%s): %s", detectorID, err)
				continue
			}
			if err != nil {
				sweeperErr := fmt.Errorf("error deleting GuardDuty Detector (%s): %w", detectorID, err)
				log.Printf("[ERROR] %s", sweeperErr)
				sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
			}
		}

		return !lastPage
	})

	if awsv1.SkipSweepError(err) {
		log.Printf("[WARN] Skipping GuardDuty Detector sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error retrieving GuardDuty Detectors: %w", err)
	}

	return sweeperErrs.ErrorOrNil()
}

func sweepPublishingDestinations(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(ctx, region)

	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	conn := client.GuardDutyClient(ctx)
	var sweeperErrs *multierror.Error

	detect_input := &guardduty.ListDetectorsInput{}

	err = conn.ListDetectorsPagesWithContext(ctx, detect_input, func(page *guardduty.ListDetectorsOutput, lastPage bool) bool {
		for _, detectorID := range page.DetectorIds {
			list_input := &guardduty.ListPublishingDestinationsInput{
				DetectorId: aws.String(detectorID),
			}

			err = conn.ListPublishingDestinationsPagesWithContext(ctx, list_input, func(page *guardduty.ListPublishingDestinationsOutput, lastPage bool) bool {
				for _, destination_element := range page.Destinations {
					input := &guardduty.DeletePublishingDestinationInput{
						DestinationId: destination_element.DestinationId,
						DetectorId:    aws.String(detectorID),
					}

					log.Printf("[INFO] Deleting GuardDuty Publishing Destination: %s", *destination_element.DestinationId)
					_, err := conn.DeletePublishingDestination(ctx, input)

					if err != nil {
						sweeperErr := fmt.Errorf("error deleting GuardDuty Publishing Destination (%s): %w", *destination_element.DestinationId, err)
						log.Printf("[ERROR] %s", sweeperErr)
						sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
					}
				}
				return !lastPage
			})
		}
		return !lastPage
	})

	if err != nil {
		sweeperErr := fmt.Errorf("Error receiving Guardduty detectors for publishing sweep : %w", err)
		log.Printf("[ERROR] %s", sweeperErr)
		sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
	}

	if awsv1.SkipSweepError(err) {
		log.Printf("[WARN] Skipping GuardDuty Publishing Destination sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error retrieving GuardDuty Publishing Destinations: %s", err)
	}

	return sweeperErrs.ErrorOrNil()
}
