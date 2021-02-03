import logging

from gateway import settings
from gateway.core import OneTimeRequest
from gateway.datapoint import DatapointList


class Subscriber:
    _datapoint_list = DatapointList(settings.DATAPOINT_LIST)

    def __init__(self, bus):
        self._bus = bus

    def on_connect(self, client, userdata, flags, rc):
        if rc == 0:
            logging.debug("Connected to broker")
            for datapoint in self._datapoint_list.datapoint_list:
                if not datapoint.write:
                    continue

                client.subscribe(settings.MQTT["TOPIC"] + "/" + datapoint.function_name + "/set")
        else:
            logging.error("Connection failed")

    def on_message(self, client, userdata, message):
        logging.debug("Message received : " + str(message.payload) + " on " + message.topic)
        # Remove "/set"
        datapoint = self._datapoint_list.get_datapoint_by_name(message.topic[:-4])

        if datapoint is not None:
            if datapoint.datatype.check(message.payload):
                request = OneTimeRequest(self._bus, settings.TARGET, settings.OPERATIONS["SET_REQUEST"],
                                         datapoint.function_name)
                request.start(int(float(message.payload.decode())))
            else:
                logging.error("Datatype is not acceptet, required %s ", str(datapoint.datatype))
        else:
            logging.error("No datapoint found for %s", message.topic)
