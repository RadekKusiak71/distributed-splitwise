from pydantic_settings import BaseSettings


class Config(BaseSettings):
    IS_DEBUG: bool
    DB_DSN: str
    JWT_SECRET_KEY: str
    JWT_ACCESS_TOKEN_TTL: int

config = Config()
