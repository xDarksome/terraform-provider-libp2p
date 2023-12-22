// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &PeerIdDataSource{}

func NewPeerIdDataSource() datasource.DataSource {
	return &PeerIdDataSource{}
}

// PeerIdDataSource defines the data source implementation.
type PeerIdDataSource struct {
	client *http.Client
}

// PeerIdDataSourceModel describes the data source data model.
type PeerIdDataSourceModel struct {
	Ed25519SecretKey types.String `tfsdk:"ed25519_secret_key"`
	Base58           types.String `tfsdk:"base58"`
}

func (d *PeerIdDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_peer_id"
}

func (d *PeerIdDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Peer ID",

		Attributes: map[string]schema.Attribute{
			"ed25519_secret_key": schema.StringAttribute{
				MarkdownDescription: "Base64 encoded ed25519 secret key (seed)",
				Sensitive:           true,
				Required:            true,
			},
			"base58": schema.StringAttribute{
				MarkdownDescription: "base58 representation of this Peer ID",
				Optional:            true,
			},
		},
	}
}

func (d *PeerIdDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *PeerIdDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PeerIdDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	bytes, err := base64.StdEncoding.DecodeString(data.Ed25519SecretKey.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Base64 decode", "Failed to base64-decode the provided secret key")
		return
	}

	if len(bytes) != 32 {
		resp.Diagnostics.AddError("Invalid secrey key length", "Should be 32")
		return
	}

	priv_key := ed25519.NewKeyFromSeed(bytes)
	_, pub_key, _ := crypto.KeyPairFromStdKey(&priv_key)
	peer_id, _ := peer.IDFromPublicKey(pub_key)

	data.Base58 = types.StringValue(peer_id.String())

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read PeerIdDataSource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
