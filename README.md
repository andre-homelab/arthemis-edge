# arthemis-edge

API Gateway do ecossistema Arthemis usando Traefik como reverse proxy de entrada.

O gateway foi preparado para:

- rotear o microserviço `brain`
- rotear o microserviço `watcher`
- hospedar BFFs no mesmo repositório e na mesma rede Docker
- centralizar CORS e middlewares HTTP na borda
- funcionar sem depender do provider Docker do Traefik

## Portas definidas

- `brain`: `8081`
- `watcher`: `8082`
- faixa sugerida para BFFs: `8091+`
- dashboard do Traefik: `8080`
- entrada HTTP do gateway: `80`
- entrada HTTPS do gateway: `443`

## Rotas públicas

- `/brain/*` -> serviço `brain`
- `/watcher/*` -> serviço `watcher`
- `/bff/<nome>/*` -> BFF correspondente

Os prefixos `/brain` e `/watcher` são removidos pelo Traefik antes de encaminhar a request para o serviço interno.

## Subida local

```bash
docker compose up -d
```

O Traefik sobe sozinho e lê as rotas a partir dos arquivos em `traefik/dynamic`.
