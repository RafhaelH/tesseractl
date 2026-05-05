# tesserato

Uma CLI em Go construída sobre o Docker e o Docker Compose, feita para centralizar a operação do dia a dia de múltiplos projetos containerizados rodando em uma mesma máquina (tipicamente um VPS).

Se você administra vários projetos baseados em Docker — cada um com seu próprio `docker-compose.yml` espalhado por diretórios diferentes — e se vê digitando comandos longos como `docker compose -f ... up -d --build` o tempo todo, o `tesserato` reduz esse fluxo a comandos curtos e cientes do projeto, dirigidos por um único arquivo YAML.

## Funcionalidades

| Comando | O que faz |
| --- | --- |
| `tesserato status` | Lista os containers em execução, agrupados pelos projetos definidos na sua configuração. Cai num modo de listagem simples quando não há configuração disponível. |
| `tesserato up <projeto> [--build]` | Roda `docker compose up -d` (opcionalmente com `--build`) dentro do diretório do projeto. |
| `tesserato down <projeto> [--volumes]` | Para e remove os containers do projeto via `docker compose down`. |
| `tesserato logs <projeto> <serviço> [-f] [--tail N]` | Mostra os logs de um único serviço, resolvido pelo mapa nome-lógico → nome-do-container do config. |
| `tesserato health` | Para cada serviço declarado em todos os projetos, reporta `Up`, `Exited (código) ...` ou `ausente`. |
| `tesserato clean [--all] [-y]` | Remove containers parados e imagens dangling. Com `--all`, roda o mais agressivo `docker system prune -a --volumes`. Pede confirmação salvo se passar `-y`. |
| `tesserato deploy <projeto> [--no-pull] [--branch X]` | Fluxo padrão de deploy em produção: `git checkout <branch>` (opcional), `git pull --ff-only` e `docker compose up -d --build`. Aborta rapidamente em erro do git. |
| `tesserato --version` | Imprime a versão do binário (padrão `dev`; pode ser injetada no momento do build). |

A flag global `--config <caminho>` sobrescreve a descoberta automática do arquivo de configuração para qualquer comando.

## Instalação

### Pré-requisitos

- Go 1.26+ (`go version` para conferir)
- Docker e Docker Compose (a CLI delega para eles)
- Git (necessário apenas para `tesserato deploy`)

### Build e instalação

```powershell
go install github.com/RafhaelH/cli_go/cmd/tesserato@latest
```

Esse comando coloca o binário em `$(go env GOPATH)\bin\tesserato.exe` (Windows) ou `$(go env GOPATH)/bin/tesserato` (macOS/Linux), que o instalador do Go já adiciona ao seu `PATH`.

Para buildar a partir do clone:

```powershell
git clone https://github.com/RafhaelH/cli_go.git
cd cli_go
go install ./cmd/tesserato
```

### Build com versão embutida

A versão reportada por `tesserato --version` é injetada via `-ldflags`:

```powershell
go install -ldflags "-X 'github.com/RafhaelH/cli_go/internal/cli.Version=v0.1.0'" ./cmd/tesserato
```

Combine com `git describe` num pipeline de release para taggar binários automaticamente.

## Configuração

O `tesserato` é dirigido por um arquivo YAML. Por padrão ele procura, nesta ordem:

1. O caminho passado via `--config <caminho>`.
2. `./tesserato.yaml` ou `./tesserato.yml` no diretório atual.
3. `~/.config/tesserato/config.yaml`.

Se nada for encontrado e o comando não precisar do config (ex.: `clean`, ou `status` sem args caindo no modo simples), a CLI roda mesmo assim. Caso contrário, ela imprime um erro acionável explicando onde procurou.

### Exemplo de configuração

```yaml
projects:
  portal_astech:
    path: /home/user/projects/portal_astech
    compose_file: docker-compose.prod.yml
    services:
      backend: portal_astech_backend
      postgres: portal_astech_postgres
      redis: portal_astech_redis

  knowledgeai:
    path: /home/user/projects/knowledgeai
    compose_file: docker-compose.yml
    services:
      api: knowledgeai_api
      frontend: knowledgeai_frontend
      db: knowledgeai_db
```

Schema:

| Campo | Obrigatório | Descrição |
| --- | --- | --- |
| `projects` | sim | Mapa de nível mais alto, indexado pelo nome lógico do projeto que você vai passar na CLI. |
| `projects.<nome>.path` | sim | Caminho absoluto do diretório do projeto (onde fica o compose). |
| `projects.<nome>.compose_file` | não | Nome do arquivo compose. Se omitido, o `docker compose` decide o default. |
| `projects.<nome>.services` | sim | Mapa do nome lógico do serviço → nome real do container. O nome real é o que o `docker ps` imprime na coluna `NAMES`. |

