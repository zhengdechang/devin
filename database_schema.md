# Database Schema Document

This document outlines the database schema for the Product Service and Order Service.

## Product Service Database

### 1. Table: `products`

| Column           | Type                | Constraints                                                              | Description / Notes                                                                                                |
|------------------|---------------------|--------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------|
| `product_id`     | BIGINT UNSIGNED     | PRIMARY KEY, AUTO_INCREMENT                                              | Unique identifier for the product.                                                                                 |
| `erp_product_id` | VARCHAR(255)        | UNIQUE, NULLABLE                                                         | ID from the ERP system, used for synchronization. Nullable if product created directly.                            |
| `sku_ref_id`     | VARCHAR(255)        | NULLABLE                                                                 | Reference ID for a standard product unit (e.g., from a catalog), if applicable.                                    |
| `merchant_id`    | BIGINT UNSIGNED     | NOT NULL, INDEX                                                          | Foreign key referencing `merchants.merchant_id`.                                                                   |
| `title`          | VARCHAR(512)        | NOT NULL                                                                 | Product title.                                                                                                     |
| `description`    | TEXT                | NULLABLE                                                                 | Detailed product description.                                                                                      |
| `images`         | JSON                | NULLABLE                                                                 | JSON array of image URLs/identifiers. Ex: `[{"url": "cdn.example.com/img1.jpg", "order": 1, "alt_text": "Front view"}]` |
| `price_currency` | VARCHAR(3)          | NOT NULL, DEFAULT 'CNY'                                                  | Currency code (e.g., CNY, USD).                                                                                    |
| `price_amount`   | DECIMAL(12, 2)      | NOT NULL                                                                 | Product price.                                                                                                     |
| `status`         | ENUM(...)           | NOT NULL, DEFAULT 'draft', INDEX                                         | ('draft', 'pending_approval', 'listed', 'locked', 'sold', 'delisted', 'archived').                                 |
| `ai_enriched`    | BOOLEAN             | NOT NULL, DEFAULT FALSE                                                  | Flag indicating if AI has processed the product.                                                                   |
| `created_at`     | TIMESTAMP           | NOT NULL, DEFAULT CURRENT_TIMESTAMP                                      | Timestamp of creation.                                                                                             |
| `updated_at`     | TIMESTAMP           | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP          | Timestamp of last update.                                                                                          |
| `version`        | INT UNSIGNED        | NOT NULL, DEFAULT 1                                                      | Optimistic locking version.                                                                                        |

**Indexes:**
*   PRIMARY KEY (`product_id`)
*   UNIQUE KEY `idx_erp_product_id` (`erp_product_id`)
*   INDEX `idx_merchant_id` (`merchant_id`)
*   INDEX `idx_status` (`status`)
*   INDEX `idx_created_at` (`created_at`)
*   INDEX `idx_sku_ref_id` (`sku_ref_id`)

---

### 2. Table: `product_ai_tags`

| Column               | Type             | Constraints                               | Description / Notes                                                                                                   |
|----------------------|------------------|-------------------------------------------|-----------------------------------------------------------------------------------------------------------------------|
| `product_ai_tag_id`  | BIGINT UNSIGNED  | PRIMARY KEY, AUTO_INCREMENT               | Unique identifier.                                                                                                    |
| `product_id`         | BIGINT UNSIGNED  | NOT NULL                                  | Foreign key referencing `products.product_id`.                                                                        |
| `tags`               | JSON             | NOT NULL                                  | JSON object/array of AI tags. Ex: `{"brand": "Chanel", "category": "Handbag"}` or `["vintage", "leather"]`.         |
| `ai_service_version` | VARCHAR(50)      | NULLABLE                                  | Version of the AI service/model that generated the tags.                                                              |
| `confidence_score`   | DECIMAL(5,4)     | NULLABLE                                  | Overall confidence score for the generated tags.                                                                      |
| `created_at`         | TIMESTAMP        | NOT NULL, DEFAULT CURRENT_TIMESTAMP       | Timestamp of creation.                                                                                                |

**Indexes:**
*   PRIMARY KEY (`product_ai_tag_id`)
*   INDEX `idx_product_id_ai_tags` (`product_id`)
*   INDEX `idx_created_at_ai_tags` (`created_at`)
*   *Consider GIN/Full-text index on `tags` if querying specific AI tags is frequent.*

