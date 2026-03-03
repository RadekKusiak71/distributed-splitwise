from src.core.exceptions import IdentityServiceException

class InvalidCredentials(IdentityServiceException):
    def __init__(self):
        super().__init__(
            message=f"Invalid credentials were provided",
            status_code=401,
            error_code="INVALID_CREDENTIALS"
        )