# WASAText

University project for WASA.

## Installation

Build the backend Docker image and run the Docker container:

```bash
docker build -f Dockerfile.backend -t wasatext-backend .
docker run --rm -p 3000:3000 --name wasatext-backend wasatext-backend
```

Build the frontend Docker image and run the Docker container:

```bash
docker build -f Dockerfile.frontend -t wasatext-frontend .
docker run --rm -p 3001:3000 --name wasatext-frontend wasatext-frontend
```