---

### 3. Table: `product_ai_descriptions`

| Column                      | Type             | Constraints                               | Description / Notes                                                                 |
|-----------------------------|------------------|-------------------------------------------|-------------------------------------------------------------------------------------|
| `product_ai_description_id` | BIGINT UNSIGNED  | PRIMARY KEY, AUTO_INCREMENT               | Unique identifier.                                                                  |
| `product_id`                | BIGINT UNSIGNED  | NOT NULL                                  | Foreign key referencing `products.product_id`.                                      |
| `language_code`             | VARCHAR(10)      | NOT NULL                                  | Language code (e.g., "en-US", "ja-JP", "zh-CN").                                    |
| `description`               | TEXT             | NOT NULL                                  | AI-generated description in the specified language.                                 |
| `ai_service_version`        | VARCHAR(50)      | NULLABLE                                  | Version of the AI service/model that generated the description.                     |
| `created_at`                | TIMESTAMP        | NOT NULL, DEFAULT CURRENT_TIMESTAMP       | Timestamp of creation.                                                              |

**Constraints:**
*   UNIQUE KEY `idx_product_lang_ai_desc` (`product_id`, `language_code`)

**Indexes:**
*   PRIMARY KEY (`product_ai_description_id`)
*   INDEX `idx_product_id_ai_desc` (`product_id`)
*   INDEX `idx_language_code_ai_desc` (`language_code`)
*   INDEX `idx_created_at_ai_desc` (`created_at`)

---

### 4. Table: `merchants`

| Column         | Type                | Constraints                                                              | Description / Notes                                                                              |
|----------------|---------------------|--------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------|
| `merchant_id`  | BIGINT UNSIGNED     | PRIMARY KEY, AUTO_INCREMENT                                              | Unique identifier for the merchant.                                                              |
| `name`         | VARCHAR(255)        | NOT NULL                                                                 | Merchant's name.                                                                                 |
| `source_type`  | VARCHAR(50)         | NULLABLE                                                                 | Origin/type of merchant (e.g., "erp_sync", "manual_entry", "consignment").                       |
| `contact_info` | JSON                | NULLABLE                                                                 | JSON object for contact details (e.g., `{"phone": "123...", "email": "contact@example.com"}`). |
| `status`       | ENUM(...)           | NOT NULL, DEFAULT 'active', INDEX                                        | ('active', 'inactive', 'pending_review').                                                        |
| `created_at`   | TIMESTAMP           | NOT NULL, DEFAULT CURRENT_TIMESTAMP                                      | Timestamp of creation.                                                                           |
| `updated_at`   | TIMESTAMP           | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP          | Timestamp of last update.                                                                        |

**Indexes:**
*   PRIMARY KEY (`merchant_id`)
*   INDEX `idx_merchant_name` (`name`)
*   INDEX `idx_merchant_status` (`status`)
*   INDEX `idx_merchant_source_type` (`source_type`)

---

### 5. Table: `tags`

| Column      | Type                | Constraints                                                              | Description / Notes                                                                           |
|-------------|---------------------|--------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------|
| `tag_id`    | BIGINT UNSIGNED     | PRIMARY KEY, AUTO_INCREMENT                                              | Unique identifier for the tag.                                                                |
| `name`      | VARCHAR(100)        | NOT NULL, UNIQUE                                                         | Tag name (e.g., "Vintage", "Limited Edition").                                                |
| `type`      | ENUM(...)           | NOT NULL, INDEX                                                          | ('manual', 'rule', 'ai', 'category'). Source/type of the tag.                                 |
| `description` | VARCHAR(255)        | NULLABLE                                                                 | Optional description for the tag's meaning or usage.                                          |
| `created_at`| TIMESTAMP           | NOT NULL, DEFAULT CURRENT_TIMESTAMP                                      | Timestamp of creation.                                                                        |
| `updated_at`| TIMESTAMP           | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP          | Timestamp of last update.                                                                     |

**Indexes:**
*   PRIMARY KEY (`tag_id`)
*   UNIQUE KEY `idx_tag_name` (`name`)
*   INDEX `idx_tag_type` (`type`)

---

