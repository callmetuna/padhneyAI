from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
import uvicron
from typing import Dict

router = APIRouter()

@router.post("/query")
async def process_query(query: str):
    if not query:
        raise HTTPException(status_code=400, detail="Query is required")
    # Placeholder for RAG logic (will be implemented later)
    return {"response": f"You asked: {query}", "sources": []}