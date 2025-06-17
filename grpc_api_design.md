# gRPC API Design Document

This document outlines the gRPC service and message definitions for the二手奢侈品商品微服务系统 (Second-hand Luxury Goods Microservice System).

**Proto Version:** `proto3`

**Common Data Types / Enums (to be defined in a common.proto file):**

```protobuf
syntax = "proto3";

package common;

option go_package = "github.com/yourorg/project/gen/go/common;commonpb";

message Empty {}

message Timestamp {
  int64 seconds = 1;
  int32 nanos = 2;
}

enum ProductStatus {
  PRODUCT_STATUS_UNSPECIFIED = 0;
  PRODUCT_STATUS_DRAFT = 1;
  PRODUCT_STATUS_PENDING_APPROVAL = 2;
  PRODUCT_STATUS_LISTED = 3;
  PRODUCT_STATUS_LOCKED = 4;
  PRODUCT_STATUS_SOLD = 5;
  PRODUCT_STATUS_DELISTED = 6;
  PRODUCT_STATUS_ARCHIVED = 7;
}

enum TagType {
  TAG_TYPE_UNSPECIFIED = 0;
  TAG_TYPE_MANUAL = 1;
  TAG_TYPE_RULE = 2;
  TAG_TYPE_AI = 3;
  TAG_TYPE_CATEGORY = 4;
}

message Price {
  string currency_code = 1; // e.g., "CNY", "USD"
  int64 amount_minor_units = 2; // Price in minor units (e.g., cents for USD, fen for CNY)
}

message AuditInfo {
  string user_id = 1; // Can be actual user ID or system ID (e.g., "ai_service")
  string ip_address = 2;
  string user_agent = 3;
}

message PaginationRequest {
  int32 page_size = 1;
  string page_token = 2; // For token-based pagination
}

message PaginationResponse {
  string next_page_token = 1;
  int32 total_size = 2; // Optional: total number of items
}
```

---

## 1. AI Service (`ai_service.proto`)

*Purpose: Provides AI-powered content generation for products.*

```protobuf
syntax = "proto3";

package ai.v1;

import "common.proto";
import "google/protobuf/struct.proto"; // For JSON-like structures

option go_package = "github.com/yourorg/project/gen/go/ai/v1;aipb";

// AIService is responsible for generating product content using AI.
service AIService {
  // Generates tags and multilingual descriptions for a given product.
  rpc GenerateProductContent(GenerateProductContentRequest) returns (GenerateProductContentResponse);
  // (Future) Suggests tags for a merchant based on their profile or products.
  // rpc SuggestMerchantTags(SuggestMerchantTagsRequest) returns (SuggestMerchantTagsResponse);
}

message ProductInput {
  string product_id_internal = 1; // Internal product ID for reference
  string title = 2;
  string existing_description = 3;
  repeated string image_urls = 4;
  // Add other relevant product fields that AI might use
}

message GenerateProductContentRequest {
  ProductInput product_input = 1;
  repeated string target_languages = 2; // e.g., ["en-US", "ja-JP"]
  common.AuditInfo audit_info = 3;
}

message AIGeneratedTag {
  string name = 1;
  double confidence_score = 2;
  string category = 3; // Optional: AI might categorize tags
}

message AIGeneratedDescription {
  string language_code = 1;
  string description_text = 2;
}

message GenerateProductContentResponse {
  string product_id_internal = 1;
  google.protobuf.Value generated_tags = 2; // JSON Value representing tags. Could be a list of strings or key-value pairs.
                                          // Example: {"tags": ["vintage", "leather"], "attributes": {"color": "red"}}
  repeated AIGeneratedDescription generated_descriptions = 3;
  string ai_service_version = 4; // Version of the AI model/service used
}

// --- Merchant Tag Suggestion (Future) ---
// message SuggestMerchantTagsRequest {
//   string merchant_id = 1;
//   // Potentially include merchant details or product samples
//   common.AuditInfo audit_info = 2;
// }
//
// message SuggestedMerchantTag {
//   string tag_name = 1;
//   double confidence = 2;
//   string rationale = 3; // Why this tag is suggested
// }
//
// message SuggestMerchantTagsResponse {
//   string merchant_id = 1;
//   repeated SuggestedMerchantTag suggested_tags = 2;
// }
```

---

## 2. Product Service (`product_service.proto`)

*Purpose: Manages product information, including core data, AI enrichment, tags, and inventory status.*

