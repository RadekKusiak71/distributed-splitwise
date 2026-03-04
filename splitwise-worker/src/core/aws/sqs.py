import json
from dataclasses import dataclass
from typing import Any

from mypy_boto3_sqs import SQSClient

from .base import AWSProvider


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

            return [
                IncomingMessage(
                    id=m['MessageId'],
                    body=json.loads(m['Body']),
                    receipt_handle=m['ReceiptHandle']
                ) for m in response.get('Messages', [])
            ]

    def delete_message(self, receipt_handle: str) -> None:
        with self._get_client('sqs') as sqs:
            sqs: SQSClient
            sqs.delete_message(QueueUrl=self.queue_url, ReceiptHandle=receipt_handle)