import os
import pymupdf
from pymupdf4llm.helpers.get_text_lines import get_raw_lines, is_white
from pymupdf4llm.helpers.multi_column import column_boxes
from typing import List, Tuple, Dict, Union, Any

if pymupdf.pymupdf_version_tuple < (1, 24, 2):
    raise NotImplementedError("PyMuPDF version 1.24.2 or later is needed.")

bullet = ("- ", "* ", chr(0xF0A7), chr(0xF0B7), chr(0xB7), chr(8226), chr(9679))
GRAPHICS_TEXT = "\n![%s](%s)\n"


class IdentifyHeaders:
    def __init__(self, doc: Union[str, pymupdf.Document], pages: List[int] = None, body_limit: float = 12) -> None:
        if isinstance(doc, pymupdf.Document):
            mydoc = doc
        else:
            mydoc = pymupdf.open(doc)

        if pages is None:
            pages = list(range(mydoc.page_count))

        fontsizes: Dict[int, int] = {}
        for pno in pages:
            page = mydoc.load_page(pno)
            blocks = page.get_text("dict", flags=pymupdf.TEXTFLAGS_TEXT)["blocks"]
            for span in [s for b in blocks for l in b["lines"] for s in l["spans"] if not is_white(s["text"])]:
                fontsz = round(span["size"])
                count = fontsizes.get(fontsz, 0) + len(span["text"].strip())
                fontsizes[fontsz] = count

        if mydoc != doc:
            mydoc.close()

        self.header_id: Dict[int, str] = {}
        temp = sorted([(k, v) for k, v in fontsizes.items()], key=lambda i: i[1], reverse=True)
        if temp:
            b_limit = max(body_limit, temp[0][0])
        else:
            b_limit = body_limit

        sizes = sorted([f for f in fontsizes.keys() if f > b_limit], reverse=True)[:6]

        for i, size in enumerate(sizes):
            self.header_id[size] = "#" * (i + 1) + " "

    def get_header_id(self, span: dict, page: pymupdf.Page = None) -> str:
        fontsize = round(span["size"])
        return self.header_id.get(fontsize, "")