```protobuf
syntax = "proto3";

package product.v1;

import "common.proto";
import "google/protobuf/struct.proto"; // For JSON fields

option go_package = "github.com/yourorg/project/gen/go/product/v1;productpb";

message Product {
  string product_id = 1;
  string erp_product_id = 2;
  string sku_ref_id = 3;
  string merchant_id = 4;
  string title = 5;
  string description = 6;
  google.protobuf.Value images = 7; // JSON array of image objects
  common.Price price = 8;
  common.ProductStatus status = 9;
  bool ai_enriched = 10;
  common.Timestamp created_at = 11;
  common.Timestamp updated_at = 12;
  int32 version = 13;
}

message ProductTag {
  string tag_id = 1;
  string name = 2;
  common.TagType type = 3;
}

message ProductAITags {
  string product_id = 1;
  google.protobuf.Value tags = 2; // JSON from product_ai_tags.tags
  string ai_service_version = 3;
  common.Timestamp created_at = 4;
}

message ProductAIDescription {
  string product_id = 1;
  string language_code = 10; // Field numbers should be unique within the message. Corrected from 1.
  string description = 11;   // Field numbers should be unique within the message. Corrected from 2.
  string ai_service_version = 12; // Field numbers should be unique within the message. Corrected from 3.
  common.Timestamp created_at = 13; // Field numbers should be unique within the message. Corrected from 4.
}

// ProductService manages all aspects of products.
service ProductService {
  // Creates a new product. Can be from ERP sync or manual entry.
  rpc CreateProduct(CreateProductRequest) returns (CreateProductResponse);
  // Gets a product by its ID.
  rpc GetProduct(GetProductRequest) returns (GetProductResponse);
  // Updates an existing product.
  rpc UpdateProduct(UpdateProductRequest) returns (UpdateProductResponse);
  // Deletes a product (logical delete, mark as archived).
  rpc DeleteProduct(DeleteProductRequest) returns (common.Empty);

  // Lists products with filtering and pagination.
  rpc ListProducts(ListProductsRequest) returns (ListProductsResponse);

  // Attempts to lock a product for an order. Called by OrderService.
  // This is an idempotent operation.
  rpc LockProduct(LockProductRequest) returns (LockProductResponse);
  // Unlocks a product if an order is cancelled or fails.
  rpc UnlockProduct(UnlockProductRequest) returns (UnlockProductResponse);
  // Marks a product as sold.
  rpc MarkProductSold(MarkProductSoldRequest) returns (common.Empty);

  // --- AI Data Integration ---
  // Adds AI-generated tags to a product.
  rpc AddProductAITags(AddProductAITagsRequest) returns (common.Empty);
  // Adds AI-generated description to a product.
  rpc AddProductAIDescription(AddProductAIDescriptionRequest) returns (common.Empty);
  // Gets AI tags for a product.
  rpc GetProductAITags(GetProductAITagsRequest) returns (GetProductAITagsResponse);
  // Gets AI descriptions for a product.
  rpc ListProductAIDescriptions(ListProductAIDescriptionsRequest) returns (ListProductAIDescriptionsResponse);


  // --- Manual Product Tags ---
  // Adds a manual tag to a product.
  rpc AddTagToProduct(AddTagToProductRequest) returns (common.Empty);
  // Removes a manual tag from a product.
  rpc RemoveTagFromProduct(RemoveTagFromProductRequest) returns (common.Empty);
  // Lists all manual tags associated with a product.
  rpc ListProductTags(ListProductTagsRequest) returns (ListProductTagsResponse);
}

// --- CreateProduct ---
message CreateProductRequest {
  // erp_product_id can be optional if not from ERP
  string erp_product_id = 1;
  string sku_ref_id = 2;
  string merchant_id = 3;
  string title = 4;
  string description = 5;
  google.protobuf.Value images = 6; // JSON array
  common.Price price = 7;
  // Initial status, usually DRAFT or PENDING_APPROVAL if AI is involved
  common.ProductStatus initial_status = 8;
  common.AuditInfo audit_info = 9;
}

message CreateProductResponse {
  Product product = 1;
}

// --- GetProduct ---
message GetProductRequest {
  string product_id = 1;
}

message GetProductResponse {
  Product product = 1;
  repeated ProductTag manual_tags = 2; // Include curated tags
}

// --- UpdateProduct ---
message UpdateProductRequest {
  string product_id = 1;
  // Fields to update - use field masks for partial updates
  optional string title = 2;
  optional string description = 3;
  optional google.protobuf.Value images = 4;
  optional common.Price price = 5;
  optional common.ProductStatus status = 6; // e.g. for manual approval
  // Add other updatable fields as needed
  int32 version = 7; // For optimistic locking
  common.AuditInfo audit_info = 8;
}

// --- DeleteProduct ---
message DeleteProductRequest {
  string product_id = 1;
  common.AuditInfo audit_info = 2;
}

// --- ListProducts ---
message ListProductsFilter {
  repeated string merchant_ids = 1;
  common.ProductStatus status = 2;
  bool ai_enriched = 3;
  repeated string product_tag_ids = 4; // Filter by manual product tags
  // (merchant_tags filter is applied by getting merchant_ids from MerchantService first)
  string title_contains = 5;
}

message ListProductsRequest {
  ListProductsFilter filter = 1;
  common.PaginationRequest pagination = 2;
  string sort_by = 3; // e.g., "created_at_desc", "price_asc"
}

message ListProductsResponse {
  repeated Product products = 1;
  common.PaginationResponse pagination = 2;
}

// --- LockProduct ---
message LockProductRequest {
  string product_id = 1;
  string order_id = 2; // For context and logging
  common.AuditInfo audit_info = 3; // audit_info.user_id would be "order_service"
}

message LockProductResponse {
  bool success = 1;
  Product product = 2; // Current state of the product
}

// --- UnlockProduct ---
message UnlockProductRequest {
  string product_id = 1;
  string order_id = 2;
  common.AuditInfo audit_info = 3;
}

// --- MarkProductSold ---
message MarkProductSoldRequest {
  string product_id = 1;
  string order_id = 2;
  common.AuditInfo audit_info = 3;
}

// --- AI Data ---
message AddProductAITagsRequest {
  string product_id = 1;
  google.protobuf.Value tags_json = 2; // JSON from AI service
  string ai_service_version = 3;
  common.AuditInfo audit_info = 4; // user_id here is "ai_service" or similar
}

message AddProductAIDescriptionRequest {
  string product_id = 1;
  string language_code = 2;
  string description = 3;
  string ai_service_version = 4;
  common.AuditInfo audit_info = 5; // user_id here is "ai_service" or similar
}

message GetProductAITagsRequest {
  string product_id = 1;
}

message GetProductAITagsResponse {
  ProductAITags ai_tags = 1;
}

message ListProductAIDescriptionsRequest {
  string product_id = 1;
}

message ListProductAIDescriptionsResponse {
  repeated ProductAIDescription ai_descriptions = 1;
}

// --- Manual Product Tags ---
message AddTagToProductRequest {
  string product_id = 1;
  string tag_id = 2; // Assumes tag already exists in the system via MerchantService or a shared Tag management
  common.AuditInfo audit_info = 3;
}

message RemoveTagFromProductRequest {
  string product_id = 1;
  string tag_id = 2;
  common.AuditInfo audit_info = 3;
}

message ListProductTagsRequest {
  string product_id = 1;
}

message ListProductTagsResponse {
  repeated ProductTag tags = 1;
}

```

