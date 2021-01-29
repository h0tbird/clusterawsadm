package main

//-----------------------------------------------------------------------------
// Imports
//-----------------------------------------------------------------------------

import (

	// stdlib
	"context"
	"io/ioutil"
	"log"

	// community
	tg "github.com/h0tbird/terrago"
	"github.com/sirupsen/logrus"

	// terraform-plugin-sdk
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	// providers
	"github.com/terraform-providers/terraform-provider-aws/aws"
)

//-----------------------------------------------------------------------------
// Init
//-----------------------------------------------------------------------------

func init() {
	// TODO: replace logrus with zap logger
	log.SetOutput(ioutil.Discard)
}

//-----------------------------------------------------------------------------
// Main
//-----------------------------------------------------------------------------

func main() {

	ctx := context.Background()
	m := tg.NewManifest()
	s := &state{}

	//----------------------------
	// Create the state directory
	//----------------------------

	if err := s.Init(); err != nil {
		logrus.Fatal(err)
	}

	//------------------------
	// Configure the provider
	//------------------------

	p := aws.Provider()
	logrus.WithFields(logrus.Fields{"region": "us-east-2"}).Info("Configuring the provider")
	diags := p.Configure(ctx, &terraform.ResourceConfig{
		Config: map[string]interface{}{
			"region": "us-east-2",
		},
	})

	if diags != nil && diags.HasError() {
		for _, d := range diags {
			if d.Severity == diag.Error {
				logrus.Fatalf("error configuring the provider: %s", d.Summary)
			}
		}
	}

	//--------------------------------------------
	// nodes.cluster-api-provider-aws.sigs.k8s.io
	//--------------------------------------------

	// AWS::IAM::Policy
	m.Resources["nodesPolicy"] = &tg.Resource{
		ResourceLogicalID: "NodesPolicy",
		ResourceType:      "aws_iam_policy",
		ResourceConfig: map[string]interface{}{
			"name":        "nodes.cluster-api-provider-aws.sigs.k8s.io",
			"description": "For the Kubernetes Cloud Provider AWS nodes",
			"policy":      nodesPolicy,
		},
	}

	// AWS::IAM::Role
	m.Resources["nodesRole"] = &tg.Resource{
		ResourceLogicalID: "NodesRole",
		ResourceType:      "aws_iam_role",
		ResourceConfig: map[string]interface{}{
			"name":               "nodes.cluster-api-provider-aws.sigs.k8s.io",
			"assume_role_policy": assumeRolePolicy,
		},
	}

	// AWS::IAM::RolePolicyAttachment
	m.Resources["nodesRoleToNodesPolicyAttachment"] = &tg.Resource{
		ResourceLogicalID: "NodesRoleToNodesPolicyAttachment",
		ResourceType:      "aws_iam_role_policy_attachment",
		ResourceConfig: map[string]interface{}{
			"role":       "nodesRole.ResourceConfig.name",
			"policy_arn": "nodesPolicy.ResourceState.ID",
		},
	}

	// AWS::IAM::InstanceProfile
	m.Resources["nodesInstanceProfile"] = &tg.Resource{
		ResourceLogicalID: "NodesInstanceProfile",
		ResourceType:      "aws_iam_instance_profile",
		ResourceConfig: map[string]interface{}{
			"name": "nodes.cluster-api-provider-aws.sigs.k8s.io",
			"role": "nodesRole.ResourceConfig.name",
		},
	}

	//--------------------------------------------------
	// controllers.cluster-api-provider-aws.sigs.k8s.io
	//--------------------------------------------------

	// AWS::IAM::Policy
	m.Resources["controllersPolicy"] = &tg.Resource{
		ResourceLogicalID: "ControllersPolicy",
		ResourceType:      "aws_iam_policy",
		ResourceConfig: map[string]interface{}{
			"name":        "controllers.cluster-api-provider-aws.sigs.k8s.io",
			"description": "For the Kubernetes Cluster API Provider AWS Controllers",
			"policy":      controllersPolicy,
		},
	}

	// AWS::IAM::Role
	m.Resources["controllersRole"] = &tg.Resource{
		ResourceLogicalID: "ControllersRole",
		ResourceType:      "aws_iam_role",
		ResourceConfig: map[string]interface{}{
			"name":               "controllers.cluster-api-provider-aws.sigs.k8s.io",
			"assume_role_policy": assumeRolePolicy,
		},
	}

	// AWS::IAM::RolePolicyAttachment
	m.Resources["controllersRoleToControllersPolicyAttachment"] = &tg.Resource{
		ResourceLogicalID: "ControllersRoleToControllersPolicyAttachment",
		ResourceType:      "aws_iam_role_policy_attachment",
		ResourceConfig: map[string]interface{}{
			"role":       "controllersRole.ResourceConfig.name",
			"policy_arn": "controllersPolicy.ResourceState.ID",
		},
	}

	// AWS::IAM::InstanceProfile
	m.Resources["controllersInstanceProfile"] = &tg.Resource{
		ResourceLogicalID: "ControllersInstanceProfile",
		ResourceType:      "aws_iam_instance_profile",
		ResourceConfig: map[string]interface{}{
			"name": "controllers.cluster-api-provider-aws.sigs.k8s.io",
			"role": "controllersRole.ResourceConfig.name",
		},
	}

	//----------------------------------------------------
	// control-plane.cluster-api-provider-aws.sigs.k8s.io
	//----------------------------------------------------

	// AWS::IAM::Policy
	m.Resources["controlPlanePolicy"] = &tg.Resource{
		ResourceLogicalID: "ControlPlanePolicy",
		ResourceType:      "aws_iam_policy",
		ResourceConfig: map[string]interface{}{
			"name":        "control-plane.cluster-api-provider-aws.sigs.k8s.io",
			"description": "For the Kubernetes Cloud Provider AWS Control Plane",
			"policy":      controlPlanePolicy,
		},
	}

	// AWS::IAM::Role
	m.Resources["controlPlaneRole"] = &tg.Resource{
		ResourceLogicalID: "ControlPlaneRole",
		ResourceType:      "aws_iam_role",
		ResourceConfig: map[string]interface{}{
			"name":               "control-plane.cluster-api-provider-aws.sigs.k8s.io",
			"assume_role_policy": assumeRolePolicy,
		},
	}

	// AWS::IAM::RolePolicyAttachment
	m.Resources["controlPlaneRoleToControlPlanePolicyAttachment"] = &tg.Resource{
		ResourceLogicalID: "ControlPlaneRoleToControlPlanePolicyAttachment",
		ResourceType:      "aws_iam_role_policy_attachment",
		ResourceConfig: map[string]interface{}{
			"role":       "controlPlaneRole.ResourceConfig.name",
			"policy_arn": "controlPlanePolicy.ResourceState.ID",
		},
	}

	// AWS::IAM::RolePolicyAttachment
	m.Resources["controlPlaneRoleToNodesPolicyAttachment"] = &tg.Resource{
		ResourceLogicalID: "ControlPlaneRoleToNodesPolicyAttachment",
		ResourceType:      "aws_iam_role_policy_attachment",
		ResourceConfig: map[string]interface{}{
			"role":       "controlPlaneRole.ResourceConfig.name",
			"policy_arn": "nodesPolicy.ResourceState.ID",
		},
	}

	// AWS::IAM::RolePolicyAttachment
	m.Resources["controlPlaneRoleToControllersPolicyAttachment"] = &tg.Resource{
		ResourceLogicalID: "ControlPlaneRoleToControllersPolicyAttachment",
		ResourceType:      "aws_iam_role_policy_attachment",
		ResourceConfig: map[string]interface{}{
			"role":       "controlPlaneRole.ResourceConfig.name",
			"policy_arn": "controllersPolicy.ResourceState.ID",
		},
	}

	// AWS::IAM::InstanceProfile
	m.Resources["controlPlaneInstanceProfile"] = &tg.Resource{
		ResourceLogicalID: "ControlPlaneInstanceProfile",
		ResourceType:      "aws_iam_instance_profile",
		ResourceConfig: map[string]interface{}{
			"name": "control-plane.cluster-api-provider-aws.sigs.k8s.io",
			"role": "controlPlaneRole.ResourceConfig.name",
		},
	}

	//--------------------
	// Apply the manifest
	//--------------------

	if err := m.Apply(ctx, p, s); err != nil {
		logrus.Fatalf("error applying the manifest: %s", err)
	}
}
