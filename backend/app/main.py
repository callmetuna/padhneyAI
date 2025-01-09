from fastapi import FastAPI
from app.routers import query, pdf

app = FastAPI()

app.include_router(query.router)
app.include_router(pdf.router)