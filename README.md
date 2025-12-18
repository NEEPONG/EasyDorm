# Dormitory Management System

A web-based application for managing dormitory operations including rooms, tenants, maintenance, and billing.

## Features

- **User Authentication**: Secure login system for administrators.
- **Dashboard**: Overview of dormitory status.
- **Room Management**:
  - Add new rooms.
  - Edit existing room details.
  - Search for rooms.
- **Tenant Management**:
  - Register new tenants.
  - Search for tenant information.
  - Handle tenant checkout.
- **Maintenance**:
  - Track maintenance requests.
  - Update maintenance status.
- **Billing**:
  - Manage billing cycles.
  - Confirm billing details.

## Tech Stack

- **Backend**: Go (Golang)
- **Database**: MySQL
- **Frontend**: HTML, CSS, JavaScript

## Setup Instructions

1.  **Clone the repository**

    ```bash
    git clone <repository-url>
    cd DormitoryMng
    ```

2.  **Prerequisites**

    - Go 1.24 or higher
    - MySQL Server

3.  **Database Configuration**

    - Create a MySQL database for the project.
    - Import the initial schema (if available).
    - _Note: Update the database connection string in the code (e.g., in `model/model.go` or `controller/service.go`) to match your local MySQL credentials._

4.  **Install Dependencies**

    ```bash
    go mod tidy
    ```

5.  **Run the Application**

    ```bash
    go run main.go
    ```

6.  **Access the Application**
    - Open your web browser and navigate to: `http://localhost:8090`
