# pgx-sqlc

Demo app for Plumb.

## stack

Designed as a monolithic Go stack. As simple as possible.

* Go
* front:
  * templ
  * templui
  * HTMX
  * Tailwind
* back:
  * Postgres
  * sqlc
  * pgx

## structure

### db

This package holds first-level wrappers for pgx + sqlc.

The .sql templates live here, which sqlc uses to generate its code.

#### db/sqlc

For generated sqlc code only.

### ui

everything related to the frontend, including templ templates and templui components.

See .templui.json
