from typing import List


class ChunkedItem:
    def __init__(self, url: str, content: str, document_id: int = None, parent_id: int = None):
        self.url = url
        self.content = content
        self.document_id = document_id
        self.parent_id = parent_id

    @classmethod
    def create_chunked_items(cls, results: List[str], url: str, document_id: int, parent_id: int) -> List['ChunkedItem']:
        return [cls(url, result, document_id, parent_id) for result in results]

    def __repr__(self):
        return f"ChunkedItem(url={self.url}, content={self.content}, document_id={self.document_id}, parent_id={self.parent_id})"
