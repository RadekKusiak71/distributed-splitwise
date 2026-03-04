from .base import AWSProvider
import io
import pandas as pd
from mypy_boto3_s3 import S3Client

class S3FileRepository(AWSProvider):

    def __init__(self, bucket_name: str, *args, **kwargs) -> None:
        super().__init__(*args, **kwargs)
        self.bucket_name = bucket_name
    
    def get_dataframe(self, file_key: str) -> pd.DataFrame | None:
        with self._get_client('s3') as s3:
            s3: S3Client
            try:
                response = s3.get_object(Bucket=self.bucket_name, Key=file_key)
                return pd.read_csv(io.BytesIO(response['Body'].read()))
            except Exception as e:
                print(f"S3 Error: {e}")
                return None