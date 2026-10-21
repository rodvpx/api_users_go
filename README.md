# API de usuários

API REST simples para criar, consultar, atualizar e remover usuários. A aplicação foi desenvolvida em Go, usa PostgreSQL para persistência e executa as migrations ao iniciar.

## Requisitos

- Go 1.26.5 ou compatível com a versão declarada em `go.mod`.
- PostgreSQL acessível localmente.
- Banco de dados `api_users_go_db` criado antes de iniciar a aplicação.

A conexão está configurada diretamente em `cmd/api/main.go`:

```text
postgres://postgres:postgres@localhost:5432/api_users_go_db?sslmode=disable
```

Altere essa configuração no código caso seu ambiente use outro endereço, usuário, senha ou banco.

## Executar

Na raiz do projeto, inicie a API:

```bash
go run ./cmd/api
```

Na inicialização, a aplicação conecta ao PostgreSQL, aplica as migrations em `migrations/` e inicia o servidor em `http://localhost:8080`.

Para executar os testes:

```bash
go test ./...
```

## Endpoints

| Método | Rota | Descrição |
| --- | --- | --- |
| `POST` | `/users` | Cria um usuário |
| `GET` | `/users` | Lista todos os usuários |
| `GET` | `/users/{id}` | Busca um usuário pelo ID |
| `PATCH` | `/users/{id}` | Atualiza parcialmente um usuário |
| `DELETE` | `/users/{id}` | Remove um usuário |

As respostas com conteúdo usam `Content-Type: application/json`. Os campos `name` e `email` são obrigatórios na criação. Na atualização, ambos são opcionais e somente os campos enviados são alterados. O banco limita cada campo a 100 caracteres e exige e-mail único.

Nos exemplos abaixo, crie uma requisição no Bruno, Postman ou Insomnia usando o método e a URL indicados. Para requisições com corpo, selecione `Body` > `JSON` (ou `raw` > `JSON`) e informe o conteúdo mostrado.

### Criar usuário

- **Método:** `POST`
- **URL:** `http://localhost:8080/users`
- **Header:** `Content-Type: application/json`
**Body (JSON):**

```json
{
	"name": "Ana Silva",
	"email": "ana@example.com"
}
```

Resposta: `201 Created`.

```json
{
	"id": 1,
	"name": "Ana Silva",
	"email": "ana@example.com"
}
```

### Listar usuários

- **Método:** `GET`
- **URL:** `http://localhost:8080/users`

Resposta: `200 OK`, com uma lista de usuários. Sem registros, a resposta é `[]`.

### Buscar usuário

- **Método:** `GET`
- **URL:** `http://localhost:8080/users/1`

Resposta: `200 OK` com o usuário encontrado; ID inválido retorna `400 Bad Request` e usuário inexistente retorna `404 Not Found`.

### Atualizar usuário

Envie apenas os campos que deseja alterar:

- **Método:** `PATCH`
- **URL:** `http://localhost:8080/users/1`
- **Header:** `Content-Type: application/json`
**Body (JSON):**

```json
{
	"name": "Ana Souza"
}
```

Resposta: `200 OK` com o usuário atualizado. Um objeto vazio (`{}`) não altera os campos.

### Remover usuário

- **Método:** `DELETE`
- **URL:** `http://localhost:8080/users/1`

Resposta: `204 No Content`, sem corpo. ID inválido retorna `400 Bad Request` e falha ao localizar/remover retorna `404 Not Found`.

## Erros e comportamento

- JSON malformado retorna `400 Bad Request`.
- Erros no cadastro, incluindo campos obrigatórios ausentes, retornam `400 Bad Request`.
- Erros ao consultar, atualizar ou remover por ID são respondidos como `404 Not Found` pelos handlers atuais.
- Falhas ao listar usuários retornam `500 Internal Server Error`.
- A validação da aplicação verifica se `name` e `email` estão vazios; não valida o formato do e-mail.

A API inclui middleware de log, recuperação de panic e CORS. O CORS permite qualquer origem e os métodos `GET`, `POST`, `PATCH`, `DELETE` e `OPTIONS`.

