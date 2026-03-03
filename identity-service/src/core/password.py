import bcrypt

class PasswordManager:
    
    def hash(self, raw_password: str) -> str:
        pwd_bytes: bytes = raw_password.encode("utf-8")
        hashed: bytes = bcrypt.hashpw(pwd_bytes, bcrypt.gensalt())
        return hashed.decode("utf-8")

    def verify(self, raw_password: str, hashed_password: str) -> bool:
        return bcrypt.checkpw(
            raw_password.encode("utf-8"), 
            hashed_password.encode("utf-8")
        )