# Product Service API Design (gRPC)

## AIService
```protobuf
syntax = "proto3";

package product;

option go_package = "luxurygoods/productpb";

// AIService is called by the Product Service to enrich product data.
service AIService {
  // Generates product tags and multilingual descriptions
  rpc GenerateProductContent(GenerateProductContentRequest) returns (GenerateProductContentResponse);
}

message GenerateProductContentRequest {
  string product_id = 1;
  // Potentially include other product details if AI service needs them
  // e.g., title, existing_description, image_urls
  string title = 2;
  string existing_description = 3;
  repeated string image_urls = 4;
}

message GenerateProductContentResponse {
  string product_id = 1;
  string tags_json = 2; // JSON string of tags
  repeated LanguageDescription descriptions = 3;
}

message LanguageDescription {
  string language_code = 1; // "en", "ja", "zh", etc.
  string description = 2;
}
```

## ProductService
```protobuf
syntax = "proto3";

package product;

option go_package = "luxurygoods/productpb";

// ProductService provides core functionalities for managing products.
service ProductService {
  // Creates a new product
  rpc CreateProduct(CreateProductRequest) returns (CreateProductResponse);
  // Retrieves a product by its ID
  rpc GetProduct(GetProductRequest) returns (GetProductResponse);
  // Lists products, supports filtering by merchant_tags and other criteria
  rpc ListProducts(ListProductsRequest) returns (ListProductsResponse);
  // Updates an existing product
  rpc UpdateProduct(UpdateProductRequest) returns (UpdateProductResponse);
  // Deletes a product
  rpc DeleteProduct(DeleteProductRequest) returns (DeleteProductResponse);
  // Locks a product when an order is placed or for other reasons
  rpc LockProduct(LockProductRequest) returns (LockProductResponse);
  // Unlocks a product if an order is cancelled or lock expires
  rpc UnlockProduct(UnlockProductRequest) returns (UnlockProductResponse);
}

message Product {
  string product_id = 1;
  string title = 2;
  string description = 3;
  repeated string images = 4; // URLs to images
  double price = 5;
  string sku_ref_id = 6; // Standard Product Unit Reference ID (标品id)
  string status = 7; // ENUM: "draft", "listed", "locked", "sold"
  string merchant_id = 8;
  bool ai_enriched = 9;
  string created_at = 10; // Timestamp
  string updated_at = 11; // Timestamp
}

message CreateProductRequest {
  string title = 1;
  string description = 2;
  repeated string images = 3;
  double price = 4;
  string sku_ref_id = 5;
  string merchant_id = 6;
  // bool trigger_ai_enrichment = 7; // Option to trigger AI enrichment on creation
}

message CreateProductResponse {
  Product product = 1;
}

message GetProductRequest {
  string product_id = 1;
}

message GetProductResponse {
  Product product = 1;
}

message ListProductsRequest {
  int32 page_size = 1;
  string page_token = 2;
  repeated string merchant_tags = 3; // Filter by merchant tags (names or IDs)
  string merchant_id = 4; // Filter by a specific merchant
  string status = 5; // Filter by product status
  // Add other filters as needed, e.g., price_min, price_max, ai_enriched
}

message ListProductsResponse {
  repeated Product products = 1;
  string next_page_token = 2;
}

message UpdateProductRequest {
  string product_id = 1;
  // Include fields that can be updated. Use FieldMask for partial updates.
  // google.protobuf.FieldMask update_mask = 2;
  string title = 3;
  string description = 4;
  repeated string images = 5;
  double price = 6;
  string status = 7;
}

message UpdateProductResponse {
  Product product = 1;
}

message DeleteProductRequest {
  string product_id = 1;
}

message DeleteProductResponse {
  bool success = 1;
}

message LockProductRequest {
  string product_id = 1;
  string order_id = 2; // To associate the lock with an order, if applicable
  string lock_reason = 3; // Optional: reason for locking
}

message LockProductResponse {
  bool success = 1;
  Product product = 2; // Returns the updated product with "locked" status
}

message UnlockProductRequest {
  string product_id = 1;
  string order_id = 2; // If lock was order-specific
}

message UnlockProductResponse {
  bool success = 1;
  Product product = 2; // Returns the updated product, likely "listed" or original status
}
```

## MerchantService (Client Interface)
This describes the gRPC client interfaces that Product Service will use to communicate with a MerchantService.
The actual MerchantService definition would reside in the Merchant Service's own `api_design.md`.

