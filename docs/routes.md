# Rotas e Contratos do Gateway

## Estratégia

O gateway expõe rotas por contexto de domínio, e não por porta. O cliente sempre conversa com o Traefik, e o Traefik encaminha para os serviços internos.

Nesta base, o Traefik usa apenas o `file provider`. Isso evita acoplamento com versões antigas do Docker Engine e mantém o gateway estável mesmo sem discovery automático.

Rotas iniciais:

- `/brain/*` -> `brain:8081`
- `/watcher/*` -> `watcher:8082`

O Traefik remove os prefixos `/brain` e `/watcher` antes de entregar a request ao serviço. Assim, um handler interno `GET /health` no `brain` pode ser acessado externamente como `GET /brain/health`.

## BFFs

Os BFFs podem ficar no mesmo repositório e também serem publicados pelo Traefik. A convenção recomendada é:

- `/bff/admin/*` -> `bff-admin:8091`
- `/bff/mobile/*` -> `bff-mobile:8092`

Cada BFF pode consumir `brain` e `watcher` pela rede Docker usando os nomes dos serviços:

- `http://brain:8081`
- `http://watcher:8082`

## Exemplo de serviço no compose

Exemplo para o `brain` quando ele existir no `docker-compose.yml`:

```yaml
  brain:
    build: ./services/brain
    container_name: arthemis-brain
    networks:
      - edge
    expose:
      - "8081"
```

Exemplo para o `watcher`:

```yaml
  watcher:
    build: ./services/watcher
    container_name: arthemis-watcher
    networks:
      - edge
    expose:
      - "8082"
```

As rotas continuam definidas em [traefik/dynamic/routes.yml](/home/andre/Documentos/Git/arthemis/arthemis-edge/traefik/dynamic/routes.yml).

## Exemplo de BFF no compose

```yaml
  bff-admin:
    build: ./bffs/admin
    container_name: arthemis-bff-admin
    networks:
      - edge
    expose:
      - "8091"
```

Depois, adicione o router/service correspondente no arquivo dinâmico do Traefik. Se o BFF precisar receber a rota interna sem o prefixo `/bff/admin`, crie mais um middleware `stripPrefix`, do mesmo jeito que foi feito para `brain` e `watcher`.
