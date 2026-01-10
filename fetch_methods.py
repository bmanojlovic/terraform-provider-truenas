#!/usr/bin/env python3
"""
Fetch TrueNAS API methods specification via core.get_methods
Saves version-specific spec file: truenas-methods-{version}.json
"""
import asyncio
import websockets
import json
import ssl
import sys
from datetime import datetime

async def fetch_methods(host, token):
    uri = f'wss://{host}/websocket'
    
    ssl_context = ssl.create_default_context()
    ssl_context.check_hostname = False
    ssl_context.verify_mode = ssl.CERT_NONE
    
    async with websockets.connect(uri, ssl=ssl_context, max_size=50*1024*1024) as ws:
        # Connect
        await ws.send(json.dumps({"msg": "connect", "version": "1", "support": ["1"]}))
        await ws.recv()
        
        # Auth
        await ws.send(json.dumps({
            "msg": "method",
            "method": "auth.login_with_api_key",
            "params": [token],
            "id": "auth"
        }))
        await ws.recv()
        
        # Get version
        await ws.send(json.dumps({
            "msg": "method",
            "method": "system.version_short",
            "params": [],
            "id": "version"
        }))
        version_resp = json.loads(await ws.recv())
        version = version_resp['result']
        
        print(f"Fetching API spec for TrueNAS {version}...", file=sys.stderr)
        
        # Get methods
        await ws.send(json.dumps({
            "msg": "method",
            "method": "core.get_methods",
            "params": [],
            "id": "methods"
        }))
        methods_resp = json.loads(await ws.recv())
        methods = methods_resp['result']
        
        # Build spec with metadata
        spec = {
            "_metadata": {
                "truenas_version": version,
                "fetched_at": datetime.utcnow().isoformat() + "Z",
                "method_count": len(methods),
                "api_endpoint": uri
            },
            "methods": methods
        }
        
        # Save to version-specific file
        filename = f"truenas-methods-{version}.json"
        with open(filename, 'w') as f:
            json.dump(spec, f, indent=2)
        
        print(f"✅ Saved {len(methods)} methods to {filename}", file=sys.stderr)
        return filename, version

if __name__ == "__main__":
    import os
    
    # Read from environment
    host = os.getenv('TRUENAS_HOST')
    token = os.getenv('TRUENAS_TOKEN')
    
    if not host or not token:
        print("ERROR: Set TRUENAS_HOST and TRUENAS_TOKEN environment variables", file=sys.stderr)
        print("Example:", file=sys.stderr)
        print("  export TRUENAS_HOST=192.168.1.100", file=sys.stderr)
        print("  export TRUENAS_TOKEN=your-api-token", file=sys.stderr)
        sys.exit(1)
    
    filename, version = asyncio.run(fetch_methods(host, token))
    print(f"{filename}")
