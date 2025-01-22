from sqlalchemy import (
    Column,
    Integer,
    String,
    Text,
    ForeignKey,
    JSON,
    DateTime,
    create_engine,
)
from sqlalchemy.orm import relationship, sessionmaker
from sqlalchemy.ext.declarative import declarative_base
import datetime

Base = declarative_base()

class User(Base):
    __tablename__ = "users"

    id = Column(Integer, primary_key=True, index=True)
    username = Column(String, unique=True, index=True, nullable=False)
    email = Column(String, unique=True, index=True, nullable=False)
    hashed_password = Column(String, nullable=False)
    bio = Column(Text)
    profile_picture = Column(String)

    # Relationships
    processed_documents = relationship("ProcessedDocument", back_populates="user", cascade="all, delete-orphan")


class ProcessedDocument(Base):
    __tablename__ = "processed_documents"

    id = Column(Integer, primary_key=True, index=True)
    user_id = Column(Integer, ForeignKey("users.id"), index=True, nullable=False)
    document_name = Column(String, nullable=False, index=True)
    content = Column(Text, nullable=False)
    file_path = Column(String)
    original_filename = Column(String)
    citation_style = Column(String)
    created_at = Column(DateTime, default=datetime.datetime.utcnow, nullable=False)
    updated_at = Column(DateTime, default=datetime.datetime.utcnow, onupdate=datetime.datetime.utcnow)

    # Relationships
    user = relationship("User", back_populates="processed_documents")
    metadata = relationship("Metadata", back_populates="document", uselist=False, cascade="all, delete-orphan")
    citations = relationship("Citation", back_populates="document", cascade="all, delete-orphan")
    embeddings = relationship("Embedding", back_populates="document", cascade="all, delete-orphan")


class Metadata(Base):
    __tablename__ = "metadata"

    id = Column(Integer, primary_key=True, index=True)
    document_id = Column(Integer, ForeignKey("processed_documents.id"), nullable=False, unique=True)
    title = Column(String, nullable=False)
    authors = Column(JSON, nullable=False)  # Store authors as a JSON array
    publication_date = Column(String)
    journal = Column(String)
    doi = Column(String, unique=True, index=True)  # Ensure DOI is unique
    other_metadata = Column(JSON)  # Store additional metadata as JSON

    # Relationships
    document = relationship("ProcessedDocument", back_populates="metadata")


class Citation(Base):
    __tablename__ = "citations"

    id = Column(Integer, primary_key=True, index=True)
    document_id = Column(Integer, ForeignKey("processed_documents.id"), nullable=False)
    citation_text = Column(Text, nullable=False)
    citation_style = Column(String)

    # Relationships
    document = relationship("ProcessedDocument", back_populates="citations")


class Embedding(Base):
    __tablename__ = "embeddings"

    id = Column(Integer, primary_key=True, index=True)
    document_id = Column(Integer, ForeignKey("processed_documents.id"), nullable=False)
    embedding_vector = Column(JSON, nullable=False)  # Store the embedding vector as JSON
    embedding_model = Column(String, nullable=False)  # Store the model that generated the embedding

    # Relationships
    document = relationship("ProcessedDocument", back_populates="embeddings")


# Database Configuration
DATABASE_URL = "sqlite:///./test.db"  # Change to your production database URL
engine = create_engine(DATABASE_URL, echo=True)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)


def init_db():
    Base.metadata.create_all(bind=engine)


if __name__ == "__main__":
    init_db()