---

## 3. Merchant Service (`merchant_service.proto`)

*Purpose: Manages merchant information and their associated tags. Also manages the global Tag definitions.*

```protobuf
syntax = "proto3";

package merchant.v1;

import "common.proto";
import "google/protobuf/struct.proto";

option go_package = "github.com/yourorg/project/gen/go/merchant/v1;merchantpb";

message Merchant {
  string merchant_id = 1;
  string name = 2;
  string source_type = 3;
  google.protobuf.Value contact_info = 4; // JSON
  string status = 5; // e.g., "active", "inactive"
  common.Timestamp created_at = 6;
  common.Timestamp updated_at = 7;
}

message Tag {
  string tag_id = 1;
  string name = 2;
  common.TagType type = 3;
  string description = 4;
  common.Timestamp created_at = 5;
  common.Timestamp updated_at = 6;
}

// MerchantService manages merchants and their tags.
service MerchantService {
  // --- Merchant CRUD ---
  rpc CreateMerchant(CreateMerchantRequest) returns (CreateMerchantResponse);
  rpc GetMerchant(GetMerchantRequest) returns (GetMerchantResponse);
  rpc UpdateMerchant(UpdateMerchantRequest) returns (UpdateMerchantResponse);
  rpc ListMerchants(ListMerchantsRequest) returns (ListMerchantsResponse);

  // --- Tag Definition CRUD ---
  rpc CreateTag(CreateTagRequest) returns (CreateTagResponse);
  rpc GetTag(GetTagRequest) returns (GetTagResponse);
  rpc UpdateTag(UpdateTagRequest) returns (UpdateTagResponse);
  rpc ListTags(ListTagsRequest) returns (ListTagsResponse); // Lists all available global tags

  // --- Merchant Tagging ---
  // Adds a tag to a specific merchant.
  rpc AddTagToMerchant(AddTagToMerchantRequest) returns (common.Empty);
  // Removes a tag from a specific merchant.
  rpc RemoveTagFromMerchant(RemoveTagFromMerchantRequest) returns (common.Empty);
  // Lists tags for a specific merchant.
  rpc ListMerchantTags(ListMerchantTagsRequest) returns (ListMerchantTagsResponse);
  // Lists merchants that have ALL specified tags.
  rpc ListMerchantsByTags(ListMerchantsByTagsRequest) returns (ListMerchantsByTagsResponse);
}

// --- Merchant CRUD Messages ---
message CreateMerchantRequest {
  string name = 1;
  string source_type = 2;
  google.protobuf.Value contact_info = 3;
  string status = 4;
  common.AuditInfo audit_info = 5;
}

message CreateMerchantResponse {
  Merchant merchant = 1;
}

message GetMerchantRequest {
  string merchant_id = 1;
}

message GetMerchantResponse {
  Merchant merchant = 1;
  repeated Tag tags = 2; // Tags associated with this merchant
}

message UpdateMerchantRequest {
  string merchant_id = 1;
  optional string name = 2;
  optional string source_type = 3;
  optional google.protobuf.Value contact_info = 4;
  optional string status = 5;
  common.AuditInfo audit_info = 6;
}

message ListMerchantsFilter {
  string name_contains = 1;
  string status = 2;
  // Potentially filter by source_type
}

message ListMerchantsRequest {
  ListMerchantsFilter filter = 1;
  common.PaginationRequest pagination = 2;
}

message ListMerchantsResponse {
  repeated Merchant merchants = 1;
  common.PaginationResponse pagination = 2;
}

// --- Tag Definition CRUD Messages ---
message CreateTagRequest {
  string name = 1;
  common.TagType type = 2;
  string description = 3;
  common.AuditInfo audit_info = 4;
}

message CreateTagResponse {
  Tag tag = 1;
}

message GetTagRequest {
  string tag_id = 1;
}

message GetTagResponse {
  Tag tag = 1;
}

message UpdateTagRequest {
  string tag_id = 1;
  optional string name = 2;
  optional string description = 3;
  // Tag type typically should not be updated once created, or with strict rules
  common.AuditInfo audit_info = 4;
}

message ListTagsFilter {
  common.TagType type = 1;
  string name_contains = 2;
}

message ListTagsRequest {
  ListTagsFilter filter = 1;
  common.PaginationRequest pagination = 2;
}

message ListTagsResponse {
  repeated Tag tags = 1;
  common.PaginationResponse pagination = 2;
}

// --- Merchant Tagging Messages ---
message AddTagToMerchantRequest {
  string merchant_id = 1;
  string tag_id = 2;
  common.AuditInfo audit_info = 3;
}

message RemoveTagFromMerchantRequest {
  string merchant_id = 1;
  string tag_id = 2;
  common.AuditInfo audit_info = 3;
}

message ListMerchantTagsRequest {
  string merchant_id = 1;
}

message ListMerchantTagsResponse {
  repeated Tag tags = 1;
}

message ListMerchantsByTagsRequest {
  repeated string tag_ids = 1; // Find merchants that have all these tags
  common.PaginationRequest pagination = 2;
}

message ListMerchantsByTagsResponse {
  repeated string merchant_ids = 1; // Returns only IDs for ProductService to use
  common.PaginationResponse pagination = 2;
}
```

