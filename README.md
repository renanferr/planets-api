# Planets API

[![Coverage Status](https://coveralls.io/repos/github/renanferr/planets-api/badge.svg?branch=master)](https://coveralls.io/github/renanferr/planets-api?branch=master)[![Build Status](https://travis-ci.org/renanferr/planets-api.svg?branch=master)](https://travis-ci.org/renanferr/planets-api)

API para adicionar e listar planetas da franquia Star Wars.

Ao adicionar um novo planeta, a aplicação buscará na [SWAPI](https://swapi.dev/) quantas vezes o mesmo apareceu nos filmes.

O _design_ se baseia em alguns conceitos do [Domain-Driver Design](https://www.amazon.com.br/Domain-Driven-Design-Eric-Evans/dp/8550800651),
Eric Evans e da [Hexagonal Architecture](https://fideloper.com/hexagonal-architecture) de Alistair Cockburn.

## Componentes

### Business Core
- Listing
- Adding

### Adapters
- Storage
- HTTP
  - RESTful server
  - Client
  
## Dependências

### Database
[MongoDB](https://docs.mongodb.com/v4.2/)

### Ambiente de Desenvolvimento
[Docker](https://docs.docker.com/) e [docker-compose](https://docs.docker.com/compose/) ou [Go](https://golang.org/).

### Go Modules
O gerenciamento de pacotes é feito através de [Go Modules](https://blog.golang.org/using-go-modules) no arquivo `go.mod`. Os módulos utilizados como dependência são:

- [govalidator](github.com/asaskevich/govalidator) como _helper_ para validação e tratamento de erros de validação
- [chi](github.com/go-chi/chi) como roteador HTTP
- [mongo-driver](go.mongodb.org/mongo-driver) como driver do MongoDB

Para baixá-las com Go instalado na máquina:
```bash
$ go mod download
```

## Como usar

Este repositório provê um ambiente de execução Docker via [Dockerfile](Dockerfile) e [docker-compose](docker-compose.yml) com todas as dependências e configurações necessárias.

Além disso, possui uma especificação [OpenAPI 3](https://swagger.io/specification/)
 [openapi.yml](openapi.yml).

### Variáveis de Ambiente

| Nome                      | Descrição                                 |
|---------------------------|-------------------------------------------|
| PORT                        | Porta do servidor HTTP                    |
| APP_LOG_LEVEL             | Nível de log estruturado da aplicação     |
| DB_CONNECTION_URI         | URI de conexão com MongoDB                |
| DB_TIMEOUT_MS             | Tempo máximo de transação no MongoDB em *milissegundos*|
| PLANETS_API_BASE_URL      | URL base da API para buscar informações extras de planetas|

### Docker-Compose

Para executar via [docker-compose](https://docs.docker.com/compose/)

```bash
$ docker-compose up --build -d
```

## Live Demo

Uma demonstração da aplicação foi implantada no plano grátis do heroku em:
https://swapi-golang-rest-api.herokuapp.com

### Exemplos de Requisição

#### Adding Planet
```bash
$ curl --request POST \
  --url https://swapi-golang-rest-api.herokuapp.com/api/planets \
  --header 'content-type: application/json' \
  --data '{
  "name": "tatooine",
  "climate": "arid",
  "terrain": "desert"
}'
```

#### Getting Planet
```bash
$ curl --request GET \
  --url https://swapi-golang-rest-api.herokuapp.com/api/planets/5f99b620ce4200066e7efed4
```

#### Listing Planet
```bash
$ curl --request GET \
  --url 'https://swapi-golang-rest-api.herokuapp.com/api/planets?page=1&limit=5'
```

## Testes

### Testando
```bash
$ go test ./...
```

### Cobetura

| Package                                                          | Coverage  |
|------------------------------------------------------------------|-----------|
| github.com/renanferr/planets-api/pkg/adding	           | 100.0%    |
| github.com/renanferr/planets-api/pkg/http/client       | 86.1%     |
| github.com/renanferr/planets-api/pkg/http/rest         | 100.0%    |
| github.com/renanferr/planets-api/pkg/http/rest/planets | 80.3%     |
| github.com/renanferr/planets-api/pkg/listing	       | 100.0%    |
| github.com/renanferr/planets-api/pkg/storage/mongo	   | 67.4%     |


## Licença
 A aplicação está sob a licença [MIT](https://choosealicense.com/licenses/mit/)
 