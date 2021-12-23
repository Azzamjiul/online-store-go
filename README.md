# Online Store

## Penyebab
Para pengguna kesal karena mereka sudah membayar produk tapi ternyata produk yang mereka beli telah habis padahal ketika mereka checkout jumlah stok masih tersedia. 
Hal ini disebabkan ketika proses checkout, **sistem kurang memvalidasi** apakah stoknya masih ada sebelum memproses pesanan ke halaman pembayaran. Karena pada prinsipnya ketika sudah masuk ke halaman pembayarn maka barang tersebut tidak boleh dipesan orang lain.

## Solusi
Sistem melakukan validasi beberapa kali untuk memastikan:
- **Order yang dilakukan tidak melebihi stok yang ada**
- **Stok berhasil diperbarui setelah ada order yang masuk**

### ERD
![ERD](/readme/erd.png "ERD")
### Berikut diagram alur yang diajukan.
![Diagaram Validasi](/readme/validasi.png "Diagaram Validasi")

## Functional Testing

## Case berhasil
### Request
```
curl --location --request POST 'localhost:8080/orders' --header 'Content-Type: application/json' --data-raw '{
    "user_id": 1,
    "total": 7000,
    "items" : [
        {
            "item_id": 1,
            "quantity": 3,
            "price": 1000
        },
        {
            "item_id": 2,
            "quantity": 2,
            "price": 2000
        }
    ]
}'
```

### Response
```
{"message":"order created successfully"}
```

## Case gagal
### Request
```
curl --location --request POST 'localhost:8080/orders' --header 'Content-Type: application/json' --data-raw '{
    "user_id": 1,
    "total": 7000,
    "items" : [
        {
            "item_id": 1,
            "quantity": 30,
            "price": 1000
        },
        {
            "item_id": 2,
            "quantity": 2,
            "price": 2000
        }
    ]
}'
```

### Response
```
{"message":"insufficient stock of items (Product 1)","status":400,"error":"bad_request"}
```
