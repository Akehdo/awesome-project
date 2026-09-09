from pathlib import Path
from tempfile import TemporaryDirectory

from app.file_storage import FileStorage
from app.transcriber import Transcriber
from app.embedder import Embedder
from app.vector_storage import VectorStorage
from app.schemas import ListChunksRequest, ListChunksResponse, StoredChunk


class ChunkService:
    def __init__(self, vector_storage: VectorStorage):
        self._vector_storage = vector_storage

    def list_chunks(self, request: ListChunksRequest) -> ListChunksResponse:
        offset = request.offset
        if offset is not None and not isinstance(offset, int):
            offset = str(offset)

        points, next_offset = self._vector_storage.scroll(
            object_key=request.object_key,
            limit=request.limit,
            offset=offset,
        )
        return ListChunksResponse(
            points=[
                StoredChunk(id=point.id, payload=point.payload or {})
                for point in points
            ],
            next_page_offset=next_offset,
        )


class ProcessingService:
    def __init__(
            self,
            file_storage: FileStorage,
            vector_storage: VectorStorage,
            transcriber: Transcriber,
            embedder: Embedder,
    ):
        self._file_storage = file_storage
        self._vector_storage = vector_storage
        self._transcriber = transcriber
        self._embedder = embedder

    def process(self, object_key: str) -> int:
        with TemporaryDirectory() as temporary_directory:
            destination = Path(temporary_directory) / Path(object_key).name

            # 1. Скачиваем аудио из MinIO
            audio_path = self._file_storage.download(
                object_key=object_key,
                destination=destination,
            )

            # 2. Транскрибируем
            result = self._transcriber.transcribe(audio_path)

            # 3. Берём непустые сегменты Whisper
            segments = [
                segment
                for segment in result["segments"]
                if segment["text"].strip()
            ]

            # 4. Получаем тексты для embedding
            texts = [
                segment["text"].strip()
                for segment in segments
            ]

            # 5. Создаём embeddings
            embeddings = self._embedder.embed(texts)

            # 6. Создаём payload для каждого embedding
            payloads = [
                {
                    "text": segment["text"].strip(),
                    "start": segment["start"],
                    "end": segment["end"],
                    "object_key": object_key,
                }
                for segment in segments
            ]

            # 7. Сохраняем в Qdrant
            self._vector_storage.upload(
                vectors=embeddings,
                payloads=payloads,
            )

            return audio_path.stat().st_size
