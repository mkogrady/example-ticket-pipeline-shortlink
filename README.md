# shortlink

A small link shortener in Go. It keeps links in memory, so it is for
development and demos, not for production.

## Run it

    go run .

Create a link, then follow it:

    curl -s -X POST localhost:8080/links -d '{"url":"https://example.com"}'
    curl -si localhost:8080/<code> | head -1

## API

`POST /links` takes `{"url": "<target>"}` and answers `201` with
`{"code": "<code>", "short_url": "http://<host>/<code>"}`. A body that
is not JSON, or a URL the service does not accept, gets `400` with
`{"error": "<reason>"}`.

`GET /<code>` answers `302` with the target in `Location`, or `404` for
a code it does not know.

## Layout

- `validate`: the rules a target URL must pass.
- `store`: the in-memory map from code to target.
- `server`: the HTTP handlers.

## License

To be decided.
