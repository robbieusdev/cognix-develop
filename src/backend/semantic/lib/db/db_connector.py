from sqlalchemy import create_engine, Column, BigInteger, UUID, TIMESTAMP, JSON, Enum, func, String
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
import enum

Base = declarative_base()


class Status(enum.Enum):
    READY_TO_PROCESS = "READY_TO_PROCESS"
    PENDING = "PENDING"
    PROCESSING = "PROCESSING"
    COMPLETED_SUCCESSFULLY = "COMPLETED_SUCCESSFULLY"
    COMPLETED_WITH_ERRORS = "COMPLETED_WITH_ERRORS"
    DISABLED = "DISABLED"
    UNABLE_TO_PROCESS = "UNABLE_TO_PROCESS"

class Connector(Base):
    __tablename__ = 'connectors'

    id = Column(BigInteger, primary_key=True, default=func.unique_rowid())
    name = Column(String, nullable=False)
    type = Column(String(50), nullable=False)
    connector_specific_config = Column(JSON, nullable=False)
    refresh_freq = Column(BigInteger, nullable=True)
    user_id = Column(UUID(as_uuid=True), nullable=False)
    tenant_id = Column(UUID(as_uuid=True), nullable=True)
    last_successful_analyzed = Column(TIMESTAMP(timezone=False), nullable=True)
    status = Column(Enum(Status), nullable=True)
    total_docs_analyzed = Column(BigInteger, nullable=False)
    creation_date = Column(TIMESTAMP(timezone=False), nullable=False)
    last_update = Column(TIMESTAMP(timezone=False), nullable=True)
    deleted_date = Column(TIMESTAMP(timezone=False), nullable=True)

    def __repr__(self):
        return (f"<Connector(id={self.id}, name={self.name}, type={self.type}, "
                f"connector_specific_config={self.connector_specific_config}, refresh_freq={self.refresh_freq}, "
                f"user_id={self.user_id}, tenant_id={self.tenant_id}, "
                f"last_successful_index_date={self.last_successful_analyzed}, last_attempt_status={self.status}, "
                f"total_docs_indexed={self.total_docs_analyzed}, creation_date={self.creation_date}, last_update={self.last_update},"
                f"deleted_date={self.deleted_date})>")


class ConnectorCRUD:
    def __init__(self, connection_string):
        self.engine = create_engine(connection_string)
        Session = sessionmaker(bind=self.engine)
        self.session = Session()
        # IMPORTANT: Cockroach by default uses isolation level SERIALIZABLE
        # Set the isolation level to READ COMMITTED
        # self.session.connection().execution_options(isolation_level="READ COMMITTED")
        Base.metadata.create_all(self.engine)


    def insert_connector(self, **kwargs) -> int:
        new_connector = Connector(**kwargs)
        self.session.add(new_connector)
        self.session.commit()
        return new_connector.id

    def select_connector(self, connector_id: int) -> Connector | None:
        if connector_id <= 0:
            raise ValueError("ID value must be positive")
        return self.session.query(Connector).filter_by(id=connector_id).first()

    def update_connector(self, connector_id: int, **kwargs) -> int:
        if connector_id <= 0:
            raise ValueError("ID value must be positive")
        updated_connectors = self.session.query(Connector).filter_by(id=connector_id).update(kwargs)
        self.session.commit()
        return updated_connectors


    def delete_connector(self, connector_id: int) -> int:
        if connector_id <= 0:
            raise ValueError("ID value must be positive")
        deleted_connectors = self.session.query(Connector).filter_by(id=connector_id).delete()
        self.session.commit()
        return deleted_connectors
