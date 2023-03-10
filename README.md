# Pixel-Sorter

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
- *help* - Print program help page.

### Flags
- *image-file-path string* - The path of the image file to be processed.
- *mask-file-path string* - The path of the image mask file to be process the image file.

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
- *mask* (-m) - Exclude the sorting effect from masked out ares of the image.
- *order* (-o) - Order of the graphic sorting stages.
    - *horizontal*
    - *vertical*
    - *horizontal-vertical*
    - *vertical-horizontal*
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

Flags:
  -a, --angle int                        The angle at which to sort the pixels.
  -c, --cycles int                       The count of sorting cycles that should be performed on the image. (default 1)
  -d, --direction string                 Pixel sorting direction in intervals. Options: [ascending, descending, random]. (default "ascending")
  -h, --help                             help for pixel-sorter
      --image-file-path string           The path of the image file to be processed.
  -i, --interval-determinant string      Parameter used to determine intervals. Options: [brightness, hue, mask, absolute, edge]. (default "brightness")
  -l, --interval-lower-threshold float   The lower threshold of the interval determination process. Options: [0.0 - 1.0]. (default 0.1)
  -k, --interval-max-length int          The max length of the interval. Zero means no length limits.
  -u, --interval-upper-threshold float   The upper threshold of the interval determination process. Options: [0.0 - 1.0]. (default 0.9)
  -m, --mask                             Exclude the sorting effect from masked out ares of the image.
      --mask-file-path string            The path of the image mask file to be process the image file.
  -o, --order string                     Order of the graphic sorting stages. Options: [horizontal, vertical, horizontal-vertical, vertical-horizontal]. (default "horizontal-vertical")
  -f, --output-format string             The output format of the graphic file. Options: [jpg, png]. (default "jpg")

Use "pixel-sorter [command] --help" for more information about a command.
```