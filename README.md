# online-store-api

Synapsis Backend Challenge.

This is a simple api server to simulate customer order system.

## How to run

To running the API service and databases using docker-compose:
(need to run db migration manually)
```
docker-compose up
```

For development:
create .env file
```
make migrate-up
make dev
```

Check the running service:
```
curl localhost:9000/v1/health
```

To generate dummy categories and products data:
```
curl localhost:9000/v1/dummy-data
```

## Database Structure

ERD

ORDER STATUS:
- payment_pending
- paid

PAYMENT STATUS:
- pending
- settlement

## API Documentation

All endpoints are available in swagger and postman collection:
- Swagger: http://localhost:9000/v1/swagger/index.html
- Postman: Synapsis_Online-Store-API.postman_collection.json

### User Register

Create new account

Endpoint: POST /v1/users/register

Json Body:
```
{
  "email": "raisa@email.com",
  "password": "Secret123",
  "full_name": "Raisa"
}
```

Created Response (201):
```
{
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGFpbXMiOnsidXNlcl9pZCI6Ijg4ZTBhZTEyLWQwNjQtNDE2NS04OGY3LWY3ODY2NmVkNDY2MCJ9LCJpc3MiOiJvbmxpbmUtc3RvcmUiLCJzdWIiOiI4OGUwYWUxMi1kMDY0LTQxNjUtODhmNy1mNzg2NjZlZDQ2NjAiLCJhdWQiOlsib25saW5lLXN0b3JlIl0sImV4cCI6MTcxNjY2NDkzMiwibmJmIjoxNzE2NjIxNzMyLCJpYXQiOjE3MTY2MjE3MzIsImp0aSI6ImY1NzcwNjVjLTAzOGYtNDQ4OS04YTg4LTYzOGYxNWVhMGE4NCJ9.P-3pi6fCT1jMDwdQg7HnQyvum132Zb0ZAlDrFxP9r8g",
    "user": {
      "id": "88e0ae12-d064-4165-88f7-f78666ed4660",
      "full_name": "Raisa",
      "created_at": "2024-05-25T07:22:12.580511Z"
    }
  }
}
```

### User Login

Login into an account

Endpoint: POST /v1/users/login

Json Body:
```
{
  "email": "raisa@email.com",
  "password": "Secret123"
}
```

Success Response (200):
```
{
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGFpbXMiOnsidXNlcl9pZCI6Ijg4ZTBhZTEyLWQwNjQtNDE2NS04OGY3LWY3ODY2NmVkNDY2MCJ9LCJpc3MiOiJvbmxpbmUtc3RvcmUiLCJzdWIiOiI4OGUwYWUxMi1kMDY0LTQxNjUtODhmNy1mNzg2NjZlZDQ2NjAiLCJhdWQiOlsib25saW5lLXN0b3JlIl0sImV4cCI6MTcxNjY2NTAwNywibmJmIjoxNzE2NjIxODA3LCJpYXQiOjE3MTY2MjE4MDcsImp0aSI6ImUzNDJhZDQ4LTI1ZjItNGY3YS1hZjc1LTZkNjNkNzY1ZmU1OSJ9.TPIr_EBk63KfWAY_lzX8umDRoHjAFeKHDBivhnhLXkY",
    "user": {
      "id": "88e0ae12-d064-4165-88f7-f78666ed4660",
      "full_name": "Raisa",
      "created_at": "2024-05-25T07:22:12.580511Z"
    }
  }
}
```

### Get Logged in User

Get active user profile

Endpoint: GET /v1/users/me

Header:
- Authorization: Bearer <access_token>

Success Response (200):
```
{
  "message": "success",
  "data": {
    "created_at": "2024-05-25T07:22:12.580511Z",
    "email": "raisa@email.com",
    "full_name": "Raisa",
    "id": "88e0ae12-d064-4165-88f7-f78666ed4660",
    "updated_at": "2024-05-25T07:22:12.580511Z"
  }
}
```

### Get List Category

Get all cateogiries (paginated)

Endpoint: GET /v1/categories

