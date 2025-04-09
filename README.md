# Assignment-2 (group err)
---

# Overview

This project is a RESTful web service in Go. It manages dashboard configurations (registrations) and dynamically populates dashboards with real-time data from multiple external APIs, including REST Countries, Open-Meteo, and a Currency API. Users can specify which features (temperature, precipitation, population, etc.) should appear on the resulting dashboards, and the service then merges and returns the appropriate information. It also provides a notifications mechanism for registering webhooks that trigger on certain events, such as creating, updating, deleting, or invoking a dashboard configuration. The service employs Firestore (Firebase) for persistent storage of configurations, ensures caching to minimize external calls, and includes a status endpoint for health checks and uptime monitoring. Additionally, the entire service is dockerized and deployed to a VM on OpenStack.

---

# External services
 - REST Countries API
   - Endpoint: http://129.241.150.113:8080/v3.1
   - Documentation: http://129.241.150.113:8080/
 - Open-Meteo APIs
   - Documentation: https://open-meteo.com/en/features#available-apis
 - Currency API
   - Endpoint: http://129.241.150.113:9090/currency/
   - Documentation: http://129.241.150.113:9090/

# Architecture & Data Flow


                         +---------------------------+
                         |         Firestore        |
                         | (Registrations, Notifications, Cache)     
+------------+           |---------------------------|
|  External  |           |                           |
|  Clients   | --->  +---v-----------------------+   |
| (Postman,  |       |  Go REST API (Handlers)   |   |
|  Browser)  | <---  | /dashboard/v1/...         |   |
+------------+       +-------------+-------------+   |
                         ^          |                |
                         |          |                |
                         |          |(Cached)        |
                         |          |                |
+------------------+     |    +-----+------------+   |
|   REST Countries |     |    |   External APIs  |   |
+------------------+     |    | - Open-Meteo     |   |
                          \---> - Currency API   +---+


---
# External Clients interact via HTTP endpoints (POST /registrations, GET /dashboards/{id}, etc.).

1- The Go REST API handles these requests, pulling and merging data from:

2- Firestore (for saved configurations, notifications, cache).

3- External APIs when not cached or when fresh data is needed.

4- Notifications are triggered upon certain events, causing POST calls to user-registered webhook URLs.

5- A Status Handler checks health of external APIs and Firestore.


# Prerequisites

1- Go (a stable version) if you plan to build/run locally (without Docker).

2- A valid Firebase service account key (placed in the project root), or a configured environment for Firestore credentials.

3- Docker installed, if you want to build and run the container.

4- OpenStack (SkyHigh) access, if you plan to deploy there.

--- 

# Local Setup & Running

## Clone the repository:
~~~
git clone https://github.com/ahmadss3/assignment-2.git
~~~

## Install dependencies
## firestore account with valid key

## Run the service locally:
~~~
go run ./cmd/main.go
~~~

## Verify:
~~~
http://localhost:8080/dashboard/v1/status/
~~~

Should return a JSON with countries_api, currency_api, meteo_api, notification_db, etc.

---
# openstack:
you can run the the project by using this url:
~~~
http://10.212.174.26:8080/dashboard/v1/status
~~~
## important: must use the NTNU network to run the provided url
---
# Running the Service with Docker:
## NOTE:
this service uses Firestore (a database service within Google’s Firebase) to store:

 - Dashboard registrations
 - Webhook notifications
 - Cached data

- You must have a Firebase project set up that includes Firestore. That way, the service can connect to a live Firestore database and persist all necessary information.

- How to Set Up Firestore:
   - Create a Firebase Project by signing up at https://firebase.google.com/
   - In your Firebase project settings, create a service account key (JSON) with permissions to read and write data in Firestore.
   - Name it something like assignment-2-firebasekey.json.
   - Put the JSON file in the project root so the application can load it (IMPORTANT: Add firebasekey.json to your .gitignore and never commit or push it).
   - Copy the JSON file into the Docker image (see the Dockerfile) or mount it as a volume so the container can access it.

## building locally (assuming you have cloned the project folder and are in its root directory):
~~~
docker build -t assignment-2 .
~~~

## Run the Container:
~~~
docker run -d -p 8080:8080 --name assignment-2 assignment-2
~~~

## Verify the Service:
Once the container is up, open your browser or use a tool like Postman to access:
~~~
http://localhost:8080/dashboard/v1/
~~~

### 
Once the container is running and responding to HTTP requests on localhost:8080, you can proceed to test all endpoints – like creating registrations, retrieving dashboards, or registering notifications.

---

---

# Endpoints & Usage

