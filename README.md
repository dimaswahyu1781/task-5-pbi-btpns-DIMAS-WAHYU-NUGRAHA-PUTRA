# task-5-pbi-btpns-DIMAS-WAHYU-NUGRAHA-PUTRA

Pada task akhir VIX Full Stack Developer ini kalian diarahkan untuk membentuk API berdasarkan kasus yang telah diberikan. Pada kasus ini, kalian diinstruksikan untuk membuat API untuk mengupload dan menghapus gambar. API yang kalian bentuk adalah POST, GET, PUT, dan DELETE.

A. Ketentuan API :

Pada bagian User Endpoint :
  1. POST : /users/register, dan gunakan atribut berikut ini :
      - ID (primary key, required)
      - Username (required)
      - Email (unique & required)
      - Password (required & minlength 6)
      - Relasi dengan model Photo (Gunakan constraint cascade)
      - Created At (timestamp)
      - Updated At (timestamp)
  2. GET: /users/login
      - Using email & password (required)
  3. PUT : /users/:userId (Update User)
  4. DELETE : /users/:userId (Delete User)

Photos Endpoint :
  1. POST : /photos
      - ID
      - Title
      - Caption
      - PhotoUrl
      - UserID
      - Relasi dengan model User
  2. GET : /photos
  3. PUT : /photoId
  4. DELETE : /:photoId

Requirement :
1. Authorization dapat menggunakan tool Go JWT
○ https://github.com/dgrijalva/jwt-go
2. Pastikan hanya user yang membuat foto yang dapat menghapus / mengubah foto Struktur dokumen / environment dari GoLang yang akan dibentuk kurang lebih sebagai berikut :
● app
Menampung pembuatan struct dalam kasus ini menggunakan struct User
untuk keperluan data dan authentication
● controllers
Berisi antara logic database yaitu models dan query
● database
Berisi konfigurasi database serta digunakan untuk menjalankan koneksi
database dan migration
● helpers
Berisi fungsi-fungsi yang dapat digunakan di setiap tempat dalam hal ini jwt,
bcrypt, headerValue
● middlewares
Berisi fungsi yang digunakan untuk proses otentikasi jwt yang digunakan untuk
proteksi api
● models
Berisi models yang digunakan untuk relasi database
● router
Berisi konfigurasi routing / endpoint yang akan digunakan untuk mengakses api
● go mod
Yang digunakan untuk manajemen package / dependency berupa library


B. Tools yang dapat kalian gunakan :
  - Gin Gonic Framework : https://github.com/gin-gonic/gin
  - Gorm : https://gorm.io/index.html
  - JWT Go : https://github.com/dgrijalva/jwt-go
  - Go Validator : http://github.com/asaskevich/govalidator

Untuk database, gunakanlah server SQL open source seperti MySQL, PostgreSQL, atau Microsoft SQL Server.


## Hasil
Dapat melakukan CREATE, READ, UPDATE, DELETE pada users dan photos

### Fitur yang belum selesai :
1. Fitur Authentication
2. Semua code di jalankan hanya satu file yaitu main.go
3. Belum adanya atribut Email, created at, updated at, dan Relasi pada users serta belum adanya relasi pada photos

### Tools yang digunakan :
- Bahasa Pemrograman GoLang & MySQL
- XAMPP
- Visual Studio Code
- Command Prompt
- Postman

### Referensi
- https://www.youtube.com/watch?v=gmP50RUd6YA
