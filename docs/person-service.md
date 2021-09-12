# Person service

This service exposes REST API endpoint for clients to create, read and update persons.
After configurable number of minutes, the person that is not updated in this time period, is copied
to S3 and then deleted. See Backup service for more details.

Person is modeled like this:

| Field        | Type          | Example value |
|--------------|---------------|---------------|
| id           | string (uuid) | 410ffb3f-bddf-409d-a397-f0e37e9f3294 |
| name         | string        | Peter         |
| address      | string        | 24 School Lane London|
| dateOfBirth  | string date   | DD/MM/YYYY    |

## Endpoints

### Create Person

**Request**

| Name           | Method | Description |
|----------------|--------|-------------|
| /api/v1/person | POST   | Persist new Person object in database and returns created Person with identifier|

**Request body example**
```json
{
  "name": "Peter",
  "address": "24 School Lane London",
  "dateOfBirth": "01/05/1991"
}
```

**Response example**

Code: 201 Created
```json
{
  "Id": "410ffb3f-bddf-409d-a397-f0e37e9f3294",
  "name": "Peter",
  "address": "24 School Lane London",
  "dateOfBirth": "01/05/1991"
}
```

### Retrieve Person

**Request**

| Name                | Method | Description |
|---------------------|--------|-------------|
| /api/v1/person/{id} | GET    | Retrieves Person from database using identifier |

**Response example**

Code: 200 OK
```json
{
  "Id": "410ffb3f-bddf-409d-a397-f0e37e9f3294",
  "name": "Marc",
  "address": "25 School Lane London",
  "dateOfBirth": "02/06/1989"
}
```

### Update Person Optimistic

**Request**

| Name           | Method | Description |
|----------------|--------|-------------|
| /api/v1/person | PATCH  | Updates Person object in database using optimistic locking |

**Request body example**
```json
{
  "name": "Peter",
  "address": "24 School Lane London",
}
```

**Response example**

Code: 200 OK
```json
{
  "Id": "410ffb3f-bddf-409d-a397-f0e37e9f3294",
  "name": "Peter",
  "address": "24 School Lane London",
  "dateOfBirth": "01/05/1991"
}
```

### Update Person Pessimistic

**Request**

| Name                       | Method | Description |
|----------------------------|--------|-------------|
| /api/v1/person/pessimistic | PATCH  | Updates Person object in database using pessimistic locking |

**Request body example**
```json
{
  "name": "Marc",
  "address": "25 School Lane London",
}
```

**Response example**

Code: 200 OK
```json
{
  "Id": "410ffb3f-bddf-409d-a397-f0e37e9f3294",
  "name": "Marc",
  "address": "25 School Lane London",
  "dateOfBirth": "02/06/1989"
}
```