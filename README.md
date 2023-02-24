<h1 align = "center"> AffiliateManager </h1>

<div align="center">

[![Test Suite](https://github.com/aaronangxz/AffiliateManager/actions/workflows/test.yml/badge.svg)](https://github.com/aaronangxz/AffiliateManager/actions/workflows/test.yml)

</div>

<div align="center">This is the backend server implementation of AffiliateManager.</div>

<h2> Development </h2>

1. Prepare an `.env` file under project root with the following credentials
    ```
    PROD_DB_HOST=
    PROD_DB_PORT=
    PROD_DB_USERNAME=
    PROD_DB_PASS=
    TEST_DB_HOST=
    TEST_DB_PORT=
    TEST_DB_USERNAME=
    TEST_DB_PASS=
    REDIS_HOST=
    REDIS_PASS=
    ```
2. `go mod download` to download dependencies.
3. `go run main.go` to start server.
4. Server can be access from `localhost:8888`

<h2> Deployment </h2>
<h3> Default Deployment </h3>

Automated deployment to GCP App Engine has been configured, and will be triggered on every push to `master`

<h3> Manual Deployment </h3>

To deploy manually, `gcloud app deploy`

<h2> Automated Test Suite </h2>

<h3> Default Test Run </h3>

Automated test suite will be triggered on every push or PR to `master`

<h3> Manual Test Run</h3>

To run locally, `go test -v ./...`

<h2> Resources </h2>

1. Server: Google Cloud Platform App Engine
2. Database: AWS RDS
3. Cache: Redis Cloud

<h2> Endpoints </h2>

| Service   | Method | Endpoint                       | 
|-----------|--------|--------------------------------|
| Affiliate | POST   | api/v1/affiliate/list          |
|           | GET    | api/v1/affiliate/:id           |
|           | GET    | api/v1/affiliate/info          |
|           | POST   | api/v1/affiliate/stats         |
|           | POST   | api/v1/affiliate/trend         |
|           | GET    | api/v1/affiliate/ranking/list  |
| Referral  | POST   | api/v1/referral/list           |
|           | POST   | api/v1/referral/stats          |
|           | POST   | api/v1/referral/trend          |
|           | POST   | api/v1/referral/recent/list    |
|           | GET    | api/v1/referral/:id            |
| Booking   | POST   | api/v1/booking/list            |
| Landing   | GET    | api/v1/booking/slots/available |
| Tracking  | POST   | api/v1/tracking/click          |

