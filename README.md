# Music API

## Overview

This project implements a RESTful API for managing music tracks and playlists. It allows users to perform CRUD operations on tracks and playlists, search for tracks and playlists, and generate Swagger documentation for easy integration and testing.

## Technical

- **Language**: Golang
- **Framework**: Gin
- **Database**: MongoDB
- **Swagger Configuration**: Generate Swagger documentation for the API endpoints
- **Docker Configuration**: Provide Dockerfile and docker-compose.yml for containerization

## Getting Started

Follow these steps to set up and run the API:

### Prerequisites

- Install Golang
- Install MongoDB
- Install Swagger
- Install Makefile

### Installation

- **1. Clone the repository**: git clone <repository_url>

- **2. Install dependencies**: run : /dockerrun.sh

### Configuration

1. System Configuration:

- Create file config.json from config.json.example at folder config of reposity.

### Running the API

- **Run the application**: make dev

### Testing the API

- **Open the url**: http://localhost:8000/swagger/index.html

## For feature in future

- Storage file audio upload in cloud storage instead of local
- Play track in playlist follow priority
- Play track in playlist random
