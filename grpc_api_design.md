# gRPC API Design Document

This document outlines the gRPC service and message definitions for the二手奢侈品商品微服务系统 (Second-hand Luxury Goods Microservice System).

**Proto Version:** `proto3`

**Common Data Types / Enums (to be defined in a common.proto file):**

```protobuf
syntax = "proto3";

package common;

option go_package = "github.com/yourorg/project/gen/go/common;commonpb";

message Empty { // -- 空消息，用于不需要参数或返回值的RPC
}

message Timestamp { // -- 时间戳
  int64 seconds = 1; // -- 秒
  int32 nanos = 2;   // -- 纳秒
}

enum ProductStatus { // -- 商品状态枚举
  PRODUCT_STATUS_UNSPECIFIED = 0;      // -- 未指定状态
  PRODUCT_STATUS_DRAFT = 1;            // -- 草稿
  PRODUCT_STATUS_PENDING_APPROVAL = 2; // -- 待审核
  PRODUCT_STATUS_LISTED = 3;           // -- 已上架
  PRODUCT_STATUS_LOCKED = 4;           // -- 已锁定
  PRODUCT_STATUS_SOLD = 5;             // -- 已售出
  PRODUCT_STATUS_DELISTED = 6;         // -- 已下架（手动）
  PRODUCT_STATUS_ARCHIVED = 7;         // -- 已归档
}

enum TagType { // -- 标签类型枚举
  TAG_TYPE_UNSPECIFIED = 0; // -- 未指定类型
  TAG_TYPE_MANUAL = 1;      // -- 手动创建
  TAG_TYPE_RULE = 2;        // -- 规则生成
  TAG_TYPE_AI = 3;          // -- AI生成
  TAG_TYPE_CATEGORY = 4;    // -- 分类标签
}

message Price { // -- 价格信息
  string currency_code = 1;      // e.g., "CNY", "USD" -- 货币代码 (例如 "CNY", "USD")
  int64 amount_minor_units = 2; // Price in minor units (e.g., cents for USD, fen for CNY) -- 价格，以最小货币单位表示 (例如 人民币的分，美元的美分)
}

message AuditInfo { // -- 审计信息
  string user_id = 1;    // Can be actual user ID or system ID (e.g., "ai_service") -- 用户ID (可以是真实用户ID或系统ID，如 "ai_service")
  string ip_address = 2; // -- IP地址
  string user_agent = 3; // -- 用户代理
}

message PaginationRequest { // -- 分页请求
  int32 page_size = 1;    // -- 每页大小
  string page_token = 2;  // For token-based pagination -- 分页令牌 (用于基于令牌的分页)
}

message PaginationResponse { // -- 分页响应
  string next_page_token = 1; // -- 下一页的令牌
  int32 total_size = 2;       // Optional: total number of items -- 总项目数 (可选)
}
```

---

## 1. AI Service (`ai_service.proto`)

*Purpose: Provides AI-powered content generation for products. -- AI服务：提供AI驱动的商品内容生成功能。*

```protobuf
syntax = "proto3";

package ai.v1;

import "common.proto";
import "google/protobuf/struct.proto"; // For JSON-like structures -- 用于JSON结构

option go_package = "github.com/yourorg/project/gen/go/ai/v1;aipb";

// AIService is responsible for generating product content using AI. -- AIService负责使用AI生成商品内容。
service AIService {
  // Generates tags and multilingual descriptions for a given product. -- 为给定商品生成标签和多语言描述。
  rpc GenerateProductContent(GenerateProductContentRequest) returns (GenerateProductContentResponse);
  // (Future) Suggests tags for a merchant based on their profile or products. -- (未来功能) 根据商家资料或商品推荐商家标签。
  // rpc SuggestMerchantTags(SuggestMerchantTagsRequest) returns (SuggestMerchantTagsResponse);
}

message ProductInput { // -- 商品输入信息 (供AI分析)
  string product_id_internal = 1; // Internal product ID for reference -- 内部商品ID，供参考
  string title = 2;                 // -- 商品标题
  string existing_description = 3;  // -- 商品现有描述
  repeated string image_urls = 4;   // -- 商品图片URL列表
  // Add other relevant product fields that AI might use -- 添加其他AI可能使用的相关商品字段
}

message GenerateProductContentRequest { // -- 生成商品内容请求
  ProductInput product_input = 1;       // -- 商品输入信息
  repeated string target_languages = 2; // e.g., ["en-US", "ja-JP"] -- 目标语言列表 (例如 ["en-US", "ja-JP"])
  common.AuditInfo audit_info = 3;      // -- 审计信息
}

message AIGeneratedTag { // -- AI生成的标签
  string name = 1;              // -- 标签名称
  double confidence_score = 2;  // -- 置信度评分
  string category = 3;          // Optional: AI might categorize tags -- 标签分类 (可选, AI可能对标签进行分类)
}

message AIGeneratedDescription { // -- AI生成的多语言描述
  string language_code = 1;    // -- 语言代码
  string description_text = 2; // -- 描述文本
}

message GenerateProductContentResponse { // -- 生成商品内容响应
  string product_id_internal = 1; // -- 内部商品ID
  google.protobuf.Value generated_tags = 2; // JSON Value representing tags. Could be a list of strings or key-value pairs. -- AI生成的标签 (JSON Value格式，可以是字符串列表或键值对)
                                          // Example: {"tags": ["vintage", "leather"], "attributes": {"color": "red"}}
  repeated AIGeneratedDescription generated_descriptions = 3; // -- AI生成的多语言描述列表
  string ai_service_version = 4; // Version of the AI model/service used -- 使用的AI模型/服务版本
}

// --- Merchant Tag Suggestion (Future) ---
// message SuggestMerchantTagsRequest {
//   string merchant_id = 1; // -- 商家ID
//   // Potentially include merchant details or product samples -- 可能包含商家详情或商品样本
//   common.AuditInfo audit_info = 2; // -- 审计信息
// }
//
// message SuggestedMerchantTag {
//   string tag_name = 1; // -- 建议的标签名称
//   double confidence = 2; // -- 置信度
//   string rationale = 3; // Why this tag is suggested -- 推荐此标签的理由
// }
//
// message SuggestMerchantTagsResponse {
//   string merchant_id = 1; // -- 商家ID
//   repeated SuggestedMerchantTag suggested_tags = 2; // -- 建议的商家标签列表
// }
```

