
class IdentityServiceException(Exception):
    def __init__(
        self, 
        message: str, 
        status_code: int = 500, 
        error_code: str = "internal_server_error"
    ):
        self.message = message
        self.status_code = status_code
        self.error_code = error_code
        super().__init__(self.message)