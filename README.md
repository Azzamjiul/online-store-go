# Online Store

## Penyebab
Para pengguna kesal karena mereka sudah membayar produk tapi ternyata produk yang mereka beli telah habis padahal ketika mereka checkout jumlah stok masih tersedia. 
Hal ini disebabkan ketika proses checkout, **sistem kurang memvalidasi** apakah stoknya masih ada sebelum memproses pesanan ke halaman pembayaran. Karena pada prinsipnya ketika sudah masuk ke halaman pembayarn maka barang tersebut tidak boleh dipesan orang lain.

## Solusi
Sistem melakukan validasi beberapa kali untuk memastikan:
- **Order yang dilakukan tidak melebihi stok yang ada**
- **Stok berhasil diperbarui setelah ada order yang masuk**

Berikut diagram alur yang diajukan.