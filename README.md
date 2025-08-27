Alur Sistem

## User mengirimkan username dan password melalui POST /login.
## Aplikasi akan cek ke Redis dengan key login_<username>.
## Data user di Redis berupa JSON (realname, email, password hash).
## Password yang dikirim user di-hash dengan SHA1 lalu dibandingkan.
## Jika cocok → login sukses, jika tidak → gagal.

