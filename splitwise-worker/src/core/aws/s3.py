import io
import logging

import pandas as pd
from mypy_boto3_s3 import S3Client

from .base import AWSProvider

logger = logging.getLogger(__name__)

class S3FileRepository(AWSProvider):

    def __init__(self, bucket_name: str, *args, **kwargs) -> None:
        super().__init__(*args, **kwargs)
        self.bucket_name = bucket_name
    
    def get_dataframe(self, file_key: str) -> pd.DataFrame | None:
        try:
            with self._get_client('s3') as s3:
                s3: S3Client
                logger.info(f"Downloading file from S3: {self.bucket_name}/{file_key}")
                response = s3.get_object(Bucket=self.bucket_name, Key=file_key)
                return pd.read_csv(io.BytesIO(response['Body'].read()))
        except Exception as e:
            logger.error(f"Failed to download or parse CSV from S3 ({file_key}): {str(e)}")
            return None

    def save_transfers(self, file_key: str, transfers: list) -> bool:
        try:
            data = [
                {
                    "debetor": t.debetor, 
                    "spender": t.spender, 
                    "value": float(t.value)
                } for t in transfers
            ]
            df_result = pd.DataFrame(data)

            csv_buffer = io.StringIO()
            df_result.to_csv(csv_buffer, index=False)

            with self._get_client('s3') as s3:
                s3: S3Client
                logger.info(f"Uploading results to S3: {self.bucket_name}/{file_key}")
                s3.put_object(
                    Bucket=self.bucket_name,
                    Key=file_key,
                    Body=csv_buffer.getvalue()
                )
            
            logger.info(f"Successfully saved transfers to {file_key}")
            return True

        except Exception as e:
            logger.error(f"S3 Upload error for key {file_key}: {str(e)}", exc_info=True)
            return False