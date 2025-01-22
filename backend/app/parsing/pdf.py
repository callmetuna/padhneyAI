import io
import logging
import os
import tempfile
from typing import Union, Dict, Optional
from dataclasses import dataclass
from enum import Enum
from pathlib import Path
from functools import partial

from fastapi import Depends, HTTPException, UploadFile, status, BackgroundTasks
from pdfminer.high_level import extract_text_to_fp
from pdfminer.layout import LAParams
from PIL import Image, UnidentifiedImageError
import pytesseract
from tika import parser
from pdfminer.pdfparser import PDFSyntaxError
from redis import Redis
from rq import Queue
from rq.job import Job
from pydantic import BaseSettings

# Configuration class using Pydantic
class Settings(BaseSettings):
    REDIS_HOST: str = "localhost"
    REDIS_PORT: int = 6379
    TESSERACT_PATH: str = "/usr/bin/tesseract"
    UPLOAD_DIR: str = "/tmp/pdf_uploads"
    MAX_FILE_SIZE: int = 10 * 1024 * 1024  # 10MB

    class Config:
        env_file = ".env"

# Create settings instance
settings = Settings()

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Initialize Redis Queue
redis_conn = Redis(host=settings.REDIS_HOST, port=settings.REDIS_PORT)
task_queue = Queue('pdf_processing', connection=redis_conn)

class ExtractionMethod(Enum):
    PDFMINER = "pdfminer"
    TESSERACT = "tesseract"
    TIKA = "tika"

@dataclass
class ExtractionResult:
    text: str
    method: ExtractionMethod
    success: bool
    error: Optional[str] = None
    job_id: Optional[str] = None

class PDFSizeError(Exception):
    """Custom exception for PDF size validation"""
    pass

class TextExtractionError(Exception):
    """Custom exception for text extraction failures"""
    pass

class PDFExtractor:
    """Class to handle PDF text extraction with multiple methods"""
    
    SUPPORTED_MIME_TYPES = {
        "application/pdf",
        "image/png",
        "image/jpeg",
        "image/tiff"
    }

    def __init__(self):
        self.tesseract_path = settings.TESSERACT_PATH
        pytesseract.pytesseract.tesseract_cmd = self.tesseract_path
        
        # Configure PDFMiner parameters
        self.laparams = LAParams(
            all_texts=True,
            detect_vertical=True,
            word_margin=0.1,
            char_margin=2.0,
            line_margin=0.5,
            boxes_flow=0.5
        )

        # Ensure upload directory exists
        os.makedirs(settings.UPLOAD_DIR, exist_ok=True)

    async def save_upload_file(self, upload_file: UploadFile) -> Path:
        """Save uploaded file to temporary location"""
        temp_file = Path(settings.UPLOAD_DIR) / f"{next(tempfile._get_candidate_names())}.pdf"
        try:
            with temp_file.open("wb") as buffer:
                while content := await upload_file.read(8192):  # Read in chunks
                    buffer.write(content)
            return temp_file
        except Exception as e:
            if temp_file.exists():
                temp_file.unlink()
            raise e

    def validate_file_size(self, file_path: Path) -> None:
        """Validate file size"""
        if file_path.stat().st_size > settings.MAX_FILE_SIZE:
            raise PDFSizeError(
                f"File size exceeds maximum limit of {settings.MAX_FILE_SIZE/1024/1024}MB"
            )

    def extract_with_pdfminer(self, file_path: Path) -> ExtractionResult:
        """Extract text using PDFMiner.six"""
        try:
            with open(file_path, 'rb') as in_file, io.StringIO() as output_string:
                extract_text_to_fp(in_file, output_string, laparams=self.laparams)
                text = output_string.getvalue().strip()
                return ExtractionResult(
                    text=text,
                    method=ExtractionMethod.PDFMINER,
                    success=bool(text)
                )
        except PDFSyntaxError as e:
            logger.warning(f"PDFMiner syntax error: {str(e)}")
            return ExtractionResult("", ExtractionMethod.PDFMINER, False, str(e))
        except Exception as e:
            logger.error(f"PDFMiner extraction failed: {str(e)}", exc_info=True)
            return ExtractionResult("", ExtractionMethod.PDFMINER, False, str(e))

    def extract_with_tesseract(self, file_path: Path) -> ExtractionResult:
        """Extract text using Tesseract OCR"""
        try:
            image = Image.open(file_path)
            if image.mode not in ('L', 'RGB'):
                image = image.convert('RGB')
            
            text = pytesseract.image_to_string(
                image,
                config='--psm 1 --oem 3 -c preserve_interword_spaces=1'
            ).strip()
            
            return ExtractionResult(
                text=text,
                method=ExtractionMethod.TESSERACT,
                success=bool(text)
            )
        except UnidentifiedImageError as e:
            logger.warning(f"Invalid image format: {str(e)}")
            return ExtractionResult("", ExtractionMethod.TESSERACT, False, str(e))
        except Exception as e:
            logger.error(f"Tesseract extraction failed: {str(e)}", exc_info=True)
            return ExtractionResult("", ExtractionMethod.TESSERACT, False, str(e))

    def extract_with_tika(self, file_path: Path) -> ExtractionResult:
        """Extract text using Apache Tika"""
        try:
            with open(file_path, 'rb') as file:
                parsed = parser.from_buffer(file.read())
            text = parsed.get("content", "").strip()
            return ExtractionResult(
                text=text,
                method=ExtractionMethod.TIKA,
                success=bool(text)
            )
        except Exception as e:
            logger.error(f"Tika extraction failed: {str(e)}", exc_info=True)
            return ExtractionResult("", ExtractionMethod.TIKA, False, str(e))

    def process_file(self, file_path: Path) -> Dict[str, Union[str, str]]:
        """
        Process file using multiple extraction methods
        This method runs in the worker process
        """
        try:
            self.validate_file_size(file_path)
            
            extraction_methods = [
                self.extract_with_pdfminer,
                self.extract_with_tesseract,
                self.extract_with_tika
            ]
            
            errors = []
            for method in extraction_methods:
                result = method(file_path)
                if result.success:
                    return {
                        "text": result.text,
                        "method": result.method.value,
                        "status": "completed"
                    }
                errors.append(f"{result.method.value}: {result.error}")
            
            error_message = " | ".join(errors)
            raise TextExtractionError(f"All extraction methods failed: {error_message}")
            
        except Exception as e:
            logger.error(f"Processing failed: {str(e)}", exc_info=True)
            return {
                "text": "",
                "method": "none",
                "status": "failed",
                "error": str(e)
            }
        finally:
            # Clean up the temporary file
            try:
                file_path.unlink()
            except Exception as e:
                logger.error(f"Failed to delete temporary file {file_path}: {str(e)}")

