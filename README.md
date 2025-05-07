# Go HTTP Server

A simple HTTP server built in Go.

## API Endpoints

This section details the available API endpoints and their functionalities.

### Health Check

* **GET /api/healthz**: Checks the health of the server.
  * **Request Body**: None.
  * **Response (200 OK)**:

### Admin

* **POST /admin/reset**: Resets the application's metrics (specifically, the file server hit counter).
  * **Request Body**: None.
  * **Response (204 No Content)**: Indicates successful reset.
* **GET /admin/metrics**: Retrieves server metrics.
  * **Request Body**: None.
  * **Response (200 OK)**: HTML page displaying the number of file server hits.

### Users

* **post /api/users**: creates a new user.
  * **request body**:

        ```json
        {
            "password": "yourpassword",
            "email": "user@example.com"
        }
        ```

  * **response body (201 created)**:

        ```json
        {
            "id": "uuid",
            "created_at": "timestamp",
            "updated_at": "timestamp",
            "email": "user@example.com",
            "is_chirpy_red": false
        }
        ```

* **post /api/login**: logs in a user.
  * **request body**:

        ```json
        {
            "email": "user@example.com",
            "password": "yourpassword"
        }
        ```

  * **response body (200 ok)**:

        ```json
        {
            "id": "uuid",
            "created_at": "timestamp",
            "updated_at": "timestamp",
            "email": "user@example.com",
            "is_chirpy_red": false,
            "token": "jwt_access_token",
            "refresh_token": "jwt_refresh_token"
        }
        ```

* **put /api/users**: updates an existing user's information.
  * **authentication**: requires bearer token in the `authorization` header.
  * **request body**:

        ```json
        {
            "password": "newpassword",
            "email": "newuser@example.com"
        }
        ```

  * **response body (200 ok)**:

        ```json
        {
            "id": "uuid",
            "created_at": "timestamp",
            "updated_at": "timestamp",
            "email": "newuser@example.com",
            "is_chirpy_red": false 
        }
        ```

### Chirps

* **POST /api/chirps**: Creates a new chirp.
  * **Authentication**: Requires Bearer Token in the `Authorization` header.
  * **Request Body**:

        ```json
        {
            "body": "This is a chirp!"
        }
        ```

  * **Response Body (201 Created)**:

        ```json
        {
            "id": "uuid",
            "created_at": "timestamp",
            "updated_at": "timestamp",
            "user_id": "uuid",
            "body": "This is a chirp!"
        }
        ```

* **GET /api/chirps**: Retrieves all chirps.
  * **Query Parameters**:
    * `author_id` (optional, uuid string): Filters chirps by the user ID of the author.
    * `sort` (optional, string: "asc" or "desc"): Sorts chirps by creation date. Defaults to ascending order.
  * **Response Body (200 OK)**:

        ```json
        [
            {
                "id": "uuid",
                "created_at": "timestamp",
                "updated_at": "timestamp",
                "user_id": "uuid",
                "body": "Chirp content"
            }
            // ... more chirps
        ]
        ```

* **GET /api/chirps/{chirpID}**: Retrieves a specific chirp by its ID.
  * **Path Parameter**: `chirpID` (uuid)
  * **Response Body (200 OK)**:

        ```json
        {
            "id": "uuid",
            "created_at": "timestamp",
            "updated_at": "timestamp",
            "user_id": "uuid",
            "body": "Chirp content"
        }
        ```

  * **Response (404 Not Found)**: If chirp with the given ID doesn't exist.
* **DELETE /api/chirps/{chirpID}**: Deletes a specific chirp by its ID.
  * **Authentication**: Requires Bearer Token in the `Authorization` header.
  * **Path Parameter**: `chirpID` (uuid)
  * **Response (204 No Content)**: On successful deletion.
  * **Response (403 Forbidden)**: If the authenticated user is not the author of the chirp.
  * **Response (404 Not Found)**: If chirp with the given ID doesn't exist.

### Authentication

* **POST /api/refresh**: Refreshes an authentication token.
  * **Authentication**: Requires Bearer Token (Refresh Token) in the `Authorization` header.
    * Example: `Authorization: Bearer your_refresh_token`
  * **Response Body (200 OK)**:

        ```json
        {
            "token": "new_jwt_access_token"
        }
        ```

* **POST /api/revoke**: Revokes an authentication token.
  * **Authentication**: Requires Bearer Token (Refresh Token) in the `Authorization` header.
    * Example: `Authorization: Bearer your_refresh_token`
  * **Response (204 No Content)**: On successful token revocation.

### Webhooks

* **POST /api/polka/webhooks**: Handles webhooks from the Polka service.
  * **Authentication**: Requires API Key in the `Authorization` header.
    * Example: `Authorization: ApiKey YOUR_POLKA_KEY`
  * **Request Body**:

        ```json
        {
            "event": "user.upgraded",
            "data": {
                "user_id": "uuid_of_upgraded_user"
            }
        }
        ```

  * **Response (204 No Content)**: If the event is not `user.upgraded` or if processing is successful for `user.upgraded`.
  * **Response (401 Unauthorized)**: If the API key is missing or invalid.
  * **Response (404 Not Found)**: If the `user_id` in the webhook data is not found (only for `user.upgraded` event).

### Static Files

* **GET /app/**: Serves static files from the `.` (root) directory. For example, `/app/index.html` would serve `index.html` from the root.

## Running the Server

1. **Environment Variables**:
    * Ensure you have a `.env` file with the following variables:
        * `DB_URL`: The connection string for your PostgreSQL database.
        * `PLATFORM`: (Optional) A string indicating the platform.
        * `SECRET`: A secret key used for JWT or other cryptographic operations.
        * `POLKA_KEY`: API key for the Polka service.
2. **Build and Run**:

    ```bash
    go build
    ./go-httpserver
    ```

3. The server will start on `http://localhost:8080`.

## Dependencies

* [github.com/joho/godotenv](https://github.com/joho/godotenv) - For loading environment variables from `.env` files.
* [github.com/lib/pq](https://github.com/lib/pq) - PostgreSQL driver.