### 6. Table: `product_tags`

(Many-to-Many relationship between `products` and `tags` - for manual/curated tags)

| Column             | Type             | Constraints                               | Description / Notes                                                                          |
|--------------------|------------------|-------------------------------------------|----------------------------------------------------------------------------------------------|
| `product_tag_id`   | BIGINT UNSIGNED  | PRIMARY KEY, AUTO_INCREMENT               | Unique identifier for the association.                                                       |
| `product_id`       | BIGINT UNSIGNED  | NOT NULL                                  | Foreign key referencing `products.product_id`.                                               |
| `tag_id`           | BIGINT UNSIGNED  | NOT NULL                                  | Foreign key referencing `tags.tag_id`.                                                       |
| `added_by_user_id` | BIGINT UNSIGNED  | NULLABLE                                  | User ID who added this tag. (References an `users` table - not defined here).                |
| `created_at`       | TIMESTAMP        | NOT NULL, DEFAULT CURRENT_TIMESTAMP       | Timestamp of association.                                                                    |

**Constraints:**
*   UNIQUE KEY `idx_product_tag_unique` (`product_id`, `tag_id`)

**Indexes:**
*   PRIMARY KEY (`product_tag_id`)
*   INDEX `idx_pt_product_id` (`product_id`)
*   INDEX `idx_pt_tag_id` (`tag_id`)

---

### 7. Table: `merchant_tags`

(Many-to-Many relationship between `merchants` and `tags`)

| Column             | Type             | Constraints                               | Description / Notes                                                                          |
|--------------------|------------------|-------------------------------------------|----------------------------------------------------------------------------------------------|
| `merchant_tag_id`  | BIGINT UNSIGNED  | PRIMARY KEY, AUTO_INCREMENT               | Unique identifier for the association.                                                       |
| `merchant_id`      | BIGINT UNSIGNED  | NOT NULL                                  | Foreign key referencing `merchants.merchant_id`.                                             |
| `tag_id`           | BIGINT UNSIGNED  | NOT NULL                                  | Foreign key referencing `tags.tag_id`.                                                       |
| `added_by_user_id` | BIGINT UNSIGNED  | NULLABLE                                  | User ID who added this tag to the merchant. (References an `users` table - not defined here). |
| `created_at`       | TIMESTAMP        | NOT NULL, DEFAULT CURRENT_TIMESTAMP       | Timestamp of association.                                                                    |

**Constraints:**
*   UNIQUE KEY `idx_merchant_tag_unique` (`merchant_id`, `tag_id`)

**Indexes:**
*   PRIMARY KEY (`merchant_tag_id`)
*   INDEX `idx_mt_merchant_id` (`merchant_id`)
*   INDEX `idx_mt_tag_id` (`tag_id`)

---

### 8. Table: `merchant_tag_logs`

| Column         | Type             | Constraints                               | Description / Notes                                                                                                |
|----------------|------------------|-------------------------------------------|--------------------------------------------------------------------------------------------------------------------|
| `log_id`       | BIGINT UNSIGNED  | PRIMARY KEY, AUTO_INCREMENT               | Unique log entry identifier.                                                                                       |
| `merchant_id`  | BIGINT UNSIGNED  | NOT NULL                                  | Foreign key referencing `merchants.merchant_id`.                                                                   |
| `tag_id`       | BIGINT UNSIGNED  | NOT NULL                                  | Foreign key referencing `tags.tag_id`.                                                                             |
| `action`       | ENUM(...)        | NOT NULL                                  | ('added', 'removed'). The action performed on the tag.                                                             |
| `operator_id`  | BIGINT UNSIGNED  | NOT NULL                                  | ID of the user or system performing the action.                                                                    |
| `source`       | VARCHAR(255)     | NULLABLE                                  | Details about source of action (e.g., "manual_ui", "ai_suggestion_approval").                                     |
| `change_details` | JSON             | NULLABLE                                  | Optional JSON for additional details about the change.                                                             |
| `created_at`   | TIMESTAMP        | NOT NULL, DEFAULT CURRENT_TIMESTAMP       | Timestamp of the log entry.                                                                                        |