async def get_job_status(job_id: str) -> Dict[str, str]:
    """Get the status of a processing job"""
    job = Job.fetch(job_id, connection=redis_conn)
    if job.is_finished:
        return job.result
    elif job.is_failed:
        return {"status": "failed", "error": str(job.exc_info)}
    else:
        return {"status": "processing"}

async def upload_document(
    file: UploadFile,
    background_tasks: BackgroundTasks
) -> Dict[str, str]:
    """Handle document upload and initiate async processing"""
    
    if file.content_type not in PDFExtractor.SUPPORTED_MIME_TYPES:
        raise HTTPException(
            status_code=status.HTTP_415_UNSUPPORTED_MEDIA_TYPE,
            detail=f"Unsupported file type. Supported types: {', '.join(PDFExtractor.SUPPORTED_MIME_TYPES)}"
        )
    
    try:
        # Initialize extractor
        extractor = PDFExtractor()
        
        # Save file to disk
        temp_file = await extractor.save_upload_file(file)
        
        # Enqueue processing job
        job = task_queue.enqueue(
            extractor.process_file,
            args=(temp_file,),
            job_timeout='10m'
        )
        
        return {
            "job_id": job.id,
            "status": "processing",
            "filename": file.filename
        }
        
    except Exception as e:
        logger.error(f"Upload failed: {str(e)}", exc_info=True)
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=str(e)
        )
    finally:
        await file.close()

# FastAPI endpoints
from fastapi import APIRouter

router = APIRouter()

@router.post("/upload/")
async def upload_endpoint(
    file: UploadFile,
    background_tasks: BackgroundTasks
) -> Dict[str, str]:
    return await upload_document(file, background_tasks)

@router.get("/status/{job_id}")
async def get_status_endpoint(job_id: str) -> Dict[str, str]:
    return await get_job_status(job_id)