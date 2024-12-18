// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sqs

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/json"
	"github.com/hashicorp/terraform-provider-aws/internal/sdkv2"
)

// @SDKResource("aws_sqs_queue_redrive_policy")
func resourceQueueRedrivePolicy() *schema.Resource {
	h := &queueAttributeHandler{
		AttributeName: types.QueueAttributeNameRedrivePolicy,
		SchemaKey:     "redrive_policy",
		ToSet: func(old, new string) (string, error) {
			if json.EqualStrings(old, new) {
				return old, nil
			}
			return new, nil
		},
	}

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"queue_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"redrive_policy": sdkv2.PolicyDocumentSchema(true),
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateWithoutTimeout: h.Upsert,
		ReadWithoutTimeout:   h.Read,
		UpdateWithoutTimeout: h.Upsert,
		DeleteWithoutTimeout: h.Delete,
	}
}
