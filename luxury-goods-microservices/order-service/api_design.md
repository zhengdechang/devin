# Order Service API Design (gRPC)

```protobuf
syntax = "proto3";

package order;

option go_package = "luxurygoods/orderpb";

import "google/protobuf/timestamp.proto";

// Interface for ProductService client (Order Service calls Product Service)
// This is a simplified view. The actual ProductService definition
// would be imported from Product Service's .proto files.
service ProductService {
  // Locks a product when an order is placed
  rpc LockProduct(LockProductRequest) returns (LockProductResponse);
  // Unlocks a product if an order is cancelled
  rpc UnlockProduct(UnlockProductRequest) returns (UnlockProductResponse);
}

message LockProductRequest {
  string product_id = 1;
  string order_id = 2; // To associate the lock with this order
}

message LockProductResponse {
  bool success = 1;
  // Could include product details if needed by Order Service
}

message UnlockProductRequest {
  string product_id = 1;
  string order_id = 2; // To specify which order caused the unlock
}

message UnlockProductResponse {
  bool success = 1;
}


// OrderService defines functionalities for managing orders.
service OrderService {
  // Places an order for a product. This involves trying to lock the product first.
  rpc PlaceOrder(PlaceOrderRequest) returns (PlaceOrderResponse);
  // Cancels an order. This involves trying to unlock the product.
  rpc CancelOrder(CancelOrderRequest) returns (CancelOrderResponse);
  // Retrieves the status and details of an order.
  rpc GetOrderStatus(GetOrderStatusRequest) returns (GetOrderStatusResponse);
  // Lists orders, with potential filters.
  rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse);
  // Marks an order as paid.
  rpc MarkOrderAsPaid(MarkOrderAsPaidRequest) returns (MarkOrderAsPaidResponse);
  // Updates order status (e.g., to shipped, delivered).
  rpc UpdateOrderStatus(UpdateOrderStatusRequest) returns (UpdateOrderStatusResponse);
}

message Order {
  string order_id = 1;
  string product_id = 2;
  string buyer_id = 3;
  string platform_id = 4; // Platform where the order originated (e.g., "xianyu", "internal_app")
  string status = 5; // ENUM: "pending_payment", "pending_lock", "confirmed", "paid", "shipped", "delivered", "cancelled", "failed"
  google.protobuf.Timestamp created_at = 6;
  google.protobuf.Timestamp paid_at = 7;    // Nullable
  double total_amount = 8;
  string currency = 9; // e.g., "USD", "JPY"
  // Potentially add shipping_address_id, billing_address_id if addresses are managed entities
  string shipping_address_snapshot = 10; // Snapshot of shipping address at time of order
  string billing_address_snapshot = 11; // Snapshot of billing address, nullable
  google.protobuf.Timestamp updated_at = 12;
}

message PlaceOrderRequest {
  string product_id = 1;
  string buyer_id = 2;
  string platform_id = 3;
  double expected_price = 4; // Price at which buyer intends to purchase
  string currency = 5;
  string shipping_address = 6; // Can be a structured message too
  string billing_address = 7;  // Can be a structured message too, nullable
  // idempotency_key might be useful here
  string idempotency_key = 8;
}

message PlaceOrderResponse {
  Order order = 1;
  // bool lock_acquired = 2; // This is implicit if order status is 'confirmed' or similar
  // If order creation involves synchronous lock attempt, status in Order message will reflect outcome.
  // e.g. status = "pending_lock" if async, or "confirmed"/"failed_lock" if sync.
}

message CancelOrderRequest {
  string order_id = 1;
  string user_id = 2; // For authorization/audit (could be buyer_id or an admin_id)
  string reason = 3; // Reason for cancellation
}

message CancelOrderResponse {
  Order order = 1; // Updated order with "cancelled" status
}

message GetOrderStatusRequest {
  string order_id = 1;
}

message GetOrderStatusResponse {
  Order order = 1;
}

message ListOrdersRequest {
  string buyer_id = 1;      // Optional filter
  string platform_id = 2;   // Optional filter
  string status = 3;        // Optional filter
  int32 page_size = 4;
  string page_token = 5;
  // google.protobuf.Timestamp created_after = 6;
  // google.protobuf.Timestamp created_before = 7;
}

message ListOrdersResponse {
  repeated Order orders = 1;
  string next_page_token = 2;
}

message MarkOrderAsPaidRequest {
  string order_id = 1;
  string payment_transaction_id = 2; // From payment gateway
  google.protobuf.Timestamp paid_at_time = 3; // Actual time of payment
}

message MarkOrderAsPaidResponse {
  Order order = 1; // Updated order with "paid" status
}

message UpdateOrderStatusRequest {
  string order_id = 1;
  string new_status = 2; // e.g., "shipped", "delivered"
  // string tracking_number = 3; // if status is "shipped"
  // string operator_id = 4; // For audit
}

message UpdateOrderStatusResponse {
  Order order = 1; // Updated order
}
```
