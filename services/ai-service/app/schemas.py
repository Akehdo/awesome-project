from typing import Any
from uuid import UUID

from pydantic import BaseModel, Field


class ProcessMeetingRequest(BaseModel):
    meeting_id: str
    object_key: str = Field(min_length=1)


class ProcessMeetingResponse(BaseModel):
    meeting_id: str
    object_key: str
    file_size: int


class ListChunksRequest(BaseModel):
    object_key: str | None = Field(default=None, min_length=1)
    limit: int = Field(default=10, ge=1, le=100)
    offset: UUID | int | None = None


class StoredChunk(BaseModel):
    id: str | int
    payload: dict[str, Any]


class ListChunksResponse(BaseModel):
    points: list[StoredChunk]
    next_page_offset: str | int | None
