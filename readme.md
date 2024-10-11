# Cloudflare ip updater

Este projeto é um automatizador de atualização de IP para o Cloudflare, funcionando de maneira similar ao No-IP. Ele verifica mudanças no endereço IP local e atualiza automaticamente o DNS no Cloudflare para garantir que o seu domínio sempre aponte para o IP correto, sem a necessidade de intervenção manual.

## Conteúdo

1. [Funcionalidades](#funcionalidades)
2. [Pré-requisitos](#pré-requisitos)
3. [Instalação](#instalação)
3. [Autores](#autores)

## Funcionalidades

-  Atualização de registros DNS no Cloudflare.
-  Suporte a múltiplos domínios e subdomínios.

## Pré-requisitos

-- Conta no Cloudflare
-- Token de API gerado no Cloudflare com permissões para modificar registros DNS.
-- Go (Golang) para compilar o projeto.

## Instalação
1. Clone o repositório:

```bash
   git clone https://github.com/seuusuario/cloudflare-ip-updater.git
```

2. Compile o projeto em Go:

```bash
   cd cloudflare-ip-updater
   go build -o cloudflare-updater

```

3. Configure suas credenciais no arquivo de configuração.

4. Configure suas credenciais no arquivo de configuração.

```bash
crontab -e
```

Exemplo de linha no crontab para execução a cada 5 minutos:

```bash
*/5 * * * * /caminho/para/seu/cloudflare-updater
```

## Autores

-   [@frankduque](https://github.com/frankduque)