from fastembed import TextEmbedding


class Embedder:
    DEFAULT_MODEL_NAME = (
        "sentence-transformers/"
        "paraphrase-multilingual-MiniLM-L12-v2"
    )

    def __init__(
            self,
            model_name: str = DEFAULT_MODEL_NAME,
            batch_size: int = 64,
    ):
        model_info = next(
            (
                item
                for item in TextEmbedding.list_supported_models()
                if item["model"] == model_name
            ),
            None,
        )

        if model_info is None:
            raise ValueError(
                f"unsupported embedding model: {model_name}"
            )

        self._model_name = model_name
        self._batch_size = batch_size

        self._model = TextEmbedding(
            model_name=model_name,
            cuda=False,
        )

    def embed(self, texts: list[str]) -> list[list[float]]:
        if not texts:
            return []

        embeddings = self._model.embed(
            texts,
            batch_size=self._batch_size,
        )

        return [
            embedding.tolist()
            for embedding in embeddings
        ]
