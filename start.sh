#!/bin/bash

cd /app/backend
uvicorn app.main:app --host 0.0.0.0 --port 8000 & # Run backend in background

cd /app/pdf-service
./pdf-service & # Run pdf service in background

wait # Wait for background processes to finish (important for Docker)