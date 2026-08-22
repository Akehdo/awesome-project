from functools import lru_cache

from minio import Minio

from app.config import settings
from app.service import ProcessingService
from app.storage import MinioStorage


@lru_cache
def get_processing_service() -> ProcessingService:
    client = Minio(
        endpoint=settings.minio_endpoint,
        access_key=settings.minio_root_user,
        secret_key=settings.minio_root_password,
        secure=settings.minio_use_ssl,
    )

    storage = MinioStorage(
        minio_client=client,
        bucket_name=settings.minio_bucket_name,
    )

    return ProcessingService(storage=storage)
