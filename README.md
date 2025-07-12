# gen-ai-lib

A Go library for aggregating generative AI models and APIs.

## Overview

gen-ai-lib is designed to provide a unified interface for interacting with multiple generative AI models and APIs. It enables developers to easily switch between different providers, aggregate results, and build robust AI-powered applications in Go.

## Features
- Unified interface for multiple generative AI providers (OpenAI, Cohere, etc.)
- Easy integration and extensibility
- Aggregation and fallback strategies
- GCP Cloud Storage helpers for uploading images and videos and returning public URLs
- Amazon S3 helpers for uploading images and videos with public access
- Designed for production use

## Getting Started

Coming soon: Installation and usage instructions.

### Video helpers

The `AppendVideos` function merges two MP4 clips using the `ffmpeg` command-line tool. You must have `ffmpeg` installed and accessible on your system `PATH`.

## License

MIT
