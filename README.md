# API para cadastro de usuários - in memory

Esse código é parte dos meus estudos em Golang e consiste na criação de um CRUD (Create, Read, Update, Delete) de usuários com a persistência dos dados em memória.
A API foi construída utilizando o framework CHI e utilizei a lib sync para lidar com race conditions.


## Instalação

1. Clone o repositório:
    ```sh
    git clone https://github.com/seu-usuario/go-crud-in-memory.git
    ```
2. Navegue até o diretório do projeto:
    ```sh
    cd go-crud-in-memory
    ```
3. Instale as dependências:
    ```sh
    go mod tidy
    ```

## Uso

Para iniciar o servidor, execute o comando:
```sh
go run src/main.go
```

O servidor estará disponível em `http://localhost:8080`.

## Endpoints

- `GET /items` - Lista todos os itens
- `POST /items` - Cria um novo item
- `GET /items/{id}` - Obtém um item pelo ID
- `PUT /items/{id}` - Atualiza um item pelo ID
- `DELETE /items/{id}` - Deleta um item pelo ID


## Anexos

- Arquivo POSTMAN para testes do projeto.


