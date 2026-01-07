# TrueNAS SCALE Version Compatibility

This provider is built against the OpenAPI specification from **TrueNAS SCALE 25.10.1 (Goldeye)**.

## Minimum Requirements

- **TrueNAS SCALE**: Version 25.10.1 (Goldeye) or later
- **Terraform**: 0.13+

## API Compatibility

The provider uses the native JSON-RPC 2.0 API over WebSocket, which provides:
- Better performance than REST APIs
- Real-time communication
- Native TrueNAS protocol support

## Version Notes

- Generated from OpenAPI spec: TrueNAS SCALE 25.10.1
- 274 resources covering the complete API surface
- Backward compatibility with newer TrueNAS versions expected
- Older versions may have missing endpoints or different schemas
