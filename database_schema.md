# Database Schema Document

This document outlines the database schema for the Product Service and Order Service.

## Product Service Database

### 1. Table: `products` -- 商品表

| Column           | Type                | Constraints                                                              | Description / Notes                                                                                                |
|------------------|---------------------|--------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------|
| `product_id`     | BIGINT UNSIGNED     | PRIMARY KEY, AUTO_INCREMENT                                              | Unique identifier for the product. -- 商品ID，主键，自增。                                                               |
| `erp_product_id` | VARCHAR(255)        | UNIQUE, NULLABLE                                                         | ID from the ERP system, used for synchronization. Nullable if the product is created directly in the system. -- ERP系统中的商品ID，用于同步，如果商品直接在系统中创建则可为空。 |
| `sku_ref_id`     | VARCHAR(255)        | NULLABLE                                                                 | Reference ID for a standard product unit (e.g., from a catalog), if applicable. -- 标品ID，关联标准产品单元（例如来自目录）。      |
| `merchant_id`    | BIGINT UNSIGNED     | NOT NULL, INDEX                                                          | Foreign key referencing `merchants.merchant_id`. -- 商家ID，外键，关联merchants表。                                        |
| `title`          | VARCHAR(512)        | NOT NULL                                                                 | Product title. -- 商品标题。                                                                                       |
| `description`    | TEXT                | NULLABLE                                                                 | Detailed product description. -- 商品详细描述。                                                                    |
| `images`         | JSON                | NULLABLE                                                                 | JSON array of image URLs or identifiers. Storing as JSON allows for flexibility in image metadata. Example: `[{"url": "cdn.example.com/img1.jpg", "order": 1, "alt_text": "Front view"}, ...]`. -- 商品图片，JSON数组格式，存储图片URL或标识符。 |
| `price_currency` | VARCHAR(3)          | NOT NULL, DEFAULT 'CNY'                                                  | Currency code (e.g., CNY, USD). -- 价格货币代码（如CNY, USD）。                                                       |
| `price_amount`   | DECIMAL(12, 2)      | NOT NULL                                                                 | Product price. -- 商品价格金额。                                                                                     |
| `status`         | ENUM(...)           | NOT NULL, DEFAULT 'draft', INDEX                                         | ('draft', 'pending_approval', 'listed', 'locked', 'sold', 'delisted', 'archived'). -- 商品状态（draft: 草稿, pending_approval: 待审核, listed: 已上架, locked: 已锁定, sold: 已售出, delisted: 已下架, archived: 已归档）。 |
| `ai_enriched`    | BOOLEAN             | NOT NULL, DEFAULT FALSE                                                  | Flag indicating if AI has processed the product. -- AI是否已处理标记，布尔值。                                               |
| `created_at`     | TIMESTAMP           | NOT NULL, DEFAULT CURRENT_TIMESTAMP                                      | Timestamp of creation. -- 创建时间。                                                                                 |
| `updated_at`     | TIMESTAMP           | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP          | Timestamp of last update. -- 更新时间。                                                                              |
| `version`        | INT UNSIGNED        | NOT NULL, DEFAULT 1                                                      | Optimistic locking version. -- 乐观锁版本号。                                                                        |

**Indexes:**
*   PRIMARY KEY (`product_id`)
*   UNIQUE KEY `idx_erp_product_id` (`erp_product_id`)
*   INDEX `idx_merchant_id` (`merchant_id`)
*   INDEX `idx_status` (`status`)
*   INDEX `idx_created_at` (`created_at`)
*   INDEX `idx_sku_ref_id` (`sku_ref_id`)

---

### 2. Table: `product_ai_tags` -- 商品AI标签表

| Column               | Type             | Constraints                               | Description / Notes                                                                                                   |
|----------------------|------------------|-------------------------------------------|-----------------------------------------------------------------------------------------------------------------------|
| `product_ai_tag_id`  | BIGINT UNSIGNED  | PRIMARY KEY, AUTO_INCREMENT               | Unique identifier. -- AI标签记录ID，主键，自增。                                                                          |
| `product_id`         | BIGINT UNSIGNED  | NOT NULL                                  | Foreign key referencing `products.product_id`. -- 商品ID，外键，关联products表。                                          |
| `tags`               | JSON             | NOT NULL                                  | JSON object or array of tags generated by AI. -- AI生成的标签，JSON格式。                                                     |
| `ai_service_version` | VARCHAR(50)      | NULLABLE                                  | Version of the AI service/model that generated the tags. -- AI服务/模型版本号。                                               |
| `confidence_score`   | DECIMAL(5,4)     | NULLABLE                                  | Overall confidence score for the generated tags, if available from AI. -- AI生成标签的置信度评分。                                |
| `created_at`         | TIMESTAMP        | NOT NULL, DEFAULT CURRENT_TIMESTAMP       | Timestamp of creation. -- 创建时间。                                                                                    |