**Indexes:**
*   PRIMARY KEY (`log_id`)
*   INDEX `idx_mtl_merchant_id` (`merchant_id`)
*   INDEX `idx_mtl_tag_id` (`tag_id`)
*   INDEX `idx_mtl_operator_id` (`operator_id`)
*   INDEX `idx_mtl_created_at` (`created_at`)

---

### 9. Table: `product_listings`

(Tracks product publication on different platforms)

| Column                     | Type             | Constraints                                                              | Description / Notes                                                                                           |
|----------------------------|------------------|--------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------|
| `listing_id`               | BIGINT UNSIGNED  | PRIMARY KEY, AUTO_INCREMENT                                              | Unique identifier for the listing.                                                                            |
| `product_id`               | BIGINT UNSIGNED  | NOT NULL                                                                 | Foreign key referencing `products.product_id`.                                                                |
| `platform_id`              | VARCHAR(50)      | NOT NULL                                                                 | Identifier for the external platform (e.g., "xianyu", "dewu", "95fen").                                       |
| `platform_product_id`      | VARCHAR(255)     | NULLABLE                                                                 | The ID of the product on the external platform.                                                               |
| `status`                   | ENUM(...)        | NOT NULL, INDEX                                                          | ('pending_publish', 'published', 'failed_publish', 'pending_unpublish', 'unpublished', 'sold_on_platform'). |
| `last_publish_attempt_at`  | TIMESTAMP        | NULLABLE                                                                 | Timestamp of the last attempt to publish.                                                                     |
| `published_at`             | TIMESTAMP        | NULLABLE                                                                 | Timestamp when successfully published.                                                                        |
| `unpublished_at`           | TIMESTAMP        | NULLABLE                                                                 | Timestamp when unpublished.                                                                                   |
| `last_sync_at`             | TIMESTAMP        | NULLABLE                                                                 | Timestamp of the last status sync with the platform.                                                          |
| `publish_details`          | JSON             | NULLABLE                                                                 | Platform-specific details (e.g., category, listing options).                                                  |
| `error_message`            | TEXT             | NULLABLE                                                                 | Stores error messages if publishing/unpublishing failed.                                                      |
| `created_at`               | TIMESTAMP        | NOT NULL, DEFAULT CURRENT_TIMESTAMP                                      | Timestamp of creation.                                                                                        |
| `updated_at`               | TIMESTAMP        | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP          | Timestamp of last update.                                                                                     |

**Constraints:**
*   UNIQUE KEY `idx_product_platform_unique` (`product_id`, `platform_id`)

**Indexes:**
*   PRIMARY KEY (`listing_id`)
*   INDEX `idx_pl_product_id` (`product_id`)
*   INDEX `idx_pl_platform_id` (`platform_id`)
*   INDEX `idx_pl_status` (`status`)
*   INDEX `idx_pl_platform_product_id` (`platform_product_id`)

---

### 10. Table: `audit_logs`

(General audit trail for significant actions)

| Column        | Type             | Constraints                               | Description / Notes                                                                                                                               |
|---------------|------------------|-------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------|
| `log_id`      | BIGINT UNSIGNED  | PRIMARY KEY, AUTO_INCREMENT               | Unique log entry identifier.                                                                                                                      |
| `user_id`     | BIGINT UNSIGNED  | NOT NULL, INDEX                           | ID of user/system performing action (e.g., user ID, reserved AI ID like 9999).                                                                    |
| `action_type` | VARCHAR(100)     | NOT NULL, INDEX                           | Type of action (e.g., "CREATE_PRODUCT", "UPDATE_PRODUCT_STATUS").                                                                                 |
| `target_type` | VARCHAR(50)      | NOT NULL, INDEX                           | Type of entity acted upon (e.g., "Product", "Merchant", "Tag").                                                                                   |
| `target_id`   | VARCHAR(255)     | NOT NULL, INDEX                           | ID of the entity acted upon (flexible for different ID types).                                                                                      |
| `details`     | JSON             | NULLABLE                                  | JSON details of the action. For changes: `{"old_value": "...", "new_value": "..."}`. For AI: AI model version, request params.                  |
| `ip_address`  | VARCHAR(45)      | NULLABLE                                  | IP address of the requester.                                                                                                                      |
| `user_agent`  | TEXT             | NULLABLE                                  | User agent of the requester.                                                                                                                      |
| `created_at`  | TIMESTAMP        | NOT NULL, DEFAULT CURRENT_TIMESTAMP, INDEX| Timestamp of the log entry.                                                                                                                       |

