# Wavis

Wavis reads WAVE PCM audio files and creates visually appealing waveforms that it can output in both SVG and ASCII formats. The project has no external dependencies.

## Usage

Passing a .wav file and using the default options outputs the audio file's properties (in a format similar to [Soxi](https://linux.die.net/man/1/soxi)'s) and a waveform:

![defaults](https://user-images.githubusercontent.com/1272713/233773663-cd70f417-c53b-414e-8cd9-d09e96d66ae6.png)


### Options

| Option | Description |
| --- | --- |
| `format` | A number representing the waveform format: <ul><li>`1`: blob SVG</li><li>`2`: single line "wavy" SVG</li><li>`3`: radial SVG</li><li>`4`: ASCII</li></ul>If no format is specified, the program outputs a file summary and an ASCII waveform.|
| `width` | Waveform's width, in characters for the ASCII format or in pixels for the other formats.  |
| `height` | Waveform's height, in lines for the ASCII format or in pixels for the other formats. |
| `padding` | Waveform's vertical padding, in lines for the ASCII format or in pixels for the other formats. |
| `resolution` | SVG formats only: the number of data points that should be represented per second. |
| `radius` | Inner circle radius; only applies to the radial format. |
| `border` | ASCII only: whether the rectangle enclosing the waveform should have a border; `0` or `1`. |
| `chars` | ASCII only: a string of 2 characters, where the first is the character the waveform is drawn with (defaults to `•`, while the other is the character used for drawind the negative space (defaults to ` `). Accepts any Unicode characters, including emojis.|

_Note: ASCII is loosely used here to refer to [text based visual art](https://en.wikipedia.org/wiki/ASCII_art) in general, which often uses non-ASCII characters._ 

The SVGs are output as plain text which you can pipe into a file and can be easily styled using CSS.

### Usage examples

```
wavis file.wav
wavis -format=1 -width=1000 -height=350 -padding=20 -resolution=5 file.wav > output.svg
wavis -format=3 -width=500 -padding=10 -resolution=10 -circle-radius=100 file.wav > output.svg
wavis -format=4 -width=100 -chars=":" -border=0 file.wav
wavis -format=4 -width=60 -height=20 -chars="✨💯" file.wav
```

### Examples of generated waveforms
| Waveform | Description |
| :---: | --- |
|<img src="https://user-images.githubusercontent.com/1272713/233843544-44f8a7b6-ee85-42c9-b946-044b8d7008f8.png" width="700"/>|`format=1`, `resolution=30`, raw|
|<img src="https://user-images.githubusercontent.com/1272713/233843566-4d011097-9e0a-4647-b40d-2034efabcd94.png" width="700"/>|`format=1`, `resolution=30`, styled with CSS|
|<img src="https://user-images.githubusercontent.com/1272713/233843667-0c4b0fbb-68c5-490b-a74d-29ae4d583895.png" width="700"/>|`format=1`, `resolution=3`, styled with CSS|
|<img src="https://user-images.githubusercontent.com/1272713/233843681-1be27bba-38b2-4fbb-afa3-81610effe6f4.png" width="700"/>|`format=1`, `resolution=3`, styled with CSS|
|<img src="https://user-images.githubusercontent.com/1272713/233843719-f8797fd6-74ef-4e42-bb3f-178889c9cf23.png" width="700"/>|`format=2`, `resolution=5`, raw|
|<img src="https://user-images.githubusercontent.com/1272713/233843733-95eaacc2-1500-42c6-898d-37607429bf74.png" width="700"/>|`format=2`, `resolution=15`, raw
|<img src="https://user-images.githubusercontent.com/1272713/233843747-54e35ff3-0c10-4fcb-a593-efa4c184fb32.png" width="700"/>|`format=2`, `resolution=25`, raw|
|<img src="https://user-images.githubusercontent.com/1272713/233843808-f3b0fcc3-398e-41ad-a69a-b8b117400696.png" width="400"/>|`format=3`, `resolution=25`, styled with CSS|
|<img src="https://user-images.githubusercontent.com/1272713/233843833-dd0f107c-ec13-48d5-ad14-294751a77735.png" width="400"/>|`format=3`, `resolution=5`, styled with CSS|
|<img src="https://user-images.githubusercontent.com/1272713/233843844-b5efd19d-5057-434f-b2fd-9fdf69430f1d.png" width="400"/>|`format=3`, `resolution=5`, styled with CSS|
|<img src="https://user-images.githubusercontent.com/1272713/233843913-29fd167b-e039-45bf-83d7-0b9a2ff0d235.png" width="700"/>|`format=4`, `chars=": "`|
|<img src="https://user-images.githubusercontent.com/1272713/233843928-7ca2d04b-54a1-46a4-82e1-6815f0a14732.png" width="700"/>|`format=4`, `chars=" :"`|
|<img src="https://user-images.githubusercontent.com/1272713/233843940-34ec4b2f-172d-4141-82fd-b25028ea4230.png" width="700"/>|`format=4`, `chars="• "`, `border=1`|
|<img src="https://user-images.githubusercontent.com/1272713/233843957-30214b72-86d5-4fb8-a6c7-2d50c524fce7.png" width="700"/>|`format=4`, `chars="🌟🖤"`|

## License

Wavis is [MIT](LICENSE.md) licensed. 
