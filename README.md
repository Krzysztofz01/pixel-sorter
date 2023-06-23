# Pixel-Sorter
[![Go Reference](https://pkg.go.dev/badge/github.com/Krzysztofz01/pixel-sorter.svg)](https://pkg.go.dev/github.com/Krzysztofz01/pixel-sorter)
[![Go Report Card](https://goreportcard.com/badge/github.com/Krzysztofz01/pixel-sorter)](https://goreportcard.com/report/github.com/Krzysztofz01/pixel-sorter)
![GitHub](https://img.shields.io/github/license/Krzysztofz01/pixel-sorter)
![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/Krzysztofz01/pixel-sorter?include_prereleases)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/Krzysztofz01/pixel-sorter)

Pixel sorting is a kind of photo editing, which is a subgenre of glitch art, the operation of which consists in reorganizing groups of pixels in a photo according to certain criteria. There are many programs whose task is to create the pixel sorting effect, but when creating my implementation in Go, I focused on optimizing this process in terms of time and creating a modular platform, thanks to which, implementation of new functionalities related to pixel sorting will be simple task.

There are two crucial stages in the sorting process, the first of which is the division of pixel rows into intervals, i.e. pixels that, according to certain criteria, have common features. The next step is to perform a sort by a certain pixel parameter in a given interval. Intervals can be determined based on the perceived brightness, hue and saturation of the HSL color space, by a mask which is a external black and white image showing which areas are to be sorted, or by performing edge detection. The sorting itself can also be performed on the basis of perceived brightness, HSL color space parameters, or instead of sorting, pixels can be arranged randomly. The sorting operation can be performed vertically and horizontally in any order once any number of times. It is possible to sort colors according to a given angle, it is aslo possible to set fixed length intervals.

# Requirements and installation
Required software:
- **git** - to download the source code from the repository
- **go (version: 1.19+)** - to build the source code locally

```sh
git clone https://github.com/Krzysztofz01/pixel-sorter

cd pixel-sorter

go build ./cli
```

# Usage

### Commands
- *brightness* - Use perceived brightness value as color sorting parameter.
- *hue* - Use hue value as color sorting parameter.
- *saturation* - Use saturation value as color sorting parameter.
- *help* - Print program help page.

### Flags
- *input-media-path* - The path of the input media file to be processed.
- *output-media-path* - The path of the output media file to be saved. The path should end with one of the supported extensions.
- *mask-image-path* - The path of the mask image file used to process the input media.

- *angle* (-a) - The angle at which to sort the pixels.
- *cycles* (-c) - The count of sorting cycles that should be performed on the image.
- *direction* (-d) - Pixel sorting direction in intervals.
    - *ascending* - Sort asceding according to the sorting determinant
    - *descending* - Sort descending according to the sorting determinant
    - *random* - Shuffle by the sorting determinant
- *interval-determinant* (-i) - Parameter used to determine intervals.
    - *brightness* - Use the perceived brightness to determine intervals
    - *hue* - Use the HSL color space hue value to determine intervals
    - *mask* - Use the external mask image to determine intervals
    - *absolute* - Use the color absolute value (old imprecise approach, but classic)
    - *edge* - Use a Canny edge detection algorithm to determine intervals
- *interval-lower-threshold* (-l) - The lower threshold of the interval determination process.
- *interval-upper-threshold* (-u) - The upper threshold of the interval determination process.
- *interval-max-length* (-k) - The max length of the interval. Zero means no length limits.
- *interval-max-length-random-factor* (-r) - The value representing the range of values that can be randomly subtracted or added to the max interval length.
- *mask* (-m) - Exclude the sorting effect from masked out ares of the image.
- *order* (-o) - Order of the graphic sorting stages.
    - *horizontal*
    - *vertical*
    - *horizontal-vertical*
    - *vertical-horizontal*
- *scale* (-s) - Image size downscale percentage factor (can be used to generate a low resolution preview).
- *blending-mode* (-b) - The blending mode algorithm to blend the original image with the sorted image.
    - *none*
    - *lighten*
    - *darken*
- *output-format* (-f) - The output format of the graphic file.
    - *jpg*
    - *png*

Output of the help command:
```sh
Pixel sorting image editing utility implemented in Go.

Usage:
  pixel-sorter [command]

Available Commands:
  brightness  Use brightness value as color sorting parameter.
  help        Help about any command
  hue         Use hue value as color sorting parameter.
  saturation  Use saturation value as color sorting parameter.

Flags:
  -a, --angle int                               The angle at which to sort the pixels.
  -b, --blending-mode string                    The blending mode algorithm to blend the sorted image into the original. Options: [none, lighten, darken]. (default "none")
  -c, --cycles int                              The count of sorting cycles that should be performed on the image. (default 1)
  -d, --direction string                        Pixel sorting direction in intervals. Options: [ascending, descending, random]. (default "ascending")
  -h, --help                                    help for pixel-sorter
      --input-media-path string                 The path of the input media file to be processed.
  -i, --interval-determinant string             Parameter used to determine intervals. Options: [brightness, hue, mask, absolute, edge]. (default "brightness")
  -l, --interval-lower-threshold float          The lower threshold of the interval determination process. Options: [0.0 - 1.0]. (default 0.1)
  -k, --interval-max-length int                 The max length of the interval. Zero means no length limits.
  -r, --interval-max-length-random-factor int   The value representing the range of values that can be randomly subtracted or added to the max interval length. Options: [0 <=]
  -u, --interval-upper-threshold float          The upper threshold of the interval determination process. Options: [0.0 - 1.0]. (default 0.9)
  -m, --mask                                    Exclude the sorting effect from masked out ares of the image.
      --mask-image-path string                  The path of the mask image file used to process the input media.
  -o, --order string                            Order of the graphic sorting stages. Options: [horizontal, vertical, horizontal-vertical, vertical-horizontal]. (default "horizontal-vertical")
      --output-media-path string                The path of the output media file to be saved. The path should end with one of the supported extensions. [jpg, png]
  -s, --scale float                             Image downscaling percentage factor. Options: [0.0 - 1.0]. (default 1)
  -v, --verbose                                 Enable verbose logging mode.

Use "pixel-sorter [command] --help" for more information about a command.
```