**Indexes:**
*   PRIMARY KEY (`product_ai_tag_id`)
*   INDEX `idx_product_id_ai_tags` (`product_id`)
*   INDEX `idx_created_at_ai_tags` (`created_at`)
*   *Consider GIN/Full-text index on `tags` if querying specific AI tags is frequent.*

---

### 3. Table: `product_ai_descriptions` -- 商品AI描述表

| Column                      | Type             | Constraints                               | Description / Notes                                                                 |
|-----------------------------|------------------|-------------------------------------------|-------------------------------------------------------------------------------------|
| `product_ai_description_id` | BIGINT UNSIGNED  | PRIMARY KEY, AUTO_INCREMENT               | Unique identifier. -- AI描述记录ID，主键，自增。                                                      |
| `product_id`                | BIGINT UNSIGNED  | NOT NULL                                  | Foreign key referencing `products.product_id`. -- 商品ID，外键，关联products表。                          |
| `language_code`             | VARCHAR(10)      | NOT NULL                                  | Language code (e.g., "en-US", "ja-JP", "zh-CN"). -- 语言代码（如 "en-US", "ja-JP", "zh-CN"）。           |
| `description`               | TEXT             | NOT NULL                                  | AI-generated description in the specified language. -- AI生成的多语言描述。                               |
| `ai_service_version`        | VARCHAR(50)      | NULLABLE                                  | Version of the AI service/model that generated the description. -- AI服务/模型版本号。                    |
| `created_at`                | TIMESTAMP        | NOT NULL, DEFAULT CURRENT_TIMESTAMP       | Timestamp of creation. -- 创建时间。                                                                  |

**Constraints:**
*   UNIQUE KEY `idx_product_lang_ai_desc` (`product_id`, `language_code`)

**Indexes:**
*   PRIMARY KEY (`product_ai_description_id`)
*   INDEX `idx_product_id_ai_desc` (`product_id`)
*   INDEX `idx_language_code_ai_desc` (`language_code`)
*   INDEX `idx_created_at_ai_desc` (`created_at`)

---

### 4. Table: `merchants` -- 商家表

| Column         | Type                | Constraints                                                              | Description / Notes                                                                              |
|----------------|---------------------|--------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------|
| `merchant_id`  | BIGINT UNSIGNED     | PRIMARY KEY, AUTO_INCREMENT                                              | Unique identifier for the merchant. -- 商家ID，主键，自增。                                              |
| `name`         | VARCHAR(255)        | NOT NULL                                                                 | Merchant's name. -- 商家名称。                                                                     |
| `source_type`  | VARCHAR(50)         | NULLABLE                                                                 | Origin or type of the merchant (e.g., "erp_sync", "manual_entry", "consignment"). -- 商家来源类型（如 "erp_sync", "manual_entry"）。 |
| `contact_info` | JSON                | NULLABLE                                                                 | JSON object for contact details (e.g., `{"phone": "123...", "email": "contact@example.com"}`). -- 联系信息，JSON格式。 |
| `status`       | ENUM(...)           | NOT NULL, DEFAULT 'active', INDEX                                        | ('active', 'inactive', 'pending_review'). -- 商家状态（active: 活跃, inactive: 非活跃, pending_review: 待审核）。 |
| `created_at`   | TIMESTAMP           | NOT NULL, DEFAULT CURRENT_TIMESTAMP                                      | Timestamp of creation. -- 创建时间。                                                               |
| `updated_at`   | TIMESTAMP           | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP          | Timestamp of last update. -- 更新时间。                                                            |

**Indexes:**
*   PRIMARY KEY (`merchant_id`)
*   INDEX `idx_merchant_name` (`name`)
*   INDEX `idx_merchant_status` (`status`)
*   INDEX `idx_merchant_source_type` (`source_type`)

