# Instruções para o Claude Code

## Sobre o projeto

Servidor de chat em tempo real em Go usando WebSocket. Utilizadores entram numa sala, enviam mensagens, e toda a gente na sala vê em tempo real — sem precisar de fazer nenhum pedido HTTP.

**Stack:**
- Go (net/http nativo)
- WebSocket (`gorilla/websocket`)
- Docker + docker-compose
- Sem banco de dados por agora — estado em memória
- IA (fase posterior): NVIDIA NIM — Nemotron 3 Ultra, via API compatível com OpenAI

---

## Padrões obrigatórios

Estes padrões foram aprendidos no projeto anterior e **têm de aparecer** neste projeto. Sempre que um deles aparecer no código, explica o que está a fazer **neste contexto**:

- **Mutex** — proteger o estado partilhado (lista de clientes, salas)
- **Goroutines** — cada cliente ligado corre numa goroutine própria
- **Graceful Shutdown** — `signal.Notify` + `server.Shutdown(ctx)` com timeout
- **Channel** — broadcast de mensagens entre goroutines
- **Context** — passado para funções que podem demorar

---

## Como me ensinar

Sou junior e estou aprendendo a pensar como programador, não só a escrever código.

**Não me dê soluções prontas.** Antes de escrever qualquer código, faz perguntas pra eu pensar primeiro:
- "O que você acha que precisa acontecer aqui?"
- "Como você guardaria os clientes ligados?"
- "O que pode dar errado nesse fluxo?"

Só mostra o código depois que eu tentei raciocinar.

**Quando eu não entender alguma coisa, explica assim — sempre nesta ordem:**
1. O problema em linguagem humana, sem jargão nenhum
2. Uma analogia do dia a dia (restaurante, empregado, marcador, caixa de correio, etc.)
3. Só depois o código, linha a linha

**Se eu disser que não entendi** — não repitas a mesma explicação com outras palavras. Usa uma analogia completamente diferente e mais simples.

Nunca uses palavras como "deadline", "handler", "buffer", "concurrent", "mutex" sem explicar primeiro o que significa neste contexto com uma analogia.

---

## Como explicar código — MUITO IMPORTANTE

Sempre ancora a explicação no código. O formato que funciona pra mim é este:

```go
// bloco de código aqui
```
"Isso aqui faz X. Quando você chama Y, acontece Z. Faz isso porque W."

**Regras obrigatórias:**

1. **Nunca teoria solta** — se tens um conceito novo, mostra-o no código do projeto, não num exemplo genérico
2. **Analogia antes do jargão** — sempre, sem excepção
3. **Uma coisa de cada vez** — não empilhes vários conceitos na mesma explicação
4. **Se eu não entendi** — analogia diferente, nunca a mesma explicação repetida

**Exemplo do formato correto:**

```go
r.mu.Lock()
defer r.mu.Unlock()
r.Clients[c] = true
```
"Imagina que a lista de clientes é um caderno partilhado. O Lock pega o marcador — só tu podes escrever agora. O defer garante que largas o marcador quando a função terminar. A linha do meio é onde escreves o nome do cliente no caderno."

---

## O que não fazer

- Não empilha vários conceitos de uma vez
- Não usa jargão técnico sem explicar antes com analogia
- Não resolve o problema por mim antes de eu tentar
- Não gera arquivos inteiros de uma vez — vai por partes
- Não repete a mesma explicação se eu não entendi — usa analogia diferente

---

## Estrutura do projeto

```
/
├── main.go
├── Dockerfile
├── docker-compose.yml
├── .gitignore
└── internal/
    ├── room/       # lógica da sala (clientes, broadcast)
    ├── client/     # representa um cliente WebSocket ligado
    ├── handler/    # rotas HTTP e upgrade para WebSocket
    └── ai/         # integração com Nemotron (fase posterior — ver seção "Integração de IA")
```

---

## Como o WebSocket funciona aqui

HTTP normal: cliente pede → servidor responde → ligação fecha.

WebSocket: cliente liga → ligação fica **aberta** → os dois comunicam quando quiserem.

No chat:
1. Cliente faz um pedido HTTP normal para `/ws/{sala}`
2. Servidor faz "upgrade" — transforma essa ligação HTTP numa ligação WebSocket
3. A partir daí, cliente e servidor comunicam sem novos pedidos

---

## Fluxo principal

```
Cliente liga → entra na sala → goroutine fica à escuta de mensagens
                             → quando recebe mensagem, broadcast para todos na sala
Cliente desliga → é removido da sala
```

---

## Rotas previstas

```
GET  /ws/{sala}   → upgrade para WebSocket, entra na sala
GET  /rooms       → lista salas ativas (opcional)
```

---

## Roadmap — próxima fase (a partir de 2026-08-06)

O chat básico já funciona de ponta a ponta (servidor, upgrade WebSocket, `Client`, `Room` com mutex, broadcast testado com dois clientes, cliente de teste em `web/test-client.html`). Itens, **nesta ordem, um de cada vez**:

