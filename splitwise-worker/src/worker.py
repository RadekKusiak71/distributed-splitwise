import logging

from src.service.process_splitwise import MessageQueue, ProcessCSVUseCase

logger = logging.getLogger(__name__)

class Worker:
    def __init__(self, queue: MessageQueue, use_case: ProcessCSVUseCase) -> None:
        self.queue = queue
        self.use_case = use_case

    def run_forever(self) -> None:
        logger.info("Worker starting... Polling SQS messages.")
        
        while True:
            try:
                messages = self.queue.receive_messages(max_messages=1)
                
                if not messages:
                    logger.debug("No messages in queue. Polling again...")
                    continue

                for msg in messages:
                    logger.info("Received message", extra={"message_id": msg.id})
                    
                    req_id = msg.body.get("request_id")
                    file_key = msg.body.get("s3_file_key")

                    if not req_id or not file_key:
                        logger.error(
                            "Invalid message format", 
                            extra={"body": msg.body, "message_id": msg.id}
                        )
                        continue

                    logger.info(
                        "Processing request", 
                        extra={"request_id": req_id, "file_key": file_key}
                    )
                    
                    if self.use_case.execute(request_id=req_id, file_key=file_key):
                        logger.info(
                            "Success! Deleting message", 
                            extra={"message_id": msg.id, "request_id": req_id}
                        )
                        self.queue.delete_message(msg.receipt_handle)
                    else:
                        logger.warning(
                            "Use case failed. Message remains in queue.", 
                            extra={"request_id": req_id}
                        )

            except Exception as e:
                logger.exception("CRITICAL ERROR in Worker Loop")
                continue