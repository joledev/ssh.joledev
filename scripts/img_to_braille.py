"""Convert image to braille art for high-resolution terminal display."""
from PIL import Image, ImageEnhance
import sys

# Braille unicode block: 2x4 dot matrix per character
# Dot positions:
# 0 3
# 1 4
# 2 5
# 6 7
BRAILLE_BASE = 0x2800
DOT_MAP = [
    (0, 0, 0x01), (1, 0, 0x08),
    (0, 1, 0x02), (1, 1, 0x10),
    (0, 2, 0x04), (1, 2, 0x20),
    (0, 3, 0x40), (1, 3, 0x80),
]


def image_to_braille(image_path, width=50, threshold=128, invert=False):
    img = Image.open(image_path)
    img = ImageEnhance.Contrast(img).enhance(1.5)
    img = img.convert("L")

    # Each braille char = 2x4 pixels
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
                    if invert:
                        is_on = val > threshold
                    else:
                        is_on = val < threshold
                    if is_on:
                        code |= bit
            line += chr(BRAILLE_BASE + code)
        lines.append(line)
    return "\n".join(lines)


if __name__ == "__main__":
    path = sys.argv[1] if len(sys.argv) > 1 else "NoLoveSong.jpg"
    width = int(sys.argv[2]) if len(sys.argv) > 2 else 45
    threshold = int(sys.argv[3]) if len(sys.argv) > 3 else 110
    print(image_to_braille(path, width, threshold, invert=False))
