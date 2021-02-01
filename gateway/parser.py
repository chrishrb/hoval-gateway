import logging

from gateway import settings
from gateway.datapoint import Device, DatapointList
from gateway.message import Message, Response, Request


class ResponseParser:
    _pending_msg = {}
    _devices = {}

    def parse(self, msg):
        if len(msg.data) < 2:
            logging.error("Message too small")
            return None

        message = Message(msg.arbitration_id, msg.data[0], msg.data[1])
        if message.operation_id == settings.OPERATIONS["RESPONSE"]:
            if message.device_id not in self._devices:
                logging.info("New device detected - Id: %d / Type: %d", message.device_id, message.device_type)
                self._devices[message.device_id] = Device(message.device_id, message.device_type)

            if message.message_id == 0x1f:
                if message.message_len == 0:
                    message.put(Response(msg.data))
                    return message.parse_data()
                else:
                    self._pending_msg[message.operation_id] = {
                        message
                    }
            else:
                msg_header = msg.data[0]
                if msg_header in self._pending_msg:
                    message = self._pending_msg[msg_header]
                    message.put(Response(msg.data))
                    if message.nb_remaining == 0:
                        del self._pending_msg[msg_header]
                        return message.parse_data()

        return None


class PeriodicRequest:
    _message_len = 1
    _prio = 8160
    _operation_id = settings.OPERATIONS["GET_REQUEST"]
    _datapoint_list = DatapointList(settings.DATAPOINT_LIST)

    def __init__(self, device):
        self._device = device
        self._arbitration_id = (self._prio << 16) | (self._device.device_type << 8) | self._device.device_id

    def start(self):
        for datapoint in self._datapoint_list.datapoint_list:
            if not datapoint.send_periodic:
                continue

            message = Message(arbitration_id=self._arbitration_id, message_len=self._message_len,
                              operation_id=self._operation_id)
            message.put(Request(datapoint.function_name))
            message.send_periodic()


class OneTimeRequest:
    _message_len = 1
    _prio = 8160
    _datapoint_list = DatapointList(settings.DATAPOINT_LIST)

    def __init__(self, device, operation, function_name):
        self._device = device
        self._arbitration_id = (self._prio << 16) | (self._device.device_type << 8) | self._device.device_id
        self._function_name = function_name
        self._operation_id = operation

    def start(self):
        datapoint = self._datapoint_list.get_datapoint_by_name(function_name=self._function_name)

        message = Message(arbitration_id=self._arbitration_id, message_len=self._message_len,
                          operation_id=self._operation_id)
        message.put(Request(datapoint.function_name))
        message.send()
