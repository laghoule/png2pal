# png2pal

`png2pal` is a command-line tool written in Go that converts a standard RGBA PNG image into a **paletted PNG** using a [GIMP palette file](https://www.gimp.org/) (`.gpl`). Each pixel is remapped to the closest color in the palette using Euclidean distance in RGB space, making it ideal for retro game development, pixel art workflows, and any project that requires a strict 256-color indexed image.

## Features

- Converts RGBA PNG images to indexed (paletted) PNG images.
- Accepts any standard 256-color GIMP palette (`.gpl`) file.
- Nearest-color matching via **Euclidean distance** in RGB space.
- Transparent pixels (alpha = 0) are automatically mapped to index 0.
- Produces a standard paletted PNG readable by any image editor or game engine.
- Multi-platform: Linux, macOS, and Windows (amd64 & arm64).
- Available as a standalone binary, Docker image, or buildable from source.

## How It Works

1. **Load the palette** – Parse the `.gpl` file and store all 256 RGB entries.
2. **Decode the source image** – Read the input PNG and verify it is in RGBA format.
3. **Remap every pixel** – For each pixel:
   - If the alpha channel is 0, the pixel is mapped to palette index 0 (transparent).
   - Otherwise, the closest palette color is found using the squared Euclidean distance formula:
     ```
     D² = (R1-R2)² + (G1-G2)² + (B1-B2)²
     ```
4. **Encode the output** – Write the resulting indexed image as a new PNG file.

## Prerequisites

- A **256-color GIMP palette** (`.gpl`) file. You can create one in GIMP via _Windows → Dockable Dialogs → Palettes_, or use any existing `.gpl` file with exactly 256 entries.
- A source PNG image in **RGBA format**.

## Installation

### Download a Pre-built Binary

Pre-built binaries for Linux, macOS, and Windows (amd64 and arm64) are available on the [Releases](https://github.com/laghoule/png2pal/releases) page.

Download the archive for your platform, extract it, and place the binary somewhere in your `$PATH`.

### Build from Source

Requires [Go](https://go.dev/) 1.26.1 or later.

```sh
git clone https://github.com/laghoule/png2pal.git
cd png2pal
go build -o png2pal ./cmd/main.go
```

### Docker

A Docker image is published to both the GitHub Container Registry and Docker Hub on every release.

**GitHub Container Registry:**

```sh
docker pull ghcr.io/laghoule/png2pal:latest
```

**Docker Hub:**

```sh
docker pull laghoule/png2pal:latest
```

Run the tool with Docker by mounting a local directory:

```sh
docker run --rm \
  -v "$(pwd)":/data \
  ghcr.io/laghoule/png2pal:latest \
  -src /data/input.png \
  -dst /data/output.png \
  -palette /data/my-palette.gpl
```

## Usage

```
png2pal -src <source file> -dst <destination file> -palette <GIMP palette file>
```

| Flag       | Description                                                       |
| ---------- | ----------------------------------------------------------------- |
| `-src`     | Path to the **source** PNG image (must be RGBA).                  |
| `-dst`     | Path to write the **output** paletted PNG image.                  |
| `-palette` | Path to a **GIMP palette** (`.gpl`) file with exactly 256 colors. |

**Example:**

```sh
png2pal -src tileset.png -dst tileset_indexed.png -palette my-palette.gpl
```

---

## GIMP Palette Format

`png2pal` expects a standard GIMP palette file (`.gpl`) with **exactly 256 color entries**. The file format looks like this:

```
GIMP Palette
Name: My Palette
Columns: 16
#
0   0   0   Index0
0   0 170   Index1
0 170   0   Index2
...
255 255 255 Index255
```

- Lines starting with `#` or empty lines are ignored.
- Each color entry is on its own line: `R G B comment`.
- **Index 0** is reserved for transparency — pixels with alpha = 0 will be mapped there.
- The palette must contain exactly 256 entries; `png2pal` will return an error otherwise.
