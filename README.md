# **Task Management System**

Task Management System is a backend service offering a collection of RESTful API endpoints designed for efficient task and user management. 
It supports task creation, updates, and user registration, all secured with JWT-based authentication.

Built using Go, the service leverages modern web frameworks and robust authentication techniques to ensure scalability, security, and ease of integration.

---

## **Table of Contents**
1. [Features](#features)  
2. [Requirements](#requirements)  
3. [Installation and Setup](#installation-and-setup)  
4. [API Endopints](#api-examples)  
5. [License](#license)  

---

## **Features**
- **Authentication**
  - All operations are authenticated with JWT token that user gets after registration.

- **User Management**:  
  - Create new users.  
  - Retrieve user details by ID.  

- **Task Management**:  
  - Create new tasks.  
  - Update task statuses (e.g., `TODO`, `IN_PROGRESS`, `DONE`).  
  - Retrieve tasks assigned to a specific user.

- Graceful server shutdown using context.  

---

## **Requirements**
- Go 1.22 or later  
- MySQL (for database support)

---

## **Installation and Setup**

### 1. Clone the repository
```bash
git clone https://github.com/your-username/task-management-system.git
cd task-management-system
```
### 2. Configure the enviroment
```
DB_HOST=localhost
DB_PORT=your_database_port
DB_USER=your_database_user
DB_PASSWORD=your_database_password
DB_NAME=task_management
SERVER_ADDRESS=:8080
```
### 3. Setup your MySQL database
```sql
CREATE DATABASE projectmanager;
```

### 4. Start the server
```bash
make run
```
## **API Endopints**

### `POST /users/register`
- **Description**: Registers a new user in the system.
- **Request Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }
  ```
- **Response**: A JWT token for the authenticated user.
- **Authentication**: None (registration does not require prior authentication).
- **Success**: On successful registration, the system will return a JWT token in the response body for user authentication.

### `GET /tasks`
- **Description**: Retrieves all tasks assigned to the currently authenticated user.
- **Authentication**: Requires a valid JWT token.
- **Response**: A list of tasks.

### `POST /tasks`
- **Description**: Creates a new task.
- **Authentication**: Requires a valid JWT token.
- **Request Body**:
  ```json
  {
    "name": "Task Name",
    "status": "TODO",
    "assigned_to_id": 1
  }
  ```
- **Response**: The newly created task.

### `GET /tasks/{id}`
- **Description**: Retrieves details of a specific task by its ID.
- **Authentication**: Requires a valid JWT token.
- **Path Parameter**:
  - id: The unique identifier of the task.
- **Response**: The details of the task.

### `POST /tasks/{id}`
- **Description**: Updates the status of a specific task by its ID to the next one:
  - TODO -> IN_PROGRESS
  - IN_PROGRESS -> DONE
- **Authentication**: Requires a valid JWT token.
- **Path Parameter**:
  - id: The unique identifier of the task.
- **Response**: The updated task details.

## License
Distributed under the MIT License. See ```LICENSE``` for more information.
