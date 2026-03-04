import logging
from src.core.logger import setup_logging
from src.core.config import Config
from src.core.aws.s3 import S3FileRepository
from src.core.aws.sqs import SQSMessageQueue
from src.service.process_splitwise import ProcessCSVUseCase
from .worker import Worker

logger = logging.getLogger(__name__)

if __name__ == "__main__":
    setup_logging(log_level="INFO")
    logger = logging.getLogger(__name__)
    logger.info("Worker process initialized")

    config = Config()

    s3_repo = S3FileRepository(
        bucket_name=config.AWS_S3_BUCKET_NAME,
        access_key=config.AWS_ACCESS_KEY,
        secret_key=config.AWS_SECRET_KEY,
        region=config.AWS_S3_REGION,
    )
    
    sqs_queue = SQSMessageQueue(
        queue_url=config.AWS_QUEUE_URL,
        access_key=config.AWS_ACCESS_KEY,
        secret_key=config.AWS_SECRET_KEY,
        region=config.AWS_S3_REGION,
    )

    processor = ProcessCSVUseCase(file_repo=s3_repo)
    worker = Worker(queue=sqs_queue, use_case=processor)

    worker.run_forever()