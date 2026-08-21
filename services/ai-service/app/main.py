from fastapi import FastAPI

from app.api import router


app = FastAPI(title="Meeting AI Service")
app.include_router(router)
