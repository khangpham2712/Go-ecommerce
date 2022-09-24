##**Ecommerce With Golang Project By Akhil Sharma**
**Code-Along video series on Youtube**

```bash
# You can start the project with below commands
docker-compose up -d
go run main.go
```

- **SIGNUP FUNCTION API CALL (POST REQUEST)**

http://localhost:8000/user/sign-up

```json
{
  "first_name": "Akhil",
  "last_name": "Sharma",
  "email": "akhil@gmail.com",
  "password": "akhilsharma",
  "phone": "+4534545435"
}
```

Response :"Successfully Signed Up!!"

- **LOGIN FUNCTION API CALL (POST REQUEST)**

  http://localhost:8000/user/log-in

```json
{
  "email": "akhil@gmail.com",
  "password": "akhilsharma"
}
```

response will be like this

```json
{
  "_id": "***********************",
  "first_name": "akhil",
  "last_name": "sharma",
  "password": "$2a$14$UIYjkTfnFnhg4qhIfhtYnuK9qsBQifPKgu/WPZAYBaaN17j0eTQZa",
  "email": "akhil@gmail.com",
  "phone": "+4534545435",
  "token": "eyJc0Bwcm90b25vbWFpbC5jb20iLCJGaXJzdF9OYW1lIjoiam9zZXBoIiwiTGFzdF9OYW1lIjoiaGVybWlzIiwiVWlkIjoiNjE2MTRmNTM5ZjI5YmU5NDJiZDlkZjhlIiwiZXhwIjoxNjMzODUzNjUxfQ.NbcpVtPLJJqRF44OLwoanynoejsjdJb5_v2qB41SmB8",
  "Refresh_Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6IiIsIkZpcnLCJVaWQiOiIiLCJleHAiOjE2MzQzNzIwNTF9.ocpU8-0gCJsejmCeeEiL8DXhFcZsW7Z3OCN34HgIf2c",
  "created_at": "2022-04-09T08:14:11Z",
  "updtaed_at": "2022-04-09T08:14:11Z",
  "user_id": "61614f539f29be942bd9df8e",
  "user_cart": [],
  "addresses": [],
  "orders": []
}
```

- **Admin add Product Function (POST REQUEST)**

  http://localhost:8000/admin/add-product

```json
{
  "product_name": "Alienware x15",
  "price": 2500,
  "rating": 10,
  "image": "alienware.jpg"
}
```

Response : "Successfully added our Product Admin!!"

- **View all the Products in db GET REQUEST**

  http://localhost:8000/user/view-products

Response

```json
[
  {
    "Product_ID": "6153ff8edef2c3c0a02ae39a",
    "product_name": "alienwarex15",
    "price": 1500,
    "rating": 10,
    "image": "alienware.jpg"
  },
  {
    "Product_ID": "616152679f29be942bd9df8f",
    "product_name": "giner ale",
    "price": 900,
    "rating": 5,
    "image": "gin.jpg"
  },
  {
    "Product_ID": "616152ee9f29be942bd9df90",
    "product_name": "iphone 13",
    "price": 1700,
    "rating": 4,
    "image": "ipho.jpg"
  },
  {
    "Product_ID": "616152fa9f29be942bd9df91",
    "product_name": "whiskey",
    "price": 100,
    "rating": 7,
    "image": "whis.jpg"
  },
  {
    "Product_ID": "616153039f29be942bd9df92",
    "product_name": "acer predator",
    "price": 3000,
    "rating": 10,
    "image": "acer.jpg"
  }
]
```

- **Search Product by regex function (GET REQUEST)**

defines the word search sorting
http://localhost:8000/user/search?name=al

response:

```json
[
  {
    "Product_ID": "616152fa9f29be942bd9df91",
    "product_name": "Alienware x15",
    "price": 1500,
    "rating": 10,
    "image": "1.jpg"
  },
  {
    "Product_ID": "616153039f29be942bd9df92",
    "product_name": "ginger Ale",
    "price": 300,
    "rating": 10,
    "image": "1.jpg"
  }
]
```

- **Adding the Products to the Cart (GET REQUEST)**

  http://localhost:8000/user/add-to-cart?id=xxxproduct_idxxx&userID=xxxxxxuser_idxxxxxx

  Corresponding mongodb query

- **Removing Item From the Cart (GET REQUEST)**

  http://localhost:8000/user/remove-item?id=xxxxxxx&userID=xxxxxxxxxxxx

- **Listing the item in the users cart (GET REQUEST) and total price**

  http://localhost:8000/user/list-cart?id=xxxxxxuser_idxxxxxxxxxx

- **Addding the Address (POST REQUEST)**

  http://localhost:8000/user/add-adress?id=user_id**\*\***\***\*\***

  The Address array is limited to two values home and work address more than two address is not acceptable

```json
{
  "house": "white house",
  "street": "white street",
  "city": "washington",
  "pin_code": "332423432"
}
```

- **Editing the Home Address(PUT REQUEST)**

  http://localhost:8000/user/edit-home-address?id=xxxxxxxxxxuser_idxxxxxxxxxxxxxxx

- **Editing the Work Address(PUT REQUEST)**

  http://localhost:8000/user/edit-work-address?id=xxxxxxxxxxuser_idxxxxxxxxxxxxxxx

- **Delete Addresses(GET REQUEST)**

  http://localhost:8000/user/delete-addresses?id=xxxxxxxxxuser_idxxxxxxxxxxxxx

  delete both addresses

- **Cart Checkout Function and placing the order(GET REQUEST)**

  After placing the order the items have to be deleted from cart functonality added

  http://localhost:8000/user/cart-checkout?id=xxuser_idxxx

- **Instantly Buying the Products(GET REQUEST)**
  http://localhost:8000/user/instant-buy?userid=xxuser_idxxx&pid=xxxxproduct_idxxxx
