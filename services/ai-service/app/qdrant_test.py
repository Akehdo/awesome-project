from pathlib import Path

import torch
import whisper
from qdrant_client import QdrantClient, models


COLLECTION_NAME = "knowledge_base"
WHISPER_MODEL_NAME = "small"

# Если аудио будет русское/казахское, лучше использовать multilingual embedding.
MODEL_NAME = "sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2"

AUDIO_PATH = str(Path(__file__).with_name("audio.mp3"))


# -------------------------
# 1. Whisper
# -------------------------

print("Loading Whisper...")

device = "cuda" if torch.cuda.is_available() else "cpu"

if device == "cuda":
    print(f"Using GPU: {torch.cuda.get_device_name(0)}")
else:
    print("CUDA is unavailable; using CPU")

whisper_model = whisper.load_model(
    WHISPER_MODEL_NAME,
    device=device,
)

print("Transcribing audio...")

result = whisper_model.transcribe(
    AUDIO_PATH,
    fp16=device == "cuda",
)

print("\nFull transcription:")
print(result["text"])


# -------------------------
# 2. Get Whisper segments
# -------------------------

documents = []
payloads = []

for index, segment in enumerate(result["segments"]):
    text = segment["text"].strip()

    if not text:
        continue

    documents.append(text)

    payloads.append(
        {
            "document": text,
            "chunk_index": index,
            "start": segment["start"],
            "end": segment["end"],
        }
    )


print(f"\nCreated {len(documents)} chunks")

for payload in payloads[:5]:
    print(
        f"[{payload['start']:.1f}s - {payload['end']:.1f}s] "
        f"{payload['document']}"
    )


# -------------------------
# 3. Qdrant
# -------------------------

client = QdrantClient(":memory:")

client.create_collection(
    collection_name=COLLECTION_NAME,
    vectors_config=models.VectorParams(
        size=client.get_embedding_size(MODEL_NAME),
        distance=models.Distance.COSINE,
    ),
)


# -------------------------
# 4. Embeddings + upload
# -------------------------

client.upload_collection(
    collection_name=COLLECTION_NAME,

    vectors=[
        models.Document(
            text=document,
            model=MODEL_NAME,
        )
        for document in documents
    ],

    payload=payloads,

    ids=list(range(len(documents))),
)


# -------------------------
# 5. Search
# -------------------------

query_text = input("\nAsk something about the audio: ")

search_results = client.query_points(
    collection_name=COLLECTION_NAME,
    query=models.Document(
        text=query_text,
        model=MODEL_NAME,
    ),
    limit=5,
).points


# -------------------------
# 6. Results
# -------------------------

print(f"\nQuery: {query_text}\n")

for result in search_results:
    payload = result.payload

    print(
        f"Score: {result.score:.4f}\n"
        f"Time: {payload['start']:.1f}s - {payload['end']:.1f}s\n"
        f"Text: {payload['document']}\n"
    )
