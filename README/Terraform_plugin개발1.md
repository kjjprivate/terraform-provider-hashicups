
### Terraform provider 개발 1

Terraform core - cli 환경에서 Terraform 명령어를 수행하는 역할을 하는 바이너리 파일

- Resource state 관리

- IaC
- Resource Graph 생성
- plan 실행

- RPC 프로토콜로 Plugin과 통신

Terraform plugin - Terraform core에 의해 실행되어진다. 각각의 플러그인들은 특별한 서비스(AWS, Bash, Provisioner, 등) 에 대한 작업을 수행하는 역할을 한다.

- API 호출에 사용되는 라이브러리 초기화
- 인프라 provider와의 인증
- 특정 서비스에 매핑되어 관리되는 리소스 및 데이터 소스를 정의
- 실무자 작업을 명시적으로 활성화하고 단순화하는 기능



Terraform plugin framework

- ./internal/provider.go - provider argument 및 구성 방식 설정

  - Schema 메소드 - Provider의 argument에 대한 스키마를 구성

    - ``` go
      // Schema defines the provider-level schema for configuration data.
      func (p *hashicupsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
          resp.Schema = schema.Schema{
              Attributes: map[string]schema.Attribute{
                  "host": schema.StringAttribute{
                      Optional: true,
                  },
                  "username": schema.StringAttribute{
                      Optional: true,
                  },
                  "password": schema.StringAttribute{
                      Optional:  true,
                      Sensitive: true,
                  },
              },
          }
      }
      ```

  - tfsdk 구조체를 이용하여 스키마의 정의를 매핑해주어야한다.

    - ``` go
      // hashicupsProviderModel maps provider schema data to a Go type.
      type hashicupsProviderModel struct {
          Host     types.String `tfsdk:"host"`
          Username types.String `tfsdk:"username"`
          Password types.String `tfsdk:"password"`
      }



#### DataSource 개발

1. internal/provider.go -> DataSources 메소드에 새로운 Data source 메소드 추가

   ```go
   // DataSources defines the data sources implemented in the provider.
   func (p *hashicupsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
     return []func() datasource.DataSource {
       NewCoffeesDataSource,
     }
   }
   ```

2. internal 경로에 go 파일을 하나 추가 및 data source의 코드 구성

   ``` go
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
   	_ datasource.DataSource = &coffeesDataSource{}
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
   type coffeesDataSourceModel struct {
   	Coffees []coffeesModel `tfsdk:"coffees"`
   }
   
   // coffeesModel maps coffees schema data
   type coffeesModel struct {
   	ID 			types.Int64 	`tfsdk:"id"`
   	Name 		types.String	`tfsdk:"name"`
   	Teaser		types.String	`tfsdk:"teaser"`
   	Description types.String 	`tfsdk:"description"`
   	Price 		types.Float64	`tfsdk:"price"`
   	Image		types.String 	`tfsdk:"image"`
   	Ingredients []coffeesIngredientsModel    `tfsdk:"ingredients"`
   }
   
   // coffeesIngredientModel maps coffee ingredients data
   type coffeesIngredientsModel struct {
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
   						},
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
   								Attributes: map[string]schema.Attribute{
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
   func (d *coffeesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
   
   	var state coffeesDataSourceModel
   	coffees, err := d.client.GetCoffees()
   	if err != nil {
   		resp.Diagnostics.AddError(
   			"Unable to Read Hashicups Coffees",
   			err.Error(),
   		)
   		return
   	}
   	for _, coffee := range coffees {
   		coffeeState := coffeesModel{
   			ID: types.Int64Value(int64(coffee.ID)),
   			Name: types.StringValue(coffee.Name),
   			Teaser: types.StringValue(coffee.Teaser),
   			Description: types.StringValue(coffee.Description),
   			Price:	types.Float64Value(coffee.Price),
   			Image: types.StringValue(coffee.Image),
   		}
   		for _, ingredient := range coffee.Ingredient {
   			coffeeState.Ingredients = append(coffeeState.Ingredients, coffeesIngredientsModel{
   				ID: types.Int64Value(int64(ingredient.ID)),
   			})
   		}
   
   		state.Coffees = append(state.Coffees, coffeeState)
   	}
   
   	diags := resp.State.Set(ctx, &state)
   	resp.Diagnostics.Append(diags...)
   	if resp.Diagnostics.HasError() {
   		return
   	}
   }
   
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
   			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
   		)
   		return
   	}
   	d.client = client
   }
   ```



3. 