---

## 2. Product Service (`product_service.proto`)

*Purpose: Manages product information, including core data, AI enrichment, tags, and inventory status. -- 商品服务：管理商品信息，包括核心数据、AI丰富内容、标签和库存状态。*

```protobuf
syntax = "proto3";

package product.v1;

import "common.proto";
import "google/protobuf/struct.proto"; // For JSON fields -- 用于JSON字段

option go_package = "github.com/yourorg/project/gen/go/product/v1;productpb";

message Product { // -- 商品核心信息
  string product_id = 1;      // -- 商品ID
  string erp_product_id = 2;  // -- ERP系统商品ID
  string sku_ref_id = 3;      // -- 标品ID
  string merchant_id = 4;     // -- 商家ID
  string title = 5;           // -- 标题
  string description = 6;     // -- 描述
  google.protobuf.Value images = 7; // JSON array of image objects -- 图片信息 (JSON数组)
  common.Price price = 8;     // -- 价格
  common.ProductStatus status = 9; // -- 商品状态
  bool ai_enriched = 10;      // -- 是否经AI丰富内容
  common.Timestamp created_at = 11; // -- 创建时间
  common.Timestamp updated_at = 12; // -- 更新时间
  int32 version = 13;         // -- 版本号 (用于乐观锁)
}

message ProductTag { // -- 商品标签信息 (通常指人工或规则校准后标签)
  string tag_id = 1;   // -- 标签ID
  string name = 2;     // -- 标签名称
  common.TagType type = 3; // -- 标签类型
}

message ProductAITags { // -- 商品的AI标签数据
  string product_id = 1;    // -- 商品ID
  google.protobuf.Value tags = 2; // JSON from product_ai_tags.tags -- AI标签内容 (JSON格式)
  string ai_service_version = 3; // -- AI服务版本
  common.Timestamp created_at = 4; // -- 创建时间
}

message ProductAIDescription { // -- 商品的AI多语言描述数据
  string product_id = 1;        // -- 商品ID
  string language_code = 10;    // Field numbers should be unique within the message. Corrected from 1. -- 语言代码
  string description = 11;      // Field numbers should be unique within the message. Corrected from 2. -- 描述文本
  string ai_service_version = 12; // Field numbers should be unique within the message. Corrected from 3. -- AI服务版本
  common.Timestamp created_at = 13; // Field numbers should be unique within the message. Corrected from 4. -- 创建时间
}

// ProductService manages all aspects of products. -- ProductService管理商品的所有方面。
service ProductService {
  // Creates a new product. Can be from ERP sync or manual entry. -- 创建新商品 (可来自ERP同步或手动录入)。
  rpc CreateProduct(CreateProductRequest) returns (CreateProductResponse);
  // Gets a product by its ID. -- 根据ID获取商品信息。
  rpc GetProduct(GetProductRequest) returns (GetProductResponse);
  // Updates an existing product. -- 更新现有商品信息。
  rpc UpdateProduct(UpdateProductRequest) returns (UpdateProductResponse);
  // Deletes a product (logical delete, mark as archived). -- 删除商品 (逻辑删除，标记为已归档)。
  rpc DeleteProduct(DeleteProductRequest) returns (common.Empty);

  // Lists products with filtering and pagination. -- 列出商品 (支持筛选和分页)。
  rpc ListProducts(ListProductsRequest) returns (ListProductsResponse);

  // Attempts to lock a product for an order. Called by OrderService. -- 尝试为订单锁定商品 (由OrderService调用)。
  // This is an idempotent operation. -- 这是一个幂等操作。
  rpc LockProduct(LockProductRequest) returns (LockProductResponse);
  // Unlocks a product if an order is cancelled or fails. -- 解锁商品 (如果订单取消或失败)。
  rpc UnlockProduct(UnlockProductRequest) returns (UnlockProductResponse);
  // Marks a product as sold. -- 标记商品为已售出。
  rpc MarkProductSold(MarkProductSoldRequest) returns (common.Empty);

  // --- AI Data Integration --- -- AI数据集成接口
  // Adds AI-generated tags to a product. -- 为商品添加AI生成的标签。
  rpc AddProductAITags(AddProductAITagsRequest) returns (common.Empty);
  // Adds AI-generated description to a product. -- 为商品添加AI生成的多语言描述。
  rpc AddProductAIDescription(AddProductAIDescriptionRequest) returns (common.Empty);
  // Gets AI tags for a product. -- 获取商品的AI标签。
  rpc GetProductAITags(GetProductAITagsRequest) returns (GetProductAITagsResponse);
  // Gets AI descriptions for a product. -- 获取商品的AI多语言描述。
  rpc ListProductAIDescriptions(ListProductAIDescriptionsRequest) returns (ListProductAIDescriptionsResponse);


  // --- Manual Product Tags --- -- 商品人工标签接口
  // Adds a manual tag to a product. -- 为商品添加人工标签。
  rpc AddTagToProduct(AddTagToProductRequest) returns (common.Empty);
  // Removes a manual tag from a product. -- 从商品移除人工标签。
  rpc RemoveTagFromProduct(RemoveTagFromProductRequest) returns (common.Empty);
  // Lists all manual tags associated with a product. -- 列出商品关联的所有人工标签。
  rpc ListProductTags(ListProductTagsRequest) returns (ListProductTagsResponse);
}

// --- CreateProduct ---
message CreateProductRequest { // -- 创建商品请求
  // erp_product_id can be optional if not from ERP -- ERP商品ID (如果非ERP来源则可选)
  string erp_product_id = 1; // -- ERP商品ID (如果非ERP来源则可选)
  string sku_ref_id = 2;            // -- 标品ID
  string merchant_id = 3;           // -- 商家ID
  string title = 4;                 // -- 标题
  string description = 5;           // -- 描述
  google.protobuf.Value images = 6; // JSON array -- 图片信息 (JSON数组)
  common.Price price = 7;           // -- 价格
  // Initial status, usually DRAFT or PENDING_APPROVAL if AI is involved -- 初始状态 (通常为DRAFT，或如果涉及AI则为PENDING_APPROVAL)
  common.ProductStatus initial_status = 8; // -- 初始状态 (通常为DRAFT，或如果涉及AI则为PENDING_APPROVAL)
  common.AuditInfo audit_info = 9;  // -- 审计信息
}

message CreateProductResponse { // -- 创建商品响应
  Product product = 1; // -- 创建的商品信息
}

// --- GetProduct ---
message GetProductRequest { // -- 获取商品请求
  string product_id = 1; // -- 商品ID
}

message GetProductResponse { // -- 获取商品响应
  Product product = 1; // -- 商品信息
  repeated ProductTag manual_tags = 2; // Include curated tags -- 商品的人工标签列表
}

// --- UpdateProduct ---
message UpdateProductRequest { // -- 更新商品请求
  string product_id = 1; // -- 商品ID
  // Fields to update - use field masks for partial updates -- 需要更新的字段 (使用field_mask进行部分更新)
  optional string title = 2;                 // -- 标题 (可选)
  optional string description = 3;           // -- 描述 (可选)
  optional google.protobuf.Value images = 4; // -- 图片信息 (可选, JSON数组)
  optional common.Price price = 5;           // -- 价格 (可选)
  optional common.ProductStatus status = 6;  // e.g. for manual approval -- 商品状态 (可选, 例如用于人工审核)
  // Add other updatable fields as needed -- 添加其他需要更新的字段
  int32 version = 7; // For optimistic locking -- 版本号 (用于乐观锁)
  common.AuditInfo audit_info = 8;           // -- 审计信息
}

// --- DeleteProduct ---
message DeleteProductRequest { // -- 删除商品请求
  string product_id = 1;          // -- 商品ID
  common.AuditInfo audit_info = 2; // -- 审计信息
}

// --- ListProducts ---
message ListProductsFilter { // -- 列出商品筛选条件
  repeated string merchant_ids = 1;   // -- 商家ID列表
  common.ProductStatus status = 2;    // -- 商品状态
  bool ai_enriched = 3;               // -- 是否经AI丰富
  repeated string product_tag_ids = 4; // Filter by manual product tags -- 商品人工标签ID列表
  // (merchant_tags filter is applied by getting merchant_ids from MerchantService first) -- (商家标签筛选首先从MerchantService获取merchant_ids)
  string title_contains = 5;          // -- 标题包含的文本
}

message ListProductsRequest { // -- 列出商品请求
  ListProductsFilter filter = 1;        // -- 筛选条件
  common.PaginationRequest pagination = 2; // -- 分页请求
  string sort_by = 3;                   // e.g., "created_at_desc", "price_asc" -- 排序字段 (例如 "created_at_desc", "price_asc")
}

message ListProductsResponse { // -- 列出商品响应
  repeated Product products = 1;          // -- 商品列表
  common.PaginationResponse pagination = 2; // -- 分页响应
}

// --- LockProduct ---
message LockProductRequest { // -- 锁定商品请求
  string product_id = 1;           // -- 商品ID
  string order_id = 2;             // For context and logging -- 订单ID (用于上下文和日志记录)
  common.AuditInfo audit_info = 3; // audit_info.user_id would be "order_service" -- 审计信息 (user_id应为 "order_service")
}

message LockProductResponse { // -- 锁定商品响应
  bool success = 1;      // -- 是否成功锁定
  Product product = 2;   // Current state of the product -- 商品当前状态
}

// --- UnlockProduct ---
message UnlockProductRequest { // -- 解锁商品请求
  string product_id = 1;           // -- 商品ID
  string order_id = 2;             // -- 订单ID
  common.AuditInfo audit_info = 3; // -- 审计信息
}

// --- MarkProductSold ---
message MarkProductSoldRequest { // -- 标记商品已售出请求
  string product_id = 1;           // -- 商品ID
  string order_id = 2;             // -- 订单ID
  common.AuditInfo audit_info = 3; // -- 审计信息
}

// --- AI Data ---
message AddProductAITagsRequest { // -- 添加商品AI标签请求
  string product_id = 1;          // -- 商品ID
  google.protobuf.Value tags_json = 2; // JSON from AI service -- AI服务返回的标签JSON
  string ai_service_version = 3;  // -- AI服务版本
  common.AuditInfo audit_info = 4; // user_id here is "ai_service" or similar -- 审计信息 (user_id为 "ai_service" 或类似)
}

message AddProductAIDescriptionRequest { // -- 添加商品AI描述请求
  string product_id = 1;          // -- 商品ID
  string language_code = 2;       // -- 语言代码
  string description = 3;         // -- 描述文本
  string ai_service_version = 4;  // -- AI服务版本
  common.AuditInfo audit_info = 5; // user_id here is "ai_service" or similar -- 审计信息 (user_id为 "ai_service" 或类似)
}

message GetProductAITagsRequest { // -- 获取商品AI标签请求
  string product_id = 1; // -- 商品ID
}

message GetProductAITagsResponse { // -- 获取商品AI标签响应
  ProductAITags ai_tags = 1; // -- 商品的AI标签数据
}

message ListProductAIDescriptionsRequest { // -- 列出商品AI描述请求
  string product_id = 1; // -- 商品ID
}

message ListProductAIDescriptionsResponse { // -- 列出商品AI描述响应
  repeated ProductAIDescription ai_descriptions = 1; // -- 商品的AI多语言描述列表
}

// --- Manual Product Tags ---
message AddTagToProductRequest { // -- 添加人工标签到商品请求
  string product_id = 1;           // -- 商品ID
  string tag_id = 2;               // Assumes tag already exists in the system via MerchantService or a shared Tag management -- 标签ID (假设标签已通过商家服务或共享标签管理在系统中存在)
  common.AuditInfo audit_info = 3; // -- 审计信息
}

message RemoveTagFromProductRequest { // -- 从商品移除人工标签请求
  string product_id = 1;           // -- 商品ID
  string tag_id = 2;               // -- 标签ID
  common.AuditInfo audit_info = 3; // -- 审计信息
}

message ListProductTagsRequest { // -- 列出商品的人工标签请求
  string product_id = 1; // -- 商品ID
}

message ListProductTagsResponse { // -- 列出商品的人工标签响应
  repeated ProductTag tags = 1; // -- 商品的人工标签列表
}

```

