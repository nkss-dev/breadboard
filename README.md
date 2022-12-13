# NKSS-backend

A RESTful API wrapper for [NKSSS](https://github.com/NIT-KKR-Student-Support-System "NIT-KKR Student Support System")'s database.

## Running a local instance

1. **Clone the repository:** `git clone https://github.com/NIT-KKR-Student-Support-System/nkss-backend`

2. **Setup initial configuration:** Populate `sample.env` with its corresponding values.

- The `HMAC_SECRET` can be any string you like as long as you have a valid JWT token for it. To generate this token, please execute `go run dev/genjwt.go <role> <rollno>` in the root directory _after_ setting this env variable.

- `DATABASE_URL` refers to the [PostgreSQL instance URL](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING) that you wish to use.

- `PORT` refers to the port number that you wish to host the backend at.

3. **Run the project:** `go run cmd/main.go`

## Features

- Endpoints follow the [REST API](https://restfulapi.net/) constraints.
- Response JSONs follow the [JSON:API](https://jsonapi.org) specifications.
- Database-independent end points to service some other requirements of [NKSSS](https://github.com/NIT-KKR-Student-Support-System "NIT-KKR Student Support System")
