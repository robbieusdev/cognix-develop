from typing import List

from sqlalchemy import create_engine, Column, Integer, BigInteger, Text, Boolean, UUID, TIMESTAMP, func
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker

Base = declarative_base()


class Document(Base):
    __tablename__ = 'documents'

    id = Column(Integer, primary_key=True, autoincrement=True)
    parent_id = Column(BigInteger, nullable=True)
    connector_id = Column(BigInteger, nullable=False)
    source_id = Column(Text, nullable=False)
    url = Column(Text, nullable=True)
    signature = Column(Text, nullable=True)
    chunking_session = Column(UUID(as_uuid=True), nullable=True)
    analyzed = Column(Boolean, nullable=False, default=False)
    creation_date = Column(TIMESTAMP(timezone=False), nullable=False, default=func.now())
    last_update = Column(TIMESTAMP(timezone=False), nullable=True)

    def __repr__(self):
        return (f"<Document(id={self.id}, parent_id={self.parent_id}, connector_id={self.connector_id}, "
                f"source_id={self.source_id}, url={self.url}, signature={self.signature}, "
                f"chunking_session={self.chunking_session}, analyzed={self.analyzed}, "
                f"creation_date={self.creation_date}, last_update={self.last_update})>")


class DocumentCRUD:
    def __init__(self, connection_string):
        self.engine = create_engine(connection_string)
        Session = sessionmaker(bind=self.engine)
        self.session = Session()
        # IMPORTANT: Cockroach by default uses isolation level SERIALIZABLE
        # Set the isolation level to READ COMMITTED
        # self.session.connection().execution_options(isolation_level="READ COMMITTED")
        Base.metadata.create_all(self.engine)

    def insert_document(self, **kwargs) -> int:
        new_document = Document(**kwargs)
        self.session.add(new_document)
        self.session.commit()
        return new_document.id

    def insert_document_object(self, document: Document) -> int:
        self.session.add(document)
        self.session.commit()
        return document.id

    def insert_documents_batch(self, documents: List[Document]) -> List[Document]:
        self.session.add_all(documents)
        self.session.commit()
        for document in documents:
            self.session.refresh(document)  # Refresh each document to get the IDs from the database
        return documents

    def select_document(self, document_id: int) -> Document | None:
        if document_id <= 0:
            raise ValueError("ID value must be positive")
        return self.session.query(Document).filter_by(id=document_id).first()

    def update_document(self, document_id: int, **kwargs) -> int:
        if document_id <= 0:
            raise ValueError("ID value must be positive")
        updated_docs = self.session.query(Document).filter_by(id=document_id).update(kwargs)
        self.session.commit()
        return updated_docs

    # def update_document_object(self, document: Document) -> int:
    #     if document.id <= 0:
    #         raise ValueError("ID value must be positive")
    #
    #     existing_document = self.session.query(Document).filter_by(id=document.id).first()
    #     if not existing_document:
    #         raise ValueError("Document not found")
    #
    #     for field in document.__dict__:
    #         if field != '_sa_instance_state':
    #             setattr(existing_document, field, getattr(document, field))
    #
    #     self.session.commit()
    #     return existing_document.id

    def update_document_object(self, document: Document):
        if document.id <= 0:
            raise ValueError("ID value must be positive")

        existing_document = self.session.query(Document).filter_by(id=document.id).first()
        if not existing_document:
            raise ValueError("Document not found")

        # really afraid of these things!!!
        existing_document.chunking_session = document.chunking_session
        existing_document.analyzed = document.analyzed
        existing_document.last_update = document.last_update

        self.session.commit()


    def delete_by_document_id(self, document_id: int) -> int:
        if document_id <= 0:
            raise ValueError("ID value must be positive")
        deleted_docs = self.session.query(Document).filter_by(id=document_id).delete()
        self.session.commit()
        return deleted_docs

    def delete_by_parent_id(self, document_id: int) -> int:
        if document_id <= 0:
            raise ValueError("ID value must be positive")
        deleted_docs = self.session.query(Document).filter_by(parent_id=document_id).delete()
        self.session.commit()
        return deleted_docs
