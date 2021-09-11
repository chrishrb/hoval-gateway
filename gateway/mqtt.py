import logging
import random

from paho.mqtt import client as mqtt_client

from gateway import datapoint
from gateway.datapoint import Datapoint
from gateway.exceptions import NoDatapointFoundError, NoRequestFoundError, NoValidMessageException
from gateway.message import SendMessage, build_arbitration_id, Operation
from gateway.request import get_subscribe_request_by_name


def connect_mqtt(mqtt_settings):
    def on_connect(client, userdata, flags, rc):
        if rc == 0:
            logging.info("Connected to MQTT Broker!")
        else:
            logging.error("Failed to connect, return code %d\n", rc)

    client_id = f'{mqtt_settings["name"]}-{random.randint(0, 1000)}'
    client = mqtt_client.Client(client_id)
    client.username_pw_set(mqtt_settings["username"], mqtt_settings["password"])
    client.on_connect = on_connect
    if mqtt_settings["port"] == 8883:
        client.tls_set_context()
    client.connect(mqtt_settings["broker"], mqtt_settings["port"])
    return client


class Subscriber:
    _operation = Operation.SET_REQUEST

    def __init__(self, can0, topic):
        self.can0 = can0
        self.topic = topic

    def on_message(self, client, userdata, msg):
        logging.debug("Message received from broker: %s on topic %s", msg.payload.decode(), msg.topic)

        datapoint_name = msg.topic[len(self.topic) + 1:]
        datapoint_name = datapoint_name[:-4]

        try:
            request = get_subscribe_request_by_name(datapoint_name)
            datapoint_of_message: Datapoint = datapoint.get_datapoint_by_name(request.datapoint_name)
        except (NoRequestFoundError, NoDatapointFoundError) as e:
            logging.error(e)
            return

        # build message
        arbitration_id = build_arbitration_id(31, request.priority, request.device_type, request.device_id)
        message = SendMessage(arbitration_id, self._operation.value, datapoint_of_message)

        # todo: support all types of all length
        # set data and send
        try:
            message.put_data(datapoint_of_message.get_datapoint_type().convert_to_bytes(msg.payload.decode()))
            can_message = message.to_can_message()
            self.can0.send(can_message)
            logging.debug("Send message: %s", message)
        except NoValidMessageException as e:
            logging.error(e)
