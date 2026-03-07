"""Convert image to colored braille art for terminal display.
Outputs braille characters with ANSI 256-color escapes for true-color rendering.
"""
from PIL import Image, ImageEnhance
import sys

BRAILLE_BASE = 0x2800
DOT_MAP = [
    (0, 0, 0x01), (1, 0, 0x08),
    (0, 1, 0x02), (1, 1, 0x10),
    (0, 2, 0x04), (1, 2, 0x20),
    (0, 3, 0x40), (1, 3, 0x80),
]


def image_to_color_braille(image_path, width=55, threshold=100):
    """Generate braille art with ANSI true-color (24-bit) per character."""
    img_orig = Image.open(image_path).convert("RGB")
    img_gray = img_orig.convert("L")
    img_gray = ImageEnhance.Contrast(img_gray).enhance(1.6)

    char_w = width
    char_h = int(char_w * (img_orig.height / img_orig.width) * (2 / 4))

    pixel_w = char_w * 2
    pixel_h = char_h * 4
    img_gray = img_gray.resize((pixel_w, pixel_h), Image.LANCZOS)
    img_color = img_orig.resize((char_w, char_h), Image.LANCZOS)

    lines = []
    for cy in range(char_h):
        line = ""
        for cx in range(char_w):
            # Compute braille pattern from grayscale
            code = 0
            for dx, dy, bit in DOT_MAP:
                px = cx * 2 + dx
                py = cy * 4 + dy
                if px < pixel_w and py < pixel_h:
                    val = img_gray.getpixel((px, py))
                    if val < threshold:
                        code |= bit

            # Get color from the color image
            r, g, b = img_color.getpixel((cx, cy))
            char = chr(BRAILLE_BASE + code)

            if code == 0:
                line += " "
            else:
                # ANSI true-color foreground
                line += f"\033[38;2;{r};{g};{b}m{char}\033[0m"
        lines.append(line)
    return "\n".join(lines)


def image_to_braille_plain(image_path, width=55, threshold=100):
    """Generate plain braille art without color."""
    img = Image.open(image_path)
    img = ImageEnhance.Contrast(img).enhance(1.5)
    img = img.convert("L")

    char_w = width
    char_h = int(char_w * (img.height / img.width) * (2 / 4))

    pixel_w = char_w * 2
    pixel_h = char_h * 4
    img = img.resize((pixel_w, pixel_h), Image.LANCZOS)

    lines = []
    for cy in range(char_h):
        line = ""
        for cx in range(char_w):
            code = 0
            for dx, dy, bit in DOT_MAP:
                px = cx * 2 + dx
                py = cy * 4 + dy
                if px < pixel_w and py < pixel_h:
                    val = img.getpixel((px, py))
                    if val < threshold:
                        code |= bit
            line += chr(BRAILLE_BASE + code)
        lines.append(line)
    return "\n".join(lines)


if __name__ == "__main__":
    path = sys.argv[1] if len(sys.argv) > 1 else "NoLoveSong.jpg"
    width = int(sys.argv[2]) if len(sys.argv) > 2 else 55
    threshold = int(sys.argv[3]) if len(sys.argv) > 3 else 100
    mode = sys.argv[4] if len(sys.argv) > 4 else "color"

    if mode == "color":
        print(image_to_color_braille(path, width, threshold))
    else:
        print(image_to_braille_plain(path, width, threshold))