---

### 5. Table: `tags` -- 标签定义表

| Column      | Type                | Constraints                                                              | Description / Notes                                                                           |
|-------------|---------------------|--------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------|
| `tag_id`    | BIGINT UNSIGNED     | PRIMARY KEY, AUTO_INCREMENT                                              | Unique identifier for the tag. -- 标签ID，主键，自增。                                              |
| `name`      | VARCHAR(100)        | NOT NULL, UNIQUE                                                         | Tag name (e.g., "Vintage", "Limited Edition", "Good Condition"). -- 标签名称，唯一。                  |
| `type`      | ENUM(...)           | NOT NULL, INDEX                                                          | ('manual', 'rule', 'ai', 'category'). Source/type of the tag. -- 标签类型（manual: 手动, rule: 规则, ai: AI生成, category: 分类）。 |
| `description` | VARCHAR(255)        | NULLABLE                                                                 | Optional description for the tag's meaning or usage. -- 标签描述信息。                                |
| `created_at`| TIMESTAMP           | NOT NULL, DEFAULT CURRENT_TIMESTAMP                                      | Timestamp of creation. -- 创建时间。                                                            |
| `updated_at`| TIMESTAMP           | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP          | Timestamp of last update. -- 更新时间。                                                         |

**Indexes:**
*   PRIMARY KEY (`tag_id`)
*   UNIQUE KEY `idx_tag_name` (`name`)
*   INDEX `idx_tag_type` (`type`)

---

### 6. Table: `product_tags` -- 商品人工标签关联表

| Column             | Type             | Constraints                               | Description / Notes                                                                          |
|--------------------|------------------|-------------------------------------------|----------------------------------------------------------------------------------------------|
| `product_tag_id`   | BIGINT UNSIGNED  | PRIMARY KEY, AUTO_INCREMENT               | Unique identifier for the association. -- 商品标签关联ID，主键，自增。                                   |
| `product_id`       | BIGINT UNSIGNED  | NOT NULL                                  | Foreign key referencing `products.product_id`. -- 商品ID，外键，关联products表。                       |
| `tag_id`           | BIGINT UNSIGNED  | NOT NULL                                  | Foreign key referencing `tags.tag_id`. -- 标签ID，外键，关联tags表。                                 |
| `added_by_user_id` | BIGINT UNSIGNED  | NULLABLE                                  | User ID who added this tag. -- 添加该标签的用户ID。                                                  |
| `created_at`       | TIMESTAMP        | NOT NULL, DEFAULT CURRENT_TIMESTAMP       | Timestamp of association. -- 创建时间。                                                        |

**Constraints:**
*   UNIQUE KEY `idx_product_tag_unique` (`product_id`, `tag_id`)

**Indexes:**
*   PRIMARY KEY (`product_tag_id`)
*   INDEX `idx_pt_product_id` (`product_id`)
*   INDEX `idx_pt_tag_id` (`tag_id`)

---

### 7. Table: `merchant_tags` -- 商家标签关联表

| Column             | Type             | Constraints                               | Description / Notes                                                                          |
|--------------------|------------------|-------------------------------------------|----------------------------------------------------------------------------------------------|
| `merchant_tag_id`  | BIGINT UNSIGNED  | PRIMARY KEY, AUTO_INCREMENT               | Unique identifier for the association. -- 商家标签关联ID，主键，自增。                                   |
| `merchant_id`      | BIGINT UNSIGNED  | NOT NULL                                  | Foreign key referencing `merchants.merchant_id`. -- 商家ID，外键，关联merchants表。                    |
| `tag_id`           | BIGINT UNSIGNED  | NOT NULL                                  | Foreign key referencing `tags.tag_id`. -- 标签ID，外键，关联tags表。                                 |
| `added_by_user_id` | BIGINT UNSIGNED  | NULLABLE                                  | User ID who added this tag to the merchant. -- 添加该标签到商家的用户ID。                                |
| `created_at`       | TIMESTAMP        | NOT NULL, DEFAULT CURRENT_TIMESTAMP       | Timestamp of association. -- 创建时间。                                                        |

**Constraints:**
*   UNIQUE KEY `idx_merchant_tag_unique` (`merchant_id`, `tag_id`)

