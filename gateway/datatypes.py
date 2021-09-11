import re
import struct
from abc import abstractmethod

from gateway.exceptions import NoValidMessageException


class Datatype:
    """Abstract Datatype"""

    @abstractmethod
    def convert_from_bytes(self, value):
        pass

    @abstractmethod
    def convert_to_bytes(self, value):
        pass


class Unsigned(Datatype):
    """Unsigned Datatype"""

    def __init__(self, length, decimal=0):
        self._length = length
        self._decimal = decimal

    def convert_from_bytes(self, value):
        val = int.from_bytes(value, byteorder='big', signed=False)
        return round(val * 10 ** (-self._decimal), 2)

    def convert_to_bytes(self, value):
        value = parse_str(value)
        if not (isinstance(value, int) or isinstance(value, float)) or value < 0:
            raise NoValidMessageException(str.format("Message is no int/float or is smaller than 0: {}", value))
        return struct.pack('!H', value)

    @staticmethod
    def get_format(length):
        if length == 8:
            return '!B'
        elif length == 16:
            return '!H'
        elif length == 32:
            return '!I'

    def __str__(self):
        return "Unsigned with length {} bit and decimal {}".format(self._length, self._decimal)


class Signed(Datatype):
    """Signed Datatype"""

    def __init__(self, length, decimal=0):
        self._length = length
        self._decimal = decimal

    def convert_from_bytes(self, value):
        val = int.from_bytes(value, byteorder='big', signed=True)
        return round(val * 10 ** (-self._decimal), 2)

    def convert_to_bytes(self, value):
        value = parse_str(value)
        if not isinstance(value, int) or isinstance(value, float):
            raise NoValidMessageException(str.format("Message is no int/float: {}", value))
        return struct.pack('!h', value)

    @staticmethod
    def get_format(length):
        if length == 8:
            return '!b'
        elif length == 16:
            return '!h'
        elif length == 32:
            return '!i'

    def __str__(self):
        return "Signed with length {} bit and decimal {}".format(self._length, self._decimal)


class List(Datatype):
    """List Datatype"""

    _length = 8

    def convert_from_bytes(self, value):
        val = int.from_bytes(value, byteorder='big', signed=False)
        return round(val)

    def convert_to_bytes(self, value):
        value = parse_str(value)
        if not isinstance(value, int) or value < 0:
            raise NoValidMessageException(str.format("Message is no int or is smaller than 0: {}", value))
        return struct.pack('!B', value)

    def __str__(self):
        return "List with length {} bit".format(self._length)


class String(Datatype):
    """String Datatype"""

    def convert_from_bytes(self, value):
        return value.decode('utf-8')

    def convert_to_bytes(self, value):
        if not isinstance(value, str):
            raise NoValidMessageException(str.format("Message is no str: {}", value))
        # todo: test if string encoding is working
        return value.encode()

    def __str__(self):
        return "String"


def parse_str(num):
    """
    Parse a string that is expected to contain a number.
    :param num: str. the number in string.
    :return: float or int. Parsed num.
    """
    if re.compile('^\s*\d+\s*$').search(num):
        return int(num)
    if re.compile('^\s*(\d*\.\d+)|(\d+\.\d*)\s*$').search(num):
        return float(num)
    raise NoValidMessageException('num is not a number. Got {}.'.format(num))  # optional
