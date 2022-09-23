# List Manager App

Um projeto para gerenciamento de listas de ítens.

## Escopo
Desenvolver uma API Rest que forneça endpoints para gerenciamento de cadastro de usuários, listas e ítens. 

**Tecnologias utilizadas :** Linguagem de programação GO e banco de dados relacional MySQL.

## Atividades
**Estruturação banco de dados:** definição das entidades do banco de dados e seus relacionamentos.

**Warmup do projeto:** inicialização do projeto GO, com estrutura inicial e inicialização de container com MySQL com entidades criadas.

**Desenvolvimento endpoint /users:** desenvolvimento das rotas para cadastro de usuário

**Desenvolvimento endpoint /lists:** desenvolvimento das rotas para cadastro de listas

**Desenvolvimento endpoint /items:** desenvolvimento das rotas para cadastro de ítens

**Desenvolvimento endpoint /authenticate:** desenvolvimento da rota para autenticação com token JWT

**Dockerização do projeto:** preparar o projeto para execução completa via docker - para testes em ambiente local

**Documentação:** documentação sobre execução e utilização do projeto

**Testes unitários:** desenvolvimento dos testes unitários da aplicação


## Execução

Para executar o projeto, execute o seguinte comando na pasta raiz:

    $ [sudo] docker-compose up -d

Ou, caso possua o make instalado:

    $ make run

O comando acima executará um container para o MySQL e um para a aplicação. Para acessar os endpoints seguros, é necessário efetuar o login através do endpoint http://localhost:8080/api/v1/authenticate, com o seguinte payload JSON:

    {
        "login": "admin",
        "password": "admin"
    }

Para fazer login via SSO, é necessário chamar o endpoint http://localhost:8080/api/v1/authenticate/sso com o seguinte payload:


    {
        "login": "admin",
        "app_token": "TOKEN_JWT"
    }

O Token deve ser um JWT HS256, não pode estar expirado e deve conter o seguinte objeto:

    {
        "id": "1", // id do usuario
        "login": "admin", // login do usuario
        "email": "admin@admin.com", // email do usuario
        "exp": 1664159645 // expiracao do token
    }

O token gerado na resposta dos endpoints de autenticação deve ser utilizado ao chamar os endpoints privados via Authorization header:

    Authorization: <TOKEN_JWT>

Os demais endpoints seguem o que foi definido no detalhamento do projeto, sendo os seguintes:

	POST /api/v1/users --> Criação de usuários (private)
	GET /api/v1/users/{id} --> Obter usuário (private)
	PUT /api/v1/users/{id} --> Atualizar usuário (private)

	POST /api/v1/lists --> Criação de lista (public)
	GET /api/v1/lists/{list_id} --> Obter lista (public)
	DELETE /api/v1/lists/{list_id} --> Deletar lista (private)

	POST /api/v1/lists/{list_id}/items --> Salvar item na lista (private)
	GET /api/v1/lists/{list_id}/items --> Obter itens da lista (private)
	PUT /api/v1/lists/{list_id}/items/{item_id} --> Atualizar item da lista (private)
	DELETE /api/v1/lists/{list_id}/items/{item_id} --> Deletar item da lista (private)

Para parar os contêineres, execute

    $ [sudo] docker-compose down

Ou, caso possua o make instalado:

    $ go generate ./...

A execução dos testes unitários só pode ser realizada caso possua o GO instalado na máquina. Para isso, são é necessária a criação dos mocks primeiro, a partir do comando:

    $ make clean

Ou, caso possua o make instalado:

    $ make generate-mocks

## Demonstração

Link do video de demonstração: https://youtu.be/kSiFGiWWg4A
