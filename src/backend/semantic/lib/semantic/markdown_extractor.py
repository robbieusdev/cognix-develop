import os
import sys
import time
import pymupdf4llm
import mistune
from typing import List, Dict, Any


class MarkdownSectionExtractor:
    def __init__(self):
        pass

    def _parse_markdown_content(self, markdown_content: str) -> List[Dict[str, Any]]:
        markdown = mistune.create_markdown(renderer='ast')
        return markdown(markdown_content)

    def extract_sections(self, markdown_content: str) -> List[Dict[str, Any]]:
        parsed_content = self._parse_markdown_content(markdown_content)
        sections = []
        current_section = []
        heading_hierarchy = []

        def add_section():
            if current_section:
                sections.append({
                    'headings': heading_hierarchy[:],  # Copy of current heading hierarchy
                    'content': current_section[:]
                })
                current_section.clear()

        for element in parsed_content:
            if element['type'] == 'heading':
                level = element['attrs']['level']
                text = element['children'][0]['raw'] if element['children'] else ''
                add_section()  # Finish the current section before starting a new heading

                # Update heading hierarchy based on the level
                if level == 1:
                    heading_hierarchy = [text]
                else:
                    if len(heading_hierarchy) >= level:
                        heading_hierarchy = heading_hierarchy[:level - 1]
                    heading_hierarchy.append(text)

                current_section.append(element)
            else:
                current_section.append(element)

        add_section()  # Add the last section

        return sections

    def extract_chunks(self, markdown: str) -> list[str]:
        sections = self.extract_sections(markdown)
        results = []
        # Process the extracted sections
        for idx, section in enumerate(sections):
            content_texts = []
            for element in section['content']:
                if 'children' in element and element['children']:
                    for child in element['children']:
                        if 'raw' in child:
                            content_texts.append(child['raw'])
                        elif 'text' in child:
                            content_texts.append(child['text'])
                elif 'raw' in element:
                    content_texts.append(element['raw'])
                elif 'text' in element:
                    content_texts.append(element['text'])

            readable_text = ' '.join(content_texts)
            headings_text = ' > '.join(section['headings'])
            result = f"{headings_text}\n{readable_text}\n"
            results.append(result)
            return results

import os
import sys
import time
import pymupdf4llm
import mistune
from typing import List, Dict, Any


class MarkdownSectionExtractor:
    def __init__(self):
        pass

    def _parse_markdown_content(self, markdown_content: str) -> List[Dict[str, Any]]:
        markdown = mistune.create_markdown(renderer='ast')
        return markdown(markdown_content)

    def extract_sections(self, markdown_content: str) -> List[Dict[str, Any]]:
        parsed_content = self._parse_markdown_content(markdown_content)
        sections = []
        current_section = []
        heading_hierarchy = []

        def add_section():
            if current_section:
                sections.append({
                    'headings': heading_hierarchy[:],  # Copy of current heading hierarchy
                    'content': current_section[:]
                })
                current_section.clear()

        for element in parsed_content:
            if element['type'] == 'heading':
                level = element['attrs']['level']
                text = element['children'][0]['raw'] if element['children'] else ''
                add_section()  # Finish the current section before starting a new heading

                # Update heading hierarchy based on the level
                if level == 1:
                    heading_hierarchy = [text]
                else:
                    if len(heading_hierarchy) >= level:
                        heading_hierarchy = heading_hierarchy[:level - 1]
                    heading_hierarchy.append(text)

                current_section.append(element)
            else:
                current_section.append(element)

        add_section()  # Add the last section

        return sections

    def extract_chunks(self, markdown: str) -> list[str]:
        sections = self.extract_sections(markdown)
        results = []
        # Process the extracted sections
        for idx, section in enumerate(sections):
            content_texts = []
            for element in section['content']:
                if 'children' in element and element['children']:
                    for child in element['children']:
                        if 'raw' in child:
                            content_texts.append(child['raw'])
                        elif 'text' in child:
                            content_texts.append(child['text'])
                elif 'raw' in element:
                    content_texts.append(element['raw'])
                elif 'text' in element:
                    content_texts.append(element['text'])

            readable_text = ' '.join(content_texts)
            headings_text = ' > '.join(section['headings'])
            result = f"{headings_text}\n{readable_text}\n"
            results.append(result)
        return results

# if __name__ == "__main__":
#     import pathlib
#
#     try:
#         filename = sys.argv[1]
#     except IndexError:
#         print(f"Usage:\npython {os.path.basename(__file__)} input.pdf")
#         sys.exit()
#
#     t0 = time.perf_counter()
#
#     markdown_content = pymupdf4llm.to_markdown(filename)
#     print(markdown_content)
#
#     extractor = MarkdownSectionExtractor()
#     # sections = extractor.extract_sections(markdown_content)
#     results = extractor.extract_chunks(markdown_content)
#
#
#     # Print the results
#     for result in results:
#         print(result)
#         print("****************** new section *********************\n")
#
#     t1 = time.perf_counter()
#     print(f"Markdown creation time for {filename=} {round(t1 - t0, 2)} sec.")