```protobuf
syntax = "proto3";

package merchant; // Assuming a different package for merchant service protos

option go_package = "luxurygoods/merchantpb";

// MerchantService defines interactions with merchant data.
// This is a client view from Product Service's perspective.
service MerchantService {
  // Adds a tag to a merchant
  rpc AddTagToMerchant(AddTagToMerchantRequest) returns (AddTagToMerchantResponse);
  // Removes a tag from a merchant
  rpc RemoveTagFromMerchant(RemoveTagFromMerchantRequest) returns (RemoveTagFromMerchantResponse);
  // Lists merchants by tags
  rpc ListMerchantsByTags(ListMerchantsByTagsRequest) returns (ListMerchantsByTagsResponse);
  // Gets merchant details (potentially used by Product Service)
  rpc GetMerchant(GetMerchantRequest) returns (GetMerchantResponse);
}

message Tag {
  string tag_id = 1;
  string name = 2;
  string type = 3; // 'manual', 'rule', 'ai'
}

message Merchant {
  string merchant_id = 1;
  string name = 2;
  string source_type = 3;
  string contact = 4;
  repeated Tag tags = 5;
}

message AddTagToMerchantRequest {
  string merchant_id = 1;
  string tag_id = 2; // Assuming tag_id is used, could be tag_name
  // string operator_id = 3; // For audit purposes
  // string source = 4; // e.g. "product_service_rule"
}

message AddTagToMerchantResponse {
  bool success = 1;
  Merchant merchant = 2; // Return updated merchant
}

message RemoveTagFromMerchantRequest {
  string merchant_id = 1;
  string tag_id = 2;
  // string operator_id = 3; // For audit
}

message RemoveTagFromMerchantResponse {
  bool success = 1;
  Merchant merchant = 2; // Return updated merchant
}

message ListMerchantsByTagsRequest {
  repeated string tag_ids = 1; // Or tag_names
  // int32 page_size = 2;
  // string page_token = 3;
}

message ListMerchantsByTagsResponse {
  repeated Merchant merchants = 1;
  // string next_page_token = 2;
}

message GetMerchantRequest {
  string merchant_id = 1;
}

message GetMerchantResponse {
  Merchant merchant = 1;
}
```

## PublishingService (Client Interface)
This describes the gRPC client interfaces that Product Service will use to communicate with a PublishingService.
The actual PublishingService definition would reside in its own `api_design.md`.

```protobuf
syntax = "proto3";

package publishing; // Assuming a different package for publishing service protos

option go_package = "luxurygoods/publishingpb";

// PublishingService is responsible for listing products on various external platforms.
// This is a client view from Product Service's perspective.
service PublishingService {
  // Publishes a product to a specific platform
  rpc PublishProduct(PublishProductRequest) returns (PublishProductResponse);
  // Unpublishes (delists) a product from a specific platform
  rpc UnpublishProduct(UnpublishProductRequest) returns (UnpublishProductResponse);
  // Notifies the PublishingService that a product has been sold,
  // so it can be unlisted from other platforms.
  rpc NotifyProductSold(NotifyProductSoldRequest) returns (NotifyProductSoldResponse);
  // Gets the status of a listing on a platform
  rpc GetListingStatus(GetListingStatusRequest) returns (GetListingStatusResponse);
}

message PublishProductRequest {
  string product_id = 1;
  string platform_id = 2; // e.g., "xianyu", "dewu", "ebay_US"
  // Potentially product details if not fetched by PublishingService itself
  // string title = 3;
  // double price = 4;
  // repeated string image_urls = 5;
}

message PublishProductResponse {
  string listing_id = 1; // ID of the listing on the external platform
  string status = 2; // e.g., "published", "pending_review", "failed"
  string platform_listing_url = 3; // Optional URL to the listing
}

message UnpublishProductRequest {
  string product_id = 1;
  string platform_id = 2;
  string listing_id = 3; // Optional: ID of the listing on the platform
}

message UnpublishProductResponse {
  bool success = 1;
}

message NotifyProductSoldRequest {
  string product_id = 1;
  string order_id = 2;
  string platform_id_sold_on = 3; // Platform where the item was sold
  // Potentially a list of other platforms it was listed on, to ensure cleanup
  // repeated string other_platform_ids = 4;
}

message NotifyProductSoldResponse {
  bool success = 1;
  // repeated UnpublishAttemptResult results = 2; // Optional: results of unpublishing from other platforms
}

// message UnpublishAttemptResult {
//   string platform_id = 1;
//   bool success = 2;
//   string error_message = 3; // if not successful
// }

message GetListingStatusRequest {
  string product_id = 1;
  string platform_id = 2;
  string listing_id = 3; // Optional: ID of the listing on the platform
}

message GetListingStatusResponse {
  string product_id = 1;
  string platform_id = 2;
  string listing_id = 3;
  string status = 4; // e.g., "active", "sold", "expired", "delisted"
  string platform_listing_url = 5; // Optional URL to the listing
}
```
