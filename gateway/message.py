import logging

from gateway import settings
from gateway.datapoint import DatapointList


class Message:
    _message_body = list()

    def __init__(self, arbitration_id, message_len, operation_id):
        self._arbitration_id = arbitration_id
        self._message_id = arbitration_id >> 24
        self._device_type = (arbitration_id >> 8) & 0xff
        self._device_id = arbitration_id & 0xff
        self._message_len = message_len >> 3
        self._nb_remaining = message_len+1
        self._operation_id = operation_id

    @property
    def message_id(self):
        return self._message_id

    @property
    def message_len(self):
        return self._message_len

    @property
    def operation_id(self):
        return self._operation_id

    @property
    def nb_remaining(self):
        return self._nb_remaining

    def put(self, body):
        if self._operation_id in settings.OPERATIONS.values():
            self._message_body.append(body)
            self._nb_remaining - 1

    def parse_data(self):
        if not self._message_body:
            return None

        data = bytearray()
        point = None
        for body in self._message_body:
            point = body.datapoint
            data += body.data

        if data and point is not None:
            return point.function_name, point.datatype.convert(data)

        return None


class Response:
    def __init__(self, data):
        datapoint_list = DatapointList(settings.DATAPOINT_LIST)
        function_group = data[2]
        function_number = data[3]
        function_datapoint = int.from_bytes(data[4:6], byteorder='big', signed=False)
        self._datapoint = datapoint_list.get_datapoint(function_group, function_number, function_datapoint)
        self._data = data[6:]
        if self._datapoint is None:
            logging.debug("No know point found for (%d, %d, %d), len %d", function_group,
                          function_number, function_datapoint, len(data))

    @property
    def data(self):
        return self._data

    @property
    def datapoint(self):
        return self._datapoint


class Request:
    def __init__(self, operation_id, function_name, data):
        pass