---

## 4. Listing Service (or Platform Listing Service) (`listing_service.proto`)

*Purpose: Manages the publishing of products to various external platforms.*

```protobuf
syntax = "proto3";

package listing.v1;

import "common.proto";
import "google/protobuf/struct.proto";

option go_package = "github.com/yourorg/project/gen/go/listing/v1;listingpb";

message Listing {
  string listing_id = 1;
  string product_id = 2;
  string platform_id = 3; // e.g., "xianyu", "dewu"
  string platform_product_id = 4; // ID on the external platform
  string status = 5; // e.g., "PENDING_PUBLISH", "PUBLISHED", "UNPUBLISHED"
  common.Timestamp published_at = 6;
  common.Timestamp unpublished_at = 7;
  google.protobuf.Value publish_details = 8; // Platform-specific details
  common.Timestamp created_at = 9;
  common.Timestamp updated_at = 10;
}

// ListingService manages product listings on external platforms.
service ListingService {
  // Publishes a product to a specified platform.
  rpc PublishProduct(PublishProductRequest) returns (PublishProductResponse);
  // Unpublishes (delists) a product from a specified platform.
  rpc UnpublishProduct(UnpublishProductRequest) returns (UnpublishProductResponse);
  // Notifies the listing service that a product has been sold (e.g., via an order)
  // so it can unlist from other platforms if necessary.
  rpc NotifyProductSold(NotifyProductSoldRequest) returns (common.Empty);
  // Gets the listing status of a product on a specific platform.
  rpc GetListingStatus(GetListingStatusRequest) returns (GetListingStatusResponse);
  // Lists all listings for a given product.
  rpc ListProductListings(ListProductListingsRequest) returns (ListProductListingsResponse);
}

message PublishProductRequest {
  string product_id = 1;
  string platform_id = 2;
  google.protobuf.Value platform_specific_options = 3; // e.g., category mapping, shipping templates
  common.AuditInfo audit_info = 4;
}

message PublishProductResponse {
  Listing listing = 1;
  string message = 2; // e.g., "Successfully submitted for publishing"
}

message UnpublishProductRequest {
  string product_id = 1;
  string platform_id = 2;
  string reason = 3; // Optional reason for unpublishing
  common.AuditInfo audit_info = 4;
}

message UnpublishProductResponse {
  string listing_id = 1;
  string status_after_unpublish = 2; // e.g., "UNPUBLISHED"
  string message = 3;
}

message NotifyProductSoldRequest {
  string product_id = 1;
  string order_id = 2;
  string platform_where_sold = 3; // The platform on which the item was sold
  common.AuditInfo audit_info = 4;
}

message GetListingStatusRequest {
  string product_id = 1;
  string platform_id = 2;
}

message GetListingStatusResponse {
  Listing listing = 1;
}

message ListProductListingsRequest {
  string product_id = 1;
}

message ListProductListingsResponse {
  repeated Listing listings = 1;
}
```

