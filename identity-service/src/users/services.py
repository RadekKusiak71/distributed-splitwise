from typing import Protocol

from src.users.exceptions import UserEmailIsAlreadyTaken
from src.users.models import User
from src.users.schemas import RegisterUserRequest, RegisterUserResponse

class UserRepository(Protocol):
    async def find_by_email(self, email: str) -> User | None: ...
    async def create(self, user: User) -> User: ...

class PasswordManager(Protocol):
    def hash(self, raw_password: str) -> str: ...

class UserService:
    def __init__(
        self, 
        user_repo: UserRepository, 
        password_manager: PasswordManager
    ) -> None:
        self.user_repo=user_repo
        self.password_manager=password_manager

    async def register_user(self, register_req: RegisterUserRequest) -> RegisterUserResponse:
        user: User | None = await self.user_repo.find_by_email(email=register_req.email)
        if user:
            raise UserEmailIsAlreadyTaken()
        hashedPassword: str = self.password_manager.hash(raw_password=register_req.password)
        newUser: User = await self.user_repo.create(user=User(email=register_req.email, password=hashedPassword))
        return RegisterUserResponse(
            id=newUser.id, 
            email=newUser.email, 
            created_at=newUser.created_at
        )