Query:
- page (int default 1, max 500)
- count (int default 10, max 100)

Success Response (200):
```
{
  "message": "success",
  "data": [
    {
      "id": "903653ba-49f4-4a95-abd1-071fd7a247d9",
      "name": "Elektronik",
      "description": "Peralatan elektronik berkualitas",
      "created_at": "2024-05-25T07:31:06.464546Z",
      "updated_at": "2024-05-25T07:31:06.464546Z"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "count": 10
  }
}
```


### Get List Product By Category

Get all cateogiries (paginated)

Endpoint: GET /v1/products/category/:category_id

Query:
- page (int default 1, max 500)
- count (int default 10, max 100)

Success Response (200):
```
{
  "message": "success",
  "data": [
    {
      "id": "7b25cc0e-314b-4a21-9738-3a4a26a7047d",
      "category_id": "903653ba-49f4-4a95-abd1-071fd7a247d9",
      "name": "TV Samsung 32 Inch",
      "description": "Bergaransi seumur jagung",
      "price": 3000000,
      "stock": 5,
      "created_at": "2024-05-25T07:31:06.464546Z",
      "updated_at": "2024-05-25T07:31:06.464546Z",
      "category_name": "Elektronik"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "count": 10
  }
}
```

### Add Cart Item

Endpoint: POST /v1/carts

Header:
- Authorization: Bearer <access_token>

Json Body:
```
{
  "product_id": "7b25cc0e-314b-4a21-9738-3a4a26a7047d",
  "notes": "yang warna ungu pls"
}
```

Created Response (201):
```
{
  "message": "success",
  "data": {
    "id": "5fc56fe8-3d1f-4754-b9bf-a38787a2ba51",
    "user_id": "3541a4dd-79c8-41e9-9239-7ae9b98343c9",
    "product_id": "7b25cc0e-314b-4a21-9738-3a4a26a7047d",
    "notes": "yang warna ungu pls",
    "created_at": "2024-05-25T07:45:17.064538Z",
    "updated_at": "2024-05-25T07:45:17.064538Z"
  }
}
```

### Remove Cart Item

Endpoint: DELETE /v1/carts/:cart_id

Header:
- Authorization: Bearer <access_token>

Success Response (200):
```
{
  "message": "success",
  "data": null
}
```

### Get List Cart Item

Endpoint: GET /v1/carts

Header:
- Authorization: Bearer <access_token>

Success Response (200):
```
{
  "message": "success",
  "data": [
    {
      "id": "d50d58e6-cdb0-46c2-9e00-e0dd668600c5",
      "user_id": "3541a4dd-79c8-41e9-9239-7ae9b98343c9",
      "product_id": "7b25cc0e-314b-4a21-9738-3a4a26a7047d",
      "notes": "yang warna ungu pls",
      "created_at": "2024-05-25T07:47:30.18255Z",
      "updated_at": "2024-05-25T07:47:30.18255Z",
      "category_id": "903653ba-49f4-4a95-abd1-071fd7a247d9",
      "product_name": "TV Samsung 32 Inch",
      "product_description": "Bergaransi seumur jagung",
      "product_price": 3000000,
      "product_stock": 5,
      "category_name": "Elektronik"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "count": 10
  }
}
```

### Create Order

Endpoint: POST /v1/orders

Header:
- Authorization: Bearer <access_token>

Json Body:
```
{
  "shipment_address": "Jl. depan rumahku",
  "payment": {
    "payment_method": "bank_transfer",
    "payment_provider": "bca"
  },
  "items": [
    {
      "product_id": "7b25cc0e-314b-4a21-9738-3a4a26a7047d",
      "qty": 1,
      "notes": "yg warna ungu pls"
    },
    {
      "product_id": "ce99b061-43fa-4b4c-8a5a-3734a500e2cd",
      "qty": 2,
      "notes": ""
    }
  ]
}
```

