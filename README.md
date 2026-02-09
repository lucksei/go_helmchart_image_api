# Go Chart Image Analyzer API

Exercise for an API that accepts a link to a helm chart, searches for container images. Downloads the docker images from their respective repositories and return a response of the list of images, their size and no. of layers in each image.

Built with Go using Gin Framework.

## Usage

Run

```sh
go run ./cmd/main
```

Build & run

```sh
go build -o ./bin/main ./cmd/main
./bin/main
```

## API Endpoints

### GET /health

Description: Healthcheck for the API

#### Responses

Returns 200 OK

```json
{
  "status": "ok"
}
```

### POST /api/helm-chart

Description: Reads and analyzes a helm chart. Can be a .tgz URI, an oci:// URI, or from a helm repo. Returns 202 if the analysis is in progress or 303 if the analysis is complete, redirects to /api/helm-chart/:id if the analysis is complete

#### Example requests:

- `repo_url` is an optional field that specifies the helm repo to download the chart from. Similar to the `helm repo add` command
- `chart_ref` is a required field that specifies the helm chart to analyze (Can be a .tgz URI, an oci:// URI, or from a helm repo)

For the example helm chart

```
POST http://localhost:8080/api/helm-chart
Content-Type: application/json
```

```json
{
  "repo_url": "https://helm.github.io/examples",
  "chart_ref": "hello-world"
}
```

For the nginx Bitnami chart

```
POST http://localhost:8080/api/helm-chart
Content-Type: application/json
```

```json
{
  "chart_ref": "oci://registry-1.docker.io/bitnamicharts/nginx"
}
```

For the redis Bitnami chart

```
POST http://localhost:8080/api/helm-chart
Content-Type: application/json
```

```json
{
  "chart_ref": "oci://registry-1.docker.io/bitnamicharts/redis"
}
```

#### Responses

Returns 202 if the analysis is in progress

```
HTTP/1.1 202 Accepted
Location: /api/helm-chart/eyJyZXBvX3VybCI6Imh0dHBzOi8vaGVsbS5naXRodWIuaW8vZXhhbXBsZXMiLCJjaGFydF9yZWYiOiJoZWxsby13b3JsZCJ9
```

Returns 303 if the analysis is complete

```
HTTP/1.1 303 SeeOther
Location: /api/helm-chart/eyJyZXBvX3VybCI6Imh0dHBzOi8vaGVsbS5naXRodWIuaW8vZXhhbXBsZXMiLCJjaGFydF9yZWYiOiJoZWxsby13b3JsZCJ9
```

### GET /api/helm-chart/:id

Description: Returns the analysis results for the helm chart. Returns 200 if the analysis is complete, 202 if the analysis is in progress

#### Requests

```
GET http://localhost:8080/api/helm-chart/:id
```

The `id` field is a base64 encoded string of the helm chart source with the fields from the helm_chart_source struct

```json
{
  "repo_url": "",
  "chart_ref": ""
}
```

#### Responses

Returns 200 OK

```json
{
  "repo_url": "https://helm.github.io/examples",
  "chart_path": "hello-world",
  "images": [
    {
      "name": "index.docker.io/library/nginx:1.16.0",
      "size": 44815103,
      "no_of_layers": 3
    }
  ]
}
```

Returns 200 OK if the analysis is in progress

```json
{
  "status": "Analysis of the helm chart is still in progress"
}
```

Returns 404 Not Found if the analysis is not found

```json
{
  "status": "Analysis not found. It has to be processed first by the POST /api/helm-chart endpoint"
}
```
