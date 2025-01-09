from fastapi import APIRouter, HTTPException, UploadFile, File
from typing import Dict
import requests
import json
import os

router = APIRouter()

PDF_SERVICE_URL = os.environ.get("PDF_SERVICE_URL", "http://localhost:8081/pdf") #Use env variable or default if not set

def process_pdf_with_service(pdf_path: str = None, pdf_file: UploadFile = None):
    try:
        if pdf_path:
            headers = {'Content-type': 'application/json'}
            data = {"pdf_path": pdf_path}
            r = requests.post(PDF_SERVICE_URL, data=json.dumps(data), headers = headers)
        elif pdf_file:
            files = {'pdf_file': (pdf_file.filename, pdf_file.file, pdf_file.content_type)}
            r = requests.post(PDF_SERVICE_URL, files=files)
        else:
            raise HTTPException(status_code=400, detail="Either pdf_path or pdf_file must be provided")

        r.raise_for_status()
        return r.json()
    except requests.exceptions.RequestException as e:
        print(f"Error communicating with PDF service: {e}")
        raise HTTPException(status_code=500, detail=f"Error processing PDF: {e}")
    except Exception as e:
        print(f"Unexpected error: {e}")
        raise HTTPException(status_code=500, detail=f"Unexpected error processing PDF: {e}")


@router.post("/process_pdf")
async def process_pdf(pdf_path: str = None, pdf_file: UploadFile = File(None)):
    pdf_response = process_pdf_with_service(pdf_path, pdf_file)
    if pdf_response:
        return pdf_response
    else:
        raise HTTPException(status_code=500, detail="Failed to process PDF")