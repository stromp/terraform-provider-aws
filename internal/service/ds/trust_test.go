package ds_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	awstypes "github.com/aws/aws-sdk-go-v2/service/directoryservice/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfds "github.com/hashicorp/terraform-provider-aws/internal/service/ds"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccDSTrust_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var v awstypes.Trust
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_directory_service_trust.test"
	domainName := acctest.RandomDomainName()
	domainNameOther := acctest.RandomDomainName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckDirectoryService(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.DSEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTrustDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTrustConfig_basic(rName, domainName, domainNameOther),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTrustExists(ctx, resourceName, &v),
					resource.TestMatchResourceAttr(resourceName, "id", regexp.MustCompile(`^t-\w{10}`)),
					resource.TestCheckResourceAttr(resourceName, "conditional_forwarder_ip_addrs.#", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "directory_id", "aws_directory_service_directory.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "remote_domain_name", domainNameOther),
					resource.TestCheckResourceAttr(resourceName, "selective_auth", string(awstypes.SelectiveAuthDisabled)),
					resource.TestCheckResourceAttr(resourceName, "trust_direction", string(awstypes.TrustDirectionTwoWay)),
					resource.TestCheckResourceAttr(resourceName, "trust_password", "Some0therPassword"),
					resource.TestCheckResourceAttr(resourceName, "trust_type", string(awstypes.TrustTypeForest)),
					acctest.CheckResourceAttrRFC3339(resourceName, "created_date_time"),
					acctest.CheckResourceAttrRFC3339(resourceName, "last_updated_date_time"),
					resource.TestCheckResourceAttr(resourceName, "trust_state", string(awstypes.TrustStateVerifyFailed)),
					resource.TestCheckResourceAttrSet(resourceName, "trust_state_reason"),
					acctest.CheckResourceAttrRFC3339(resourceName, "state_last_updated_date_time"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTrustStateIdFunc(resourceName),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"trust_password",
				},
			},
		},
	})
}

func TestAccDSTrust_bidirectionalBasic(t *testing.T) {
	ctx := acctest.Context(t)
	var v1, v2 awstypes.Trust
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_directory_service_trust.test"
	resourceOtherName := "aws_directory_service_trust.other"
	domainName := acctest.RandomDomainName()
	domainNameOther := acctest.RandomDomainName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckDirectoryService(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.DSEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTrustDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTrustConfig_bidirectionalBasic(rName, domainName, domainNameOther),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTrustExists(ctx, resourceName, &v1),
					resource.TestCheckResourceAttr(resourceName, "trust_state", string(awstypes.TrustStateVerified)),
					resource.TestCheckNoResourceAttr(resourceName, "trust_state_reason"),

					testAccCheckTrustExists(ctx, resourceOtherName, &v2),
					resource.TestCheckResourceAttr(resourceOtherName, "trust_state", string(awstypes.TrustStateVerified)),
					resource.TestCheckNoResourceAttr(resourceOtherName, "trust_state_reason"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTrustStateIdFunc(resourceName),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"trust_password",
				},
			},
		},
	})
}

func TestAccDSTrust_SelectiveAuth(t *testing.T) {
	ctx := acctest.Context(t)
	var v awstypes.Trust
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_directory_service_trust.test"
	domainName := acctest.RandomDomainName()
	domainNameOther := acctest.RandomDomainName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckDirectoryService(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.DSEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTrustDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTrustConfig_SelectiveAuth(rName, domainName, domainNameOther, awstypes.SelectiveAuthEnabled),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTrustExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "selective_auth", string(awstypes.SelectiveAuthEnabled)),
					resource.TestCheckResourceAttr(resourceName, "trust_state", string(awstypes.TrustStateVerifyFailed)), // Updating single-sided config
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTrustStateIdFunc(resourceName),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"trust_password",
				},
			},
			{
				Config: testAccTrustConfig_SelectiveAuth(rName, domainName, domainNameOther, awstypes.SelectiveAuthDisabled),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTrustExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "selective_auth", string(awstypes.SelectiveAuthDisabled)),
					resource.TestCheckResourceAttr(resourceName, "trust_state", string(awstypes.TrustStateVerifyFailed)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTrustStateIdFunc(resourceName),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"trust_password",
				},
			},
		},
	})
}

