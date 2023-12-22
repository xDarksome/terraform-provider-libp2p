// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccExampleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.libp2p_peer_id.this", "base58", "12D3KooWBnTyEyBVeYpZJobw78rb85nNamrYQR3Tc6gJmfQ76pG4"),
				),
			},
		},
	})
}

const testAccExampleDataSourceConfig = `
data "libp2p_peer_id" "this" {
  ed25519_secret_key = base64encode("00000000000000000000000000000001")
}
`
