# orderResource의 Create 함수 설명

이 문서는 `terraform-provider-hashicups`의 `orderResource`에서 사용되는 `Create` 함수의 역할과 동작을 단계별로 설명합니다.

---

## 함수 전체 역할

이 함수는 **Terraform Provider**에서 "order" 리소스를 생성할 때 호출됩니다. 
즉, 사용자가 Terraform에서 주문을 만들면, 실제로 주문을 생성하고 그 결과를 Terraform에 알려주는 함수입니다.

---

## 코드 단계별 설명

### 1. 입력값(plan) 읽기

```go
var plan orderResourceModel
diags := req.Plan.Get(ctx, &plan)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
    return
}
```
- **plan**: 사용자가 Terraform에서 입력한 주문 정보(커피 종류, 수량 등)를 담는 구조체
- **req.Plan.Get**: Terraform이 전달한 입력값을 plan에 복사 (외부 라이브러리: terraform-plugin-framework)

---

### 2. 주문 데이터 생성

```go
var items []hashicups.OrderItem
for _, item := range plan.Items {
    items = append(items, hashicups.OrderItem{
        Coffee: hashicups.Coffee{
            ID: int(item.Coffee.ID.ValueInt64()),
        },
        Quantity: int(item.Quantity.ValueInt64()),
    })
}
```
- **items**: 실제 주문 API에 전달할 주문 항목 배열
- **hashicups.OrderItem, hashicups.Coffee**: HashiCups API와 통신할 때 사용하는 구조체 (외부 라이브러리: hashicups-client-go)

---

### 3. 주문 생성 API 호출

```go
order, err := r.client.CreateOrder(items)
if err != nil {
    resp.Diagnostics.AddError(
        "Error creating order",
        "Could not create order, unexpected error: "+err.Error(),
    )
    return
}
```
- **r.client.CreateOrder**: HashiCups 서버에 주문 생성 요청 (외부 라이브러리: hashicups-client-go)
- **order**: 생성된 주문 정보

---

### 4. 결과(plan) 갱신

```go
plan.ID = types.StringValue(strconv.Itoa(order.ID))
for orderItemIndex, orderItem := range order.Items {
    plan.Items[orderItemIndex] = orderItemModel{
        Coffee: orderItemCoffeeModel{
            ID:          types.Int64Value(int64(orderItem.Coffee.ID)),
            Name:        types.StringValue(orderItem.Coffee.Name),
            Teaser:      types.StringValue(orderItem.Coffee.Teaser),
            Description: types.StringValue(orderItem.Coffee.Description),
            Price:       types.Float64Value(orderItem.Coffee.Price),
            Image:       types.StringValue(orderItem.Coffee.Image),
        },
        Quantity: types.Int64Value(int64(orderItem.Quantity)),
    }
}
plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
```
- **plan.ID**: 생성된 주문의 ID 저장
- **plan.Items**: 실제 주문된 항목 정보 저장
- **types.StringValue 등**: Go 타입을 Terraform 타입으로 변환 (외부 라이브러리: terraform-plugin-framework)

---

### 5. 상태(State) 저장

```go
diags = resp.State.Set(ctx, plan)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
    return
}
```
- **resp.State.Set**: 최종 상태를 Terraform에 전달 (외부 라이브러리: terraform-plugin-framework)

---

## 외부 라이브러리 정리

- **terraform-plugin-framework**  
  - Plan/State 관리, 타입 변환, 진단 메시지 처리 등
- **hashicups-client-go**  
  - HashiCups API와 통신, 주문 생성, 데이터 구조체 제공

---

## 함수에서 사용하는 주요 메소드와 호출 위치

- **req.Plan.Get**: 입력값(plan) 읽기 (terraform-plugin-framework)
- **r.client.CreateOrder**: 주문 생성 API 호출 (hashicups-client-go)
- **resp.State.Set**: 상태 저장 (terraform-plugin-framework)
- **resp.Diagnostics.AddError/Append/HasError**: 에러 및 진단 메시지 관리 (terraform-plugin-framework)