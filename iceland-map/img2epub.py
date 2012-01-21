import string
import getopt
import sys
import zipfile
import uuid
import os
import os.path
import tempfile
from PIL import Image

CONTAINER_XML = """<?xml version="1.0"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
  <rootfiles>
    <rootfile full-path="OEBPS/content.opf" media-type="application/oebps-package+xml"/>
  </rootfiles>
</container>"""

STYLESHEET_CSS = """
body {
  margin: 0 0 0 0;
}
"""

TOC_NCX = string.Template("""<?xml version="1.0" encoding="UTF-8"?>

<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1">
  <head>
    <meta name="dtb:uid" content="$uid"/>
    <meta name="dtb:depth" content="1"/>
    <meta name="dtb:totalPageCount" content="0"/>
    <meta name="dtb:maxPageNumber" content="0"/>
  </head>
  <docTitle>
    <text>$title</text>
  </docTitle>
  <navMap>
    $nav_map
  </navMap>
</ncx>""")

CONTENT_OPF = string.Template("""<?xml version="1.0"?>

<package version="2.0" xmlns="http://www.idpf.org/2007/opf" unique-identifier="BookId">
 <metadata xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:opf="http://www.idpf.org/2007/opf">
   <dc:title>$title</dc:title> 
   <dc:creator opf:role="aut">$creator</dc:creator>
   <dc:language>en-US</dc:language> 
   <dc:rights>Public Domain</dc:rights> 
   <dc:publisher>$publisher</dc:publisher> 
   <dc:identifier id="BookId">urn:uuid:$uid</dc:identifier>
 </metadata>

 <manifest>
  <item id="ncx" href="toc.ncx" media-type="application/x-dtbncx+xml" />
  <item id="style" href="stylesheet.css" media-type="text/css" />
  $manifest
 </manifest>

 <spine toc="ncx">
  $spine
 </spine>
</package>""")

OVERVIEW_XHTML = string.Template("""<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
  <title>Overview</title>
  <link href="stylesheet.css" type="text/css" rel="stylesheet" />
</head>
<body>
  <svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" 
       width="100%" height="100%" viewBox="0 0 $page_width $page_height" preserveAspectRatio="xMidYMid meet">
    <image width="$page_width" height="$page_height" xlink:href="images/overview.$img_format" />

    $links
  </svg>
</body>
</html>
""")

PAGE_XHTML = string.Template("""<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
  <title>$page</title>
  <link href="stylesheet.css" type="text/css" rel="stylesheet" />
</head>
<body>
  <svg version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" 
       width="100%" height="100%" viewBox="0 0 $page_width $page_height" preserveAspectRatio="xMidYMid meet">
    <image width="$page_width" height="$page_height" xlink:href="images/$page.$img_format" />

    $overview_link
    $top_link
    $right_link
    $bottom_link
    $left_link
  </svg>
</body>
</html>
""")

def create_overview(f, page_width, page_height):
  img = Image.open(f)
  (w, h) = img.size
  if w % page_width != 0:
    w = page_width * (w / page_width + 1)
  if h % page_height != 0:
    h = page_height * (h / page_height + 1)

  overview_img = Image.new("P", (w, h), 255)
  overview_img.paste(img, (0, 0))

  return overview_img

def create_nav_map(ids):
  t = """<navPoint id="%s" playOrder="%d">
      <navLabel>
        <text>%s</text>
      </navLabel>
      <content src="%s.xhtml"/>
    </navPoint>"""

  return "\n".join([t % (i, k+1, i, i) for (k, i) in enumerate(ids)])

def create_manifest(ids, fmt):
  return "\n".join(['<item id="%s" href="%s.xhtml" media-type="application/xhtml+xml" />' % (i, i) for i in ids]) + "\n" + \
         "\n".join(['<item id="%s-img" href="images/%s.%s" media-type="image/%s" />' % (i, i, fmt, fmt) for i in ids])

def create_spine(ids):
  return "\n".join(['<itemref idref="%s" />' % i for i in ids])

def create_page_name(x, y):
  return chr(ord('A') + x) + str(y)

