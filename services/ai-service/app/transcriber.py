import torch
import whisper
from pathlib import Path


class Transcriber:
    WHISPER_MODEL_NAME = "small"

    def __init__(self, whisper_model_name: str):
        self._device = "cuda" if torch.cuda.is_available() else "cpu"

        if self._device == "cuda":
            print(f"Using GPU: {torch.cuda.get_device_name(0)}")
        else:
            print("CUDA is unavailable; using CPU")

        self._whisper_model = whisper.load_model(
            whisper_model_name,
            device=self._device,
        )

    def transcribe(self, path: Path):
        return self._whisper_model.transcribe(
            str(path),
            fp16=self._device == "cuda",
        )