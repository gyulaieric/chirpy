# Endpoints

## /api/healthz
## GET  
#### Description:  
Returns status code 200 if the server is running.  
Response Body:
```json
{
    "OK"
}
```

## /api/users
## POST  
#### Description:  
Registers a new user to the database.  
#### Request Body:
```json
{
    "email": "your-email",
    "password": "your-password"
}
```

#### Respoonse Body:
```json
{
  "id": "f713a4b7-551a-4083-9a9f-def33afe508d",
  "created_at": "2026-01-17T16:51:40.212611Z",
  "updated_at": "2026-01-17T16:51:40.212611Z"
  "email": "your-email",
  "token": "",
  "refresh_token": "",
  "is_chirpy_red": false,
}
```

## PUT  
#### Description:  
Allows a user to change their email and password.  
#### Request Headers:
```bash
"Authorization": "Bearer your-access-token"
```
#### Request Body:
```json
{
    "email": "your-new-email",
    "password": "your-new-password"
}
```

#### Response Body:
```json
{
  "id": "f713a4b7-551a-4083-9a9f-def33afe508d",
  "created_at": "2026-01-17T16:51:40.212611Z",
  "updated_at": "2026-01-17T16:51:40.212611Z"
  "email": "your-new-email",
  "token": "",
  "refresh_token": "",
  "is_chirpy_red": false,
}
```

## /api/login
## POST  
#### Description:  
Allows a user to log in and get their access and refresh tokens.  
#### Request Body:
```json
{
    "email": "your-email",
    "password": "your-password"
}
```

#### Response Body:
```json
{
  "id": "f713a4b7-551a-4083-9a9f-def33afe508d",
  "created_at": "2026-01-17T16:51:40.212611Z",
  "updated_at": "2026-01-17T16:51:40.212611Z"
  "email": "your-email",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiJmNzEzYTRiNy01NTFhLTQwODMtOWE5Zi1kZWYzM2FmZTUwOGQiLCJleHAiOjE3Njg2NjUxMDAsImlhdCI6MTc2ODY2MTUwMH0.j3gcv7HDJx5rCDWmamMmi8UmfWHuK_UPOz54M_5rWQA",
  "refresh_token": "30474e79bffe131d7d85bad04a3143f6749fb9e94670ddbe2dc2a94a3b2185c8",
  "is_chirpy_red": false,
}
```

## /api/refresh
## POST  
#### Description:  
Allows a user to refresh their access token.  
#### Request Headers:
```bash
"Authorization": "Bearer your-refresh-token"
```

#### Request Body:
```json
{
}
```

#### Response Body:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiJmNzEzYTRiNy01NTFhLTQwODMtOWE5Zi1kZWYzM2FmZTUwOGQiLCJleHAiOjE3Njg2NjUxMDAsImlhdCI6MTc2ODY2MTUwMH0.j3gcv7HDJx5rCDWmamMmi8UmfWHuK_UPOz54M_5rWQA",
}
```

## /api/revoke
## POST  
#### Description:  
Revokes a users' refresh token.
#### Request Headers:
```bash
"Authorization": "Bearer your-refresh-token"
```

#### Request Body:
```json
{
}
```

#### Response Body:
```json
{
}
```

## /api/chirps
## GET  
#### Parameters:
```bash
?author_id=your_id
```
#### Description:  
If the author_id parameter is provided, retrieves the author's chirps from the database.  
Otherwise, retrieves all chirps from the database.

#### Request Body:
```json
{
}
```

#### Response Body:
```json
{
    [
        {
            "body": "Darn that fly, I just wanna cook",
            "created_at": "2026-01-17T16:51:40.23628Z",
            "id": "2dc109ff-3904-4aa0-8ffd-a90093dff0f1",
            "updated_at": "2026-01-17T16:51:40.23628Z",
            "user_id": "f713a4b7-551a-4083-9a9f-def33afe508d"
        },
        {
            "body": "Cmon Pinkman",
            "created_at": "2026-01-17T16:51:40.233778Z",
            "id": "031c19c4-11e7-42a6-b246-6a3725fbf45f",
            "updated_at": "2026-01-17T16:51:40.233778Z",
            "user_id": "f713a4b7-551a-4083-9a9f-def33afe508d"
        },
        {
            "body": "Gale!",
            "created_at": "2026-01-17T16:51:40.231381Z",
            "id": "919145af-c7b0-472a-8068-314e6cd91a97",
            "updated_at": "2026-01-17T16:51:40.231381Z",
            "user_id": "f713a4b7-551a-4083-9a9f-def33afe508d"
        },
        {
            "body": "I'm the one who knocks!",
            "created_at": "2026-01-17T16:51:40.228984Z",
            "id": "82745829-e4db-4061-ae0e-41d044a4af11",
            "updated_at": "2026-01-17T16:51:40.228984Z",
            "user_id": "f713a4b7-551a-4083-9a9f-def33afe508d"
        }
    ]
}
```

## GET /{chirpID} 

#### Description:  
Retrieves a single chirp from the database, based on the given ID from the path variable.

#### Request Body:
```json
{
}
```

#### Response Body:
```json
{
    "body": "I'm the one who knocks!",
    "created_at": "2026-01-17T16:51:40.228984Z",
    "id": "82745829-e4db-4061-ae0e-41d044a4af11",
    "updated_at": "2026-01-17T16:51:40.228984Z",
    "user_id": "f713a4b7-551a-4083-9a9f-def33afe508d"
}
```


## POST

#### Description:  
Creates a new chirp.

#### Request Headers:
```bash
"Authorization": "Bearer your-access-token"
```

#### Request Body:
##### Restrictions:
body parameter should not be longer than 140 characters.
```json
{
    "body": "I'm the one who knocks!",
}
```

#### Response Body:
```json
{
    "body": "I'm the one who knocks!",
    "created_at": "2026-01-17T16:51:40.228984Z",
    "id": "82745829-e4db-4061-ae0e-41d044a4af11",
    "updated_at": "2026-01-17T16:51:40.228984Z",
    "user_id": "f713a4b7-551a-4083-9a9f-def33afe508d"
}
```
## DELETE /{chirpID}

#### Description:  
Allows a user to delete one of the chirps they've created.

#### Request Headers:
```bash
"Authorization": "Bearer your-access-token"
```

#### Request Body:
```json
{
}
```

#### Response Body:
```json
{
}
```

## /admin/reset
## POST  
#### Description:  
Resets fileserver hit count and wipes database.  
Returns 403 Forbidden unless .env variable "PLATFORM" is set to "dev".
#### Request Body:
```json
{
}
```

#### Response Body:
```json
{
}
```

## /admin/metrics
## GET  
#### Description:  
Returns HTML containing the fileserver hit count.
#### Request Body:
```json
{
}
```

#### Response Body:
```html
<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
</html>
```
