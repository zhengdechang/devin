# Product Service

This service is responsible for managing all aspects of product information for the二手奢侈品商品微服务系统 (Second-hand Luxury Goods Microservice System).

## Core Responsibilities:
- Receiving and synchronizing product data from ERP systems.
- Invoking AI services to enrich products with tags and multi-language descriptions.
- Managing product status through its lifecycle (draft, listed, locked, sold).
- Handling product上架 (listing) to multiple platforms.
- Storing and managing relationships between products, merchants, and tags.
- Providing a comprehensive audit trail for all product-related operations.

## Design Documents:
- [Database Design](./database_design.md)
- [API Design](./api_design.md)

## Key Sub-directories:
- `/service`: Contains the core business logic, including AI service integration and tag processing.
- `/repo`: Handles database operations for product master data, AI-generated content, and tag relationships.
- `/audit`: Provides an abstraction layer for logging all significant actions.
- `/aiclient`: Encapsulates gRPC calls to the external AI service.
- `/merchantclient`: Encapsulates gRPC calls to the Merchant service (for tag-based merchant filtering).
