from collections.abc import AsyncGenerator

from fastapi import Depends
from sqlalchemy.ext.asyncio import AsyncSession
from src.authentication.service import AuthenticationService
from src.core.config import Config, config
from src.core.jwt import JWTManager
from src.core.password import PasswordManager
from src.core.postgres import postgres_client
from src.users.repository import UserRepository
from src.users.services import UserService


async def get_db() -> AsyncGenerator:
    async with postgres_client.get_session() as session:
        yield session

def get_config() -> Config:
    return config

def get_user_repository(db: AsyncSession = Depends(get_db)) -> UserRepository:
    return UserRepository(db)

def get_password_manager() -> PasswordManager:
    return PasswordManager()

def get_jwt_manager(cfg: Config = Depends(get_config)) -> JWTManager:
    return JWTManager(
        secret_key=cfg.JWT_SECRET_KEY, 
        access_token_ttl=cfg.JWT_ACCESS_TOKEN_TTL
    )


def get_user_service(
    repo: UserRepository = Depends(get_user_repository), 
    pwd_manager: PasswordManager = Depends(get_password_manager)
) -> UserService:
    return UserService(repo, pwd_manager)


def get_auth_service(
    user_repo: UserRepository = Depends(get_user_repository),
    pwd_manager: PasswordManager = Depends(get_password_manager),
    jwt_manager: JWTManager = Depends(get_jwt_manager)
) -> AuthenticationService:
    return AuthenticationService(user_repo, jwt_manager, pwd_manager)