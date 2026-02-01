eBUS registry, schemas, and vendor providers.

## TinyGo

- CI builds TinyGo for main packages only; when none exist, the build is skipped.
- Some TinyGo installs require a board-specific target (e.g., `esp32-coreboard-v2`) if `-target esp32` reports missing pin definitions.
