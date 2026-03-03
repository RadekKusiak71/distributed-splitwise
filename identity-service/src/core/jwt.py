import jwt
from datetime import datetime, timedelta, timezone
from typing import Any

class InvalidTokenException(Exception):
    pass

class TokenExpiredException(InvalidTokenException):
    pass

class JWTManager:
    def __init__(self, secret_key: str, algorithm: str = "HS256", access_token_ttl: int = 1440) -> None:
        self.secret_key = secret_key
        self.algorithm = algorithm
        self.access_token_ttl = access_token_ttl

    def generate(self, user_id: int) -> str:
        payload: dict = {
            "sub": str(user_id),
            "iat": datetime.now(timezone.utc),
            "iss": "identity-service-splitwise",
            "exp": datetime.now(timezone.utc) + timedelta(seconds=self.access_token_ttl),
        }
        return jwt.encode(payload, self.secret_key, algorithm=self.algorithm)

    def verify(self, token: str) -> dict[str, Any]:
        try:
            decoded_data: dict = jwt.decode(
                token, 
                self.secret_key, 
                algorithms=[self.algorithm],
                issuer="identity-service-splitwise"
            )
            return decoded_data
        
        except jwt.ExpiredSignatureError:
            raise TokenExpiredException("Token has expired")
        except (jwt.InvalidTokenError, jwt.DecodeError) as e:
            raise InvalidTokenException(f"Invalid token: {str(e)}")
