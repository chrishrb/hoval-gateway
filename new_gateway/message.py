from enum import Enum

import can

from new_gateway import datapoint
from new_gateway.exceptions import NoValidMessageException


class Operation(Enum):
    """Operations in CAN-Message"""
    RESPONSE = 0x42
    GET_REQUEST = 0x40
    SET_REQUEST = 0x46

    @staticmethod
    def list():
        return list(map(lambda c: c.value, Operation))


class Message:
    """Message - Basic Class for all Messages"""

    def __init__(self, arbitration_id, operation_id, message_len=0):
        self.arbitration_id = arbitration_id
        self.message_id = arbitration_id >> 24
        self.priority = (arbitration_id >> 16) & 0xff
        self.device_type = (arbitration_id >> 8) & 0xff
        self.device_id = arbitration_id & 0xff
        self.message_len = 0
        self.nb_remaining = message_len + 1
        self.operation_id = operation_id
        self.data = bytearray()

    def put(self, data):
        if not self._is_valid(data):
            return
        self.data.extend(data)
        self.nb_remaining -= 1
        self.message_len += 1

    def put_extended_msg(self, data):
        if not self._is_valid(data):
            return
        self.data.extend(data[:-2])
        self.nb_remaining -= 1
        self.message_len += 1

    @staticmethod
    def _is_valid(data):
        return len(data) > 2

    def parse(self):
        if len(self.data) < 6:
            raise NoValidMessageException("Message too short for parsing")
        function_group = self.data[2]
        function_number = self.data[3]
        datapoint_id = int.from_bytes(self.data[4:6], byteorder='big', signed=False)
        read_datapoint = datapoint.get_datapoint_by_id(function_group, function_number, datapoint_id)
        return read_datapoint, Operation(self.operation_id), read_datapoint.get_datapoint_type().convert(self.data[6:])

    def to_can_message(self, datapoint_of_message):
        datapoint_by_bytes = datapoint_of_message.get_datapoint_by_bytes()
        return can.Message(arbitration_id=self.arbitration_id, data=[
            self.message_len, self.operation_id, datapoint_of_message.function_group,
            datapoint_of_message.function_number, datapoint_by_bytes[0],
            datapoint_by_bytes[1], self.data], is_extended_id=True)

    def __str__(self):
        return str.format("Message: id: {}, priority: {}, operation: {}, nb_remaining: {}, device_type: {}, "
                          "device_id: {}, data: {}",
                          self.message_id,
                          self.priority,
                          Operation(self.operation_id),
                          self.nb_remaining,
                          self.device_type,
                          self.device_id,
                          self.data)


def get_message_id(arbitration_id):
    return arbitration_id >> 24


def get_message_header(data):
    return data[0]


def get_message_len(data):
    return data[0] >> 3


def get_operation_id(data):
    return data[1]


def build_arbitration_id(priority, device_type, device_id):
    return (priority << 16) | (device_type << 8) | device_id