---

## 3. Merchant Service (`merchant_service.proto`)

*Purpose: Manages merchant information and their associated tags. Also manages the global Tag definitions. -- 商家服务：管理商家信息及其关联标签，同时管理全局标签定义。*

```protobuf
syntax = "proto3";

package merchant.v1;

import "common.proto";
import "google/protobuf/struct.proto";

option go_package = "github.com/yourorg/project/gen/go/merchant/v1;merchantpb";

message Merchant { // -- 商家信息
  string merchant_id = 1;    // -- 商家ID
  string name = 2;           // -- 商家名称
  string source_type = 3;    // -- 商家来源类型
  google.protobuf.Value contact_info = 4; // JSON -- 联系信息 (JSON)
  string status = 5;         // e.g., "active", "inactive" -- 状态 (例如 "active", "inactive")
  common.Timestamp created_at = 6; // -- 创建时间
  common.Timestamp updated_at = 7; // -- 更新时间
}

message Tag { // -- 标签定义信息
  string tag_id = 1;      // -- 标签ID
  string name = 2;        // -- 标签名称
  common.TagType type = 3; // -- 标签类型
  string description = 4; // -- 标签描述
  common.Timestamp created_at = 5; // -- 创建时间
  common.Timestamp updated_at = 6; // -- 更新时间
}

// MerchantService manages merchants and their tags. -- MerchantService管理商家及其标签。
service MerchantService {
  // --- Merchant CRUD --- -- 商家 CRUD 接口
  rpc CreateMerchant(CreateMerchantRequest) returns (CreateMerchantResponse); // -- 创建商家
  rpc GetMerchant(GetMerchantRequest) returns (GetMerchantResponse);       // -- 获取商家信息
  rpc UpdateMerchant(UpdateMerchantRequest) returns (UpdateMerchantResponse); // -- 更新商家信息
  rpc ListMerchants(ListMerchantsRequest) returns (ListMerchantsResponse);   // -- 列出商家

  // --- Tag Definition CRUD --- -- 标签定义 CRUD 接口
  rpc CreateTag(CreateTagRequest) returns (CreateTagResponse);     // -- 创建标签定义
  rpc GetTag(GetTagRequest) returns (GetTagResponse);           // -- 获取标签定义
  rpc UpdateTag(UpdateTagRequest) returns (UpdateTagResponse);     // -- 更新标签定义
  rpc ListTags(ListTagsRequest) returns (ListTagsResponse);       // Lists all available global tags -- 列出所有可用的全局标签定义

  // --- Merchant Tagging --- -- 商家标签管理接口
  // Adds a tag to a specific merchant. -- 为指定商家添加标签。
  rpc AddTagToMerchant(AddTagToMerchantRequest) returns (common.Empty);
  // Removes a tag from a specific merchant. -- 从指定商家移除标签。
  rpc RemoveTagFromMerchant(RemoveTagFromMerchantRequest) returns (common.Empty);
  // Lists tags for a specific merchant. -- 列出指定商家的标签。
  rpc ListMerchantTags(ListMerchantTagsRequest) returns (ListMerchantTagsResponse);
  // Lists merchants that have ALL specified tags. -- 列出拥有所有指定标签的商家。
  rpc ListMerchantsByTags(ListMerchantsByTagsRequest) returns (ListMerchantsByTagsResponse);
}

// --- Merchant CRUD Messages ---
message CreateMerchantRequest { // -- 创建商家请求
  string name = 1;                  // -- 商家名称
  string source_type = 2;           // -- 来源类型
  google.protobuf.Value contact_info = 3; // -- 联系信息 (JSON)
  string status = 4;                // -- 状态
  common.AuditInfo audit_info = 5;  // -- 审计信息
}

message CreateMerchantResponse { // -- 创建商家响应
  Merchant merchant = 1; // -- 创建的商家信息
}

message GetMerchantRequest { // -- 获取商家请求
  string merchant_id = 1; // -- 商家ID
}

message GetMerchantResponse { // -- 获取商家响应
  Merchant merchant = 1;   // -- 商家信息
  repeated Tag tags = 2;   // Tags associated with this merchant -- 与此商家关联的标签列表
}

message UpdateMerchantRequest { // -- 更新商家请求
  string merchant_id = 1;           // -- 商家ID
  optional string name = 2;         // -- 商家名称 (可选)
  optional string source_type = 3;  // -- 来源类型 (可选)
  optional google.protobuf.Value contact_info = 4; // -- 联系信息 (可选, JSON)
  optional string status = 5;       // -- 状态 (可选)
  common.AuditInfo audit_info = 6;  // -- 审计信息
}

message ListMerchantsFilter { // -- 列出商家筛选条件
  string name_contains = 1; // -- 名称包含的文本
  string status = 2;        // -- 状态
  // Potentially filter by source_type -- 可能按来源类型筛选
}

message ListMerchantsRequest { // -- 列出商家请求
  ListMerchantsFilter filter = 1;          // -- 筛选条件
  common.PaginationRequest pagination = 2; // -- 分页请求
}

message ListMerchantsResponse { // -- 列出商家响应
  repeated Merchant merchants = 1;        // -- 商家列表
  common.PaginationResponse pagination = 2; // -- 分页响应
}

// --- Tag Definition CRUD Messages ---
message CreateTagRequest { // -- 创建标签定义请求
  string name = 1;              // -- 标签名称
  common.TagType type = 2;      // -- 标签类型
  string description = 3;       // -- 描述
  common.AuditInfo audit_info = 4; // -- 审计信息
}

message CreateTagResponse { // -- 创建标签定义响应
  Tag tag = 1; // -- 创建的标签定义信息
}

message GetTagRequest { // -- 获取标签定义请求
  string tag_id = 1; // -- 标签ID
}

message GetTagResponse { // -- 获取标签定义响应
  Tag tag = 1; // -- 标签定义信息
}

message UpdateTagRequest { // -- 更新标签定义请求
  string tag_id = 1;                // -- 标签ID
  optional string name = 2;         // -- 标签名称 (可选)
  optional string description = 3;  // -- 描述 (可选)
  // Tag type typically should not be updated once created, or with strict rules -- 标签类型一旦创建通常不应更新，或遵循严格规则
  common.AuditInfo audit_info = 4;  // -- 审计信息
}

message ListTagsFilter { // -- 列出标签定义筛选条件
  common.TagType type = 1;    // -- 标签类型
  string name_contains = 2;   // -- 名称包含的文本
}

message ListTagsRequest { // -- 列出标签定义请求
  ListTagsFilter filter = 1;              // -- 筛选条件
  common.PaginationRequest pagination = 2; // -- 分页请求
}

message ListTagsResponse { // -- 列出标签定义响应
  repeated Tag tags = 1;                  // -- 标签定义列表
  common.PaginationResponse pagination = 2; // -- 分页响应
}

// --- Merchant Tagging Messages ---
message AddTagToMerchantRequest { // -- 添加标签到商家请求
  string merchant_id = 1;          // -- 商家ID
  string tag_id = 2;               // -- 标签ID
  common.AuditInfo audit_info = 3; // -- 审计信息
}

message RemoveTagFromMerchantRequest { // -- 从商家移除标签请求
  string merchant_id = 1;          // -- 商家ID
  string tag_id = 2;               // -- 标签ID
  common.AuditInfo audit_info = 3; // -- 审计信息
}

message ListMerchantTagsRequest { // -- 列出商家标签请求
  string merchant_id = 1; // -- 商家ID
}

message ListMerchantTagsResponse { // -- 列出商家标签响应
  repeated Tag tags = 1; // -- 商家的标签列表
}

message ListMerchantsByTagsRequest { // -- 根据标签列出商家请求
  repeated string tag_ids = 1;             // Find merchants that have all these tags -- 标签ID列表 (查找拥有所有这些标签的商家)
  common.PaginationRequest pagination = 2; // -- 分页请求
}

message ListMerchantsByTagsResponse { // -- 根据标签列出商家响应
  repeated string merchant_ids = 1;       // Returns only IDs for ProductService to use -- 返回商家ID列表 (供ProductService使用)
  common.PaginationResponse pagination = 2; // -- 分页响应
}
```

