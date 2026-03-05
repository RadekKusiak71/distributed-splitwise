import psycopg
import logging

logger = logging.getLogger(__name__)

class PostgresRequestRepository:
    def __init__(self, dsn: str):
        self.dsn = dsn

    def mark_completed(self, request_id: str, output_s3_key: str):
        with psycopg.connect(self.dsn) as conn:
            with conn.cursor() as cur:
                cur.execute("""
                    UPDATE requests 
                    SET status = 'completed', 
                        output_s3_key = %s, 
                        updated_at = NOW() 
                    WHERE id = %s
                """, (output_s3_key, request_id))

    def mark_failed(self, request_id: str):
        with psycopg.connect(self.dsn) as conn:
            with conn.cursor() as cur:
                cur.execute("""
                    UPDATE requests SET status = 'failed', updated_at = NOW() WHERE id = %s
                """, (request_id,))