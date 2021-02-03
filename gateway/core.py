import logging

from gateway import settings
from gateway.datapoint import Device, DatapointList
from gateway.message import Message, MessageResponse, MessageRequest


class ResponseParser:
    """
    Parser for all messages
    """
    _pending_msg = {}
    _devices = {}

    def parse(self, msg):
        """
        Parse all messages

        :param msg:
        :return:
        """
        if len(msg.data) < 2:
            logging.error("Message too small")
            return None

        message = Message(msg.arbitration_id, (msg.data[0] >> 3), msg.data[1])
        if message.operation_id in settings.OPERATIONS.values():
            if message.device_id not in self._devices:
                logging.info("New device detected - Id: %d / Type: %d", message.device_id, message.device_type)
                self._devices[message.device_id] = Device(message.device_id, message.device_type)

            if message.message_id == 0x1f:
                if message.message_len == 0:
                    message.put(MessageResponse(msg.data))
                    if message.operation_id == settings.OPERATIONS["RESPONSE"]:
                        return message.parse_data()
                    else:
                        logging.debug("Message data: " + str(message.parse_data()))
                        logging.debug("arbitration_id: " + str(message.arbitration_id))
                else:
                    self._pending_msg[message.operation_id] = {
                        message
                    }
            else:
                msg_header = msg.data[0]
                if msg_header in self._pending_msg:
                    message = self._pending_msg[msg_header]
                    message.put(MessageResponse(msg.data))
                    if message.nb_remaining == 0:
                        del self._pending_msg[msg_header]

                        if message.operation_id == settings.OPERATIONS["RESPONSE"]:
                            return message.parse_data()
                        else:
                            logging.debug(message.parse_data())
        return None


class PeriodicRequest:
    """
    Periodic Requests to ask different datapoints (e.g. temperature, ..)
    """
    _message_len = 1
    # todo: is this id okay??
    _prio = 8160
    _operation_id = settings.OPERATIONS["GET_REQUEST"]
    _datapoint_list = DatapointList(settings.DATAPOINT_LIST)
    _tasks = []

    def __init__(self, bus, device):
        self._device = device
        self._bus = bus

    def start(self):
        """
        Start Periodic Requests

        :return:
        """
        for datapoint in self._datapoint_list.datapoint_list:
            if not datapoint.read:
                continue

            arbitration_id = (self._prio << 16) | (self._device.device_type << 8) | self._device.device_id
            message = Message(arbitration_id=arbitration_id, message_len=self._message_len,
                              operation_id=self._operation_id)
            message.put(MessageRequest(datapoint.function_name))
            task = message.send_periodic(self._bus)
            self._prio += 1
            self._tasks += task

    def stop(self):
        """
        Stopp all tasks

        :return:
        """
        for task in self._tasks:
            task.stop()


class SimpleRequest:
    """
    Simple request for one time requests
    """
    _message_len = 1
    _prio = 8160
    _datapoint_list = DatapointList(settings.DATAPOINT_LIST)

    def __init__(self, bus, device, operation, function_name):
        self._device = device
        self._arbitration_id = (self._prio << 16) | (self._device.device_type << 8) | self._device.device_id
        self._function_name = function_name
        self._operation_id = operation
        self._bus = bus

    def send(self, data):
        """
        Send a simple request

        :param data:
        :return:
        """
        datapoint = self._datapoint_list.get_datapoint_by_name(function_name=self._function_name)
        if not datapoint:
            return

        message = Message(arbitration_id=self._arbitration_id, message_len=self._message_len,
                          operation_id=self._operation_id)
        message.put(MessageRequest(datapoint.function_name, data))
        message.send(self._bus)
        logging.debug("Request sent with data %s", str(data))
