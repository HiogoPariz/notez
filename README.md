# Notez

Monorepo do notez: aplicação de notas com editor WYSIWYG, autenticação via Google
e armazenamento de metadados (Postgres) + conteúdo (filesystem).

## Estrutura

```
notez/
├── apps/
│   ├── api/     API REST em Go (Gin) — metadados das notas no Postgres
│   ├── files/   Serviço de arquivos em Go (Gin) — conteúdo das notas em .files/
│   └── web/     Front-end em SvelteKit (editor Editor.js + Auth.js/Google)
├── docker-compose.yml   Postgres local
├── Makefile             orquestração da stack
└── go.work              workspace Go (api + files)
```

## Como rodar

Pré-requisitos: Go 1.21+, Bun, Docker (com o plugin `docker compose`) e `goose`
(opcional, apenas para criar migrations).

```bash
# 1. cria os arquivos .env a partir dos exemplos
make setup

# 2. sobe o Postgres
make database

# 3. instala as dependências (go mod + bun install)
make install

# 4. roda api + files + web juntos
make dev
```

Depois, acesse o front em <http://localhost:5173>. As migrations do banco rodam
automaticamente quando a API inicia.

## Comandos úteis

| Comando | Descrição |
| --- | --- |
| `make dev` | Roda api + files + web juntos |
| `make dev-api` / `make dev-files` / `make dev-web` | Roda um serviço isolado |
| `make build` | Compila api, files e o front |
| `make test` | Roda os testes dos serviços Go |
| `make fmt` | Formata o código Go |
| `make migration name=criar_tabela` | Cria uma migration (goose) |
| `make down` | Derruba o Postgres |
| `make help` | Lista todos os comandos |

## Configuração

- `.env` (raiz): variáveis dos serviços Go (`api`, `files`).
- `apps/web/.env`: variáveis do front (`API_URL`, credenciais do Google, `AUTH_SECRET`).

Ambos são criados por `make setup` a partir dos respectivos `.env.example`.
Os arquivos `.env` reais **não** são versionados.