O mapa `services` é o que possibilita as visões agrupadas por projeto: `tesserato status` e `tesserato health` correlacionam a saída do `docker ps` com esse mapa para saber qual container pertence a qual projeto.

## Uso

Com o binário instalado e a configuração no lugar:

```powershell
# O que está rodando, por projeto?
tesserato status

# Sobe um projeto
tesserato up portal_astech --build

# Acompanha os logs de um serviço
tesserato logs portal_astech backend -f --tail 100

# Para um projeto
tesserato down portal_astech

# Deploy em produção
tesserato deploy portal_astech

# Snapshot de saúde de todos os projetos
tesserato health

# Libera espaço em disco (com confirmação)
tesserato clean

# Limpeza agressiva
tesserato clean --all -y
```

Passe `--config /etc/tesserato.yaml` (ou qualquer caminho) para qualquer comando para sobrescrever a descoberta automática.

## Estrutura do projeto

```
cli_go/
├── cmd/
│   └── tesserato/
│       └── main.go              Ponto de entrada fino
├── internal/
│   ├── cli/
│   │   ├── root.go              Comando raiz cobra, flags globais, versão
│   │   ├── status.go            tesserato status
│   │   ├── up.go                tesserato up
│   │   ├── down.go              tesserato down
│   │   ├── logs.go              tesserato logs
│   │   ├── clean.go             tesserato clean (+ helper de confirmação)
│   │   ├── health.go            tesserato health
│   │   ├── deploy.go            tesserato deploy
│   │   ├── helpers.go           lookupProject, indexByName, etc.
│   │   └── print.go             printPerProject (helper compartilhado de tabela)
│   ├── config/
│   │   ├── config.go            Schema YAML, Load, Resolve, ErrNotFound
│   │   └── config_test.go
│   ├── docker/
│   │   ├── docker.go            ListContainers, ListAll, RunCompose, RunDocker
│   │   └── docker_test.go
│   └── proc/
│       └── proc.go              Runner genérico de processo com streaming
├── docker-compose.yml           Compose de exemplo para testes locais
├── tesserato.yml                Configuração tesserato de exemplo
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

A separação entre `cmd/` e `internal/` segue o layout padrão de Go: `cmd/<nome>/main.go` por binário, `internal/` para código que **não pode** ser importado de fora deste módulo (regra imposta pelo próprio compilador).

## Desenvolvimento

### Rodar a partir do código

```powershell
go run ./cmd/tesserato status
```

### Testes

```powershell
go test ./...
go test -v ./...        # verboso, mostrando cada subtest
```

A suíte atual cobre o pacote `internal/config` (casos de Load + sentinel do Resolve) e o pacote `internal/docker` (o parser `parseContainers`, exercitado através de um `io.Reader` injetado para não exigir Docker rodando). Ambos os pacotes usam testes table-driven com `t.Run` para subtests — o padrão idiomático em Go.

### Verificações estáticas

```powershell
gofmt -l ./internal ./cmd     # lista arquivos fora do padrão
go vet ./...                  # linter nativo
```

### Reinstalação após mudanças

```powershell
go install ./cmd/tesserato
```

Pode rodar a qualquer momento após editar — o `go install` reconstrói e substitui o binário no `PATH` em segundos.

## Roadmap

Pronto:

- Fase 1 (MVP): `status`, `up`, `down`, `logs`.
- Fase 2: `clean`, `health`.
- Fase 3: `deploy` (git pull + rebuild).
- Polimento: binário real via `go install`, versão injetada via ldflags, mensagens de erro multilinhas com sugestão de projetos disponíveis, primeiros testes, refactor para `internal/proc` e `printPerProject` compartilhado.

Planejado:

- Inspeções concorrentes em `health` usando goroutines + `errgroup`.
- Testes ponta a ponta para o pacote `cli` usando `cobra.Command.SetArgs` e uma camada Docker mockável.
- CI no GitHub Actions: testes, `go vet`, `golangci-lint`, e binários cross-platform com versão embutida.
- Saída colorida / estilizada para `status` e `health` (provavelmente `fatih/color` ou `lipgloss`).
- Flag `--json` em `status` e `health` para consumo via shell scripts.

## Por que "tesserato"?

Um *tesseract* é o análogo em quatro dimensões de um cubo — uma estrutura que contém vários cubos num único objeto coerente. A ideia se aplica aqui: muitos projetos Docker, um único lugar para operá-los.