---

## 5. Order Service (`order_service.proto`)

*Purpose: Handles order creation, status management, and fulfillment coordination.*
*(This is a separate microservice but its interface is defined here for completeness of the system design).*

```protobuf
syntax = "proto3";

package order.v1;

import "common.proto";
import "google/protobuf/struct.proto";

option go_package = "github.com/yourorg/project/gen/go/order/v1;orderpb";

enum OrderStatus {
  ORDER_STATUS_UNSPECIFIED = 0;
  ORDER_STATUS_PENDING_LOCK = 1;
  ORDER_STATUS_LOCK_FAILED = 2;
  ORDER_STATUS_PENDING_PAYMENT = 3; // Product locked, awaiting payment
  ORDER_STATUS_PAID = 4;
  ORDER_STATUS_PAYMENT_FAILED = 5;
  ORDER_STATUS_SHIPPED = 6;
  ORDER_STATUS_DELIVERED = 7;
  ORDER_STATUS_CANCELLED = 8;
  ORDER_STATUS_REFUNDED = 9;
}

message Order {
  string order_id = 1;
  string external_order_id = 2; // Platform's order ID
  string product_id = 3;
  string buyer_id = 4;
  string merchant_id = 5; // Denormalized from product for easier access
  string platform_id = 6;
  OrderStatus status = 7;
  common.Price order_price = 8; // Price at the time of order
  google.protobuf.Value shipping_address = 9; // JSON
  google.protobuf.Value payment_details = 10; // JSON, e.g., transaction ID
  string cancellation_reason = 11;
  common.Timestamp created_at = 12;
  common.Timestamp paid_at = 13;
  common.Timestamp updated_at = 14;
  int32 version = 15;
}

// OrderService manages the lifecycle of orders.
service OrderService {
  // Creates an order and attempts to lock the product.
  rpc PlaceOrder(PlaceOrderRequest) returns (PlaceOrderResponse);
  // Cancels an order. May involve unlocking the product.
  rpc CancelOrder(CancelOrderRequest) returns (CancelOrderResponse);
  // Gets the status of a specific order.
  rpc GetOrderStatus(GetOrderStatusRequest) returns (GetOrderStatusResponse);
  // Lists orders, possibly filtered.
  rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse);

  // Callback for payment gateway or internal payment confirmation
  rpc ConfirmPayment(ConfirmPaymentRequest) returns (ConfirmPaymentResponse);
  // Marks an order as shipped
  rpc MarkOrderShipped(MarkOrderShippedRequest) returns (common.Empty);
  // Marks an order as delivered
  rpc MarkOrderDelivered(MarkOrderDeliveredRequest) returns (common.Empty);
}

message OrderItemInput {
  string product_id = 1;
  common.Price price_at_purchase = 2; // Ensure price is fixed at order time
  // Potentially quantity if applicable, though luxury goods are often single items
}

message PlaceOrderRequest {
  string external_order_id = 1; // Optional, if order originates from an external platform
  string buyer_id = 2;
  string platform_id = 3;
  OrderItemInput item = 4;
  google.protobuf.Value shipping_address = 5;
  common.AuditInfo audit_info = 6; // user_id might be platform system or actual user
}

message PlaceOrderResponse {
  Order order = 1;
  bool product_lock_acquired = 2;
  string message = 3; // e.g., "Order created, awaiting payment" or "Failed to lock product"
}

message CancelOrderRequest {
  string order_id = 1;
  string reason = 2;
  common.AuditInfo audit_info = 3;
}

message CancelOrderResponse {
  Order order = 1; // Updated order
  string message = 2;
}

message GetOrderStatusRequest {
  string order_id = 1;
}

message GetOrderStatusResponse {
  Order order = 1;
}

message ListOrdersFilter {
  string buyer_id = 1;
  string product_id = 2;
  string merchant_id = 3;
  string platform_id = 4;
  OrderStatus status = 5;
}

message ListOrdersRequest {
  ListOrdersFilter filter = 1;
  common.PaginationRequest pagination = 2;
}

message ListOrdersResponse {
  repeated Order orders = 1;
  common.PaginationResponse pagination = 2;
}

message ConfirmPaymentRequest {
  string order_id = 1;
  string payment_transaction_id = 2;
  string payment_method = 3;
  common.Timestamp paid_at = 4;
  common.AuditInfo audit_info = 5;
}

message ConfirmPaymentResponse {
  Order order = 1;
  string message = 2;
}

message MarkOrderShippedRequest {
  string order_id = 1;
  string tracking_number = 2;
  string carrier = 3;
  common.AuditInfo audit_info = 4;
}

message MarkOrderDeliveredRequest {
  string order_id = 1;
  common.Timestamp delivered_at = 2;
  common.AuditInfo audit_info = 3;
}
```

