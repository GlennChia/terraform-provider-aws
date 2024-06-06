// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package guardduty

import (
	"github.com/aws/aws-sdk-go-v2/service/guardduty/types"
)

var (
	errCodeAccessDeniedException = (*types.AccessDeniedException)(nil).ErrorCode()
	errCodeBadRequestException   = (*types.BadRequestException)(nil).ErrorCode()
)