## Registrations
- Create, retrieve, update, patch, and delete dashboard configurations that specify the data features (temperature, population, area, etc.) for a given country or ISO code.
- Stored persistently in Firestore (Firebase).

### `POST /dashboard/v1/registrations/`
Creates a new dashboard configuration.

#### **Request**
- **Method**: `POST`
- **Path**: `/dashboard/v1/registrations/`
- **Headers**:
  - `Content-Type: application/json`
- **Body** (example):
~~~
{
  "country": "Norway",
  "isoCode": "",
  "features": {
    "temperature": true,
    "precipitation": true,
    "capital": true,
    "coordinates": true,
    "population": true,
    "area": true,
    "targetCurrencies": ["EUR", "USD", "SEK"]
  }
}
~~~

#### **Response**
- **Status**: 201 Created (or relevant error code)
- **Body** (example):
~~~
{
  "id": "abc123def",
  "lastChange": "20250410 12:10"
}
~~~

---

### `GET /dashboard/v1/registrations/`
Retrieves all dashboard configurations.

#### **Request**
- **Method**: `GET`
- **Path**: `/dashboard/v1/registrations/`

#### **Response**
- **Status**: 200 OK
- **Body** (example):
~~~
[
  {
    "id": "abc123def",
    "country": "Norway",
    "isoCode": "NO",
    "features": {
      "temperature": true,
      "precipitation": true,
      "capital": true,
      "coordinates": true,
      "population": true,
      "area": false,
      "targetCurrencies": ["EUR", "USD", "SEK"]
    },
    "lastChange": "20250410 12:10"
  },
  {
    "id": "doc456xyz",
    "country": "Sweden",
    "isoCode": "SE",
    "features": {
      "temperature": false,
      "precipitation": true,
      "capital": true,
      "coordinates": false,
      "population": true,
      "area": true,
      "targetCurrencies": ["USD", "EUR"]
    },
    "lastChange": "20250410 14:07"
  }
]
~~~

**HEAD Method**:  
You can also call `HEAD /dashboard/v1/registrations/` to return headers without a body.  
The service responds with `200 OK` if successful.

---

### `GET /dashboard/v1/registrations/{id}`
Retrieves a single configuration by its ID.

#### **Request**
- **Method**: `GET`
- **Path**: `/dashboard/v1/registrations/{id}`

#### **Response**
- **Status**: 200 OK if found; 404 Not Found if it does not exist
- **Body** (example):
~~~
{
  "id": "abc123def",
  "country": "Norway",
  "isoCode": "NO",
  "features": {
    "temperature": true,
    "precipitation": true,
    "capital": true,
    "coordinates": true,
    "population": true,
    "area": true,
    "targetCurrencies": ["EUR", "USD", "SEK"]
  },
  "lastChange": "20250410 12:10"
}
~~~

---

### `PUT /dashboard/v1/registrations/{id}`
Overwrites an existing configuration.

#### **Request**
- **Method**: `PUT`
- **Path**: `/dashboard/v1/registrations/{id}`
- **Headers**:
  - `Content-Type: application/json`
- **Body** (example):
~~~
{
  "country": "Norway",
  "isoCode": "NO",
  "features": {
    "temperature": false,
    "precipitation": true,
    "capital": true,
    "coordinates": true,
    "population": true,
    "area": false,
    "targetCurrencies": ["EUR", "SEK"]
  }
}
~~~

#### **Response**
- **Status**: 204 No Content if successful
- **Body**: (empty)

---

### `PATCH /dashboard/v1/registrations/{id}`
Applies a partial update to the existing configuration.

#### **Request**
- **Method**: `PATCH`
- **Path**: `/dashboard/v1/registrations/{id}`
- **Headers**:
  - `Content-Type: application/json`
- **Body** (example):
~~~
{
  "country": "Norway",
  "features": {
    "temperature": true,
    "area": false,
    "targetCurrencies": ["GBP", "USD"]
  }
}
~~~

#### **Response**
- **Status**: 204 No Content
- **Body**: (empty)

---

### `DELETE /dashboard/v1/registrations/{id}`
Deletes a dashboard configuration.

#### **Request**
- **Method**: `DELETE`
- **Path**: `/dashboard/v1/registrations/{id}`

#### **Response**
- **Status**: 204 No Content if deleted; 404 Not Found otherwise
- **Body**: (empty)

---

## Dashboards
- Dynamically merges real-time data from:
  - REST Countries (capital, population, area, lat/long, base currency, etc.)
  - Open-Meteo (temperature & precipitation)
  - Currency API (exchange rates)
- Returns a populated dashboard with the requested features.

### `GET /dashboard/v1/dashboards/{id}`
Retrieves a populated dashboard with data from external APIs (REST Countries, Open-Meteo, Currency API).

