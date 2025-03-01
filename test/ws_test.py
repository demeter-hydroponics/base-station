#!/usr/bin/env python3
import json
import asyncio
import websockets
import argparse

async def connect_to_websocket(uri):
    """
    Connect to a WebSocket endpoint and print received JSON messages.
    """
    print(f"Connecting to WebSocket at {uri}...")
    try:
        async with websockets.connect(uri) as websocket:
            print(f"Connected to {uri}")
            # Keep listening for messages
            while True:
                try:
                    # Receive message
                    message = await websocket.recv()
                    print("recieved message")
                except websockets.exceptions.ConnectionClosed:
                    print("Connection closed")
                    break
    except Exception as e:
        print(f"Error: {e}")


def main():
    parser = argparse.ArgumentParser(description='WebSocket client to receive and print JSON messages')
    parser.add_argument('uri', help='WebSocket URI (e.g., ws://localhost:8080/socket)')
    args = parser.parse_args()
    # Run the async function
    asyncio.run(connect_to_websocket(args.uri))

if __name__ == "__main__":
    main()