**Indexes:**
*   PRIMARY KEY (`log_id`)
*   INDEX `idx_al_user_id` (`user_id`)
*   INDEX `idx_al_action_type` (`action_type`)
*   INDEX `idx_al_target_type_id` (`target_type`, `target_id`)
*   INDEX `idx_al_created_at` (`created_at`)

---

## Order Service Database

### 1. Table: `orders`

| Column              | Type                | Constraints                                                              | Description / Notes                                                                                                  |
|---------------------|---------------------|--------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------|
| `order_id`          | BIGINT UNSIGNED     | PRIMARY KEY, AUTO_INCREMENT                                              | Unique identifier for the order.                                                                                     |
| `external_order_id` | VARCHAR(255)        | NULLABLE, UNIQUE                                                         | Order ID from the external platform (e.g., Xianyu order ID).                                                         |
| `product_id`        | BIGINT UNSIGNED     | NOT NULL, INDEX                                                          | ID of the product being ordered (refers to `products.product_id` in Product Service).                                |
| `buyer_id`          | VARCHAR(255)        | NOT NULL, INDEX                                                          | Identifier for the buyer (platform-specific user ID).                                                                |
| `merchant_id`       | BIGINT UNSIGNED     | NOT NULL, INDEX                                                          | ID of the merchant selling the product (denormalized from Product Service).                                          |
| `platform_id`       | VARCHAR(50)         | NOT NULL, INDEX                                                          | Identifier for the platform where the order originated.                                                              |
| `status`            | ENUM(...)           | NOT NULL, INDEX                                                          | ('pending_lock', 'lock_failed', 'pending_payment', 'paid', 'payment_failed', 'shipped', 'delivered', 'cancelled', 'refunded'). |
| `currency`          | VARCHAR(3)          | NOT NULL                                                                 | Currency of the order.                                                                                               |
| `price_amount`      | DECIMAL(12, 2)      | NOT NULL                                                                 | Price at the time of order.                                                                                          |
| `shipping_address`  | JSON                | NULLABLE                                                                 | Buyer's shipping address.                                                                                            |
| `payment_details`   | JSON                | NULLABLE                                                                 | Details about the payment (e.g., transaction ID, payment method).                                                    |
| `cancellation_reason`| TEXT                | NULLABLE                                                                 | Reason for cancellation, if applicable.                                                                              |
| `created_at`        | TIMESTAMP           | NOT NULL, DEFAULT CURRENT_TIMESTAMP                                      | Timestamp of order creation.                                                                                         |
| `paid_at`           | TIMESTAMP           | NULLABLE                                                                 | Timestamp when payment was confirmed.                                                                                |
| `updated_at`        | TIMESTAMP           | NOT NULL, DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP          | Timestamp of last update.                                                                                            |
| `version`           | INT UNSIGNED        | NOT NULL, DEFAULT 1                                                      | Optimistic locking version.                                                                                          |

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

### 2. Table: `order_status_history`

| Column       | Type             | Constraints                               | Description / Notes                                                                                 |
|--------------|------------------|-------------------------------------------|-----------------------------------------------------------------------------------------------------|
| `history_id` | BIGINT UNSIGNED  | PRIMARY KEY, AUTO_INCREMENT               | Unique identifier for the history entry.                                                            |
| `order_id`   | BIGINT UNSIGNED  | NOT NULL                                  | Foreign key referencing `orders.order_id`.                                                          |
| `status_from`| VARCHAR(50)      | NULLABLE                                  | Previous status.                                                                                    |
| `status_to`  | VARCHAR(50)      | NOT NULL                                  | New status.                                                                                         |
| `changed_by` | VARCHAR(255)     | NULLABLE                                  | Who/what initiated change (e.g., "user_action", "system_event", "payment_gateway_callback").        |
| `notes`      | TEXT             | NULLABLE                                  | Additional notes about the status change.                                                           |
| `created_at` | TIMESTAMP        | NOT NULL, DEFAULT CURRENT_TIMESTAMP       | Timestamp of the status change.                                                                     |

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
