from uuid import uuid4

from qdrant_client import QdrantClient, models

COLLECTION_NAME = "knowledge_base"
MODEL_NAME = "sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2"


class VectorStorage:
    def __init__(
            self,
            client: QdrantClient,
            collection_name: str,
            model_name: str,
    ):
        self._client = client
        self._collection_name = collection_name
        self._model_name = model_name

    def create_collection(self):
        if self._client.collection_exists(self._collection_name):
            return

        return self._client.create_collection(
            self._collection_name,
            vectors_config=models.VectorParams(
                size=self._client.get_embedding_size(self._model_name),
                distance=models.Distance.COSINE,
            ),
        )

    def upload(
        self,
        vectors: list[list[float]],
        payloads: list[dict],
    ) -> None:
        if len(vectors) != len(payloads):
            raise ValueError(
                "vectors and payloads must have the same length"
            )

        if not vectors:
            return

        self._client.upload_collection(
            collection_name=self._collection_name,
            vectors=vectors,
            payload=payloads,
            ids=[
                str(uuid4())
                for _ in vectors
            ],
            wait=True,
        )

    def search(
            self,
            query: str,
            limit: int = 5,
    ):
        return self._client.query_points(
            collection_name=self._collection_name,
            query=models.Document(
                text=query,
                model=self._model_name,
            ),
            limit=limit,
        ).points

    def close(self) -> None:
        self._client.close()
