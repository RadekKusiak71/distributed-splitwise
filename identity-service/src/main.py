from contextlib import asynccontextmanager

from fastapi import FastAPI
from src.api.v1 import v1_auth_router, v1_user_router
from src.core.config import config
from src.core.postgres import postgres_client
from collections.abc import AsyncGenerator

@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncGenerator:
    postgres_client._connect()
    yield
    await postgres_client._disconnect()

app = FastAPI(
    debug=config.IS_DEBUG,
    name="Identity-Service for Splitwise System",
    lifespan=lifespan
)

app.include_router(v1_user_router, prefix='/api/v1')
app.include_router(v1_auth_router, prefix='/api/v1')

@app.get("/health")
def get_service_health() -> dict:
    return {"status": "ok"}

