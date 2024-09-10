# Task Manager Server

## Overview

This project implements a robust Task Manager Server using Go, designed for handling real image processing tasks.
The server supports a distributed system architecture, allowing for horizontal scalability to manage and process
image-related tasks efficiently.  It includes a comprehensive API for task submission and status tracking,
with the capability to offload tasks to a dedicated processing service and store results in a database.

## Getting Started

### Building and Running the Docker Image

To manage the Task Manager Server, you can use the following `make` commands:

- **Build the Docker image and run tests:**

   ```bash
   make tests
   ```

- **Start the server without running tests:**

   ```bash
   make run
   ```

- **Stop the server and close any running images:**

   ```bash
   make stop
   ```

### Shooter API

The Shooter API allows you to apply filters to images. Here’s how to use it:

1. **Navigate to the `shooter` directory and install dependencies:**

    ```bash
    cd shooter
    pip install -r requirements.txt
    ```

2. **Run the `shoot.py` script to apply a filter to an image:**

    ```bash
    python shoot.py <FilterName> [<Parameters>]
    ```

    - **`<FilterName>`**: Name of the filter (e.g., `Grayscale`, `Negative`, `Blur`, `Sharpen`).
    - **`<Parameters>`**: Optional filter parameters in JSON format (e.g., `{"sigma": 5.0}`).

The processed images will be saved in the 'shooter/results' directory.

**Examples:**

   ```bash
   python shoot.py Grayscale
   ```

   ```bash
   python shoot.py Negative
   ```

   ```bash
   python shoot.py Blur '{"sigma": 5.0}'
   ```

   ```bash
   python shoot.py Sharpen '{"sigma": 5.0}'
   ```

### Accessing the API Documentation

Once the server is up and running, you can view the API documentation by navigating to:

http://localhost:8000/swagger/index.html#/

This link will direct you to the Swagger UI, where you can explore the available API endpoints.