**Indexes:**
*   PRIMARY KEY (`merchant_tag_id`)
*   INDEX `idx_mt_merchant_id` (`merchant_id`)
*   INDEX `idx_mt_tag_id` (`tag_id`)

---

### 8. Table: `merchant_tag_logs` -- 商家标签操作日志表

| Column         | Type             | Constraints                               | Description / Notes                                                                                                |
|----------------|------------------|-------------------------------------------|--------------------------------------------------------------------------------------------------------------------|
| `log_id`       | BIGINT UNSIGNED  | PRIMARY KEY, AUTO_INCREMENT               | Unique log entry identifier. -- 日志ID，主键，自增。                                                                   |
| `merchant_id`  | BIGINT UNSIGNED  | NOT NULL                                  | Foreign key referencing `merchants.merchant_id`. -- 商家ID，外键，关联merchants表。                                     |
| `tag_id`       | BIGINT UNSIGNED  | NOT NULL                                  | Foreign key referencing `tags.tag_id`. -- 标签ID，外键，关联tags表。                                                     |
| `action`       | ENUM(...)        | NOT NULL                                  | ('added', 'removed'). The action performed on the tag. -- 操作类型（added: 添加, removed: 移除）。                       |
| `operator_id`  | BIGINT UNSIGNED  | NOT NULL                                  | ID of the user or system performing the action. -- 操作人ID（用户ID或系统ID）。                                           |
| `source`       | VARCHAR(255)     | NULLABLE                                  | More details about the source of the action. -- 操作来源详情。                                                          |
| `change_details` | JSON             | NULLABLE                                  | Optional JSON field to store additional details about the change. -- 变更详情，JSON格式。                                  |
| `created_at`   | TIMESTAMP        | NOT NULL, DEFAULT CURRENT_TIMESTAMP       | Timestamp of the log entry. -- 创建时间。                                                                            |

**Indexes:**
*   PRIMARY KEY (`log_id`)
*   INDEX `idx_mtl_merchant_id` (`merchant_id`)
*   INDEX `idx_mtl_tag_id` (`tag_id`)
*   INDEX `idx_mtl_operator_id` (`operator_id`)
*   INDEX `idx_mtl_created_at` (`created_at`)

---

### 9. Table: `product_listings` -- 商品上架信息表

| Column                     | Type             | Constraints                                                              | Description / Notes                                                                                           |
|----------------------------|------------------|--------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------|
| `listing_id`               | BIGINT UNSIGNED  | PRIMARY KEY, AUTO_INCREMENT                                              | Unique identifier for the listing. -- 上架记录ID，主键，自增。                                                        |
| `product_id`               | BIGINT UNSIGNED  | NOT NULL                                                                 | Foreign key referencing `products.product_id`. -- 商品ID，外键，关联products表。                                      |
| `platform_id`              | VARCHAR(50)      | NOT NULL                                                                 | Identifier for the external platform (e.g., "xianyu", "dewu", "95fen"). -- 外部平台标识（如 "xianyu", "dewu"）。        |
| `platform_product_id`      | VARCHAR(255)     | NULLABLE                                                                 | The ID of the product on the external platform, once listed. -- 商品在外部平台的ID。                                  |
| `status`                   | ENUM(...)        | NOT NULL, INDEX                                                          | ('pending_publish', 'published', 'failed_publish', 'pending_unpublish', 'unpublished', 'sold_on_platform'). -- 上架状态（pending_publish: 待发布, published: 已发布, failed_publish: 发布失败, pending_unpublish: 待下架, unpublished: 已下架, sold_on_platform: 已在平台售出）。 |
| `last_publish_attempt_at`  | TIMESTAMP        | NULLABLE                                                                 | Timestamp of the last attempt to publish. -- 最后尝试发布时间。                                                     |
| `published_at`             | TIMESTAMP        | NULLABLE                                                                 | Timestamp when successfully published. -- 成功发布时间。                                                          |
| `unpublished_at`           | TIMESTAMP        | NULLABLE                                                                 | Timestamp when unpublished. -- 下架时间。                                                                       |
| `last_sync_at`             | TIMESTAMP        | NULLABLE                                                                 | Timestamp of the last status sync with the platform. -- 最后与平台同步状态时间。                                      |
| `publish_details`          | JSON             | NULLABLE                                                                 | Any specific details related to publishing on this platform. -- 发布详情，JSON格式，存储平台特定信息。                           |
| `error_message`            | TEXT             | NULLABLE                                                                 | Stores error messages if publishing/unpublishing failed. -- 错误信息（如果发布/下架失败）。                               |
| `created_at`               | TIMESTAMP        | NOT NULL, DEFAULT CURRENT_TIMESTAMP                                      | Timestamp of creation. -- 创建时间。                                                                            |
| `updated_at`               | TIMESTAMP        | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP          | Timestamp of last update. -- 更新时间。                                                                         |