func TestAccDSTrust_bidirectionalSelectiveAuth(t *testing.T) {
	ctx := acctest.Context(t)
	var v awstypes.Trust
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_directory_service_trust.test"
	domainName := acctest.RandomDomainName()
	domainNameOther := acctest.RandomDomainName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckDirectoryService(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.DSEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTrustDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTrustConfig_bidirectionalSelectiveAuth(rName, domainName, domainNameOther, awstypes.SelectiveAuthEnabled),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTrustExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "selective_auth", string(awstypes.SelectiveAuthEnabled)),
					resource.TestCheckResourceAttr(resourceName, "trust_state", string(awstypes.TrustStateVerified)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTrustStateIdFunc(resourceName),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"trust_password",
				},
			},
			{
				Config: testAccTrustConfig_bidirectionalSelectiveAuth(rName, domainName, domainNameOther, awstypes.SelectiveAuthDisabled),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTrustExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "selective_auth", string(awstypes.SelectiveAuthDisabled)),
					resource.TestCheckResourceAttr(resourceName, "trust_state", string(awstypes.TrustStateVerified)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTrustStateIdFunc(resourceName),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"trust_password",
				},
			},
		},
	})
}

func TestAccDSTrust_TrustType(t *testing.T) {
	ctx := acctest.Context(t)
	var v awstypes.Trust
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_directory_service_trust.test"
	domainName := acctest.RandomDomainName()
	domainNameOther := acctest.RandomDomainName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckDirectoryService(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.DSEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTrustDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTrustConfig_TrustType(rName, domainName, domainNameOther, awstypes.TrustTypeExternal),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTrustExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "trust_type", string(awstypes.TrustTypeExternal)),
					resource.TestCheckResourceAttr(resourceName, "trust_state", string(awstypes.TrustStateVerifyFailed)), // Updating single-sided config
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTrustStateIdFunc(resourceName),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"trust_password",
				},
			},
		},
	})
}

func TestAccDSTrust_TrustTypeSpecifyDefault(t *testing.T) {
	ctx := acctest.Context(t)
	var v awstypes.Trust
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_directory_service_trust.test"
	domainName := acctest.RandomDomainName()
	domainNameOther := acctest.RandomDomainName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckDirectoryService(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.DSEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTrustDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTrustConfig_basic(rName, domainName, domainNameOther),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTrustExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "trust_type", string(awstypes.TrustTypeForest)),
				),
			},
			{
				Config:   testAccTrustConfig_TrustType(rName, domainName, domainNameOther, awstypes.TrustTypeForest),
				PlanOnly: true,
			},
		},
	})
}

func TestAccDSTrust_ConditionalForwarderIPs(t *testing.T) {
	ctx := acctest.Context(t)
	var v awstypes.Trust
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_directory_service_trust.test"
	domainName := acctest.RandomDomainName()
	domainNameOther := acctest.RandomDomainName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckDirectoryService(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.DSEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTrustDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTrustConfig_basic(rName, domainName, domainNameOther),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTrustExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "conditional_forwarder_ip_addrs.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTrustStateIdFunc(resourceName),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"trust_password",
				},
			},
			{
				Config: testAccTrustConfig_ConditionalForwarderIPs(rName, domainName, domainNameOther),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTrustExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "conditional_forwarder_ip_addrs.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTrustStateIdFunc(resourceName),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"trust_password",
				},
			},
		},
	})
}

// TODO: Test one-directional trusts

func testAccCheckTrustExists(ctx context.Context, n string, v *awstypes.Trust) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Directory Service Trust ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).DSClient()

		output, err := tfds.FindTrustByID(ctx, conn, rs.Primary.Attributes["directory_id"], rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckTrustDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).DSClient()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_directory_service_trust" {
				continue
			}

			_, err := tfds.FindTrustByID(ctx, conn, rs.Primary.Attributes["directory_id"], rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("Directory Service Trust %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccTrustStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["directory_id"], rs.Primary.ID), nil
	}
}

