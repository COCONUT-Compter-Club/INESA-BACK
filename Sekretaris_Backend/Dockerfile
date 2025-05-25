# Gunakan image base Go versi terbaru
FROM golang:latest

# Buat direktori kerja di dalam container
WORKDIR /app

# Salin file go.mod dan go.sum untuk menginstal dependensi
COPY go.mod go.sum ./

# Instal dependensi
RUN go mod download

# Salin seluruh kode sumber ke dalam container
COPY . .

# Build aplikasi Go
RUN go build -o main .

# Expose port yang digunakan oleh aplikasi (sesuaikan dengan kebutuhan)
EXPOSE 8088

# Perintah untuk menjalankan aplikasi
CMD ["./main"]


