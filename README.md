# NKSS-backend

A RESTful API wrapper for [NKSSS](https://github.com/NIT-KKR-Student-Support-System "NIT-KKR Student Support System")'s database.

## Running a local instance

1. **Clone the repository:** `git clone https://github.com/NIT-KKR-Student-Support-System/nkss-backend`

2. **Setup initial configuration:**
- Populate `sample.env` with its corresponding values.
- The `HMAC_SECRET` can be any string you like as long as you have a valid JWT token for it. To generate this token, please execute `go run dev/genjwt.go <role> <rollno>` in the root directory _after_ setting this env variable.

3. **Run the project:** `go run cmd/main.go`

## Features

- Endpoints follow the [REST API](https://restfulapi.net/) constraints.
- Response JSONs follow the [JSON:API](https://jsonapi.org) specifications.
- Database-independent end points to service some other requirements of [NKSSS](https://github.com/NIT-KKR-Student-Support-System "NIT-KKR Student Support System")
