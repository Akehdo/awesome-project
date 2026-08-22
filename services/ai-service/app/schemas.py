from pydantic import BaseModel, Field


class ProcessMeetingRequest(BaseModel):
    meeting_id: str
    object_key: str = Field(min_length=1)


class ProcessMeetingResponse(BaseModel):
    meeting_id: str
    object_key: str
    file_size: int
