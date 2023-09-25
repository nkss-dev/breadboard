NKSS Breadboard
===============

.. image:: https://img.shields.io/github/go-mod/go-version/nkss-dev/breadboard?logo=Go
    :target: https://go.dev
    :alt: Go version info

A RESTful API wrapper for the database our project, `NKSS <https://github.com/nkss-dev>`_, relies on!

Running a local instance
------------------------

1. **Clone the repository:** ``git clone https://github.com/nkss-dev/breadboard``

2. **Setup initial configuration:** Populate ``sample.env`` with its corresponding values.
    - The ``HMAC_SECRET`` can be any string you like as long as you have a valid JWT token for it. To generate this token, please execute ``go run dev/genjwt.go <role> <rollno>`` in the root directory *after* setting this env variable.

    - ``DATABASE_URL`` refers to the `PostgreSQL instance URL <https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING>`_ that you wish to use.

    - ``PORT`` refers to the port number that you wish to host the backend at.

3. **Install dependencies:** You need the following deps to run this project:
    - `go` 1.18

    You can either install these deps externally or simply run `nix develop` to enter a shell environment with all the dependencies installed.

3. **Generate the required query files:** ``sqlc generate``

4. **Run the project:** ``go run cmd/main.go``

Features
--------

Well, it's a `RESTful <https://restfulapi.net>`_ API. What else?
