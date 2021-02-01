import can
import asyncio
import logging
import paho.mqtt.client as mqtt

from gateway import settings
from gateway.datapoint import Device
from gateway.parser import ResponseParser, Request, PeriodicRequest


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

    # Message Parser
    message_parser = ResponseParser()

    # Periodic Requests
    periodic_request = PeriodicRequest(Device(10, 8))
    periodic_request.start()

    # MQTT Client
    client = mqtt.Client("hoval-client")
    client.username_pw_set(username=settings.MQTT["BROKER_USERNAME"], password=settings.MQTT["BROKER_PASSWORD"])
    client.connect(settings.MQTT["BROKER"])

    while True:
        # Wait for next message from AsyncBufferedReader
        msg = await reader.get_message()
        parsed = message_parser.parse(msg)
        if parsed:
            logging.info(parsed)
            client.publish("hoval-gw/" + str(parsed[0]) + str(parsed[1]) + "/status")

    # Clean-up
    notifier.stop()
    can0.shutdown()


# Get the default event loop
loop = asyncio.get_event_loop()
# Run until main coroutine finishes
loop.run_until_complete(main())
loop.close()
