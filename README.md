# Commands

- Build image: (from `./back`) `docker build -t typhoon-back-go .`

- Start in dev: (from `./back`) `docker-compose up`
- Start in prod: (from `./back/docker_compose_production`) `docker-compose up`

- Open console: `docker-compose run code bash`

# Env variable

| valid `.env` filenames | `GO_ENV=\*` | `GO_ENV=test` |
| ---------------------- | ----------- | ------------- |
| .env                   | ✔️          | ✔️            |
| .env.{GO_ENV}          | ✔️          | ✔️            |
| .env.local             | ✔️          | ✖️            |
| .env.{GO_ENV}.local    | ✔️          | ✖️            |

Notably:
* `GO_ENV` defaults to `development`, can be `development`, `test`, `production`
* `.env.local` and `.env.test.local` aren't loaded when `GO_ENV=test` since tests should produce the same results for everyone

# Doc

Oauths:
- [GitHub](https://developer.github.com/apps/building-oauth-apps/authorizing-oauth-apps/)
- [ViaRezo](https://auth.viarezo.fr/docs)
