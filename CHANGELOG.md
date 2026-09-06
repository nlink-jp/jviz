# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- **Third-party attributions moved from `LICENSE` to `NOTICE.md`.** The
  appended notice defeated GitHub's licence-template match, so the repository
  was reported as "Other" rather than MIT. `LICENSE` is now the plain MIT text
  and the attributions live in `NOTICE.md`, following the convention already
  used by markdown-viewer. Nothing about the licensing changed.

## [1.2.0] - 2026-07-12

### Removed

- **darwin/amd64 (Intel) pre-built binary.** macOS releases now ship
  **arm64 only**, per the org-wide policy (darwin is Apple-Silicon only; no
  universal binaries). Intel Mac users can build from source.

### Changed

- **Linux release archives are now `.tar.gz`** (darwin/windows remain `.zip`),
  per `nlink-jp/.github` CONVENTIONS.md §Release Archive Standard.
- **darwin code-signature identifier** is now the canonical `jviz`
  (was `jviz-darwin-arm64`), set via `codesign -i` so it stays stable
  after the archived binary is renamed to its canonical name.

No change to the binary's behaviour — a packaging / build-config release.

## [1.1.1] - 2026-05-23

### Changed

- **Darwin releases are now Developer ID signed and Apple-notarized.**
  `jviz-v1.1.1-darwin-{amd64,arm64}.zip` carry full Apple Developer
  ID Application signatures and notarization tickets from Apple. End
  users on macOS no longer need to bypass Gatekeeper with right-click
  → Open or `xattr -d com.apple.quarantine` on first launch; local
  users who place `jviz` under Dropbox-synced (or any other
  FileProvider-managed) paths are no longer killed by macOS's
  ad-hoc + provenance distrust policy. Pipeline:
  `scripts/codesign-darwin.sh` + `scripts/notarize-darwin.sh`,
  driven by `make package`. Adopts the org-wide convention in
  `nlink-jp/.github` CONVENTIONS.md §Code Signing.
- **Release asset naming changed** from `jviz_<os>_<arch>.zip`
  (underscore-separated, version-less) to
  `jviz-vX.Y.Z-<os>-<arch>.zip` (hyphen-separated, versioned),
  aligning with the rest of util-series. LICENSE is now bundled
  alongside README.md in each zip.

No behaviour change to the binary itself — feature-wise this is
identical to v1.1.0.

## [1.1.0] - 2026-03-28

### Added

- Light / dark theme support.
  - Defaults to the OS `prefers-color-scheme` setting.
  - Toggle button (☀ / 🌙) in the header for manual override.
  - Selection is persisted in `localStorage`.
  - Chart.js axis and grid colours update automatically on theme change.

## [1.0.1] - 2026-03-28

### Internal

- Added 11 unit tests covering `sseMessage`, `Hub` (snapshot, subscribe/unsubscribe, fan-out), and `decodeJSONStream`.
- Extracted `decodeJSONStream` helper for testability.
- Fixed potential infinite loop in `decodeJSONStream` when invalid non-JSON input precedes a valid JSON array: decoder is now rebuilt from buffered bytes after skipping to the next newline.

## [1.0.0] - 2026-03-28

### Added

- Initial release.
- Read JSON array from stdin or file (`--watch`).
- Local HTTP server with SSE-based live updates.
- Interactive browser UI with Chart.js: bar, line, pie, and table views.
- Column selectors for X-axis label and Y-axis value.
- Auto-open browser on start (`--no-open` to disable).
- Cross-platform builds: Linux amd64/arm64, macOS amd64/arm64, Windows amd64.