**Constraints:**
*   UNIQUE KEY `idx_product_platform_unique` (`product_id`, `platform_id`)

**Indexes:**
*   PRIMARY KEY (`listing_id`)
*   INDEX `idx_pl_product_id` (`product_id`)
*   INDEX `idx_pl_platform_id` (`platform_id`)
*   INDEX `idx_pl_status` (`status`)
*   INDEX `idx_pl_platform_product_id` (`platform_product_id`)

---

### 10. Table: `audit_logs` -- 审计日志表

| Column        | Type             | Constraints                               | Description / Notes                                                                                                                               |
|---------------|------------------|-------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------|
| `log_id`      | BIGINT UNSIGNED  | PRIMARY KEY, AUTO_INCREMENT               | Unique log entry identifier. -- 日志ID，主键，自增。                                                                                                    |
| `user_id`     | BIGINT UNSIGNED  | NOT NULL, INDEX                           | ID of the user or system performing the action. -- 操作用户ID（实际用户ID或AI等系统保留ID）。                                                              |
| `action_type` | VARCHAR(100)     | NOT NULL, INDEX                           | Type of action. -- 操作类型（如 "CREATE_PRODUCT", "AI_GENERATE_TAGS"）。                                                                               |
| `target_type` | VARCHAR(50)      | NOT NULL, INDEX                           | The type of entity being acted upon. -- 操作目标实体类型（如 "Product", "Merchant"）。                                                                       |
| `target_id`   | VARCHAR(255)     | NOT NULL, INDEX                           | The ID of the entity being acted upon. -- 操作目标实体ID。                                                                                              |
| `details`     | JSON             | NULLABLE                                  | JSON object containing details of the action. -- 操作详情，JSON格式，可包含操作前后数据对比。                                                                      |
| `ip_address`  | VARCHAR(45)      | NULLABLE                                  | IP address of the requester, if applicable. -- 请求者IP地址。                                                                                           |
| `user_agent`  | TEXT             | NULLABLE                                  | User agent of the requester, if applicable. -- 请求者用户代理。                                                                                         |
| `created_at`  | TIMESTAMP        | NOT NULL, DEFAULT CURRENT_TIMESTAMP, INDEX| Timestamp of the log entry. -- 创建时间。                                                                                                           |

**Indexes:**
*   PRIMARY KEY (`log_id`)
*   INDEX `idx_al_user_id` (`user_id`)
*   INDEX `idx_al_action_type` (`action_type`)
*   INDEX `idx_al_target_type_id` (`target_type`, `target_id`)
*   INDEX `idx_al_created_at` (`created_at`)

---

## Order Service Database -- 订单服务数据库（独立微服务）

### 1. Table: `orders` -- 订单表