func testAccTrustConfig_basic(rName, domain, domainOther string) string {
	return acctest.ConfigCompose(
		acctest.ConfigVPCWithSubnets(rName, 2),
		fmt.Sprintf(`
resource "aws_directory_service_trust" "test" {
  directory_id = aws_directory_service_directory.test.id

  remote_domain_name = aws_directory_service_directory.other.name
  trust_direction    = "Two-Way"
  trust_password     = "Some0therPassword"

  conditional_forwarder_ip_addrs = aws_directory_service_directory.other.dns_ip_addresses
}

resource "aws_directory_service_directory" "test" {
  name     = %[1]q
  password = "SuperSecretPassw0rd"
  type     = "MicrosoftAD"
  edition  = "Standard"

  vpc_settings {
    vpc_id     = aws_vpc.test.id
    subnet_ids = aws_subnet.test[*].id
  }
}

resource "aws_directory_service_directory" "other" {
  name     = %[2]q
  password = "SuperSecretPassw0rd"
  type     = "MicrosoftAD"
  edition  = "Standard"

  vpc_settings {
    vpc_id     = aws_vpc.test.id
    subnet_ids = aws_subnet.test[*].id
  }
}

resource "aws_security_group_rule" "test" {
  security_group_id = aws_directory_service_directory.test.security_group_id

  type                     = "egress"
  protocol                 = "all"
  from_port                = 0
  to_port                  = 65535
  source_security_group_id = aws_directory_service_directory.other.security_group_id
}

resource "aws_security_group_rule" "other" {
  security_group_id = aws_directory_service_directory.other.security_group_id

  type                     = "egress"
  protocol                 = "all"
  from_port                = 0
  to_port                  = 65535
  source_security_group_id = aws_directory_service_directory.test.security_group_id
}
`, domain, domainOther),
	)
}

func testAccTrustConfig_bidirectionalBasic(rName, domain, domainOther string) string {
	return acctest.ConfigCompose(
		acctest.ConfigVPCWithSubnets(rName, 2),
		fmt.Sprintf(`
resource "aws_directory_service_trust" "test" {
  directory_id = aws_directory_service_directory.test.id

  remote_domain_name = aws_directory_service_directory.other.name
  trust_direction    = "Two-Way"
  trust_password     = "Some0therPassword"

  conditional_forwarder_ip_addrs = aws_directory_service_directory.other.dns_ip_addresses
}

resource "aws_directory_service_trust" "other" {
  directory_id = aws_directory_service_directory.other.id

  remote_domain_name = aws_directory_service_directory.test.name
  trust_direction    = "Two-Way"
  trust_password     = "Some0therPassword"

  conditional_forwarder_ip_addrs = aws_directory_service_directory.test.dns_ip_addresses
}

resource "aws_directory_service_directory" "test" {
  name     = %[1]q
  password = "SuperSecretPassw0rd"
  type     = "MicrosoftAD"
  edition  = "Standard"

  vpc_settings {
    vpc_id     = aws_vpc.test.id
    subnet_ids = aws_subnet.test[*].id
  }
}

resource "aws_directory_service_directory" "other" {
  name     = %[2]q
  password = "SuperSecretPassw0rd"
  type     = "MicrosoftAD"
  edition  = "Standard"

  vpc_settings {
    vpc_id     = aws_vpc.test.id
    subnet_ids = aws_subnet.test[*].id
  }
}

resource "aws_security_group_rule" "test" {
  security_group_id = aws_directory_service_directory.test.security_group_id

  type                     = "egress"
  protocol                 = "all"
  from_port                = 0
  to_port                  = 65535
  source_security_group_id = aws_directory_service_directory.other.security_group_id
}

resource "aws_security_group_rule" "other" {
  security_group_id = aws_directory_service_directory.other.security_group_id

  type                     = "egress"
  protocol                 = "all"
  from_port                = 0
  to_port                  = 65535
  source_security_group_id = aws_directory_service_directory.test.security_group_id
}
`, domain, domainOther),
	)
}

