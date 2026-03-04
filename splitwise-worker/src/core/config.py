from decouple import config


class Config:

    POSTGRES_DNS = config("POSTGRES_DNS")
    AWS_ACCESS_KEY = config("AWS_ACCESS_KEY")
    AWS_SECRET_KEY = config("AWS_SECRET_KEY")
    AWS_S3_BUCKET_NAME = config("AWS_S3_BUCKET_NAME")
    AWS_S3_REGION = config("AWS_S3_REGION")
    AWS_QUEUE_URL = config("AWS_QUEUE_URL")
    