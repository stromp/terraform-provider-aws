// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package networkfirewall

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/networkfirewall"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_networkfirewall_resource_policy")
func ResourceResourcePolicy() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceResourcePolicyPut,
		ReadWithoutTimeout:   resourceResourcePolicyRead,
		UpdateWithoutTimeout: resourceResourcePolicyPut,
		DeleteWithoutTimeout: resourceResourcePolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			names.AttrPolicy: {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: verify.SuppressEquivalentPolicyDiffs,
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
			},
			names.AttrResourceARN: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
		},
	}
}

func resourceResourcePolicyPut(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).NetworkFirewallConn(ctx)
	resourceArn := d.Get(names.AttrResourceARN).(string)

	policy, err := structure.NormalizeJsonString(d.Get(names.AttrPolicy).(string))

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "policy (%s) is invalid JSON: %s", policy, err)
	}

	input := &networkfirewall.PutResourcePolicyInput{
		ResourceArn: aws.String(resourceArn),
		Policy:      aws.String(policy),
	}

	log.Printf("[DEBUG] Putting NetworkFirewall Resource Policy for resource: %s", resourceArn)

	_, err = conn.PutResourcePolicyWithContext(ctx, input)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "putting NetworkFirewall Resource Policy (for resource: %s): %s", resourceArn, err)
	}

	d.SetId(resourceArn)

	return append(diags, resourceResourcePolicyRead(ctx, d, meta)...)
}

func resourceResourcePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).NetworkFirewallConn(ctx)
	resourceArn := d.Id()

	log.Printf("[DEBUG] Reading NetworkFirewall Resource Policy for resource: %s", resourceArn)

	policy, err := FindResourcePolicy(ctx, conn, resourceArn)
	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, networkfirewall.ErrCodeResourceNotFoundException) {
		log.Printf("[WARN] NetworkFirewall Resource Policy (for resource: %s) not found, removing from state", resourceArn)
		d.SetId("")
		return diags
	}
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading NetworkFirewall Resource Policy (for resource: %s): %s", resourceArn, err)
	}

	if policy == nil {
		return sdkdiag.AppendErrorf(diags, "reading NetworkFirewall Resource Policy (for resource: %s): empty output", resourceArn)
	}

	d.Set(names.AttrResourceARN, resourceArn)

	policyToSet, err := verify.PolicyToSet(d.Get(names.AttrPolicy).(string), aws.StringValue(policy))

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "setting policy %s: %s", aws.StringValue(policy), err)
	}

	d.Set(names.AttrPolicy, policyToSet)

	return diags
}

func resourceResourcePolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	const (
		timeout = 2 * time.Minute
	)
	conn := meta.(*conns.AWSClient).NetworkFirewallConn(ctx)

	log.Printf("[DEBUG] Deleting NetworkFirewall Resource Policy: %s", d.Id())
	_, err := tfresource.RetryWhenAWSErrMessageContains(ctx, timeout, func() (interface{}, error) {
		return conn.DeleteResourcePolicyWithContext(ctx, &networkfirewall.DeleteResourcePolicyInput{
			ResourceArn: aws.String(d.Id()),
		})
	}, networkfirewall.ErrCodeInvalidResourcePolicyException, "The supplied policy does not match RAM managed permissions")

	if tfawserr.ErrCodeEquals(err, networkfirewall.ErrCodeResourceNotFoundException) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting NetworkFirewall Resource Policy (%s): %s", d.Id(), err)
	}

	return diags
}
