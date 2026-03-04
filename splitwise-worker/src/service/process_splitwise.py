import logging
from typing import Protocol
import pandas as pd
from src.core.aws.sqs import IncomingMessage
from .splitwise_processor import SplitwiseProcessor

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
        log_context = {"request_id": request_id, "file_key": file_key}
        
        logger.info("Executing CSV process use case", extra=log_context)

        try:
            df = self.file_repo.get_dataframe(file_key)
            
            if df is None:
                logger.warning("File repository returned no data", extra=log_context)
                return False

            logger.info("CSV successfully loaded", extra={**log_context, "rows": len(df)})

            processor = SplitwiseProcessor(df=df)
            transfers = processor.calculate()

            result_key = file_key.replace("uploads/", "results/").replace(".csv", "_settled.csv")
            success = self.file_repo.save_transfers(result_key, transfers)
            if success:
                logger.info("Process finished and results uploaded", extra={
                    **log_context, 
                    "result_key": result_key,
                    "transfers_count": len(transfers)
                })
            
            return success

        except Exception as e:
            logger.exception("Failed to process CSV file", extra=log_context)
            return False