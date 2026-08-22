from pathlib import Path
from tempfile import TemporaryDirectory

from app.storage import MinioStorage


class ProcessingService:
    def __init__(self, storage: MinioStorage):
        self._storage = storage

    def process(self, object_key: str) -> int:
        with TemporaryDirectory() as temporary_directory:
            destination = Path(temporary_directory) / Path(object_key).name

            audio_path = self._storage.download(
                object_key=object_key,
                destination=destination,
            )

            return audio_path.stat().st_size
