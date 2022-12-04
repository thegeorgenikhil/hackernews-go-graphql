# HackerNews Clone

This is a HackerNews clone built with GraphQL(gglgen), Go and MySQL.

### Frontend

- [Live Demo]()
- [Frontend Repo]()

It is based on the [How to GraphQL Tutorial](https://www.howtographql.com/graphql-go/0-introduction/).

## Running the app

1. Clone the repo

```bash
git clone https://github.com/thegeorgenikhil/hackernews-go-graphql
```

2. Install dependencies

```bash
go mod download
```

3. Start a MySQL Docker container using the following command:

```bash
docker run -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=dbpass -e MYSQL_DATABASE=hackernews -d mysql:latest
```

4. Run the app

```bash
go run server.go
```

5. Open [http://localhost:8080](http://localhost:8080) to view the GraphQL Playground it in the browser.

6. Use [http://localhost:8000/query](http://localhost:8000/query) to query the GraphQL API.

## GraphQL Queries and Mutations

### Queries

- Get all links

```graphql
query GetLinks {
  links {
    title
    address
    id
    user {
      id
      name
    }
  }
}
```

Response:

```json
{
  "data": {
    "links": [
      {
        "title": "GraphQL",
        "address": "https://graphql.org/",
        "id": "1",
        "user": {
          "id": "1",
          "name": "thegeorgenikhil"
        }
      },
      {
        "title": "How to GraphQL",
        "address": "https://www.howtographql.com/",
        "id": "2",
        "user": {
          "id": "1",
          "name": "thegeorgenikhil"
        }
      }
    ]
  }
}
```

## Mutations

- Post a new link

```graphql
mutation CreateLink($input: NewLink!) {
  createLink(input: $input) {
    id
    title
    address
    user {
      id
      name
    }
  }
}
```

Variables:

```json
{
  "input": {
    "title": "GraphQL",
    "address": "https://graphql.org/"
  }
}
```

Headers:

```json
{
  "Authorization": "<user-token-here>"
}
```

Response:

```json
{
  "data": {
    "createLink": {
      "id": "1",
      "title": "GraphQL",
      "address": "https://graphql.org/",
      "user": {
        "id": "1",
        "name": "thegeorgenikhil"
      }
    }
  }
}
```

- Create a new user

```graphql
mutation CreateUser($input: NewUser!) {
  createUser(input: $input)
}
```

Variables:

```json
{
  "input": {
    "username": "user",
    "password": "password"
  }
}
```

Response:

```json
{
  "data": {
    "createUser": "<user-token-here>"
  }
}
```

---

- Login a user

```graphql
mutation LoginUser($input: Login!) {
  login(input: $input)
}
```

Variables:

```json
{
  "input": {
    "username": "user",
    "password": "password"
  }
}
```

Response:

```json
{
  "data": {
    "login": "<user-token-here>"
  }
}
```

---

- Refresh a user token

```graphql
mutation RefereshToken($input: RefreshTokenInput!) {
  refreshToken(input: $input)
}
```

Variables:

```json
{
  "input": {
    "token": "<expired-token-here>"
  }
}
```

Response:

```json
{
  "data": {
    "refreshToken": "<new-user-token-here>"
  }
}
```
