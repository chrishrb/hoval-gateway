import can
import asyncio
import logging
import paho.mqtt.client as mqtt

from gateway.parser import ResponseParser


async def main():
    can0 = can.Bus(channel='can0', bustype='socketcan', receive_own_messages=False)
    reader = can.AsyncBufferedReader()
    # logger = can.Logger('logfile.log')

    listeners = [
        reader,  # AsyncBufferedReader() listener
        # logger          # Regular Listener object
    ]
    # Create Notifier with an explicit loop to use for scheduling of callbacks
    event_loop = asyncio.get_event_loop()
    notifier = can.Notifier(can0, listeners, loop=event_loop)

    message_parser = ResponseParser()
    while True:
        # Wait for next message from AsyncBufferedReader
        msg = await reader.get_message()
        parsed = message_parser.parse(msg)
        if parsed:
            logging.info(parsed)
            print("hoval-gw/" + parsed[0], parsed[1])

    # Clean-up
    notifier.stop()
    can0.shutdown()


# Get the default event loop
loop = asyncio.get_event_loop()
# Run until main coroutine finishes
loop.run_until_complete(main())
loop.close()
