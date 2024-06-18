from lib.semantic.semantic_base import BaseSemantic
from lib.semantic.semantic_url import URLSemantic
from lib.semantic.semantic_pdf import PDFSemantic
from lib.semantic.semantic_doc import DOCXSemantic
from lib.semantic.semantic_txt import TXTSemantic
from lib.semantic.semantic_md import MDSemantic
from lib.gen_types.semantic_data_pb2 import FileType
from lib.semantic.semantic_youtube import YTSemantic


# Plaintext	.eml, .html, .md, .msg, .rst, .rtf, .txt, .xml
# Documents	.csv, .doc, .docx, .epub, .odt, .pdf, .ppt, .pptx, .tsv, .xlsx

class SemanticFactory:
    factories = {
        FileType.URL: URLSemantic,
        FileType.PDF: PDFSemantic,
        FileType.DOC: DOCXSemantic,
        FileType.TXT: TXTSemantic,
        FileType.MD: MDSemantic,
        FileType.YT: YTSemantic,
        # Add other mappings here
    }

    @staticmethod
    def create_semantic_analyzer(file_type: FileType) -> BaseSemantic:
        method_name = SemanticFactory.factories.get(file_type)
        if not method_name:
            raise ValueError(f"No factory method found for file type: {file_type}")
        semantic_class = SemanticFactory.factories.get(file_type)
        if not semantic_class:
            raise ValueError(f"Unsupported file type")

        return semantic_class()

