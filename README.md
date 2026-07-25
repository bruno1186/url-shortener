# url-shortener

[![CI](https://github.com/bruno1186/url-shortener/actions/workflows/ci.yml/badge.svg)](https://github.com/bruno1186/url-shortener/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

Servico de encurtamento de URLs escrito em Go usando apenas a biblioteca padrao. Expoe uma API HTTP para criar codigos curtos e redirecionar para a URL original. O storage e definido por uma interface, permitindo trocar a implementacao (memoria, Redis, SQL) sem mudar a logica de negocio.

## Caracteristicas

- HTTP API com `net/http` e roteamento por metodo/rota do Go 1.22
- Geracao de codigos aleatorios com `crypto/rand`
- Validacao de URL (apenas http/https)
- Storage plugavel via interface (`MemoryStore` incluso, thread-safe)
- Testes unitarios e de integracao HTTP (`net/http/httptest`)
- CI (GitHub Actions): gofmt, go vet, build e testes com `-race`
- Imagem Docker multi-stage (distroless)

## Estrutura

```
cmd/server/main.go              # entrypoint (configuracao e bootstrap HTTP)
internal/shortener/
  store.go                      # interface Store + MemoryStore
  shortener.go                  # logica: validacao, geracao de codigo, resolve
  shortener_test.go             # testes de dominio
internal/server/
  server.go                     # handlers HTTP e roteamento
  server_test.go                # testes de API (httptest)
Dockerfile
.github/workflows/ci.yml
```

## Como rodar

```
go run ./cmd/server
```

O servico sobe em `:8080` por padrao. Configuravel via `ADDR` e `BASE_URL`.

## Testes

```
go test -race -cover ./...
```

## API

### Health

```
GET /healthz  ->  200 {"status":"ok"}
```

### Encurtar

```
POST /api/shorten
Content-Type: application/json

{"url":"https://example.com/uma/pagina"}
```

Resposta `201`:

```
{"code":"aB3xYz9","short_url":"http://localhost:8080/aB3xYz9"}
```

### Redirecionar

```
GET /{code}  ->  302 Location: <url original>
```

Exemplo com curl:

```
curl -s -X POST localhost:8080/api/shorten \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://go.dev"}'

curl -i localhost:8080/<code>
```

## Docker

```
docker build -t url-shortener .
docker run -p 8080:8080 url-shortener
```

## Licenca

MIT
