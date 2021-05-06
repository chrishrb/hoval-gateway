import logging
from enum import Enum

from new_gateway import datapoint


class NoMessageError(Exception):
    pass


class Operation(Enum):
    RESPONSE = 0x42
    GET_REQUEST = 0x40
    SET_REQUEST = 0x46

    @staticmethod
    def list():
        return list(map(lambda c: c.value, Operation))


class Message:
    """
    Message - Basic Class for all Messages
    """

    def __init__(self, arbitration_id, message_len, operation_id, message_header):
        self.arbitration_id = arbitration_id
        self.message_id = arbitration_id >> 24
        self.device_type = (arbitration_id >> 8) & 0xff
        self.device_id = arbitration_id & 0xff
        self.message_len = message_len
        self.nb_remaining = message_len + 1
        self.operation_id = operation_id
        self.message_header = message_header
        self.data = bytearray()

    def put(self, msg):
        if not self.is_valid(msg):
            return
        self.data.extend(msg.data)
        self.nb_remaining -= 1

    def put_extended_msg(self, msg):
        if not self.is_valid(msg):
            return
        self.data.extend(msg.data[:-2])
        self.nb_remaining -= 1

    def is_valid(self, msg):
        return len(msg.data) > 2 and self.operation_id in Operation.list()

    def __str__(self):
        return str.format("Message: id: {}, nb_remaining: {}, data: {}", self.message_id, self.nb_remaining, self.data)


class MessageResponse(Message):
    def __init__(self, arbitration_id, message_len, operation_id, message_header):
        super().__init__(arbitration_id, message_len, operation_id, message_header)

    def parse(self):
        logging.debug("message %s", self.data)
        function_group = self.data[2]
        function_number = self.data[3]
        function_datapoint = int.from_bytes(self.data[4:6], byteorder='big', signed=False)
        read_datapoint = datapoint.get_datapoint_by_id(function_group, function_number, function_datapoint)
        return read_datapoint, Operation(self.operation_id), read_datapoint.get_datapoint_type().convert(self.data[6:])


def get_message_id(arbitration_id):
    return arbitration_id >> 24


def get_message_header(data):
    return data[0]


def get_message_len(data):
    return data[0] >> 3


def get_operation_id(data):
    return data[1]