---

## 4. Listing Service (or Platform Listing Service) (`listing_service.proto`)

*Purpose: Manages the publishing of products to various external platforms. -- 上架服务：管理商品到各个外部平台的发布。*

```protobuf
syntax = "proto3";

package listing.v1;

import "common.proto";
import "google/protobuf/struct.proto";

option go_package = "github.com/yourorg/project/gen/go/listing/v1;listingpb";

message Listing { // -- 商品上架信息
  string listing_id = 1;    // -- 上架ID
  string product_id = 2;    // -- 商品ID
  string platform_id = 3;   // e.g., "xianyu", "dewu" -- 平台ID (例如 "xianyu", "dewu")
  string platform_product_id = 4; // ID on the external platform -- 在外部平台上的商品ID
  string status = 5;        // e.g., "PENDING_PUBLISH", "PUBLISHED", "UNPUBLISHED" -- 状态 (例如 "待发布", "已发布", "未发布")
  common.Timestamp published_at = 6;   // -- 发布时间
  common.Timestamp unpublished_at = 7; // -- 下架时间
  google.protobuf.Value publish_details = 8; // Platform-specific details -- 平台特定发布详情 (JSON)
  common.Timestamp created_at = 9;     // -- 创建时间
  common.Timestamp updated_at = 10;    // -- 更新时间
}

// ListingService manages product listings on external platforms. -- ListingService管理外部平台上的商品上架信息。
service ListingService {
  // Publishes a product to a specified platform. -- 将商品发布到指定平台。
  rpc PublishProduct(PublishProductRequest) returns (PublishProductResponse);
  // Unpublishes (delists) a product from a specified platform. -- 从指定平台下架商品。
  rpc UnpublishProduct(UnpublishProductRequest) returns (UnpublishProductResponse);
  // Notifies the listing service that a product has been sold (e.g., via an order)
  // so it can unlist from other platforms if necessary. -- 通知上架服务商品已售出 (例如通过订单)，以便必要时从其他平台下架。
  rpc NotifyProductSold(NotifyProductSoldRequest) returns (common.Empty);
  // Gets the listing status of a product on a specific platform. -- 获取商品在特定平台的上架状态。
  rpc GetListingStatus(GetListingStatusRequest) returns (GetListingStatusResponse);
  // Lists all listings for a given product. -- 列出给定商品的所有上架信息。
  rpc ListProductListings(ListProductListingsRequest) returns (ListProductListingsResponse);
}

message PublishProductRequest { // -- 发布商品请求
  string product_id = 1;    // -- 商品ID
  string platform_id = 2;   // -- 平台ID
  google.protobuf.Value platform_specific_options = 3; // e.g., category mapping, shipping templates -- 平台特定选项 (JSON, 例如类目映射、运费模板)
  common.AuditInfo audit_info = 4; // -- 审计信息
}

message PublishProductResponse { // -- 发布商品响应
  Listing listing = 1; // -- 上架信息
  string message = 2;  // e.g., "Successfully submitted for publishing" -- 消息 (例如 "已成功提交发布请求")
}

message UnpublishProductRequest { // -- 下架商品请求
  string product_id = 1;  // -- 商品ID
  string platform_id = 2; // -- 平台ID
  string reason = 3;      // Optional reason for unpublishing -- 下架原因 (可选)
  common.AuditInfo audit_info = 4; // -- 审计信息
}

message UnpublishProductResponse { // -- 下架商品响应
  string listing_id = 1;             // -- 上架ID
  string status_after_unpublish = 2; // e.g., "UNPUBLISHED" -- 下架后状态 (例如 "UNPUBLISHED")
  string message = 3;                // -- 消息
}

message NotifyProductSoldRequest { // -- 通知商品已售出请求
  string product_id = 1;          // -- 商品ID
  string order_id = 2;            // -- 订单ID
  string platform_where_sold = 3; // The platform on which the item was sold -- 商品售出所在的平台
  common.AuditInfo audit_info = 4; // -- 审计信息
}

message GetListingStatusRequest { // -- 获取上架状态请求
  string product_id = 1;  // -- 商品ID
  string platform_id = 2; // -- 平台ID
}

message GetListingStatusResponse { // -- 获取上架状态响应
  Listing listing = 1; // -- 上架信息
}

message ListProductListingsRequest { // -- 列出商品的上架信息请求
  string product_id = 1; // -- 商品ID
}

message ListProductListingsResponse { // -- 列出商品的上架信息响应
  repeated Listing listings = 1; // -- 上架信息列表
}
```

