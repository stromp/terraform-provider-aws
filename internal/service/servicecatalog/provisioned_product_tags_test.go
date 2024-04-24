// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicecatalog_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/servicecatalog"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccServiceCatalogProvisionedProduct_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var v servicecatalog.ProvisionedProductDetail
	resourceName := "aws_servicecatalog_provisioned_product.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ServiceCatalogServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProvisionedProductDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccProvisionedProductConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
				},
			},
			{
				Config: testAccProvisionedProductConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
				},
			},
			// {
			// 	Config: testAccProvisionedProductConfig_tags1(rName, "key2", "value2"),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		testAccCheckProvisionedProductExists(ctx, resourceName, &v),
			// 		resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
			// 		resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
			// 	),
			// },
			// {
			// 	ResourceName:      resourceName,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	ImportStateVerifyIgnore: []string{
			// 		"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
			// 	},
			// },
			// {
			// 	Config: testAccProvisionedProductConfig_tags0(rName),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		testAccCheckProvisionedProductExists(ctx, resourceName, &v),
			// 		resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
			// 	),
			// },
			// {
			// 	ResourceName:      resourceName,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	ImportStateVerifyIgnore: []string{
			// 		"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
			// 	},
			// },
		},
	})
}

func TestAccServiceCatalogProvisionedProduct_tags_null(t *testing.T) {
	ctx := acctest.Context(t)
	var v servicecatalog.ProvisionedProductDetail
	resourceName := "aws_servicecatalog_provisioned_product.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ServiceCatalogServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProvisionedProductDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccProvisionedProductConfig_tagsNull(rName, "key1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
				},
			},
			{
				Config:             testAccProvisionedProductConfig_tags0(rName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccServiceCatalogProvisionedProduct_tags_AddOnUpdate(t *testing.T) {
	ctx := acctest.Context(t)
	var v servicecatalog.ProvisionedProductDetail
	resourceName := "aws_servicecatalog_provisioned_product.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ServiceCatalogServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProvisionedProductDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccProvisionedProductConfig_tags0(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				Config: testAccProvisionedProductConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
				},
			},
		},
	})
}

func TestAccServiceCatalogProvisionedProduct_tags_EmptyTag_OnCreate(t *testing.T) {
	t.Skip("Resource ProvisionedProduct does not support empty tags")

	ctx := acctest.Context(t)
	var v servicecatalog.ProvisionedProductDetail
	resourceName := "aws_servicecatalog_provisioned_product.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ServiceCatalogServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProvisionedProductDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccProvisionedProductConfig_tags1(rName, "key1", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
				},
			},
			{
				Config: testAccProvisionedProductConfig_tags0(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
				},
			},
		},
	})
}

func TestAccServiceCatalogProvisionedProduct_tags_EmptyTag_OnUpdate_Add(t *testing.T) {
	t.Skip("Resource ProvisionedProduct does not support empty tags")

	ctx := acctest.Context(t)
	var v servicecatalog.ProvisionedProductDetail
	resourceName := "aws_servicecatalog_provisioned_product.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ServiceCatalogServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProvisionedProductDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccProvisionedProductConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				Config: testAccProvisionedProductConfig_tags2(rName, "key1", "value1", "key2", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
				},
			},
			{
				Config: testAccProvisionedProductConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
				},
			},
		},
	})
}

func TestAccServiceCatalogProvisionedProduct_tags_EmptyTag_OnUpdate_Replace(t *testing.T) {
	t.Skip("Resource ProvisionedProduct does not support empty tags")

	ctx := acctest.Context(t)
	var v servicecatalog.ProvisionedProductDetail
	resourceName := "aws_servicecatalog_provisioned_product.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ServiceCatalogServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProvisionedProductDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccProvisionedProductConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				Config: testAccProvisionedProductConfig_tags1(rName, "key1", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
				},
			},
		},
	})
}

