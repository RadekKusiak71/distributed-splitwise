from typing import Protocol

from src.authentication.exceptions import InvalidCredentials
from src.authentication.schemas import LoginUserRequest, LoginUserResponse
from src.users.models import User


class UserRepository(Protocol):
    async def find_by_email(self, email: str) -> User | None: ...

class JWTManager(Protocol):
    def generate(self, user_id: str) -> str: ...
    def verify(self, token: str) -> bool: ...

class PasswordManager(Protocol):
    def verify(self, raw_password: str, hashed_password: str) -> bool: ...

class AuthenticationService:
    def __init__(
        self, 
        user_repo: UserRepository, 
        jwt_manager: JWTManager, 
        password_manager: PasswordManager
    ) -> None:
        self.user_repo=user_repo
        self.jwt_manager=jwt_manager
        self.password_manager=password_manager

    async def login_user(self, login_request: LoginUserRequest) -> LoginUserResponse:
        user: User | None = await self.user_repo.find_by_email(email=login_request.email)
        if not user:
            raise InvalidCredentials()
        if not self.password_manager.verify(login_request.password, user.password):
            raise InvalidCredentials()
        return LoginUserResponse(
            access_token=self.jwt_manager.generate(user_id=str(user.id))
        )