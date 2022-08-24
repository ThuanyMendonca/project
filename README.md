![Logo](https://pubnative.net/wp-content/uploads/2018/01/Go.png)

# Transfer API

API de transferência de saldos entre usuários.

## Clonando o projeto

Clone o projeto

```bash
  git clone https://github.com/ThuanyMendonca/project.git
```

## Variáveis de Ambiente

Para rodar esse projeto, você vai precisar adicionar as seguintes variáveis de ambiente no seu .env

`DB_HOST`
`DB_USER`
`DB_NAME`
`DB_PASSWORD`
`DB_PORT`
`DB_TIME_ZONE`
`PORT`
`GIN_MODE`
`AUTHORIZATOR_URL`

Obs: Já configurado para rodar na sua máquina. 😉

## Executar script

É necessário rodar o script que está na pasta scripts para inserir os tipos de usuários

## Como executar o projeto:

Entre no diretório

```bash
  cd project
```

Inicie o docker

```bash
  docker-compose up
```

Rodar o projeto

```bash
  go run main.go
```

## Rodando os testes

Para rodar os testes, rode o seguinte comando

```bash
  go test ./... -coverprofile=coverage.out
  go tool cover -html=coverage.out
```

## Autores

- [@thuanymendonca](https://www.github.com/ThuanyMendonca)