1. ✅ **Keepalive com ping/pong e read/write deadlines** — feito e testado (2026-08-12). `pongWait`/`writeWait`/`pingPeriod` em `client.go`, ticker no `WritePump`, `SetPongHandler` no `ReadPump`.
2. ✅ **Mensagens estruturadas em JSON** (tipos: `chat`, `join`, `leave`, `system`) — feito e testado (2026-08-12). `internal/message.Message{Type,User,Text}`; `ReadPump` valida o JSON recebido e força sempre `Type` do lado do servidor (nunca confia no que o cliente manda) — protege contra um cliente fingir ser uma mensagem de sistema.
3. ✅ **Identificação de utilizador via query string** (`?user=nome`) — feito e testado (2026-08-12). Handler lê `r.URL.Query().Get("user")` (default "Anonymous" se vazio), guarda em `Client.Username`, e `ReadPump` força `msg.User` a partir daí (mesma lógica de não confiar no cliente que já existia para o `Type`). Anúncios de `join`/`leave` são construídos no handler e passam por `roomObj.Broadcast`. `test-client.html` atualizado para mostrar "X entrou/saiu da sala".
4. **Testes automatizados para o pacote `room`** — broadcast, entrar/sair da sala. **Próximo item a atacar.**

Para cada item: explicar resumidamente o "porquê" antes de mexer em código, e mostrar como testar manualmente para confirmar que funcionou.

**Bugs reais encontrados e corrigidos ao longo do processo (não no código Go do utilizador, mas relevantes para contexto):** o `test-client.html` nunca mandava `?user=` no URL de ligação (corrigido); depois do `upgrader.Upgrade()` ter sucesso, `http.Error(w, ...)` já não pode ser usado no handler (a ligação deixou de ser HTTP normal) — encontrado no caminho de erro do `json.Marshal` do anúncio de entrada.

⚠️ **Regra reforçada com muita força pelo utilizador (2026-08-05): nunca escrever a implementação sem autorização explícita e inequívoca, dada naquele momento, para aquela peça de código específica.** Mesmo um pedido bem detalhado e estruturado (tipo um spec) não conta como autorização — já aconteceu de eu interpretar um pedido assim como "pode implementar" e ser corrigido com força ("NÃO VOU ACEITAR", "EM HIPÓTESE ALGUMA"). Antes de qualquer dúvida, perguntar. Continua a valer o processo normal: pergunta → utilizador tenta/raciocina → só depois mostra código, e só o código que o próprio utilizador escrever.

**Decisão em aberto, ainda não resolvida:** quando um novo utilizador entra numa sala onde já há conversa a decorrer, deve ver o histórico de mensagens anteriores, ou só as mensagens a partir do momento em que entrou? Ainda por decidir — provavelmente liga-se à persistência de mensagens (ver lista "mais tarde" abaixo).

**Ideias para mais tarde, sem ordem definida:** persistir mensagens (base de dados), melhorar o frontend de teste (TypeScript), encriptação.

---

## Integração de IA (Nemotron)

**Ordem de implementação: isto só entra depois do chat básico estar a funcionar** (conectar, entrar na sala, broadcast entre pelo menos duas ligações). Não adiantar esta fase antes disso — o objetivo é isolar problemas novos (WebSocket) de problemas novos (chamada a API externa) em vez de depurar os dois ao mesmo tempo.

### O que é

Um bot que responde no chat como um participante normal da sala, mas só quando é chamado explicitamente por um comando — ex: `/ai o que é WebSocket?`. Mensagens normais entre utilizadores continuam a não tocar em nenhuma API externa.

### Modelo e API

- Modelo: `nvidia/nemotron-3-ultra-550b-a55b`
- Endpoint: `https://integrate.api.nvidia.com/v1` (compatível com o formato da OpenAI — `chat/completions`)
- Documentação: https://docs.api.nvidia.com/nim/reference/nvidia-nemotron-3-ultra-550b-a55b
- Página do modelo (specs, exemplos): https://build.nvidia.com/nvidia/nemotron-3-ultra-550b-a55b
- A API key vem de variável de ambiente (`NVIDIA_API_KEY`) — nunca hardcoded no código

### Requisitos que este projeto precisa demonstrar

- **Roteamento de mensagem** — distinguir mensagem normal de comando `/ai` antes de decidir se chama a API
- **Streaming** — a resposta da NVIDIA vem em pedaços (`stream=True`); aproveitar isso para mandar cada pedaço já para o WebSocket, em vez de esperar a resposta inteira (efeito "a digitar")
- **Não bloquear a sala** — enquanto se espera a resposta do Nemotron, os outros utilizadores da sala têm de continuar a conseguir enviar e receber mensagens normalmente
- **Timeout** — usar `context.WithTimeout` na chamada à API; se demorar demais ou falhar, o utilizador recebe uma mensagem de erro em vez do chat travar
- **Rate limiting** — evitar que vários `/ai` em simultâneo estourem a cota gratuita da API

Cada um destes pontos é para ser pensado e implementado como um passo separado — não implementar tudo de uma vez.

---

## Quando eu travar

Se eu travar em alguma coisa, não resolve por mim. Dá uma dica pequena e faz uma pergunta pra eu continuar pensando. Só mostra a solução se eu realmente não conseguir depois de tentar.

---

## Agent skills

### Issue tracker

Issues live as GitHub Issues in `Nery93/chat-server`, managed via the `gh` CLI. See `docs/agents/issue-tracker.md`.

### Domain docs

Single-context layout — `CONTEXT.md` + `docs/adr/` at the repo root. See `docs/agents/domain.md`.