#### **Request**
- **Method**: `GET`
- **Path**: `/dashboard/v1/dashboards/{id}`

#### **Response**
- **Status**: 200 OK if found; 404 Not Found otherwise
- **Body** (example):
~~~
{
  "country": "Norway",
  "isoCode": "NO",
  "features": {
    "temperature": 5.2,
    "precipitation": 1.4,
    "capital": "Oslo",
    "coordinates": {
      "latitude": 60.0,
      "longitude": 10.0
    },
    "population": 5372000,
    "area": 385207,
    "targetCurrencies": {
      "EUR": 0.09,
      "USD": 0.1
    }
  },
  "lastRetrieval": "20250410 18:15"
}
~~~

Also triggers an `INVOKE` event for any matching webhooks.

---

## Notifications (Webhooks)
- Users can register webhooks that trigger on specific events:
  - REGISTER (new configuration created)
  - CHANGE (configuration updated or patched)
  - DELETE (configuration deleted)
  - INVOKE (dashboard retrieved)
- The webhooks are themselves stored persistently, so they survive service restarts.

### `POST /dashboard/v1/notifications/`
Registers a new webhook for specified events (`REGISTER`, `CHANGE`, `DELETE`, `INVOKE`).

#### **Request**
- **Method**: `POST`
- **Path**: `/dashboard/v1/notifications/`
- **Headers**:
  - `Content-Type: application/json`
- **Body** (example):
~~~
{
  "url": "https://example.com/hook",
  "country": "NO",
  "event": "REGISTER"
}
~~~

#### **Response**
- **Status**: 201 Created
- **Body** (example):
~~~
{
  "id": "notif-abc123"
}
~~~

---

### `GET /dashboard/v1/notifications/`
Retrieves all registered webhooks.

#### **Request**
- **Method**: `GET`
- **Path**: `/dashboard/v1/notifications/`

#### **Response**
- **Status**: 200 OK
- **Body** (example):
~~~
[
  {
    "id": "notif-1",
    "url": "https://example.com/hook",
    "country": "NO",
    "event": "REGISTER",
    "created": "20250410T10:23:42Z"
  },
  {
    "id": "notif-2",
    "url": "https://webhook.site/test",
    "country": "",
    "event": "INVOKE",
    "created": "20250410T12:05:11Z"
  }
]
~~~

---

### `GET /dashboard/v1/notifications/{id}`
Retrieves a single webhook registration.

#### **Request**
- **Method**: `GET`
- **Path**: `/dashboard/v1/notifications/{id}`

#### **Response**
- **Status**: 200 OK if found; 404 Not Found otherwise
- **Body** (example):
~~~
{
  "id": "notif-abc123",
  "url": "https://example.com/hook",
  "country": "NO",
  "event": "REGISTER",
  "created": "20250410T10:23:42Z"
}
~~~

---

### `DELETE /dashboard/v1/notifications/{id}`
Deletes an existing webhook registration.

#### **Request**
- **Method**: `DELETE`
- **Path**: `/dashboard/v1/notifications/{id}`

#### **Response**
- **Status**: 204 No Content if deleted; 404 Not Found otherwise
- **Body**: (empty)

---

# Status Endpoint
- Indicates the health of external services (REST Countries, Open-Meteo, Currency API), the status of Firestore, the uptime, and how many webhooks are registered.

### `GET /dashboard/v1/status/`
Checks the health of external APIs, notification database, and shows service uptime.

#### **Request**
- **Method**: `GET`
- **Path**: `/dashboard/v1/status/`

#### **Response**
- **Status**: 200 OK if dependencies are healthy; 503 if not
- **Body** (example):
~~~
{
  "countries_api": 200,
  "meteo_api": 200,
  "currency_api": 200,
  "notification_db": 200,
  "webhooks": 3,
  "version": "v1.0.0",
  "uptime": 3600
}
~~~

---

## Webhook Invocation Format

When an event triggers (e.g., `REGISTER`, `CHANGE`, `DELETE`, or `INVOKE`), the service sends a `POST` request to each matching webhook, using JSON like:
~~~
{
  "id": "notif-abc123",
  "country": "NO",
  "event": "REGISTER",
  "time": "20250410 10:25"
}
~~~

- **id**: The webhook registration ID
- **country**: The relevant country or ISO code
- **event**: One of `REGISTER`, `CHANGE`, `DELETE`, or `INVOKE`
- **time**: Timestamp indicating when the event occurred

---
---
# Caching & Periodic Purging
- Country data and other external responses can be cached in Firestore to reduce overhead.
- A background goroutine purges stale cache entries older than 24 hours.
