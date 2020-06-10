# Super API

[![Go Report Card](https://goreportcard.com/badge/github.com/Dr-Pudim/super-api)](https://goreportcard.com/report/github.com/Dr-Pudim/super-api)

API de super heróis e vilões que usa a [SuperHero API](https://www.superheroapi.com/) como base de dados.

## Requisitos

* [Go](https://golang.org/) versão **>= 1.14**
* [BuffaloCLI](https://gobuffalo.io/en/docs/getting-started/installation/)
* Um banco de dados instalado

## Configuração de Ambiente

É necessario ter as sequintes variaveis de ambiente:

* *SUPERHEROAPI_ACCESS_TOKEN* - Com o token para acessar a [SuperHero API](https://www.superheroapi.com/)
* *SUPER_API_KEY* - Com a key de autenticação de acesso

### Variaveis Opcionais

* *TEST_DATABASE_URL* - Se for usar uma url diferente da padrão para o DB de teste
* *DATABASE_URL* - Se for usar uma url diferente da padrão para o DB da produção**(Recomendado)**

### Banco de Dados

A configuração padrão do DB de desenvolvimento espera que você esteja usando [PostGreSQL](https://www.postgresql.org/) e tenha um usuario *postgres* cuja senha é *postgres*.

As configurações de banco de dados ficam em *database.yml*. Para ver as opções dispoveis visite a [documentação do gobuffalo](https://gobuffalo.io/en/docs/db/configuration/).

## Criar Banco de Dados

Com o *database.yml* devidamente configurado e seu banco de dados rodando, execute

	$ buffalo pop create -a

E aplique as migrações para atualizar o schema

 $ buffalo db migrate

Esse comando aplica as migrações no DB de desenvolvimento, para aplicar na produção execute

	$ buffalo db migrate -e production

## Testes

Os testes podem ser feitos atraves da Buffalo CLI usando o comando

	$ buffalo test

Para executar um teste especifico execute

	$ buffalo test -m "Nome_do_Teste"

## Iniciando a API

Para iniciar a API em modo de desenvolvimento:

	$ buffalo dev

Você pode interagir com a API em [http://127.0.0.1:3000](http://127.0.0.1:3000)

## Usando a API

Listar todos os supers cadastrados:

> /

ou

> /all

Cadastrar supers:

> /{chave}/add?name={nome}

Onde {chave} é o valor de SUPER_API_KEY e {nome} é o nome do super para ser cadastrado, a API procura por personagens na [SuperHero API](https://www.superheroapi.com/) que contem o valor de {nome}. Retornas os supers cadastrados.

Para buscas supers:

> /{chave}/search?{campo}={valor}

Onde {chave} é o valor de SUPER_API_KEY, {campo} o campo alvo da busca e {valor} o valor para buscar

Os campos disponiveis para busca:

* uuid
* intelligence
* strength
* speed
* durability
* power
* combat
* min_height
* max_height
* min_weight
* max_weight
* image
* name
* full_name
* alias
* group
* relative
* occupation
* gender
* race
* eye_color
* hair_color

É possivel fazer uma pesquisa com multiplos campos usando **&**

> /{chave}/search?{campo1}={valor1}&{campo2}={valor2}&{campo3}={valor3}

Listar todos os herois cadastratos

> /heros

Listar todos os vilões cadastrados

> /villains

Remover um super cadastrado

> /{chave}/destroy?super_id={id}

Onde {chave} é o valor de SUPER_API_KEY e {id} o UUID do super a remover
