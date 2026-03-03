from pydantic import BaseModel, EmailStr

class LoginUserRequest(BaseModel):
    email: EmailStr
    password: str

class LoginUserResponse(BaseModel):
    access_token: str