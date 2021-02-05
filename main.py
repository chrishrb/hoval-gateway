import can
import asyncio
import logging
import paho.mqtt.client as mqtt

from gateway import settings
from gateway.core import ResponseParser, PeriodicRequest
from gateway.subscriber import Subscriber


async def main():
    can0 = settings.CAN_BUS
    reader = can.AsyncBufferedReader()
    # logger = can.Logger('logfile.log')

    listeners = [
        reader,  # AsyncBufferedReader() listener
        # logger          # Regular Listener object
    ]
    # Create Notifier with an explicit loop to use for scheduling of callbacks
    event_loop = asyncio.get_event_loop()
    notifier = can.Notifier(can0, listeners, loop=event_loop)

    # Periodic Requests
    periodic_request = PeriodicRequest(can0, settings.TARGET)
    periodic_request.start()

    # Message Parser
    message_parser = ResponseParser()

    # MQTT Client
    client = mqtt.Client("hoval-client")
    client.username_pw_set(username=settings.MQTT["BROKER_USERNAME"], password=settings.MQTT["BROKER_PASSWORD"])
    if settings.MQTT["PORT"] == 8883:
        client.tls_set_context()
    client.connect(settings.MQTT["BROKER"], port=settings.MQTT["PORT"])

    # Receive Requests
    subscriber = Subscriber(can0)
    client.on_connect = subscriber.on_connect
    client.on_message = subscriber.on_message
    client.loop_start()

    while True:
        # Wait for next message from AsyncBufferedReader
        msg = await reader.get_message()
        parsed = message_parser.parse(msg)
        if parsed:
            client.publish(settings.MQTT["TOPIC"] + "/" + str(parsed[0]) + "/status", parsed[1])
            logging.info(settings.MQTT["TOPIC"] + "/" + str(parsed[0]) + "/status " + str(parsed[1]))

    # Clean-up
    # todo: does this work??
    periodic_request.stop()
    notifier.stop()
    can0.shutdown()


loop = asyncio.get_event_loop()
try:
    loop.run_until_complete(main())
finally:
    print("Exit..")
    loop.close()
