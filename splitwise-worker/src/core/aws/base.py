import contextlib
from typing import Any, Generator

import boto3


class AWSProvider:
    def __init__(self, access_key: str, secret_key: str, region: str) -> None:
        self.access_key = access_key
        self.secret_key = secret_key
        self.region = region

    @contextlib.contextmanager
    def _get_client(self, service: str) -> Generator[Any, None, None]:
        client = boto3.client(
            service,
            aws_access_key_id=self.access_key,
            aws_secret_access_key=self.secret_key,
            region_name=self.region
        )
        try:
            yield client
        finally:
            client.close()