func TestAccServiceCatalogProvisionedProduct_tags_DefaultTags_providerOnly(t *testing.T) {
	ctx := acctest.Context(t)
	var v servicecatalog.ProvisionedProductDetail
	resourceName := "aws_servicecatalog_provisioned_product.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ServiceCatalogServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProvisionedProductDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: acctest.ConfigCompose(
					acctest.ConfigDefaultTags_Tags1("key1", "value1"),
					testAccProvisionedProductConfig_tags0(rName),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
				},
			},
			{
				Config: acctest.ConfigCompose(
					acctest.ConfigDefaultTags_Tags2("key1", "value1updated", "key2", "value2"),
					testAccProvisionedProductConfig_tags0(rName),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.key2", "value2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
				},
			},
			// {
			// 	Config: acctest.ConfigCompose(
			// 		acctest.ConfigDefaultTags_Tags1("key2", "value2"),
			// 		testAccProvisionedProductConfig_tags0(rName),
			// 	),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		testAccCheckProvisionedProductExists(ctx, resourceName, &v),
			// 		resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
			// 		resource.TestCheckResourceAttr(resourceName, "tags_all.%", "1"),
			// 		resource.TestCheckResourceAttr(resourceName, "tags_all.key2", "value2"),
			// 	),
			// },
			// {
			// 	ResourceName:      resourceName,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	ImportStateVerifyIgnore: []string{
			// 		"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
			// 	},
			// },
			// {
			// 	Config: acctest.ConfigCompose(
			// 		acctest.ConfigDefaultTags_Tags0(),
			// 		testAccProvisionedProductConfig_tags0(rName),
			// 	),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		testAccCheckProvisionedProductExists(ctx, resourceName, &v),
			// 		resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
			// 		resource.TestCheckResourceAttr(resourceName, "tags_all.%", "0"),
			// 	),
			// },
			// {
			// 	ResourceName:      resourceName,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	ImportStateVerifyIgnore: []string{
			// 		"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
			// 	},
			// },
		},
	})
}

func TestAccServiceCatalogProvisionedProduct_tags_DefaultTags_nonOverlapping(t *testing.T) {
	ctx := acctest.Context(t)
	var v servicecatalog.ProvisionedProductDetail
	resourceName := "aws_servicecatalog_provisioned_product.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ServiceCatalogServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProvisionedProductDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: acctest.ConfigCompose(
					acctest.ConfigDefaultTags_Tags1("providerkey1", "providervalue1"),
					testAccProvisionedProductConfig_tags1(rName, "resourcekey1", "resourcevalue1"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.resourcekey1", "resourcevalue1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.providerkey1", "providervalue1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.resourcekey1", "resourcevalue1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
				},
			},
			{
				Config: acctest.ConfigCompose(
					acctest.ConfigDefaultTags_Tags1("providerkey1", "providervalue1updated"),
					testAccProvisionedProductConfig_tags2(rName, "resourcekey1", "resourcevalue1updated", "resourcekey2", "resourcevalue2"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.resourcekey1", "resourcevalue1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.resourcekey2", "resourcevalue2"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.providerkey1", "providervalue1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.resourcekey1", "resourcevalue1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.resourcekey2", "resourcevalue2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
				},
			},
			// {
			// 	Config: acctest.ConfigCompose(
			// 		acctest.ConfigDefaultTags_Tags0(),
			// 		testAccProvisionedProductConfig_tags0(rName),
			// 	),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		testAccCheckProvisionedProductExists(ctx, resourceName, &v),
			// 		resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
			// 		resource.TestCheckResourceAttr(resourceName, "tags_all.%", "0"),
			// 	),
			// },
			// {
			// 	ResourceName:      resourceName,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	ImportStateVerifyIgnore: []string{
			// 		"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
			// 	},
			// },
		},
	})
}

