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
* `.env.local` and `.env.test.local` are not loaded when `GO_ENV=test` since tests should produce the same results for everyone

# Doc

Oauths:
- [GitHub](https://developer.github.com/apps/building-oauth-apps/authorizing-oauth-apps/)
- [ViaRezo](https://auth.viarezo.fr/docs)

# API specification

## Projects management

#### Get my projects, or get all projects if I am admin

`/projects(?all)` - **GET**
- headers : `{ "Authorization: "Bearer <token>" }`
- return: `[Project]`

#### Get details by project id

`/projects/:id` - **GET**
- headers : `{ "Authorization: "Bearer <token>" }`
- return: `Project`

#### Update project with new info (full override, id in url and body must match)

`/projects/:id` - **PUT**
- headers : `{ "Authorization: "Bearer <token>", "Content-Type": "application/json" }`
- body: Project
- return: `[Project]`

#### Post a new project (will be added to database but not built)

`/projects` - **POST**
- headers : `{ "Authorization: "Bearer <token>", "Content-Type": "application/json" }`
- body: Project
- return: `[Project]`

### Remove a project by id

`/projects/:id` - **DELETE**
- headers : `{ "Authorization: "Bearer <token>" }`
- return: `[Project]`

## Docker management

#### Build and run project

`/docker/apply/:id` - **POST**
- headers : `{ "Authorization: "Bearer <token>" }`
- return: `{ "project": <project structure>, "dockerfile_1": <dockerfile>, "docker_compose": <docker_compose>}`

#### Run project

`/docker/up/:id` - **POST**
- headers : `{ "Authorization: "Bearer <token>" }`
- return: `"OK" || <error>`

#### Undeploy project

`/docker/down/:id` - **POST**
- headers : `{ "Authorization: "Bearer <token>" }`
- return: `"OK" || <error>`

#### Get project status

`/docker/status/:id` - **GET**
- headers : `{ "Authorization: "Bearer <token>" }`
- return: `raw $(docker-compose ps) output`

#### Get logs

`/docker/logs/:id?lines=<lines>` - **GET**

- headers : `{ "Authorization: "Bearer <token>" }`
- return: `raw logs`

Parameter lines is optional. Default is 30.

## User and Admin management

#### List all users, or all admins

`/admin/list(?admin)` - **GET**
- headers : `{ "Authorization: "Bearer <token>" }`
- return: `[User]`

#### Change user scope

`/admin/scope/:id?scope=<scope>` - **PUT**
- headers : `{ "Authorization: "Bearer <token>" }`
- return: `"OK" || <error>`

#### Edit user

`/admin/user/:id` - **PUT**
- body: User
- return: `User`

#### Delete user

`/admin/user/:id` - **DELETE**
- return: `user_id`

## Misc.

### Check if the server is up

`/healthCheck` - **GET**
- return: `"OK"`

### Check if a project has the given name

`/checkProject?name=<name>` - **GET**
- return: `true` if project with that name exists, else `false`

### See my token info

`/showme` - **GET**
- headers : `{ "Authorization: "Bearer <token>" }`
- return: `<JWT info>`
