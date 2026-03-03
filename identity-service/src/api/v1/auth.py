from fastapi import APIRouter, Depends, HTTPException, status
from src.api.dependencies import get_auth_service
from src.authentication.exceptions import InvalidCredentials
from src.authentication.schemas import LoginUserRequest, LoginUserResponse
from src.authentication.service import AuthenticationService

v1_router = APIRouter(prefix="/auth", tags=["Authentication"])

@v1_router.post(
    "/login",
    response_model=LoginUserResponse,
    status_code=status.HTTP_200_OK,
    summary="Authenticate user and return JWT",
    responses={
        200: {"description": "Successful login"},
        401: {"description": "Invalid username or password"},
        500: {"description": "Internal server error"}
    }
)
async def login_user(
    login_request: LoginUserRequest,
    auth_service: AuthenticationService = Depends(get_auth_service)
) -> LoginUserResponse:
    try:
        return await auth_service.login_user(login_request=login_request)
    except InvalidCredentials:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid email or password",
            headers={"WWW-Authenticate": "Bearer"},
        )
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="An unexpected error occurred."
        )