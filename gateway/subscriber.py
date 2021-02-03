import logging

from gateway import settings
from gateway.core import SimpleRequest
from gateway.datapoint import DatapointList


class Subscriber:
    """
    Subscriber class for all messages that are sent to the gateway
    """
    _datapoint_list = DatapointList(settings.DATAPOINT_LIST)

    def __init__(self, bus):
        self._bus = bus

    def on_connect(self, client, userdata, flags, rc):
        """
        Gateway connects to MQTT and subscribes the write datapoints

        :param client:
        :param userdata:
        :param flags:
        :param rc:
        :return:
        """
        if rc == 0:
            logging.debug("Connected to broker..")
            for datapoint in self._datapoint_list.datapoint_list:
                if not datapoint.write:
                    continue

                logging.debug("subscribed to %s", settings.MQTT["TOPIC"] + "/" + datapoint.function_name + "/set")
                client.subscribe(settings.MQTT["TOPIC"] + "/" + datapoint.function_name + "/set")
        else:
            logging.error("Connection failed")

    def on_message(self, client, userdata, message):
        """
        Gateway recieved message from MQTT Server and starts request to Hoval Device

        :param client:
        :param userdata:
        :param message:
        :return:
        """
        logging.debug("Message received : " + str(message.payload) + " on " + message.topic)
        # Remove TOPIC and "/set"
        datapoint_name = message.topic[len(settings.MQTT["TOPIC"])+1:]
        datapoint_name = datapoint_name[:-4]
        datapoint = self._datapoint_list.get_datapoint_by_name(datapoint_name)

        # todo: add more checks
        if datapoint is not None:
            request = SimpleRequest(self._bus, settings.TARGET, settings.OPERATIONS["SET_REQUEST"],
                                    datapoint.function_name)
            request.send(int(float(message.payload.decode())))
        else:
            logging.error("No datapoint found for %s", message.topic)
