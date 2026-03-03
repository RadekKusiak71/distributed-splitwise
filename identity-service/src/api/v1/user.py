from starlette import status
from fastapi import APIRouter, Depends, HTTPException
from src.api.dependencies import get_user_service
from src.users.services import UserService
from src.users.exceptions import UserEmailIsAlreadyTaken
from src.users.schemas import RegisterUserRequest, RegisterUserResponse

v1_router = APIRouter(prefix="/users", tags=["Users"])

@v1_router.post(
    "/register", 
    status_code=status.HTTP_201_CREATED,
    response_model=RegisterUserResponse,
    summary="Register a new user",
    description="Creates a new user account in the system and hashes the password.",
    responses={
        201: {"description": "User successfully created"},
        409: {"description": "Conflict: Email already exists"},
        500: {"description": "Internal server error"}
    }
)
async def register_user(
    request: RegisterUserRequest, 
    service: UserService = Depends(get_user_service)
) -> RegisterUserResponse:
    """
    Registers a new user into the LoL Sentinel system.
    
    - **email**: Must be a valid and unique email address.
    - **password**: Will be salted and hashed using bcrypt.
    """
    try:
        return await service.register_user(register_req=request)
    
    except UserEmailIsAlreadyTaken as e:
        raise HTTPException(
            status_code=status.HTTP_409_CONFLICT, 
            detail=str(e)
        )
    except Exception:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="An unexpected error occurred on the server."
        )