def create_overview_links(size_x, size_y, page_width, page_height):
  t = string.Template("""
    <rect x="$left" y="$top" width="$width" height="$height" fill="none" stroke="black" stroke-dasharray="10,10" stroke-width="1"/>
    <a xlink:href="$page.xhtml">
      <text x="$text_center" y="$text_bottom" text-anchor="middle" 
         font-family="Verdana" font-size="$font_size" font-weight="bold" fill="#888888" fill-opacity="0.1">$page</text>
    </a>""")

  w = page_width / size_x
  h = page_height / size_y
  font_size = min(100, h - 20)
  env = {
    "width": w,
    "height": h,
    "font_size": font_size
  }

  return "\n".join([t.substitute(env, page=create_page_name(x, y), 
                                      left=x*w, top=y*h,
                                      text_center=x*w+w/2, text_bottom=y*h+(h+font_size)/2) \
      for x in range(0, size_x) for y in range(0, size_y)])

def create_page_link(x, y, sizes, direction):
  (size_x, size_y, page_width, page_height) = sizes
  if x < 0 or x >= size_x or y < 0 or y >= size_y:
    return ""

  page = create_page_name(x, y)
  font_size = 16
  if direction == 't':
    text_center = page_width / 2
    text_bottom = 10 + font_size
    text = '^'
  if direction == 'r':
    text_center = page_width - font_size * len(page) / 2
    text_bottom = page_height / 2 + font_size / 2
    text = '>'
  if direction == 'b':
    text_center = page_width / 2
    text_bottom = page_height - 10
    text = 'v'
  if direction == 'l':
    text_center = font_size * len(page) / 2
    text_bottom = page_height / 2 + font_size / 2
    text = '<'
  if direction == 'o':
    page = "overview"
    text = 'O'
    text_center = page_width - font_size * len(page) / 2
    text_bottom = 10 + font_size

  return """   
    <a xlink:href="%s.xhtml">
      <text x="%d" y="%d" text-anchor="middle" 
         font-family="Verdana" font-size="%d" font-weight="bold" fill="#111111" fill-opacity="0.5">%s</text>
    </a>""" % (page, text_center, text_bottom, font_size, text)

               
def convert_img_to_epub(img_file, epub_file, width, height, title):
  overview_img = create_overview(img_file, width, height)

  img_format = "png"
  pages = ["overview"]
  images = [overview_img.resize((width, height))]
  page_numbers = [(None, None)]

  size_x = overview_img.size[0] / width
  size_y = overview_img.size[1] / height

  for x in range(0, size_x):
    for y in range(0, size_y):
      pages.append(create_page_name(x, y))
      page_numbers.append((x, y))
      images.append(overview_img.crop((x * width, y * height, (x+1) * width, (y+1) * height)))

  env = {
    "uid": uuid.uuid1().hex,
    "title": title,
    "creator": "me",
    "publisher": "me",
    "img_format": img_format,
    "page_width": width,
    "page_height": height,
  }

  with zipfile.ZipFile(epub_file, 'w', zipfile.ZIP_DEFLATED) as epub:
    epub.writestr("mimetype", "application/epub+zip")
    epub.writestr("META-INF/container.xml", CONTAINER_XML)
    epub.writestr("OEBPS/stylesheet.css", STYLESHEET_CSS)
    epub.writestr("OEBPS/toc.ncx", TOC_NCX.substitute(env,
      nav_map=create_nav_map(pages)))
    epub.writestr("OEBPS/content.opf", CONTENT_OPF.substitute(env, 
      manifest=create_manifest(pages, img_format),
      spine=create_spine(pages)))

    for (page, img, (x, y)) in zip(pages, images, page_numbers):
      with tempfile.TemporaryFile() as f:
        img.save(f, img_format)
        f.seek(0)
        epub.writestr("OEBPS/images/%s.%s" % (page, img_format), f.read())

      if page == "overview":
        epub.writestr("OEBPS/%s.xhtml" % page, OVERVIEW_XHTML.substitute(env, page=page,
          links=create_overview_links(size_x, size_y, width, height)))
      else:
        sizes = (size_x, size_y, width, height)
        epub.writestr("OEBPS/%s.xhtml" % page, PAGE_XHTML.substitute(env, page=page,
          overview_link=create_page_link(x, y, sizes, 'o'),
          top_link=create_page_link(x, y-1, sizes, 't'), right_link=create_page_link(x+1, y, sizes, 'r'),
          bottom_link=create_page_link(x, y+1, sizes, 'b'), left_link=create_page_link(x-1, y, sizes, 'l')))

if __name__ == "__main__":
  (opts, args) = getopt.gnu_getopt(sys.argv[1:], "w:h:t:")

  width = 592
  height = 780
  title = os.path.basename(args[0])
  for (o, a) in opts:
    if o == "-w":
      width = int(a)
    if o == "-h":
      height = int(a)
    if o == "-t":
      title = a
  
  convert_img_to_epub(args[0], args[1], width, height, title)
