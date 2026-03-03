import contextlib
import logging

from sqlalchemy.ext.asyncio import AsyncAttrs
from sqlalchemy.ext.asyncio import AsyncEngine, create_async_engine
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker
from sqlalchemy.orm import DeclarativeBase

from collections.abc import AsyncGenerator
from src.core.config import config

logger = logging.getLogger(__name__)

class Base(AsyncAttrs, DeclarativeBase):
    ...

class PostgresConnection:
    def __init__(self, dsn: str, echo: bool) -> None:
        self._dsn: str = dsn
        self._echo: bool = echo
        self._engine: AsyncEngine | None = None
        self._sessionmaker: async_sessionmaker | None = None

    @contextlib.asynccontextmanager
    async def get_session(self) -> AsyncGenerator[AsyncSession, None, None]:
        if not self._sessionmaker:
            logger.error("Database is not connected yet.")
            raise RuntimeError("Database is not connected yet.")

        async with self._sessionmaker() as session:
            try:
                yield session
            except Exception as exc:
                logger.exception(f"Exception occured while operating with db session: {str(exc)}")
                await session.rollback()
                raise
            finally:
                await session.close()

    def _connect(self) -> None:
        if self._engine and self._sessionmaker:
            logger.exception("Postgres database is already connected")
            raise RuntimeError("Postgres database is already connected")
        self._engine = create_async_engine(url=self._dsn, echo=self._echo)
        self._sessionmaker = async_sessionmaker(bind=self._engine, expire_on_commit=False)


    async def _disconnect(self):
        if not self._engine:
            logger.exception("Postgres database is already discconnected")
            raise RuntimeError("Postgres database is already disconnected")
        
        await self._engine.dispose()
        self._engine = None
        self._sessionmaker = None

postgres_client = PostgresConnection(dsn=config.DB_DSN, echo=config.IS_DEBUG)