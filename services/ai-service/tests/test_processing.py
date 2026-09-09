import unittest
from pathlib import Path
from unittest.mock import Mock

from qdrant_client import QdrantClient, models

from app.embedder import Embedder
from app.schemas import ProcessMeetingResponse
from app.service import ProcessingService
from app.transcriber import Transcriber
from app.vector_storage import VectorStorage


class ProcessingTests(unittest.TestCase):
    def test_process_downloads_transcribes_and_uploads_before_cleanup(self):
        storage = Mock()
        transcriber = Mock()
        embedder = Mock()
        vectors = Mock()
        downloaded_paths = []

        def download(*, object_key, destination):
            self.assertEqual(object_key, "meetings/one/source")
            destination.write_bytes(b"audio data")
            downloaded_paths.append(destination)
            return destination

        def transcribe(path):
            self.assertTrue(path.is_file())
            return {"segments": [
                {"text": " Hello ", "start": 0.0, "end": 1.0},
                {"text": "  ", "start": 1.0, "end": 2.0},
            ]}

        storage.download.side_effect = download
        transcriber.transcribe.side_effect = transcribe
        embedder.embed.return_value = [[0.1, 0.2, 0.3]]
        service = ProcessingService(storage, vectors, transcriber, embedder)

        file_size = service.process("meetings/one/source")

        response = ProcessMeetingResponse(
            meeting_id="one", object_key="meetings/one/source", file_size=file_size,
        )
        self.assertEqual(response.file_size, 10)
        embedder.embed.assert_called_once_with(["Hello"])
        vectors.upload.assert_called_once_with(
            vectors=[[0.1, 0.2, 0.3]],
            payloads=[{
                "text": "Hello", "start": 0.0, "end": 1.0,
                "object_key": "meetings/one/source",
            }],
        )
        self.assertFalse(downloaded_paths[0].exists())

    def test_transcriber_converts_path_to_string(self):
        transcriber = Transcriber.__new__(Transcriber)
        transcriber._device = "cpu"
        transcriber._whisper_model = Mock()
        path = Path("recording.mp3")
        transcriber.transcribe(path)
        transcriber._whisper_model.transcribe.assert_called_once_with(
            str(path), fp16=False,
        )

    def test_empty_embedding_input_does_not_call_model(self):
        embedder = Embedder.__new__(Embedder)
        embedder._model = Mock()
        self.assertEqual(embedder.embed([]), [])
        embedder._model.embed.assert_not_called()


class VectorStorageTests(unittest.TestCase):
    def test_scroll_filters_and_paginates_without_vectors(self):
        client = QdrantClient(":memory:")
        self.addCleanup(client.close)
        storage = VectorStorage(client, "test", "unused")
        self.assertEqual(storage.scroll(), ([], None))
        client.create_collection(
            collection_name="test",
            vectors_config=models.VectorParams(size=3, distance=models.Distance.COSINE),
        )
        storage.upload(
            [[1.0, 0.0, 0.0]] * 3,
            [{"object_key": "first"}, {"object_key": "second"}, {"object_key": "first"}],
        )
        first, offset = storage.scroll(object_key="first", limit=1)
        self.assertEqual(len(first), 1)
        self.assertIsNotNone(offset)
        second, next_offset = storage.scroll(object_key="first", limit=1, offset=offset)
        self.assertEqual(len(second), 1)
        self.assertIsNone(next_offset)
        self.assertNotEqual(first[0].id, second[0].id)
        for point in first + second:
            self.assertEqual(point.payload["object_key"], "first")
            self.assertIsNone(point.vector)
        self.assertEqual(storage.scroll(object_key="missing"), ([], None))

    def test_two_audio_files_keep_separate_points(self):
        client = QdrantClient(":memory:")
        self.addCleanup(client.close)
        client.create_collection(
            collection_name="test",
            vectors_config=models.VectorParams(size=3, distance=models.Distance.COSINE),
        )
        storage = VectorStorage(client, "test", "unused")
        storage.create_collection()
        storage.upload([[1.0, 0.0, 0.0]], [{"object_key": "first"}])
        storage.upload([[0.0, 1.0, 0.0]], [{"object_key": "second"}])
        points, _ = client.scroll("test", with_vectors=True)
        self.assertEqual(len(points), 2)
        self.assertEqual({point.payload["object_key"] for point in points}, {"first", "second"})
        self.assertEqual(len({point.id for point in points}), 2)
        self.assertEqual({tuple(point.vector) for point in points}, {
            (1.0, 0.0, 0.0), (0.0, 1.0, 0.0),
        })

    def test_empty_upload_and_mismatched_payloads(self):
        client = Mock()
        storage = VectorStorage(client, "test", "unused")
        storage.upload([], [])
        client.upload_collection.assert_not_called()
        with self.assertRaises(ValueError):
            storage.upload([[1.0]], [])
        client.upload_collection.assert_not_called()


if __name__ == "__main__":
    unittest.main()
