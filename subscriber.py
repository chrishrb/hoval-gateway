import logging

import paho.mqtt.client as mqtt

from gateway import settings
from gateway.datapoint import Device
from gateway.core import OneTimeRequest


def on_connect(client, userdata, flags, rc):
    logging.debug("Connected with result code " + str(rc))
    client.subscribe("hoval-gw/lueftungs_modulation/set")


def on_message(client, userdata, msg):
    request = OneTimeRequest(Device(10, 8), settings.OPERATIONS["SET_REQUEST"], "normal_lueftungs_modulation")
    request.start(int(float(msg.payload.decode())))
    logging.debug("Request sent!")


client = mqtt.Client("hoval-client")
client.username_pw_set(username=settings.MQTT["BROKER_USERNAME"], password=settings.MQTT["BROKER_PASSWORD"])
client.connect(settings.MQTT["BROKER"])

client.on_connect = on_connect
client.on_message = on_message

client.loop_forever()
