# Typhoon

Here is the repository of the Typhoon project. With Typhoon we allow logged in users to deploy their own websites on our infrastructure, on a given domain name, and with HTTPS.

## How to deploy the Typhoon project in production

Hallo! Here's how to deploy:

You'll need a machine (preferably debian), and be able to run ansible scripts in that machine.
Git clone `https://github.com/typhoon-docker/ansible.git`.
Follow the `README` (and make sure to update `users.yml` in `ansible/roles/setup_users/vars` with your admin's RSA keys).

This will install `docker`, `go`, `oh-my-zsh`, and create the required directories.

Then, connect into the VM via SSH. Make sure docker is working: `docker run hello-world`. You may have to deal with group issues to be able to run docker from the debian user (use `sudo groupadd docker` to create the group and `sudo usermod -aG docker debian` to add `debian` to it).

Once docker works fine, git clone the typhoon repo in `~`: `https://github.com/typhoon-docker/typhoon.git`. It contains the frontend and backend code.

### Setting up the backend

This projects uses a nginx proxy to deploy the websites. You will need to (from `~`):

```sh
git clone https://github.com/typhoon-docker/docker-nginx-proxy
cd docker-nginx-proxy
mkdir certs
docker network create nginx-proxy
docker-compose up
```

Once the proxy is up and running it's time to run the backend.
Go to the `back` folder of the typhoon project.

First, use the `.env.template` to create the `.env` (copy it and add the correct values):

```
VIAREZO_CLIENT_ID= # App ID for ViaRezo Oauth
VIAREZO_CLIENT_SECRET= # App secret for ViaRezo Oauth
GITHUB_CLIENT_ID=# App ID for Github Oauth
GITHUB_CLIENT_SECRET=# App secret for Github Oauth
JWT_SECRET= # To sign you tokens, put a random (long) string
```

Then `cd docker_compose_production`. This contains the production docker-compose yml file.
You will have to (manually) modify the `docker-compose.yml`. Replace:
```
      - VIRTUAL_HOST=typhoon-back.typhoon.viarezo.fr
      - LETSENCRYPT_HOST=typhoon-back.typhoon.viarezo.fr
      - LETSENCRYPT_EMAIL=aymeric.bernard@student.ecp.fr
```

By changing `typhoon.viarezo.fr` to your domain name.

Then modify `.env.production` in the same way regarding the domain name.

Finally `docker-compose up`.

Do make sure your DNS config (that pointing to the VM) also allows for wildcards as your projects will be deployed as `[project_name].typhoon.viarezo.fr`.

You can use `typhoon-back.[your domain name]/healthCheck` to see it the back is successfully running (and the nginx).

### Setting up the frontend

In a similar fashion, go to the `front` folder from the root. Modify

```
      - VIRTUAL_HOST=typhoon.viarezo.fr
      - LETSENCRYPT_HOST=typhoon.viarezo.fr
```

to suit your domain name, as well as `.env.production`:

```
BACKEND_URL=https://typhoon-back.[your domain name]
FRONTEND_URL=https://[your domain name]
```

Then:

- build image: `docker build -t typhoon-front -f docker/Dockerfile .`
- start container: `docker-compose up`

# Commands to deploy the backend code

- Build image: (from `./back`) `docker build -t typhoon-back-go .`

- Start in dev: (from `./back`) `docker-compose up`
- Start in prod: (from `./back/docker_compose_production`) `docker-compose up`

- Open console (for debug): `docker-compose run code bash`

# Commands to deploy the frontend code

git pull && docker build -t 2015koltesb-typhoon -f docker/Dockerfile . && cd docker && docker-compose up -d && cd ..

- Build image: (from `./back`) `docker build -t typhoon-back-go .`

- Start in dev: (from `./back`) `docker-compose up`
- Start in prod: (from `./back/docker_compose_production`) `docker-compose up`

- Open console (for debug): `docker-compose run code bash`

# Env variable files that are loaded

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

### Structs

#### Container

```go
type Container struct {
	Id     string `json:"id"`
	Image  string `json:"name"`
	Status string `json:"status"` // ex: "Up for 14 min."
	State  string `json:"state"`  // ex: "running"
}
```

### Routes

#### Get my projects, or get all projects if I am admin

`/projects(?all)` - **GET**
- headers : `{ "Authorization: "Bearer <token>" }`
- return: `[Project]`

#### Get details by project id

`/projects/:id` - **GET**
- headers : `{ "Authorization: "Bearer <token>" }`
- return: `Project`

#### Post a new project (will be added to database but not built)

`/projects` - **POST**
- headers : `{ "Authorization: "Bearer <token>", "Content-Type": "application/json" }`
- body: Project
- return: `[Project]`

#### Update project with new info (full override, id in url and body must match)

`/projects/:id` - **PUT**
- headers : `{ "Authorization: "Bearer <token>", "Content-Type": "application/json" }`
- body: Project
- return: `[Project]`

### Remove a project by id

`/projects/:id` - **DELETE**
- headers : `{ "Authorization: "Bearer <token>" }`
- return: `[Project]`

## Docker management

#### Return docker files of project (doesn't write them)

`/docker/files/:id` - **GET**
- headers : `{ "Authorization: "Bearer <token>" }`
- return:

```
{
  "project": <project structure>,
  "dockerfile_0": <dockerfile> || "error_dockerfile_0": <error>,
  "docker_compose": <docker_compose> || "error_docker_compose": <error>
}
```

#### Clone, build and run project

`/docker/apply/:id` - **POST**
- headers : `{ "Authorization: "Bearer <token>" }`
- return:

```
{
  "project": <project structure>,
  "dockerfile_0": <dockerfile> || "error_dockerfile_0": <error>,
  "docker_compose": <docker_compose> || "error_docker_compose": <error>
}
```

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
- return: `[Container]`

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
