# Order Service Database Design

## Order Table (`Order`)
- `order_id`: BIGINT PRIMARY KEY AUTO_INCREMENT
- `product_id`: BIGINT (Should ideally reference Product in Product Service, but direct FK might be across service boundaries. Eventual consistency or local cache/replication might be used.)
- `buyer_id`: BIGINT (Identifier for the buyer)
- `platform_id`: VARCHAR(255) (Identifier for the platform where the order originated)
- `status`: VARCHAR(50) (e.g., 'pending_payment', 'paid', 'shipped', 'delivered', 'cancelled')
- `created_at`: TIMESTAMP DEFAULT CURRENT_TIMESTAMP
- `paid_at`: TIMESTAMP (Nullable, set when payment is confirmed)
- `total_amount`: DECIMAL(10, 2)
- `currency`: VARCHAR(3) (e.g., 'USD', 'JPY')
- `shipping_address`: TEXT
- `billing_address`: TEXT (Nullable, if same as shipping)