---

## 5. Order Service (`order_service.proto`)

*Purpose: Handles order creation, status management, and fulfillment coordination. -- 订单服务：处理订单创建、状态管理和履约协调。*
*(This is a separate microservice but its interface is defined here for completeness of the system design). -- (这是一个独立的微服务，但其接口在此定义以保证系统设计的完整性)。*

```protobuf
syntax = "proto3";

package order.v1;

import "common.proto";
import "google/protobuf/struct.proto";

option go_package = "github.com/yourorg/project/gen/go/order/v1;orderpb";

enum OrderStatus { // -- 订单状态枚举
  ORDER_STATUS_UNSPECIFIED = 0;      // -- 未指定状态
  ORDER_STATUS_PENDING_LOCK = 1;     // -- 待锁定商品
  ORDER_STATUS_LOCK_FAILED = 2;      // -- 商品锁定失败
  ORDER_STATUS_PENDING_PAYMENT = 3;  // Product locked, awaiting payment -- 商品已锁定，等待支付
  ORDER_STATUS_PAID = 4;             // -- 已支付
  ORDER_STATUS_PAYMENT_FAILED = 5;   // -- 支付失败
  ORDER_STATUS_SHIPPED = 6;          // -- 已发货
  ORDER_STATUS_DELIVERED = 7;        // -- 已送达
  ORDER_STATUS_CANCELLED = 8;        // -- 已取消
  ORDER_STATUS_REFUNDED = 9;         // -- 已退款
}

message Order { // -- 订单信息
  string order_id = 1;          // -- 订单ID
  string external_order_id = 2; // Platform's order ID -- 外部平台订单ID
  string product_id = 3;        // -- 商品ID
  string buyer_id = 4;          // -- 买家ID
  string merchant_id = 5;       // Denormalized from product for easier access -- 商家ID (从商品信息中冗余，便于查询)
  string platform_id = 6;       // -- 平台ID
  OrderStatus status = 7;       // -- 订单状态
  common.Price order_price = 8; // Price at the time of order -- 订单价格 (下单时价格)
  google.protobuf.Value shipping_address = 9; // JSON -- 收货地址 (JSON)
  google.protobuf.Value payment_details = 10; // JSON, e.g., transaction ID -- 支付详情 (JSON, 例如交易ID)
  string cancellation_reason = 11; // -- 取消原因
  common.Timestamp created_at = 12;    // -- 创建时间
  common.Timestamp paid_at = 13;       // -- 支付时间
  common.Timestamp updated_at = 14;    // -- 更新时间
  int32 version = 15;             // -- 版本号 (用于乐观锁)
}

// OrderService manages the lifecycle of orders. -- OrderService管理订单的生命周期。
service OrderService {
  // Creates an order and attempts to lock the product. -- 创建订单并尝试锁定商品。
  rpc PlaceOrder(PlaceOrderRequest) returns (PlaceOrderResponse);
  // Cancels an order. May involve unlocking the product. -- 取消订单 (可能涉及解锁商品)。
  rpc CancelOrder(CancelOrderRequest) returns (CancelOrderResponse);
  // Gets the status of a specific order. -- 获取特定订单的状态。
  rpc GetOrderStatus(GetOrderStatusRequest) returns (GetOrderStatusResponse);
  // Lists orders, possibly filtered. -- 列出订单 (可筛选)。
  rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse);

  // Callback for payment gateway or internal payment confirmation -- 支付网关回调或内部支付确认
  rpc ConfirmPayment(ConfirmPaymentRequest) returns (ConfirmPaymentResponse);
  // Marks an order as shipped -- 标记订单为已发货
  rpc MarkOrderShipped(MarkOrderShippedRequest) returns (common.Empty);
  // Marks an order as delivered -- 标记订单为已送达
  rpc MarkOrderDelivered(MarkOrderDeliveredRequest) returns (common.Empty);
}

message OrderItemInput { // -- 订单项输入信息
  string product_id = 1;           // -- 商品ID
  common.Price price_at_purchase = 2; // Ensure price is fixed at order time -- 购买时价格 (确保价格在下单时固定)
  // Potentially quantity if applicable, though luxury goods are often single items -- 如果适用，则可能包含数量，但奢侈品通常是单件商品
}

message PlaceOrderRequest { // -- 下单请求
  string external_order_id = 1; // Optional, if order originates from an external platform -- 外部订单ID (可选, 如果订单源自外部平台)
  string buyer_id = 2;          // -- 买家ID
  string platform_id = 3;       // -- 平台ID
  OrderItemInput item = 4;      // -- 订单项
  google.protobuf.Value shipping_address = 5; // -- 收货地址 (JSON)
  common.AuditInfo audit_info = 6; // user_id might be platform system or actual user -- 审计信息 (user_id可能是平台系统用户或真实用户)
}

message PlaceOrderResponse { // -- 下单响应
  Order order = 1;                 // -- 订单信息
  bool product_lock_acquired = 2;  // -- 商品是否成功锁定
  string message = 3;              // e.g., "Order created, awaiting payment" or "Failed to lock product" -- 消息 (例如 "订单已创建，等待支付" 或 "商品锁定失败")
}

message CancelOrderRequest { // -- 取消订单请求
  string order_id = 1;            // -- 订单ID
  string reason = 2;              // -- 取消原因
  common.AuditInfo audit_info = 3; // -- 审计信息
}

message CancelOrderResponse { // -- 取消订单响应
  Order order = 1;   // Updated order -- 更新后的订单信息
  string message = 2; // -- 消息
}

message GetOrderStatusRequest { // -- 获取订单状态请求
  string order_id = 1; // -- 订单ID
}

message GetOrderStatusResponse { // -- 获取订单状态响应
  Order order = 1; // -- 订单信息
}

message ListOrdersFilter { // -- 列出订单筛选条件
  string buyer_id = 1;    // -- 买家ID
  string product_id = 2;  // -- 商品ID
  string merchant_id = 3; // -- 商家ID
  string platform_id = 4; // -- 平台ID
  OrderStatus status = 5; // -- 订单状态
}

message ListOrdersRequest { // -- 列出订单请求
  ListOrdersFilter filter = 1;            // -- 筛选条件
  common.PaginationRequest pagination = 2; // -- 分页请求
}

message ListOrdersResponse { // -- 列出订单响应
  repeated Order orders = 1;              // -- 订单列表
  common.PaginationResponse pagination = 2; // -- 分页响应
}

message ConfirmPaymentRequest { // -- 确认支付请求
  string order_id = 1;               // -- 订单ID
  string payment_transaction_id = 2; // -- 支付交易ID
  string payment_method = 3;         // -- 支付方式
  common.Timestamp paid_at = 4;      // -- 支付时间
  common.AuditInfo audit_info = 5;   // -- 审计信息
}

message ConfirmPaymentResponse { // -- 确认支付响应
  Order order = 1;   // -- 订单信息
  string message = 2; // -- 消息
}

message MarkOrderShippedRequest { // -- 标记订单已发货请求
  string order_id = 1;        // -- 订单ID
  string tracking_number = 2; // -- 物流追踪号
  string carrier = 3;         // -- 物流公司
  common.AuditInfo audit_info = 4; // -- 审计信息
}

message MarkOrderDeliveredRequest { // -- 标记订单已送达请求
  string order_id = 1;               // -- 订单ID
  common.Timestamp delivered_at = 2; // -- 送达时间
  common.AuditInfo audit_info = 3;   // -- 审计信息
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