Created Response (201):
```
{
  "message": "success",
  "data": {
    "id": "44ccbaa9-fd7f-4115-9c8b-c1346e4dba15",
    "user_id": "3541a4dd-79c8-41e9-9239-7ae9b98343c9",
    "status": "payment_pending",
    "other_cost": 0,
    "total_cost": 3360000,
    "shipment_address": "Jl. depan rumahku",
    "created_at": "2024-05-25T07:49:29.954533Z",
    "updated_at": "2024-05-25T07:49:29.954533Z"
  }
}
```

### Get List Order

Endpoint: GET /v1/orders

Header:
- Authorization: Bearer <access_token>

Success Response (200):
```
{
  "message": "success",
  "data": [
    {
      "id": "44ccbaa9-fd7f-4115-9c8b-c1346e4dba15",
      "user_id": "3541a4dd-79c8-41e9-9239-7ae9b98343c9",
      "status": "payment_pending",
      "other_cost": 0,
      "total_cost": 3360000,
      "shipment_address": "Jl. depan rumahku",
      "created_at": "2024-05-25T07:49:29.954533Z",
      "updated_at": "2024-05-25T07:49:29.954533Z"
    }
  ],
  "meta": {
    "total": 1,
    "page": 1,
    "count": 10
  }
}
```

### Get Order

Endpoint: GET /v1/orders/:order_id

Header:
- Authorization: Bearer <access_token>

Success Response (200):
```
{
  "message": "success",
  "data": {
    "items": [
      {
        "id": "02a1a6a3-1c9c-4f46-ae18-162e2b0d7a9a",
        "order_id": "44ccbaa9-fd7f-4115-9c8b-c1346e4dba15",
        "product_id": "7b25cc0e-314b-4a21-9738-3a4a26a7047d",
        "qty": 1,
        "product_price": 3000000,
        "notes": "yg warna ungu pls",
        "created_at": "2024-05-25T07:49:29.954533Z",
        "product_name": "TV Samsung 32 Inch",
        "product_description": "Bergaransi seumur jagung"
      },
      {
        "id": "07c9e5dd-0f12-47bb-8ab5-74a4894fd9d4",
        "order_id": "44ccbaa9-fd7f-4115-9c8b-c1346e4dba15",
        "product_id": "ce99b061-43fa-4b4c-8a5a-3734a500e2cd",
        "qty": 2,
        "product_price": 180000,
        "notes": "",
        "created_at": "2024-05-25T07:49:29.954533Z",
        "product_name": "Kipas angin kesedot sampah",
        "product_description": "Wuss"
      }
    ],
    "order": {
      "id": "44ccbaa9-fd7f-4115-9c8b-c1346e4dba15",
      "user_id": "3541a4dd-79c8-41e9-9239-7ae9b98343c9",
      "status": "payment_pending",
      "other_cost": 0,
      "total_cost": 3360000,
      "shipment_address": "Jl. depan rumahku",
      "created_at": "2024-05-25T07:49:29.954533Z",
      "updated_at": "2024-05-25T07:49:29.954533Z"
    },
    "payment": {
      "id": "21381a91-5be7-49ce-97bd-d5eccbb76d2d",
      "order_id": "44ccbaa9-fd7f-4115-9c8b-c1346e4dba15",
      "payment_method": "bank_transfer",
      "payment_provider": "bca",
      "bill_amount": 3360000,
      "paid_amount": 0,
      "status": "pending",
      "transaction_id": "",
      "paid_at": null,
      "log": null,
      "created_at": "2024-05-25T07:49:29.954533Z",
      "updated_at": "2024-05-25T07:49:29.954533Z"
    }
  }
}
```

### Payment Webhook

Endpoint: POST /v1/payments/webhook

Header:
- X-API-KEY: <API Key>

Json Body:
```
{
  "transaction_id": "eae7efe4-tx-id-from-payment-gateway-44ccbaa9",
  "payment_amount": 3360123,
  "status": "settlement",
  "transaction_details": {
    "order_id": "44ccbaa9-fd7f-4115-9c8b-c1346e4dba15",
    "gross_amount": 3360000
  },
  "user_details": {
    "full_name": "Raisa"
  }
}
```

Success Response (200):
```
{
  "message": "success",
  "data": null
}
```
