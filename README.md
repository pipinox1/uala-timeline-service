
# Timeline Service

Service to manage timelines of users. This service allow to get timeline user between two dates. This timelines are stored on postsgres db.

To optimze this solution we use a dynamodb to store a day "snapshot" that allows us to scale.

The core mechanisim of this flows in on the add post method whis is trigger when users publish a posts. In the next diagram we represent the add_post_flow

![Add Post Flow](https://i.ibb.co/GQtWdh6x/uala-posts-architecture.jpg)

Here we explain how we store the posts on dynamo.

![Dynamo](https://i.ibb.co/wZznN8V0/test.jpg)

To reduce and scale more we will reduce the size of post stored using gzip compression.

## API Reference

#### Get user timeline

```http
  POST /api/v1/user_timeline/${user_id}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `user_id` | `string` | **Required**. ID of the user to get timeline for |
| `Content-Type` | `string` | **Required**. application/json |

**Request Body:**

| Field | Type     | Description                       |
| :---- | :------- | :-------------------------------- |
| `from_day` | `number` | **Required**. Starting day of the date range |
| `from_month` | `number` | **Required**. Starting month of the date range |
| `from_year` | `number` | **Required**. Starting year of the date range |
| `to_day` | `number` | **Required**. Ending day of the date range |
| `to_month` | `number` | **Required**. Ending month of the date range |
| `to_year` | `number` | **Required**. Ending year of the date range |

## How to Run?


In order to run you have to clone:
- `uala-posts-service`
- `uala-followers-service`
- `uala-timeline-service`

In the `uala-posts-service` repository, you'll find the `docker-compose` file used to start all services.

### Steps:


1. After cloning each repository, run the following command inside each one:

   ```bash
   make build
   ```

2. Once all images are built, navigate to the uala-posts-service repository and run:

   ```bash
   docker-compose up -d
   ```

### Service Ports:
- posts-service: 8080
- timeline-service: 8081
- followers-service: 8082