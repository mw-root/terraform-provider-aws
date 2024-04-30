// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ec2

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @FrameworkResource("aws_vpc_endpoint_service_private_dns_verification", name="Endpoint Service Private DNS Verification")
func newResourceEndpointServicePrivateDNSVerification(_ context.Context) (resource.ResourceWithConfigure, error) {
	return &resourceEndpointServicePrivateDNSVerification{}, nil
}

const (
	ResNameEndpointServicePrivateDNSVerification = "Endpoint Service Private DNS Verification"
)

type resourceEndpointServicePrivateDNSVerification struct {
	framework.ResourceWithConfigure
}

func (r *resourceEndpointServicePrivateDNSVerification) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "aws_vpc_endpoint_service_private_dns_verification"
}

func (r *resourceEndpointServicePrivateDNSVerification) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"service_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *resourceEndpointServicePrivateDNSVerification) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	conn := r.Meta().EC2Client(ctx)

	var plan resourceEndpointServicePrivateDNSVerificationData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	in := &ec2.StartVpcEndpointServicePrivateDnsVerificationInput{
		ServiceId: aws.String(plan.ServiceID.ValueString()),
	}

	out, err := conn.StartVpcEndpointServicePrivateDnsVerification(ctx, in)
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.EC2, create.ErrActionCreating, ResNameEndpointServicePrivateDNSVerification, plan.ServiceID.String(), err),
			err.Error(),
		)
		return
	}
	if out == nil || out.ReturnValue == nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.EC2, create.ErrActionCreating, ResNameEndpointServicePrivateDNSVerification, plan.ServiceID.String(), nil),
			errors.New("empty output").Error(),
		)
		return
	}
	if !aws.ToBool(out.ReturnValue) {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.EC2, create.ErrActionCreating, ResNameEndpointServicePrivateDNSVerification, plan.ServiceID.String(), nil),
			errors.New("request failed").Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *resourceEndpointServicePrivateDNSVerification) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read is a no-op
}

func (r *resourceEndpointServicePrivateDNSVerification) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Update is a no-op
}

func (r *resourceEndpointServicePrivateDNSVerification) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Delete is a no-op
}

type resourceEndpointServicePrivateDNSVerificationData struct {
	ServiceID types.String `tfsdk:"service_id"`
}