| Column              | Type                | Constraints                                                              | Description / Notes                                                                                                  |
|---------------------|---------------------|--------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------|
| `order_id`          | BIGINT UNSIGNED     | PRIMARY KEY, AUTO_INCREMENT                                              | Unique identifier for the order. -- 订单ID，主键，自增。                                                                 |
| `external_order_id` | VARCHAR(255)        | NULLABLE, UNIQUE                                                         | Order ID from the external platform. -- 外部平台订单ID，唯一。                                                               |
| `product_id`        | BIGINT UNSIGNED     | NOT NULL, INDEX                                                          | The ID of the product being ordered. -- 商品ID，关联商品服务中的商品。                                                           |
| `buyer_id`          | VARCHAR(255)        | NOT NULL, INDEX                                                          | Identifier for the buyer. -- 买家ID。                                                                                  |
| `merchant_id`       | BIGINT UNSIGNED     | NOT NULL, INDEX                                                          | The ID of the merchant selling the product. -- 商家ID，冗余字段，方便查询。                                                      |
| `platform_id`       | VARCHAR(50)         | NOT NULL, INDEX                                                          | Identifier for the platform where the order originated. -- 订单来源平台ID。                                                  |
| `status`            | ENUM(...)           | NOT NULL, INDEX                                                          | ('pending_lock', 'lock_failed', 'pending_payment', 'paid', 'payment_failed', 'shipped', 'delivered', 'cancelled', 'refunded'). -- 订单状态（pending_lock: 待锁定, lock_failed: 锁定失败, pending_payment: 待支付, paid: 已支付, payment_failed: 支付失败, shipped: 已发货, delivered: 已送达, cancelled: 已取消, refunded: 已退款）。 |
| `currency`          | VARCHAR(3)          | NOT NULL                                                                 | Currency of the order. -- 货币代码。                                                                                   |
| `price_amount`      | DECIMAL(12, 2)      | NOT NULL                                                                 | Price at the time of order. -- 订单金额（下单时价格）。                                                                    |
| `shipping_address`  | JSON                | NULLABLE                                                                 | Buyer's shipping address. -- 收货地址，JSON格式。                                                                      |
| `payment_details`   | JSON                | NULLABLE                                                                 | Details about the payment. -- 支付详情，JSON格式。                                                                       |
| `cancellation_reason`| TEXT                | NULLABLE                                                                 | Reason for cancellation, if applicable. -- 取消原因。                                                                  |
| `created_at`        | TIMESTAMP           | NOT NULL, DEFAULT CURRENT_TIMESTAMP                                      | Timestamp of order creation. -- 创建时间。                                                                             |
| `paid_at`           | TIMESTAMP           | NULLABLE                                                                 | Timestamp when payment was confirmed. -- 支付时间。                                                                    |
| `updated_at`        | TIMESTAMP           | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP          | Timestamp of last update. -- 更新时间。                                                                                |
| `version`           | INT UNSIGNED        | NOT NULL, DEFAULT 1                                                      | Optimistic locking version. -- 乐观锁版本号。                                                                            |

**Indexes:**
*   PRIMARY KEY (`order_id`)
*   UNIQUE KEY `idx_external_order_id` (`external_order_id`)
*   INDEX `idx_ord_product_id` (`product_id`)
*   INDEX `idx_ord_buyer_id` (`buyer_id`)
*   INDEX `idx_ord_merchant_id` (`merchant_id`)
*   INDEX `idx_ord_platform_id` (`platform_id`)
*   INDEX `idx_ord_status` (`status`)
*   INDEX `idx_ord_created_at` (`created_at`)

---

### 2. Table: `order_status_history` -- 订单状态历史表

| Column       | Type             | Constraints                               | Description / Notes                                                                                 |
|--------------|------------------|-------------------------------------------|-----------------------------------------------------------------------------------------------------|
| `history_id` | BIGINT UNSIGNED  | PRIMARY KEY, AUTO_INCREMENT               | Unique identifier for the history entry. -- 历史记录ID，主键，自增。                                          |
| `order_id`   | BIGINT UNSIGNED  | NOT NULL                                  | Foreign key referencing `orders.order_id`. -- 订单ID，外键，关联orders表。                                  |
| `status_from`| VARCHAR(50)      | NULLABLE                                  | Previous status. -- 原状态。                                                                          |
| `status_to`  | VARCHAR(50)      | NOT NULL                                  | New status. -- 新状态。                                                                             |
| `changed_by` | VARCHAR(255)     | NULLABLE                                  | Who/what initiated the change. -- 状态变更发起者。                                                        |
| `notes`      | TEXT             | NULLABLE                                  | Additional notes about the status change. -- 备注信息。                                                 |
| `created_at` | TIMESTAMP        | NOT NULL, DEFAULT CURRENT_TIMESTAMP       | Timestamp of the status change. -- 创建时间。                                                         |

**Indexes:**
*   PRIMARY KEY (`history_id`)
*   INDEX `idx_osh_order_id_created_at` (`order_id`, `created_at` DESC)

---

## Relationships between tables

