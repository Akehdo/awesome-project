from typing import Annotated

from fastapi import APIRouter, Depends

from app.dependencies import get_processing_service
from app.schemas import (
    ProcessMeetingRequest,
    ProcessMeetingResponse,
)
from app.service import ProcessingService


router = APIRouter()

ProcessingServiceDependency = Annotated[
    ProcessingService,
    Depends(get_processing_service),
]


@router.post(
    "/meetings/process",
    response_model=ProcessMeetingResponse,
    tags=["meetings"],
)
def process_meeting(
    body: ProcessMeetingRequest,
    processing_service: ProcessingServiceDependency,
) -> ProcessMeetingResponse:
    file_size = processing_service.process(body.object_key)

    return ProcessMeetingResponse(
        meeting_id=body.meeting_id,
        object_key=body.object_key,
        file_size=file_size,
    )


@router.get("/health", tags=["system"])
async def health() -> dict[str, str]:
    return {
        "status": "ok",
        "service": "meeting-ai-service",
    }
