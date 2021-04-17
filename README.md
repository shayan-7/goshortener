# goshortener

A URL shortener in Go.

#### Prerequisites:

You should have `postgres` and `redis` installed which listening on ports 5432
and 6379 respectively.

#### Usage:

To serve the project, use:

```sh
make serve
```

To make subscriber working:

```sh
make subscribe
```

Subscriber listen for messages about viewed url, So it can calculate and insert
the number of views in `Stats` table.

#### List of API:

1. Creating a shortened url

   path: / \
   verb: POST \
   form: url -> `string` (ex: `example.com`) \
   response: {"id": "wEJQBwz", "url": "example.com"} (example) \
   Header: `Authorization: Bearer <JWT token>` (Case Sensitive)

2. Getting a shortened url

   path: /{id:`shortened url`} (ex: `/wEJQBwz`) \
   verb: GET \
   response: {"id": "wEJQBwz", "url": "example.com"} (example)

3. Signup

   path: /signup \
   verb: POST \
   form: username -> `string` (ex: `username`)
   password -> `string` (ex: `password`) \
   response: {"username": "username"} (example)

4. Login

   path: /login \
   verb: POST \
   form: username -> `string` (ex: `username`)
   password -> `string` (ex: `password`) \
   response: {"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6Im1lbWJlcjIiLCJleHAiOjE2MTg3Mjk3MzAsImlzcyI6IkF1dGhTZXJ2aWNlIn0.EBQDC9cxOdV2ob8Sujy0iSnzmJi5gYLNfOtdQSvp_gw"} (example)

**Following API list should be implemented in V2**

5. Online stats

   path: /stats/online \
   verb: GET \
   response: {"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6Im1lbWJlcjIiLCJleHAiOjE2MTg3Mjk3MzAsImlzcyI6IkF1dGhTZXJ2aWNlIn0.EBQDC9cxOdV2ob8Sujy0iSnzmJi5gYLNfOtdQSvp_gw"} (example)

6. Daily stats

   path: /stats/daily \
   verb: GET \
   response: {"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6Im1lbWJlcjIiLCJleHAiOjE2MTg3Mjk3MzAsImlzcyI6IkF1dGhTZXJ2aWNlIn0.EBQDC9cxOdV2ob8Sujy0iSnzmJi5gYLNfOtdQSvp_gw"} (example)

7. Weekly stats

   path: /stats/weekly \
   verb: GET \
   response: {"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6Im1lbWJlcjIiLCJleHAiOjE2MTg3Mjk3MzAsImlzcyI6IkF1dGhTZXJ2aWNlIn0.EBQDC9cxOdV2ob8Sujy0iSnzmJi5gYLNfOtdQSvp_gw"} (example)

#### Benchmarking:

You can use one of [benchmarking tools](https://gist.github.com/denji/8333630).

#### TODO:

    1. Provide tests for all the packages
    2. Dockerize the project
