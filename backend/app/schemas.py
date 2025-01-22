from pydantic import BaseModel, EmailStr
from typing import Optional, List


# User Schemas
class UserBase(BaseModel):
    username: str
    email: EmailStr
    bio: Optional[str] = None
    profile_picture: Optional[str] = None


class UserCreate(UserBase):
    password: str  # Password field for registration


class UserResponse(UserBase):
    id: int

    class Config:
        orm_mode = True  # Enables ORM serialization for SQLAlchemy models


# Processed Document Schemas
class MetadataSchema(BaseModel):
    title: str
    authors: List[str]
    publication_date: Optional[str] = None
    journal: Optional[str] = None
    doi: Optional[str] = None
    other_metadata: Optional[dict] = None


class CitationSchema(BaseModel):
    citation_text: str
    citation_style: Optional[str] = None


class ProcessedDocumentSchema(BaseModel):
    id: int
    document_name: str
    content: str
    citations: List[CitationSchema]
    metadata: MetadataSchema

    class Config:
        orm_mode = True
