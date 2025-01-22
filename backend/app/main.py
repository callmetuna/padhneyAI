from fastapi import FastAPI, File, UploadFile, HTTPException, Depends
from typing import List
import os
import io
from pypdf import PdfReader
import docx2txt
from unstructured.partition.auto import partition
from sqlalchemy.orm import Session
from . import models, database

app = FastAPI()

# Database dependency
def get_db():
    db = database.SessionLocal()
    try:
        yield db
    finally:
        db.close()

@app.post("/upload", response_model=List[dict])
async def upload_files(files: List[UploadFile] = File(...), db: Session = Depends(get_db)):
    extracted_data = []
    for file in files:
        try:
            contents = await file.read()
            filename = file.filename
            file_extension = os.path.splitext(filename)[1].lower()

            if file_extension == ".pdf":
                reader = PdfReader(io.BytesIO(contents))
                text = ""
                for page in reader.pages:
                    text += page.extract_text()
                extracted_data.append({"filename": filename, "content": text, "file_type": "pdf"})

            elif file_extension == ".docx":
                text = docx2txt.process(io.BytesIO(contents))
                extracted_data.append({"filename": filename, "content": text, "file_type": "docx"})

            elif file_extension in [".txt", ".rtf", ".html", ".htm", ".csv", ".tsv", ".eml"]:
                elements = partition(filename=filename, file=io.BytesIO(contents)) #unstructured
                text = "\n\n".join([str(el) for el in elements])
                extracted_data.append({"filename": filename, "content": text, "file_type": file_extension})
            else:
                raise HTTPException(status_code=400, detail=f"Unsupported file type: {file_extension}")
            
            #store data in database
            db_document = models.ProcessedDocument(document_name=filename, content=text)
            db.add(db_document)
            db.commit()
            db.refresh(db_document)

        except Exception as e:
            raise HTTPException(status_code=500, detail=f"Error processing {file.filename}: {str(e)}")

    return extracted_data

app.include_router(query.router)
app.include_router(pdf.router)
app.include_router(auth.router)
app.include_router(user.router)
