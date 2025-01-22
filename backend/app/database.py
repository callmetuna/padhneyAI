import os
from dotenv import load_dotenv
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker
from sqlalchemy.exc import OperationalError

from .models import Base, User, ProcessedDocument

# Load environment variables
load_dotenv()

# Fetch the database URL from environment variables
DATABASE_URL = os.getenv("DATABASE_URL")

if not DATABASE_URL:
    raise ValueError("DATABASE_URL is not set in the environment variables")

# Async SQLAlchemy engine
engine = create_async_engine(
    DATABASE_URL,
    echo=False,  # Set to True only for debugging in development
    pool_pre_ping=True,  # Ensures database connections are alive
)

# Async session maker
SessionLocal = sessionmaker(
    autocommit=False,
    autoflush=False,
    bind=engine,
    class_=AsyncSession,
    expire_on_commit=False,  # Prevents unloading objects after a commit
)

async def create_tables():
    """
    Creates all tables defined in the SQLAlchemy models.
    """
    try:
        async with engine.begin() as conn:
            print("Creating tables...")
            await conn.run_sync(Base.metadata.create_all)
            print("Tables created successfully")
    except OperationalError as e:
        print(f"Error creating tables: {e}")
        raise

async def get_db():
    """
    Provides an asynchronous session for database interactions.
    Automatically closes the session after use.
    """
    async with SessionLocal() as session:
        try:
            yield session
        except Exception as e:
            await session.rollback()
            print(f"Error during database session: {e}")
            raise
        finally:
            await session.close()

def test_database_connection():
    """
    Tests the database connection and raises an error if the connection fails.
    """
    try:
        with engine.sync_engine.connect() as conn:
            conn.execute("SELECT 1")  # Simple query to test connection
            print("Database connection successful")
    except OperationalError as e:
        raise RuntimeError(f"Database connection failed: {e}")
