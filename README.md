# Sistema de Gerenciamento de Usuários (CRUD)

Projeto full-stack para gerenciar usuários: é possível cadastrar,
consultar, editar e excluir registros. O backend foi desenvolvido em Go,
com PostgreSQL para armazenar os dados, e o frontend utiliza Vue 3.

## Tecnologias

-   **Backend:** Go, `net/http`, `pgx/v5` e `golang.org/x/crypto/bcrypt`
-   **Banco de dados:** PostgreSQL 16, executado com Docker Compose
-   **Frontend:** Vue 3 (Composition API), Vite e Bootstrap 5

## Estrutura do projeto

O código está separado em duas partes principais: `backend`, que
concentra a API e o acesso ao banco, e `frontend`, responsável pela
interface.

``` text
DESAFIO_DEVPOOL/
├── backend/
│   ├── cmd/api/main.go                  # Inicialização do servidor e definição das rotas
│   ├── internal/
│   │   ├── config/database.go           # Configuração do pool de conexões
│   │   ├── handler/user_handler.go      # Recebe as requisições e monta as respostas HTTP
│   │   ├── middleware/cors.go           # Configuração de CORS
│   │   ├── models/
│   │   │   ├── errors.go                # Erros utilizados pela aplicação
│   │   │   └── users.go                 # Estruturas de usuário e dados de entrada/saída
│   │   ├── repository/
│   │   │   ├── users/users.go           # Consultas SQL e operações com usuários
│   │   │   └── repositories.go          # Organização dos repositórios
│   │   └── service/services.go          # Validações, regras de negócio e hash de senhas
│   ├── migrations/001_create_users.sql  # Criação da tabela de usuários
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── index.html
│   ├── package.json
│   └── src/
│       ├── components/UserFormModal.vue # Formulário de cadastro e edição
│       ├── services/api.js              # Funções de comunicação com a API
│       ├── App.vue                      # Tela principal
│       └── main.js                      # Inicialização do Vue e do Bootstrap
├── .env                                 # Variáveis de ambiente
└── docker-compose.yml                   # Configuração do PostgreSQL
```

### Algumas decisões do projeto

-   **Separação por responsabilidades:** o fluxo da API passa pelo
    handler, service e repository. Assim, cada parte fica responsável
    por uma tarefa: receber a requisição, aplicar as regras da aplicação
    e conversar com o banco.
-   **Armazenamento de senhas:** as senhas não são salvas em texto puro.
    O backend gera um hash com BCrypt, e o campo de senha também é
    omitido das respostas em JSON por meio da tag `json:"-"`.
-   **Exclusão de usuários:** a operação `DELETE` remove o registro
    diretamente do banco. Como o projeto não prevê histórico ou
    auditoria de usuários excluídos, essa abordagem mantém a
    implementação mais simples.
-   **Chamadas do frontend:** as requisições ficam centralizadas em
    `api.js`. Esse módulo também define um timeout de 5 segundos e trata
    erros de conexão, evitando repetir essa lógica nos componentes.

## Pré-requisitos

Antes de executar o projeto, você precisa ter instalado:

-   Go 1.22 ou superior
-   Node.js 18 ou superior e `npm`
-   Docker e Docker Compose

## Como executar

### 1. Configurar as variáveis de ambiente

Na raiz do projeto, confira se o arquivo `.env` contém as variáveis
usadas pelo Docker Compose e pelo backend:

``` env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=devpool
DATABASE_URL=postgres://postgres:postgres@localhost:5432/devpool?sslmode=disable
```

Se estiver usando outros valores de usuário, senha ou nome do banco,
ajuste as variáveis de acordo.

### 2. Iniciar o PostgreSQL e criar a tabela

Na raiz do projeto (`DESAFIO_DEVPOOL/`), suba o container:

``` bash
docker compose up -d
```

Depois, execute a migration que cria a tabela `users`:

``` bash
docker exec -i <nome_do_container_postgres> psql -U postgres -d devpool < backend/migrations/001_create_users.sql
```

Substitua `<nome_do_container_postgres>` pelo nome do container criado
pelo Docker Compose. Se preferir, também é possível executar o arquivo
`backend/migrations/001_create_users.sql` usando um cliente SQL ou o
`psql` instalado na máquina.

### 3. Iniciar o backend

Em um terminal, entre na pasta `backend`, baixe as dependências e
execute a API:

``` bash
cd backend
go mod download
go run cmd/api/main.go
```

A API ficará disponível em `http://localhost:8080`.

### 4. Iniciar o frontend

Abra outro terminal, entre na pasta `frontend` e instale as
dependências:

``` bash
cd frontend
npm install
npm run dev
```

O Vite vai mostrar no terminal o endereço para acessar a aplicação. Por
padrão, é `http://localhost:5173`.

## API

A API utiliza a seguinte URL base:

`http://localhost:8080`

| Método | Rota | Descrição | Status de Sucesso |
| :--- | :--- | :--- | :--- |
| `POST` | `/users` | Cadastra um novo usuário | `201 Created` |
| `GET` | `/users` | Lista todos os usuários cadastrados | `200 OK` |
| `GET` | `/users/{id}` | Busca um usuário específico pelo UUID | `200 OK` |
| `PUT` | `/users/{id}` | Atualiza os dados e/ou senha de um usuário | `200 OK` |
| `DELETE` | `/users/{id}` | Remove um usuário pelo UUID | `204 No Content` |
  
### Exemplos de uso

Os exemplos abaixo usam `curl` e podem ser executados no terminal com a
API em funcionamento.

#### Cadastrar um usuário

``` bash
curl -i -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Maria Silva",
    "email": "maria.silva@exemplo.com",
    "login": "mariasilva",
    "password": "senhaSegura123",
    "active": true
  }'
```

#### Listar usuários

``` bash
curl -i http://localhost:8080/users
```

#### Buscar um usuário pelo ID

``` bash
curl -i http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000
```

#### Atualizar um usuário

Para atualizar somente os dados cadastrais, envie os campos desejados
sem `current_password` e `new_password`:

``` bash
curl -i -X PUT http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Maria Oliveira Silva",
    "email": "maria.silva@exemplo.com",
    "login": "mariasilva",
    "active": false
  }'
```

#### Excluir um usuário

``` bash
curl -i -X DELETE http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000
```
