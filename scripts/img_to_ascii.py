"""Convert album cover image to ASCII art for terminal display."""
from PIL import Image, ImageEnhance
import sys


def image_to_ascii(image_path, width=70, invert=True):
    ascii_chars = " .:-=+*#%@"
    if invert:
        ascii_chars = ascii_chars[::-1]

    img = Image.open(image_path)
    aspect_ratio = img.height / img.width
    height = int(width * aspect_ratio * 0.45)
    img = img.resize((width, height))
    img = ImageEnhance.Contrast(img).enhance(1.4)
    img = img.convert("L")

    ascii_art = []
    for y in range(height):
        line = ""
        for x in range(width):
            pixel = img.getpixel((x, y))
            line += ascii_chars[pixel * (len(ascii_chars) - 1) // 255]
        ascii_art.append(line)
    return "\n".join(ascii_art)


if __name__ == "__main__":
    path = sys.argv[1] if len(sys.argv) > 1 else "NoLoveSong.jpg"
    width = int(sys.argv[2]) if len(sys.argv) > 2 else 60
    print(image_to_ascii(path, width))