func testAccTrustConfig_SelectiveAuth(rName, domain, domainOther string, selectiveAuth awstypes.SelectiveAuth) string {
	return acctest.ConfigCompose(
		acctest.ConfigVPCWithSubnets(rName, 2),
		fmt.Sprintf(`
resource "aws_directory_service_trust" "test" {
  directory_id = aws_directory_service_directory.test.id

  remote_domain_name = aws_directory_service_directory.other.name
  trust_direction    = "Two-Way"
  trust_password     = "Some0therPassword"

  conditional_forwarder_ip_addrs = aws_directory_service_directory.other.dns_ip_addresses

  selective_auth = %[3]q
}

resource "aws_directory_service_directory" "test" {
  name     = %[1]q
  password = "SuperSecretPassw0rd"
  type     = "MicrosoftAD"
  edition  = "Standard"

  vpc_settings {
    vpc_id     = aws_vpc.test.id
    subnet_ids = aws_subnet.test[*].id
  }
}

resource "aws_directory_service_directory" "other" {
  name     = %[2]q
  password = "SuperSecretPassw0rd"
  type     = "MicrosoftAD"
  edition  = "Standard"

  vpc_settings {
    vpc_id     = aws_vpc.test.id
    subnet_ids = aws_subnet.test[*].id
  }
}

resource "aws_security_group_rule" "test" {
  security_group_id = aws_directory_service_directory.test.security_group_id

  type                     = "egress"
  protocol                 = "all"
  from_port                = 0
  to_port                  = 65535
  source_security_group_id = aws_directory_service_directory.other.security_group_id
}

resource "aws_security_group_rule" "other" {
  security_group_id = aws_directory_service_directory.other.security_group_id

  type                     = "egress"
  protocol                 = "all"
  from_port                = 0
  to_port                  = 65535
  source_security_group_id = aws_directory_service_directory.test.security_group_id
}
`, domain, domainOther, selectiveAuth),
	)
}

func testAccTrustConfig_bidirectionalSelectiveAuth(rName, domain, domainOther string, selectiveAuth awstypes.SelectiveAuth) string {
	return acctest.ConfigCompose(
		acctest.ConfigVPCWithSubnets(rName, 2),
		fmt.Sprintf(`
resource "aws_directory_service_trust" "test" {
	directory_id = aws_directory_service_directory.test.id

	remote_domain_name = aws_directory_service_directory.other.name
	trust_direction    = "Two-Way"
	trust_password     = "Some0therPassword"

	conditional_forwarder_ip_addrs = aws_directory_service_directory.other.dns_ip_addresses

	selective_auth = %[3]q
}

resource "aws_directory_service_trust" "other" {
	directory_id = aws_directory_service_directory.other.id

	remote_domain_name = aws_directory_service_directory.test.name
	trust_direction    = "Two-Way"
	trust_password     = "Some0therPassword"

	conditional_forwarder_ip_addrs = aws_directory_service_directory.test.dns_ip_addresses
}

resource "aws_directory_service_directory" "test" {
  name     = %[1]q
  password = "SuperSecretPassw0rd"
  type     = "MicrosoftAD"
  edition  = "Standard"

  vpc_settings {
    vpc_id     = aws_vpc.test.id
    subnet_ids = aws_subnet.test[*].id
  }
}

resource "aws_directory_service_directory" "other" {
  name     = %[2]q
  password = "SuperSecretPassw0rd"
  type     = "MicrosoftAD"
  edition  = "Standard"

  vpc_settings {
    vpc_id     = aws_vpc.test.id
    subnet_ids = aws_subnet.test[*].id
  }
}

resource "aws_security_group_rule" "test" {
	security_group_id = aws_directory_service_directory.test.security_group_id

	type = "egress"
	protocol = "all"
	from_port = 0
	to_port = 65535
	source_security_group_id =aws_directory_service_directory.other.security_group_id
}

resource "aws_security_group_rule" "other" {
	security_group_id = aws_directory_service_directory.other.security_group_id

	type = "egress"
	protocol = "all"
	from_port = 0
	to_port = 65535
	source_security_group_id =aws_directory_service_directory.test.security_group_id
}
`, domain, domainOther, selectiveAuth),
	)
}

