import json
import logging
from dataclasses import dataclass
from typing import Any

from mypy_boto3_sqs import SQSClient

from .base import AWSProvider

logger = logging.getLogger(__name__)

@dataclass
class IncomingMessage:
    id: str
    body: dict[str, Any]
    receipt_handle: str

class SQSMessageQueue(AWSProvider):
    def __init__(self, queue_url: str, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.queue_url = queue_url

    def receive_messages(self, max_messages: int = 1) -> list[IncomingMessage]:
        with self._get_client('sqs') as sqs:
            sqs: SQSClient
            response = sqs.receive_message(
                QueueUrl=self.queue_url,
                MaxNumberOfMessages=max_messages,
                WaitTimeSeconds=20
            )

            messages = []
            for m in response.get('Messages', []):
                raw_body = m.get('Body', '')
                
                if not raw_body or not raw_body.strip():
                    logger.warning("Skipping empty SQS message", extra={"message_id": m['MessageId']})
                    self.delete_message(m['ReceiptHandle'])
                    continue

                try:
                    parsed_body = json.loads(raw_body)
                    messages.append(
                        IncomingMessage(
                            id=m['MessageId'],
                            body=parsed_body,
                            receipt_handle=m['ReceiptHandle']
                        )
                    )
                except json.JSONDecodeError:
                    logger.error(
                        "Failed to decode SQS message JSON", 
                        extra={
                            "message_id": m['MessageId'],
                            "raw_body": raw_body
                        }
                    )
                    self.delete_message(m['ReceiptHandle'])
                    continue
            
            return messages

    def delete_message(self, receipt_handle: str) -> None:
        with self._get_client('sqs') as sqs:
            sqs: SQSClient
            sqs.delete_message(QueueUrl=self.queue_url, ReceiptHandle=receipt_handle)