package utils

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
)

/*
Nama cookie dan path-nya dipakai di beberapa tempat (login, refresh, logout,
middleware), jadi dikumpulkan di sini supaya tidak ada yang salah ketik lalu
menghasilkan cookie yang di-set tapi tidak pernah terbaca.
*/
const (
	AccessCookieName  = "trackora_access"
	RefreshCookieName = "trackora_refresh"

	/*
		Refresh token dibatasi ke satu endpoint saja. Browser hanya mengirim
		cookie yang Path-nya cocok, jadi token berumur panjang ini tidak ikut
		menempel di ratusan request board/card setiap harinya — makin jarang
		lewat kabel, makin sedikit kesempatan bocor.
	*/
	RefreshCookiePath = "/v1/auth/refresh"
)

/*
Secure hanya bisa dinyalakan kalau situsnya HTTPS: cookie Secure yang dikirim
lewat http:// akan diabaikan browser, dan login jadi gagal total di localhost.
Karena itu flag-nya ikut env, bukan dipaku.
*/
func secureCookie() bool {
	return os.Getenv("COOKIE_SECURE") == "true"
}

/*
SameSite=Lax adalah penangkal CSRF utama di sini. Cookie ikut terkirim saat
user mengetik URL atau mengklik tautan ke situs ini, tapi TIDAK ikut pada
POST/PUT/DELETE lintas situs — yang justru semua operasi yang mengubah data.

Kalau suatu saat frontend dan backend beda domain (bukan sekadar beda port),
ini harus jadi "None" dan Secure wajib true, karena "Lax" tidak akan terkirim
sama sekali dari origin lain.
*/
func sameSiteMode() string {
	if mode := os.Getenv("COOKIE_SAMESITE"); mode != "" {
		return mode
	}
	return fiber.CookieSameSiteLaxMode
}

/*
SetAuthCookies menaruh kedua token sebagai cookie HttpOnly.

HttpOnly berarti document.cookie tidak bisa membacanya. Itu inti dari seluruh
perubahan ini: sebelumnya token dikirim di body JSON lalu disimpan frontend di
localStorage, yang bisa dibaca skrip mana pun yang berhasil jalan di halaman —
satu celah XSS sama dengan session dicuri utuh. Sekarang token tidak pernah
tersedia untuk JavaScript sejak awal.

accessMaxAge sengaja pendek. Karena refresh berjalan otomatis lewat cookie,
umur pendek tidak terasa oleh user, tapi memperkecil jendela pakai token yang
sempat bocor lewat cara lain (log proxy, backup, dsb).
*/
func SetAuthCookies(c fiber.Ctx, accessToken, refreshToken string, accessMaxAge, refreshMaxAge time.Duration) {
	c.Cookie(&fiber.Cookie{
		Name:     AccessCookieName,
		Value:    accessToken,
		Path:     "/",
		MaxAge:   int(accessMaxAge.Seconds()),
		HTTPOnly: true,
		Secure:   secureCookie(),
		SameSite: sameSiteMode(),
	})

	/* Refresh token kosong saat sekadar memperpanjang access token. */
	if refreshToken == "" {
		return
	}

	c.Cookie(&fiber.Cookie{
		Name:     RefreshCookieName,
		Value:    refreshToken,
		Path:     RefreshCookiePath,
		MaxAge:   int(refreshMaxAge.Seconds()),
		HTTPOnly: true,
		Secure:   secureCookie(),
		SameSite: sameSiteMode(),
	})
}

/*
ClearAuthCookies menimpa kedua cookie dengan nilai kosong dan MaxAge negatif.

Logout harus dikerjakan server. Setelah cookie jadi HttpOnly, frontend secara
teknis tidak bisa menghapusnya sendiri — tombol keluar yang cuma membersihkan
state React akan tampak berhasil padahal session masih hidup di browser.

Path-nya wajib sama persis dengan saat di-set. Cookie dengan nama sama tapi
path berbeda dianggap cookie lain, jadi salah path berarti yang lama tetap
tertinggal.
*/
func ClearAuthCookies(c fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     AccessCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   secureCookie(),
		SameSite: sameSiteMode(),
	})

	c.Cookie(&fiber.Cookie{
		Name:     RefreshCookieName,
		Value:    "",
		Path:     RefreshCookiePath,
		MaxAge:   -1,
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   secureCookie(),
		SameSite: sameSiteMode(),
	})
}
