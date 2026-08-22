from pathlib import Path

from minio import Minio


class MinioStorage:
    def __init__(self, minio_client: Minio, bucket_name: str):
        self._client = minio_client
        self._bucket_name = bucket_name

    def download(self, object_key: str, destination: Path) -> Path:
        destination.parent.mkdir(parents=True, exist_ok=True)

        self._client.fget_object(
            bucket_name=self._bucket_name,
            object_name=object_key,
            file_path=str(destination),
        )

        return destination