func TestAccServiceCatalogProvisionedProduct_tags_DefaultTags_overlapping(t *testing.T) {
	ctx := acctest.Context(t)
	var v servicecatalog.ProvisionedProductDetail
	resourceName := "aws_servicecatalog_provisioned_product.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ServiceCatalogServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProvisionedProductDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: acctest.ConfigCompose(
					acctest.ConfigDefaultTags_Tags1("overlapkey1", "providervalue1"),
					testAccProvisionedProductConfig_tags1(rName, "overlapkey1", "resourcevalue1"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.overlapkey1", "resourcevalue1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.overlapkey1", "resourcevalue1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
				},
			},
			{
				Config: acctest.ConfigCompose(
					acctest.ConfigDefaultTags_Tags2("overlapkey1", "providervalue1", "overlapkey2", "providervalue2"),
					testAccProvisionedProductConfig_tags2(rName, "overlapkey1", "resourcevalue1", "overlapkey2", "resourcevalue2"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.overlapkey1", "resourcevalue1"),
					resource.TestCheckResourceAttr(resourceName, "tags.overlapkey2", "resourcevalue2"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.overlapkey1", "resourcevalue1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.overlapkey2", "resourcevalue2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
				},
			},
			{
				Config: acctest.ConfigCompose(
					acctest.ConfigDefaultTags_Tags1("overlapkey1", "providervalue1"),
					testAccProvisionedProductConfig_tags1(rName, "overlapkey1", "resourcevalue2"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.overlapkey1", "resourcevalue2"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.overlapkey1", "resourcevalue2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
				},
			},
		},
	})
}

func TestAccServiceCatalogProvisionedProduct_tags_DefaultTags_updateToProviderOnly(t *testing.T) {
	ctx := acctest.Context(t)
	var v servicecatalog.ProvisionedProductDetail
	resourceName := "aws_servicecatalog_provisioned_product.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ServiceCatalogServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProvisionedProductDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccProvisionedProductConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.key1", "value1"),
				),
			},
			{
				Config: acctest.ConfigCompose(
					acctest.ConfigDefaultTags_Tags1("key1", "value1"),
					testAccProvisionedProductConfig_tags0(rName),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
				},
			},
		},
	})
}

func TestAccServiceCatalogProvisionedProduct_tags_DefaultTags_updateToResourceOnly(t *testing.T) {
	ctx := acctest.Context(t)
	var v servicecatalog.ProvisionedProductDetail
	resourceName := "aws_servicecatalog_provisioned_product.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ServiceCatalogServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProvisionedProductDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: acctest.ConfigCompose(
					acctest.ConfigDefaultTags_Tags1("key1", "value1"),
					testAccProvisionedProductConfig_tags0(rName),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.key1", "value1"),
				),
			},
			{
				Config: testAccProvisionedProductConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
				},
			},
		},
	})
}

func TestAccServiceCatalogProvisionedProduct_tags_DefaultTags_emptyResourceTag(t *testing.T) {
	t.Skip("Resource ProvisionedProduct does not support empty tags")

	ctx := acctest.Context(t)
	var v servicecatalog.ProvisionedProductDetail
	resourceName := "aws_servicecatalog_provisioned_product.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ServiceCatalogServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProvisionedProductDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: acctest.ConfigCompose(
					acctest.ConfigDefaultTags_Tags1("key1", "value1"),
					testAccProvisionedProductConfig_tags1(rName, "key1", ""),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", ""),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.key1", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
				},
			},
		},
	})
}

func TestAccServiceCatalogProvisionedProduct_tags_DefaultTags_nullOverlappingResourceTag(t *testing.T) {
	ctx := acctest.Context(t)
	var v servicecatalog.ProvisionedProductDetail
	resourceName := "aws_servicecatalog_provisioned_product.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ServiceCatalogServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProvisionedProductDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: acctest.ConfigCompose(
					acctest.ConfigDefaultTags_Tags1("key1", "providervalue1"),
					testAccProvisionedProductConfig_tagsNull(rName, "key1"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.key1", "providervalue1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
				},
			},
		},
	})
}

func TestAccServiceCatalogProvisionedProduct_tags_DefaultTags_nullNonOverlappingResourceTag(t *testing.T) {
	ctx := acctest.Context(t)
	var v servicecatalog.ProvisionedProductDetail
	resourceName := "aws_servicecatalog_provisioned_product.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.ServiceCatalogServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckProvisionedProductDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: acctest.ConfigCompose(
					acctest.ConfigDefaultTags_Tags1("providerkey1", "providervalue1"),
					testAccProvisionedProductConfig_tagsNull(rName, "resourcekey1"),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProvisionedProductExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.providerkey1", "providervalue1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"accept_language", "ignore_errors", "provisioning_artifact_name", "provisioning_parameters", "retain_physical_resources",
				},
			},
		},
	})
}
