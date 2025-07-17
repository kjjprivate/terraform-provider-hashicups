package provider


import (
	"context"
	"fmt"
  
	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
  )

// Ensure the implementation satisfies the expexted interfaces
var (
	_ datasource.DataSource = &coffeeDatasources{}
	_ datasource.DataSourceWithConfigure = &coffeesDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation
func NewCoffeesDataSource() datasource.DataSource {
	return &coffeesDataSource{}
}

// cooffeesDataSource is the data source implementation
type coffeesDataSource struct {
	client *hashicups.Client
}


//coffeeDataSourceModel maps the data source schema data
type coffeeDataSourceModel struct {
	Coffees []coffeesModel `tfsdk:"coffeees"`
}

// coffeesModel maps coffees schema data
type coffeesModel struct {
	ID 			types.Int64 	`tfsdk:"id"`
	Name 		types.String	`tfsdk:"name"`
	Teaser		types.String	`tfsdk:"teaser"`
	Description types.String 	`tfsdk:"description"`
	Price 		types.String	`tfsdk:"price"`
	Image		types.String 	`tfsdk:"image"`
	Ingredients types.String    `tfsdk:"ingredients"`
}

// coffeesIngredientModel maps coffee ingredients data
type coffeesIngredientModel struct {
	ID types.Int64 `tfsdk:"id"`
}


// Metadata returns the data source type name
func (d *coffeesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_coffees"
}

// Schmea defines the schema for the data source
func (d *coffeesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"coffees": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						}
						"teaser": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"price": schema.Float64Attribute{
							Computed: true,
						},
						"image": schema.StringAttribute{
							Computed: true,
						},
						"ingredients": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attribute: map[string]schema.Attribute{
									"id": schema.Int64Attribute{
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}



// Read refreshes the Terraform state with the latest data
func (d *coffeesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {}

//Configure adds the provider configured client to the data source.
func (d *coffeesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*hashicups.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.SPrintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}