
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationsignals

import (
    "context"
    "errors"

    "github.com/YakDriver/regexache"
    "github.com/aws/aws-sdk-go-v2/service/applicationsignals"
    awstypes "github.com/aws/aws-sdk-go-v2/service/applicationsignals/types"
    "github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
    "github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
    "github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
    "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-provider-aws/internal/create"
    "github.com/hashicorp/terraform-provider-aws/internal/framework"
    "github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
    fwtypes "github.com/hashicorp/terraform-provider-aws/internal/framework/types"
    tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
    "github.com/hashicorp/terraform-provider-aws/internal/tfresource"
    "github.com/hashicorp/terraform-provider-aws/names"
)

// @FrameworkResource(name="Service Level Objective")
// @Tags(identifierAttribute="arn")
func newResourceServiceLevelObjective(_ context.Context) (resource.ResourceWithConfigure, error) {
    r := &resourceServiceLevelObjective{}

    return r, nil
}

const (
    ResNameServiceLevelObjective = "Service Level Objective"
)

type resourceServiceLevelObjective struct {
    framework.ResourceWithConfigure
}

func (r *resourceServiceLevelObjective) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = "aws_applicationsignals_service_level_objective"
}

func (r *resourceServiceLevelObjective) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
    }
}

func (r *resourceServiceLevelObjective) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    conn := r.Meta().ApplicationSignalsClient(ctx)
    var plan serviceLevelObjectiveResourceModel

    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

    input := &applicationsignals.CreateServiceLevelObjectiveInput{}
    resp.Diagnostics.Append(flex.Expand(ctx, plan, input)...)

    if resp.Diagnostics.HasError() {
        return
    }

    input.Tags = getTagsIn(ctx)

    var out *applicationsignals.CreateServiceLevelObjectiveOutput

    state := plan
    state.ID = flex.StringToFramework(ctx, out.Slo.Arn)

    // Read after create to get computed attributes omitted from the create response
    readOut, err := FindServiceLevelObjectiveByID(ctx, conn, state.ID.ValueString())
    if err != nil {
        resp.Diagnostics.AddError(
            create.ProblemStandardMessage(names.SSOAdmin, create.ErrActionCreating, ResNameServiceLevelObjective, plan.ID.String(), err),
            err.Error(),
        )
        return
    }
    resp.Diagnostics.Append(flex.Flatten(ctx, readOut, &state)...)

    if resp.Diagnostics.HasError() {
        return
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *resourceServiceLevelObjective) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    conn := r.Meta().ApplicationSignalsClient(ctx)
    var state serviceLevelObjectiveResourceModel

    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

    if resp.Diagnostics.HasError() {
        return
    }

    out, err := FindServiceLevelObjectiveByID(ctx, conn, state.ID.ValueString())

    if tfresource.NotFound(err) {
        resp.State.RemoveResource(ctx)
        return
    }

    if err != nil {
        resp.Diagnostics.AddError(
            create.ProblemStandardMessage(names.ApplicationSignals, create.ErrActionSetting, ResNameServiceLevelObjective, state.ID.ValueString(), err),
            err.Error(),
        )
        return
    }

    resp.Diagnostics.Append(flex.Flatten(ctx, out, &state)...)

    if resp.Diagnostics.HasError() {
        return
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *resourceServiceLevelObjective) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    conn := r.Meta().ApplicationSignalsClient(ctx)
    var state, plan serviceLevelObjectiveResourceModel

    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

    if resp.Diagnostics.HasError() {
        return
    }

    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

    if resp.Diagnostics.HasError() {
        return
    }

    if serviceLevelObjectiveHasChanges(ctx, plan, state) {
        input := &applicationsignals.UpdateServiceLevelObjectiveInput{}

        resp.Diagnostics.Append(flex.Expand(ctx, plan, input)...)

        if resp.Diagnostics.HasError() {
            return
        }

        input.Id = flex.StringFromFramework(ctx, state.ID)

        _, err := conn.UpdateServiceLevelObjective(ctx, input)

        if err != nil {
            resp.Diagnostics.AddError(
                create.ProblemStandardMessage(names.ApplicationSignals, create.ErrActionUpdating, ResNameServiceLevelObjective, state.ID.ValueString(), err),
                err.Error(),
            )
            return
        }
    }

    out, err := FindServiceLevelObjectiveByID(ctx, conn, state.ID.ValueString())

    if err != nil {
        resp.Diagnostics.AddError(
            create.ProblemStandardMessage(names.ApplicationSignals, create.ErrActionUpdating, ResNameServiceLevelObjective, state.ID.ValueString(), err),
            err.Error(),
        )
        return
    }

    resp.Diagnostics.Append(flex.Flatten(ctx, out, &plan)...)

    if resp.Diagnostics.HasError() {
        return
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *resourceServiceLevelObjective) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    conn := r.Meta().ApplicationSignalsClient(ctx)

    var state serviceLevelObjectiveResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    _, err := conn.DeleteServiceLevelObjective(ctx, &applicationsignals.DeleteServiceLevelObjectiveInput{
        Id: state.ID.ValueStringPointer(),
    })
    if err != nil {
        var nfe *awstypes.ResourceNotFoundException
        if errors.As(err, &nfe) {
            return
        }
        resp.Diagnostics.AddError(
            create.ProblemStandardMessage(names.ApplicationSignals, create.ErrActionDeleting, ResNameServiceLevelObjective, state.ID.String(), nil),
            err.Error(),
        )
    }
}

func (r *resourceServiceLevelObjective) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root(names.AttrID), req, resp)
}

func (r *resourceServiceLevelObjective) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
    r.SetTagsAll(ctx, req, resp)
}

type serviceLevelObjectiveResourceModel struct {
    
}

func serviceLevelObjectiveHasChanges(_ context.Context, plan, state serviceLevelObjectiveResourceModel) bool {
    return !plan.Actions.Equal(state.Actions) ||
        !plan.ProtectedResource.Equal(state.ProtectedResource) ||
        !plan.Role.Equal(state.Role)
}
