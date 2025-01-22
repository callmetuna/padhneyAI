from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from fastapi.security import OAuth2PasswordBearer

from .. import models, database
from .auth import get_current_active_user, verify_password, get_password_hash, create_access_token

# Authentication router
router = APIRouter(prefix="/auth", tags=["Authentication"])

# Dependency for OAuth2 Password Bearer
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="auth/login")

# Registration Route
@router.post("/register", status_code=status.HTTP_201_CREATED, response_model=models.UserResponse)
async def register(user: models.UserCreate, db: Session = Depends(database.get_db)):
    """
    Register a new user.
    """
    # Check if the username or email already exists
    existing_user = db.query(models.User).filter(
        (models.User.username == user.username) | (models.User.email == user.email)
    ).first()
    if existing_user:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Username or email already registered.",
        )

    # Hash the user's password
    hashed_password = get_password_hash(user.password)
    
    # Create a new user instance
    new_user = models.User(
        username=user.username,
        email=user.email,
        hashed_password=hashed_password,
        bio=user.bio,
        profile_picture=user.profile_picture,
    )
    db.add(new_user)
    db.commit()
    db.refresh(new_user)

    return new_user


# Login Route
@router.post("/login", status_code=status.HTTP_200_OK)
async def login(
    form_data: OAuth2PasswordBearer = Depends(), db: Session = Depends(database.get_db)
):
    """
    Authenticate a user and issue a JWT access token.
    """
    user = (
        db.query(models.User)
        .filter(models.User.username == form_data.username)
        .first()
    )
    if not user or not verify_password(form_data.password, user.hashed_password):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid username or password.",
            headers={"WWW-Authenticate": "Bearer"},
        )

    # Create JWT token
    access_token = create_access_token(data={"user_id": user.id})
    return {"access_token": access_token, "token_type": "bearer"}


# Protected Route Example
@router.get("/protected", response_model=models.UserResponse)
async def get_protected_data(current_user: models.User = Depends(get_current_active_user)):
    """
    A protected route that requires a valid JWT access token.
    """
    return current_user  # Return the authenticated user data
