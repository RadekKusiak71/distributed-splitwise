import uuid
from datetime import datetime
from pydantic import Field, BaseModel, EmailStr, model_validator
from typing import Self

class RegisterUserRequest(BaseModel):
    email: EmailStr
    password: str = Field(min_length=8, max_length=120)
    password2: str

    @model_validator(mode="after")
    def validate_password(self) -> Self:
        if self.password != self.password2:
            raise ValueError("password and password confirmation does not match")
        return self

class RegisterUserResponse(BaseModel):
    id: uuid.UUID
    email: EmailStr
    created_at: datetime    
