from src.core.exceptions import IdentityServiceException

class UserEmailIsAlreadyTaken(IdentityServiceException):
    def __init__(self) -> None:
        super().__init__(
            message=f"The email is already registered.",
            status_code=409,
            error_code="EMAIL_ALREADY_EXISTS"
        )