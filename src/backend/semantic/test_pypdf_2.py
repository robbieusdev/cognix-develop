import os
import os
import string
import pathlib
import sys
import time
import pymupdf


def extract_text_from_pdf(pdf_path):
    # open document
    doc = pymupdf.open(pdf_path)

    text = ""  # we will return this string
    row_count = 0  # counts table rows
    header = ""  # overall table header: output this only once!

    # iterate over the pages
    for page in doc:
        # only read the table rows on each page, ignore other content
        tables = page.find_tables()  # a "TableFinder" object
        for table in tables:

            # on first page extract external column names if present
            if page.number == 0 and table.header.external:
                # build the overall table header string
                # technical note: incomplete / complex tables may have
                # "None" in some header cells. Just use empty string then.
                header = (
                        ";".join(
                            [
                                name if name is not None else ""
                                for name in table.header.names
                            ]
                        )
                        + "\n"
                )
                text += header
                row_count += 1  # increase row counter

            # output the table body
            for row in table.extract():  # iterate over the table rows

                # again replace any "None" in cells by an empty string
                row_text = (
                        ";".join([cell if cell is not None else "" for cell in row]) + "\n"
                )
                if row_text != header:  # omit duplicates of header row
                    text += row_text
                    row_count += 1  # increase row counter
    doc.close()  # close document
    print(f"Loaded {row_count} table rows from file '{doc.name}'.\n")
    return text


if __name__ == "__main__":
    import pathlib
    import sys
    import time

    try:
        filename = sys.argv[1]
    except IndexError:
        print(f"Usage:\npython {os.path.basename(__file__)} input.pdf")
        sys.exit()

    print(extract_text_from_pdf(filename))