---

## Service Interactions Summary:

1.  **ERP -> Product Service (`CreateProduct`):** ERP syncs product data.
2.  **Product Service -> AI Service (`GenerateProductContent`):** Product Service calls AI Service to get tags/descriptions.
3.  **AI Service (Callback/Response) -> Product Service (`AddProductAITags`, `AddProductAIDescription`):** Product Service stores AI-generated content.
4.  **User/Admin -> Product Service (`AddTagToProduct`, `UpdateProduct` etc.):** Manual data curation.
5.  **User/Admin -> Merchant Service (`CreateTag`, `AddTagToMerchant` etc.):** Merchant and global tag management.
6.  **Product Service -> Merchant Service (`ListMerchantsByTags`):** When filtering products by merchant tags, Product Service queries Merchant Service for relevant `merchant_id`s.
7.  **Product Service -> Listing Service (`PublishProduct`):** Product Service requests Listing Service to publish a product.
8.  **Platform User -> Platform -> Order Service (`PlaceOrder`):** User places an order on an external platform.
9.  **Order Service -> Product Service (`LockProduct`):** Order Service attempts to lock the product.
10. **Product Service (Event/Callback) -> Listing Service (`NotifyProductSold` or Listing Service subscribes to `ProductLocked` event):** If product is locked/sold, Listing Service is notified to delist from other platforms.
11. **Order Service -> Product Service (`UnlockProduct`):** If order is cancelled or payment fails after lock.

This API design provides a comprehensive set of RPCs for managing the lifecycle of luxury goods, from creation and AI enrichment to multi-platform listing and order fulfillment.