### Product Service:
*   `products.merchant_id` -> `merchants.merchant_id` (Many-to-One)
*   `product_ai_tags.product_id` -> `products.product_id` (Many-to-One)
*   `product_ai_descriptions.product_id` -> `products.product_id` (One-to-Many)
*   `product_tags.product_id` -> `products.product_id` (Many-to-Many junction)
*   `product_tags.tag_id` -> `tags.tag_id` (Many-to-Many junction)
*   `merchant_tags.merchant_id` -> `merchants.merchant_id` (Many-to-Many junction)
*   `merchant_tags.tag_id` -> `tags.tag_id` (Many-to-Many junction)
*   `merchant_tag_logs.merchant_id` -> `merchants.merchant_id`
*   `merchant_tag_logs.tag_id` -> `tags.tag_id`
*   `product_listings.product_id` -> `products.product_id` (One-to-Many)
*   `audit_logs.target_id` conceptually links to various primary keys based on `target_type` (No enforced FK).

### Order Service:
*   `order_status_history.order_id` -> `orders.order_id` (One-to-Many)
*   `orders.product_id` logically refers to `products.product_id` in Product Service (No enforced FK across microservices).
*   `orders.merchant_id` logically refers to `merchants.merchant_id` in Product Service (No enforced FK, denormalized).

---

## Rationale for Design Choices

*   **JSON Fields (`products.images`, `product_ai_tags.tags`, `merchants.contact_info`, `audit_logs.details`, `product_listings.publish_details`, `orders.shipping_address`, `orders.payment_details`):**
    *   Provides flexibility for semi-structured data that might evolve without requiring schema migrations (e.g., adding new contact fields, different image attributes).
    *   Databases like PostgreSQL and MySQL (newer versions) offer good JSON indexing and querying capabilities if needed.
    *   For `tags` in `product_ai_tags`, JSON allows storing key-value pairs or arrays, which can be useful if AI provides structured tag information (e.g., `{"brand": "Rolex", "model": "Submariner"}`).
*   **ENUM Types:** Used for fields with a fixed set of possible values (`status` fields, `tags.type`, `merchant_tag_logs.action`). This improves data integrity and can be more efficient than VARCHARs.
*   **`created_at` and `updated_at` Timestamps:** Standard practice for tracking record creation and modifications. `ON UPDATE CURRENT_TIMESTAMP` for `updated_at` automates this for MySQL.
*   **`version` column (Optimistic Locking):** Included in `products` and `orders` tables to help prevent lost updates in concurrent environments. The application layer would check and increment this version during updates.
*   **Separate AI Data Tables (`product_ai_tags`, `product_ai_descriptions`):**
    *   Keeps the main `products` table cleaner.
    *   Allows for potentially multiple versions of AI data (e.g., if AI reruns or different models are used) by looking at `created_at`.
    *   Can be scaled or managed independently.
*   **Separate `product_tags` for manual/curated tags:** This clearly distinguishes AI-generated tags from human-curated or rule-based tags, which might have different trust levels or uses.
*   **`merchant_tag_logs` vs `audit_logs`:**
    *   `merchant_tag_logs` is specific to changes in merchant-tag associations, providing a focused history for this critical relationship. It's directly queryable for merchant tag history.
    *   `audit_logs` is a more generic, system-wide log for various significant actions, including those that might also be in `merchant_tag_logs` but with broader context. This duplication is acceptable for different querying needs (specific vs. system-wide audit).
*   **`product_listings` Table:** Decouples the core product information from its state on various external sales platforms. This is crucial for managing multi-platform sales.
*   **`target_id` in `audit_logs` as VARCHAR:** While many IDs are BIGINT, using VARCHAR provides flexibility if some target entities use UUIDs or other string-based identifiers in the future.
*   **Denormalization in `orders` table (`merchant_id`, `price_amount`, `currency`):**
    *   `merchant_id` is copied to avoid joins to Product Service for common order queries.
    *   `price_amount` and `currency` are captured at the time of order to ensure historical accuracy, as product prices might change later.
*   **Indexing Strategy:**
    *   Primary keys are inherently indexed.
    *   Foreign keys are indexed to improve join performance.
    *   Fields frequently used in WHERE clauses (e.g., `status`, `type`, `created_at`, `platform_id`) are indexed.
    *   Unique constraints are created where necessary (e.g., `tags.name`, `product_listings` on `product_id` + `platform_id`).
    *   Consider composite indexes for queries involving multiple columns (e.g., `audit_logs` on `target_type`, `target_id`).
    *   For JSON fields, if specific keys within the JSON are frequently queried, consider generated columns and indexing them (PostgreSQL GIN indexes are also powerful for JSON).

This completes the database design part.
