import logging
from typing import Protocol
import pandas as pd
from src.core.aws.sqs import IncomingMessage

logger = logging.getLogger(__name__)

class FileRepository(Protocol):
    def get_dataframe(self, file_key: str) -> pd.DataFrame | None: ...

class MessageQueue(Protocol):
    def receive_messages(self, max_messages: int = 1) -> list[IncomingMessage]: ...
    def delete_message(self, receipt_handle: str) -> None: ...

class ProcessCSVUseCase:
    def __init__(self, file_repo: FileRepository) -> None:
        self.file_repo = file_repo

    def execute(self, request_id: str, file_key: str) -> bool:
        logger.info("Executing CSV process use case", extra={
            "request_id": request_id,
            "file_key": file_key
        })

        try:
            df = self.file_repo.get_dataframe(file_key)
            
            if df is not None:
                logger.info("CSV successfully loaded and parsed", extra={
                    "request_id": request_id,
                    "rows_count": len(df),
                    "columns": list(df.columns)
                })
                
                return True

            logger.warning("File repository returned no data", extra={
                "request_id": request_id,
                "file_key": file_key
            })
            return False

        except Exception:
            logger.exception("Failed to process CSV file", extra={
                "request_id": request_id,
                "file_key": file_key
            })
            return False