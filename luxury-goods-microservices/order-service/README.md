# Order Service

This service manages the lifecycle of orders within the 二手奢侈品商品微服务系统 (Second-hand Luxury Goods Microservice System). It is decoupled from the Product Service to handle order-specific responsibilities.

## Core Responsibilities:
- Processing new order placements.
- Coordinating with the Product Service to lock items when an order is initiated.
- Managing order status transitions (e.g., pending lock, locked, paid, cancelled, fulfilled).
- Handling payment-related logic and updates (though payment gateway integration is a separate concern).
- Providing an audit trail for all order-related operations.

## Design Documents:
- [Database Design](./database_design.md)
- [API Design](./api_design.md)

## Key Sub-directories:
- `/service`: Contains the core business logic for order processing, status management, and state transitions.
- `/repo`: Handles database operations for the orders table.
- `/client`: Encapsulates gRPC calls to the Product Service (specifically for locking/unlocking products).
- `/audit`: Provides an abstraction layer for logging order-related actions, potentially reusing a common framework.
