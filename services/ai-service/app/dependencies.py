from functools import lru_cache

from minio import Minio
from qdrant_client import QdrantClient

from app.config import settings
from app.embedder import Embedder
from app.file_storage import FileStorage
from app.service import ChunkService, ProcessingService
from app.transcriber import Transcriber
from app.vector_storage import VectorStorage


@lru_cache
def get_processing_service() -> ProcessingService:
    file_storage_client = Minio(
        endpoint=settings.minio_endpoint,
        access_key=settings.minio_root_user,
        secret_key=settings.minio_root_password,
        secure=settings.minio_use_ssl,
    )

    file_storage = FileStorage(
        client=file_storage_client,
        bucket_name=settings.minio_bucket_name,
    )

    vector_storage = get_vector_storage()
    vector_storage.create_collection()

    return ProcessingService(
        file_storage=file_storage,
        vector_storage=vector_storage,
        transcriber=get_transcriber(),
        embedder=get_embedder(),
    )


@lru_cache
def get_vector_storage() -> VectorStorage:
    vector_storage_client = QdrantClient(
        host=settings.qdrant_host,
        port=settings.qdrant_http_port,
        grpc_port=settings.qdrant_grpc_port,
        prefer_grpc=settings.qdrant_transport == "grpc",
    )

    return VectorStorage(
        client=vector_storage_client,
        collection_name="knowledge_base",
        model_name=Embedder.DEFAULT_MODEL_NAME,
    )

@lru_cache
def get_chunk_service() -> ChunkService:
    return ChunkService(vector_storage=get_vector_storage())


@lru_cache
def get_transcriber() -> Transcriber:
    return Transcriber(whisper_model_name="small")


@lru_cache
def get_embedder() -> Embedder:
    return Embedder()
