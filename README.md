# ebiten_gvvideo

GV video player for [Ebitengine](https://ebitengine.org/), using [go-gv-video](https://github.com/funatsufumiya/go-gv-video)

( Partial port from [bevy_movie_player](https://github.com/funatsufumiya/bevy_movie_player) )

> [!WARNING]
> Go port was almost done by GitHub Copilot. Use with care.

## Limitation

Currently CPU decoding `DXT1/DXT3/DXT5` compressed data into `RGBA`, because Ebitengine now doesn't provide GPU compressed texture assignment.

## Example

```bash
$ go run ./example/main.go
```
