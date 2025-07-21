package provider

import {
	"context",
	"github.com/hashicorp/terraform-plugin-framework/resource",
	"github.com/hashicorp/terraform-plugin-framework/resource/schema",
	"github.com/hashicorp-demoapp/hashicups-client-go"
}

var (
	_ resource.Resource = &orderResource{}
	_ resource.ResourceWithConfigure = &orderResource{}
)
func NewOrderResource() resource.Resource {
	return &orderResource{}
}

type orderResource struct{
	client *hashicups.Client
}

func (r *orderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_order"
}

func (r *orderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

// Create resource
func (r *orderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

}

// Read resource
func (r *orderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

// Update resource
func (r *orderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (r *orderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}
// Configure adds the provider configured client to the resource.
func (r *orderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    // Add a nil check when handling ProviderData because Terraform
    // sets that data after it calls the ConfigureProvider RPC.
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*hashicups.Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Data Source Configure Type",
            fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    r.client = client
}
