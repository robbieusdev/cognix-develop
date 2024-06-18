from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker, scoped_session
from contextlib import contextmanager

class ConnectionManager:
    _instance = None

    def __new__(cls, connection_string):
        if cls._instance is None:
            cls._instance = super(ConnectionManager, cls).__new__(cls)
            cls._instance._engine = create_engine(connection_string, pool_size=20, max_overflow=0)
            cls._instance._session_factory = scoped_session(sessionmaker(bind=cls._instance._engine))
        return cls._instance

    @contextmanager
    def get_session(self):
        session = self._session_factory()
        try:
            yield session
            session.commit()
        except Exception:
            session.rollback()
            raise
        finally:
            session.close()
