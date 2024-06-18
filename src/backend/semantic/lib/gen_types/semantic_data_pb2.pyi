from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class FileType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    UNKNOWN: _ClassVar[FileType]
    URL: _ClassVar[FileType]
    PDF: _ClassVar[FileType]
    RTF: _ClassVar[FileType]
    DOC: _ClassVar[FileType]
    XLS: _ClassVar[FileType]
    PPT: _ClassVar[FileType]
    TXT: _ClassVar[FileType]
    MD: _ClassVar[FileType]
    YT: _ClassVar[FileType]
UNKNOWN: FileType
URL: FileType
PDF: FileType
RTF: FileType
DOC: FileType
XLS: FileType
PPT: FileType
TXT: FileType
MD: FileType
YT: FileType

class SemanticData(_message.Message):
    __slots__ = ("url", "url_recursive", "site_map", "search_for_sitemap", "document_id", "connector_id", "file_type", "collection_name", "model_name", "model_dimension")
    URL_FIELD_NUMBER: _ClassVar[int]
    URL_RECURSIVE_FIELD_NUMBER: _ClassVar[int]
    SITE_MAP_FIELD_NUMBER: _ClassVar[int]
    SEARCH_FOR_SITEMAP_FIELD_NUMBER: _ClassVar[int]
    DOCUMENT_ID_FIELD_NUMBER: _ClassVar[int]
    CONNECTOR_ID_FIELD_NUMBER: _ClassVar[int]
    FILE_TYPE_FIELD_NUMBER: _ClassVar[int]
    COLLECTION_NAME_FIELD_NUMBER: _ClassVar[int]
    MODEL_NAME_FIELD_NUMBER: _ClassVar[int]
    MODEL_DIMENSION_FIELD_NUMBER: _ClassVar[int]
    url: str
    url_recursive: bool
    site_map: str
    search_for_sitemap: bool
    document_id: int
    connector_id: int
    file_type: FileType
    collection_name: str
    model_name: str
    model_dimension: int
    def __init__(self, url: _Optional[str] = ..., url_recursive: bool = ..., site_map: _Optional[str] = ..., search_for_sitemap: bool = ..., document_id: _Optional[int] = ..., connector_id: _Optional[int] = ..., file_type: _Optional[_Union[FileType, str]] = ..., collection_name: _Optional[str] = ..., model_name: _Optional[str] = ..., model_dimension: _Optional[int] = ...) -> None: ...
