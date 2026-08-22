from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    minio_endpoint: str
    minio_bucket_name: str
    minio_root_user: str
    minio_root_password: str
    minio_use_ssl: bool

    model_config = SettingsConfigDict(
        env_file="../../.env",
        env_file_encoding="utf-8",
    )


settings = Settings()  # pyright: ignore[reportCallIssue]