class MarkdownConverter:
    def __init__(self, doc_path: str, pages: List[int] = None, body_limit: float = 12, write_images: bool = False,
                 page_chunks: bool = False,
                 margins: Union[float, Tuple[float, float], Tuple[float, float, float, float]] = (
                 0, 50, 0, 50)) -> None:
        self.doc_path = doc_path
        self.doc = pymupdf.open(doc_path)
        self.pages = pages if pages is not None else list(range(self.doc.page_count))
        self.body_limit = body_limit
        self.write_images = write_images
        self.page_chunks = page_chunks
        self.margins = self._validate_margins(margins)
        self.hdr_info = IdentifyHeaders(self.doc, pages=self.pages, body_limit=body_limit)

    def _validate_margins(self, margins: Union[float, Tuple[float, float], Tuple[float, float, float, float]]) -> Tuple[
        float, float, float, float]:
        if hasattr(margins, "__float__"):
            return (margins, margins, margins, margins)
        if len(margins) == 2:
            return (0, margins[0], 0, margins[1])
        if len(margins) != 4:
            raise ValueError("Margins must have length 2 or 4 or be a number.")
        if not all([hasattr(m, "__float__") for m in margins]):
            raise ValueError("Margin values must be numbers")
        return tuple(margins)

    def convert_to_markdown(self) -> Union[str, List[Dict[str, Any]]]:
        document_output: Union[str, List[Dict[str, Any]]] = "" if not self.page_chunks else []
        toc = self.doc.get_toc()
        textflags = pymupdf.TEXT_DEHYPHENATE | pymupdf.TEXT_MEDIABOX_CLIP

        for pno in self.pages:
            page_output, images, tables, graphics = self._get_page_output(pno, textflags)
            if not self.page_chunks:
                document_output += page_output
            else:
                page_tocs = [t for t in toc if t[-1] == pno + 1]
                metadata = self._get_metadata(pno)
                document_output.append({
                    "metadata": metadata,
                    "toc_items": page_tocs,
                    "tables": tables,
                    "images": images,
                    "graphics": graphics,
                    "text": page_output,
                })

        return document_output

    def _get_metadata(self, pno: int) -> Dict[str, Union[str, int]]:
        meta = self.doc.metadata.copy()
        meta["file_path"] = self.doc.name
        meta["page_count"] = self.doc.page_count
        meta["page"] = pno + 1
        return meta

    def _get_page_output(self, pno: int, textflags: int) -> Tuple[
        str, List[Dict[str, Any]], List[Dict[str, Any]], List[Dict[str, Any]]]:
        page = self.doc[pno]
        md_string = ""
        left, top, right, bottom = self.margins
        clip = page.rect + (left, top, -right, -bottom)
        links = [l for l in page.get_links() if l["kind"] == 2]
        textpage = page.get_textpage(flags=textflags, clip=clip)
        img_info = [img for img in page.get_image_info() if img["bbox"] in clip]
        images = img_info[:]
        tables = []
        graphics = []
        tabs = page.find_tables(clip=clip, strategy="lines_strict")
        tab_rects = self._get_tab_rects(tabs)
        tab_rects0 = list(tab_rects.values())
        paths, vg_clusters0 = self._get_paths_and_graphics(page, tab_rects0, clip, img_info)
        vg_clusters = {i: r for i, r in enumerate(vg_clusters0)}
        text_rects = column_boxes(page, paths=paths, no_image_text=self.write_images, textpage=textpage,
                                  avoid=tab_rects0 + vg_clusters0)

        for text_rect in text_rects:
            md_string += self._output_tables(tabs, text_rect, tab_rects)
            md_string += self._output_images(page, text_rect, vg_clusters)
            md_string += self._write_text(page, textpage, text_rect, tabs, tab_rects, vg_clusters, links)

        md_string += self._output_tables(tabs, None, tab_rects)
        md_string += self._output_images(page, None, vg_clusters)
        md_string += "\n-----\n\n"
        while md_string.startswith("\n"):
            md_string = md_string[1:]

        return md_string, images, tables, graphics

    def _get_tab_rects(self, tabs: List[Any]) -> Dict[int, pymupdf.Rect]:
        tab_rects: Dict[int, pymupdf.Rect] = {}
        for i, t in enumerate(tabs):
            tab_rects[i] = pymupdf.Rect(t.bbox) | pymupdf.Rect(t.header.bbox)
        return tab_rects

    def _get_paths_and_graphics(self, page: pymupdf.Page, tab_rects0: List[pymupdf.Rect], clip: pymupdf.Rect,
                                img_info: List[Dict[str, Any]]) -> Tuple[List[Dict[str, Any]], List[pymupdf.Rect]]:
        page_clip = page.rect + (36, 36, -36, -36)
        paths = [p for p in page.get_drawings() if
                 not self._intersects_rects(p["rect"], tab_rects0) and p["rect"] in page_clip and p[
                     "rect"].width < page_clip.width and p["rect"].height < page_clip.height]
        vg_clusters = [r for r in page.cluster_drawings(drawings=paths) if
                       not self._intersects_rects(r, tab_rects0) and r.height > 20]

        if self.write_images:
            vg_clusters += [pymupdf.Rect(i["bbox"]) for i in img_info]

        return paths, vg_clusters

    def _intersects_rects(self, rect: pymupdf.Rect, rect_list: List[pymupdf.Rect]) -> bool:
        for r in rect_list:
            if (rect.tl + rect.br) / 2 in r:
                return True
        return False

    def _output_tables(self, tabs: List[Any], text_rect: pymupdf.Rect, tab_rects: Dict[int, pymupdf.Rect]) -> str:
        this_md = ""
        if text_rect is not None:
            for i, trect in sorted([(j[0], j[1]) for j in tab_rects.items() if j[1].y1 <= text_rect.y0],
                                   key=lambda j: (j[1].y1, j[1].x0)):
                this_md += tabs[i].to_markdown(clean=False)
                del tab_rects[i]
        else:
            for i, trect in sorted(tab_rects.items(), key=lambda j: (j[1].y1, j[1].x0)):
                this_md += tabs[i].to_markdown(clean=False)
                del tab_rects[i]
        return this_md

    def _output_images(self, page: pymupdf.Page, text_rect: pymupdf.Rect, img_rects: Dict[int, pymupdf.Rect]) -> str:
        if img_rects is None:
            return ""
        this_md = ""
        if text_rect is not None:
            for i, img_rect in sorted([(j[0], j[1]) for j in img_rects.items() if j[1].y1 <= text_rect.y0],
                                      key=lambda j: (j[1].y1, j[1].x0)):
                pathname = self._save_image(page, img_rect, i)
                if pathname:
                    this_md += GRAPHICS_TEXT % (pathname, pathname)
                del img_rects[i]
        else:
            for i, img_rect in sorted(img_rects.items(), key=lambda j: (j[1].y1, j[1].x0)):
                pathname = self._save_image(page, img_rect, i)
                if pathname:
                    this_md += GRAPHICS_TEXT % (pathname, pathname)
                del img_rects[i]
        return this_md

    def _save_image(self, page: pymupdf.Page, rect: pymupdf.Rect, i: int) -> str:
        filename = page.parent.name.replace("\\", "/")
        image_path = f"{filename}-{page.number}-{i}.png"
        if self.write_images:
            pix = page.get_pixmap(clip=rect)
            pix.save(image_path)
            del pix
            return os.path.basename(image_path)
        return ""

    def _write_text(self, page: pymupdf.Page, textpage: pymupdf.TextPage, clip: pymupdf.Rect, tabs: List[Any],
                    tab_rects: Dict[int, pymupdf.Rect], img_rects: Dict[int, pymupdf.Rect],
                    links: List[Dict[str, Any]]) -> str:
        out_string = ""
        nlines = get_raw_lines(textpage, clip=clip, tolerance=3)
        tab_rects0 = list(tab_rects.values())
        img_rects0 = list(img_rects.values())
        prev_lrect = None
        prev_bno = -1
        code = False
        prev_hdr_string = None

        for lrect, spans in nlines:
            if self._intersects_rects(lrect, tab_rects0) or self._intersects_rects(lrect, img_rects0):
                continue

            for i, tab_rect in sorted(
                    [j for j in tab_rects.items() if j[1].y1 <= lrect.y0 and not (j[1] & clip).is_empty],
                    key=lambda j: (j[1].y1, j[1].x0)):
                out_string += "\n" + tabs[i].to_markdown(clean=False) + "\n"
                del tab_rects[i]

            for i, img_rect in sorted(
                    [j for j in img_rects.items() if j[1].y1 <= lrect.y0 and not (j[1] & clip).is_empty],
                    key=lambda j: (j[1].y1, j[1].x0)):
                pathname = self._save_image(page, img_rect, i)
                if pathname:
                    out_string += GRAPHICS_TEXT % (pathname, pathname)
                del img_rects[i]

            text = " ".join([s["text"] for s in spans])
            all_mono = all([s["flags"] & 8 for s in spans])

            if all_mono:
                if not code:
                    out_string += "```\n"
                    code = True
                delta = int((lrect.x0 - clip.x0) / (spans[0]["size"] * 0.5))
                indent = " " * delta
                out_string += indent + text + "\n"
                continue

            span0 = spans[0]
            bno = span0["block"]
            if bno != prev_bno:
                out_string += "\n"
                prev_bno = bno

            if (prev_lrect and lrect.y1 - prev_lrect.y1 > lrect.height * 1.5 or span0["text"].startswith("[") or span0[
                "text"].startswith(bullet) or span0["flags"] & 1):
                out_string += "\n"
            prev_lrect = lrect

            hdr_string = self.hdr_info.get_header_id(span0, page=page)
            if hdr_string and hdr_string == prev_hdr_string:
                out_string = out_string[:-1] + " " + text + "\n"
                continue

            prev_hdr_string = hdr_string
            if hdr_string.startswith("#"):
                out_string += hdr_string + text + "\n"
                continue

            if code:
                out_string += "```\n"
                code = False

            for i, s in enumerate(spans):
                mono = s["flags"] & 8
                bold = s["flags"] & 16
                italic = s["flags"] & 2

                if mono:
                    out_string += f"`{s['text'].strip()}` "
                else:
                    prefix = ""
                    suffix = ""
                    if hdr_string == "":
                        if bold:
                            prefix = "**"
                            suffix += "**"
                        if italic:
                            prefix += "_"
                            suffix = "_" + suffix

                    ltext = self._resolve_links(links, s)
                    if ltext:
                        text = f"{hdr_string}{prefix}{ltext}{suffix} "
                    else:
                        text = f"{hdr_string}{prefix}{s['text'].strip()}{suffix} "

                    if text.startswith(bullet):
                        text = "-  " + text[1:]
                    out_string += text

            if not code:
                out_string += "\n"
        out_string += "\n"
        if code:
            out_string += "```\n"
            code = False

        return out_string.replace(" \n", "\n").replace("  ", " ").replace("\n\n\n", "\n\n")

    def _resolve_links(self, links: List[Dict[str, Any]], span: Dict[str, Any]) -> Union[str, None]:
        bbox = pymupdf.Rect(span["bbox"])
        bbox_area = 0.7 * abs(bbox)
        for link in links:
            hot = link["from"]
            if abs(hot & bbox) >= bbox_area:
                return f'[{span["text"].strip()}]({link["uri"]})'
        return None


if __name__ == "__main__":
    import pathlib
    import sys
    import time

    try:
        filename = sys.argv[1]
    except IndexError:
        print(f"Usage:\npython {os.path.basename(__file__)} input.pdf")
        sys.exit()

    t0 = time.perf_counter()
    converter = MarkdownConverter(filename)
    md_string = converter.convert_to_markdown()

    if isinstance(md_string, str):
        print(md_string)
        outname = converter.doc_path.replace(".pdf", ".md")
        pathlib.Path(outname).write_bytes(md_string.encode())
    else:
        for page in md_string:
            print(page["text"])
            outname = converter.doc_path.replace(".pdf", f"_page_{page['metadata']['page']}.md")
            pathlib.Path(outname).write_bytes(page["text"].encode())

    t1 = time.perf_counter()
    print(f"Markdown creation time for {converter.doc.name=} {round(t1 - t0, 2)} sec.")