func testAccTrustConfig_TrustType(rName, domain, domainOther string, trustType awstypes.TrustType) string {
	return acctest.ConfigCompose(
		acctest.ConfigVPCWithSubnets(rName, 2),
		fmt.Sprintf(`
resource "aws_directory_service_trust" "test" {
	directory_id = aws_directory_service_directory.test.id

	remote_domain_name = aws_directory_service_directory.other.name
	trust_direction    = "Two-Way"
	trust_password     = "Some0therPassword"

	conditional_forwarder_ip_addrs = aws_directory_service_directory.other.dns_ip_addresses

	trust_type = %[3]q
}

resource "aws_directory_service_directory" "test" {
  name     = %[1]q
  password = "SuperSecretPassw0rd"
  type     = "MicrosoftAD"
  edition  = "Standard"

  vpc_settings {
    vpc_id     = aws_vpc.test.id
    subnet_ids = aws_subnet.test[*].id
  }
}

resource "aws_directory_service_directory" "other" {
  name     = %[2]q
  password = "SuperSecretPassw0rd"
  type     = "MicrosoftAD"
  edition  = "Standard"

  vpc_settings {
    vpc_id     = aws_vpc.test.id
    subnet_ids = aws_subnet.test[*].id
  }
}

resource "aws_security_group_rule" "test" {
	security_group_id = aws_directory_service_directory.test.security_group_id

	type = "egress"
	protocol = "all"
	from_port = 0
	to_port = 65535
	source_security_group_id =aws_directory_service_directory.other.security_group_id
}

resource "aws_security_group_rule" "other" {
	security_group_id = aws_directory_service_directory.other.security_group_id

	type = "egress"
	protocol = "all"
	from_port = 0
	to_port = 65535
	source_security_group_id =aws_directory_service_directory.test.security_group_id
}
`, domain, domainOther, trustType),
	)
}

func testAccTrustConfig_ConditionalForwarderIPs(rName, domain, domainOther string) string {
	return acctest.ConfigCompose(
		acctest.ConfigVPCWithSubnets(rName, 2),
		fmt.Sprintf(`
resource "aws_directory_service_trust" "test" {
	directory_id = aws_directory_service_directory.test.id

	remote_domain_name = aws_directory_service_directory.other.name
	trust_direction    = "Two-Way"
	trust_password     = "Some0therPassword"

	conditional_forwarder_ip_addrs = toset(slice(tolist(aws_directory_service_directory.other.dns_ip_addresses),0,1))
}

resource "aws_directory_service_directory" "test" {
  name     = %[1]q
  password = "SuperSecretPassw0rd"
  type     = "MicrosoftAD"
  edition  = "Standard"

  vpc_settings {
    vpc_id     = aws_vpc.test.id
    subnet_ids = aws_subnet.test[*].id
  }
}

resource "aws_directory_service_directory" "other" {
  name     = %[2]q
  password = "SuperSecretPassw0rd"
  type     = "MicrosoftAD"
  edition  = "Standard"

  vpc_settings {
    vpc_id     = aws_vpc.test.id
    subnet_ids = aws_subnet.test[*].id
  }
}

resource "aws_security_group_rule" "test" {
	security_group_id = aws_directory_service_directory.test.security_group_id

	type = "egress"
	protocol = "all"
	from_port = 0
	to_port = 65535
	source_security_group_id =aws_directory_service_directory.other.security_group_id
}

resource "aws_security_group_rule" "other" {
	security_group_id = aws_directory_service_directory.other.security_group_id

	type = "egress"
	protocol = "all"
	from_port = 0
	to_port = 65535
	source_security_group_id =aws_directory_service_directory.test.security_group_id
}
`, domain, domainOther),
	)
}
