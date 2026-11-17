package helper

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// bcryptCost menentukan berapa kali proses hashing diulang secara internal.
// Semakin besar nilainya, semakin lambat percobaan menebak password satu per satu.
// Nilai 12 adalah kompromi yang lazim: cukup lambat bagi penyerang,
// masih cukup cepat untuk sekali login.
const bcryptCost = 12

// dummyHash dipakai ketika username tidak ditemukan, agar waktu tanggap
// login tetap mirip dengan kasus password salah. Tanpa ini, selisih waktu
// respons membocorkan username mana yang terdaftar (timing attack).
var dummyHash = []byte("$2a$12$abcdefghijklmnopqrstuuLKa3Bt1TCmU/6zvhZ8x4nq1yBiuGvS")

// HashPassword mengubah password plaintext menjadi hash bcrypt yang tidak
// dapat dikembalikan ke bentuk aslinya. bcrypt menyisipkan salt acak ke dalam
// hasilnya, sehingga dua user dengan password sama tetap menghasilkan hash
// yang berbeda.
func HashPassword(plain string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// VerifyPassword membandingkan password plaintext dengan hash bcrypt-nya.
// Mengembalikan true jika cocok, false jika tidak.
func VerifyPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

// VerifyDummyPassword sengaja membuang waktu seperti VerifyPassword normal,
// dipakai ketika username tidak ditemukan agar response time tetap konsisten
// dan tidak membocorkan informasi keberadaan username (timing attack prevention).
func VerifyDummyPassword(plain string) {
	_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(plain))
}

// RandomToken menghasilkan string hex acak yang aman secara kriptografis
// sepanjang numBytes byte (hasil hex-nya 2x panjang numBytes).
// math/rand TIDAK boleh dipakai untuk keperluan ini karena hasilnya
// dapat diperkirakan bila seed-nya diketahui.
func RandomToken(numBytes int) (string, error) {
	buf := make([]byte, numBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// SHA256Hex menghasilkan hash SHA-256 dari sebuah string dalam bentuk hex.
// Dipakai untuk menyimpan refresh token dalam bentuk hash di database.
// bcrypt tidak dipakai di sini karena token sudah acak dan panjang,
// sehingga tidak perlu diperlambat seperti password buatan manusia.
func SHA256Hex(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}