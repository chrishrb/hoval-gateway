import struct
from abc import abstractmethod
from fastnumbers import fast_real

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
        value = fast_real(value)
        if not (isinstance(value, int) or isinstance(value, float)) or value < 0:
            raise NoValidMessageException(str.format("Message is no int/float or is smaller than 0: {}", value))
        value = int(value * 10 ** self._decimal)
        return struct.pack(self.get_format(), value)

    def get_format(self):
        if self._length == 8:
            return '!B'
        elif self._length == 16:
            return '!H'
        elif self._length == 32:
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
        value = fast_real(value)
        if not (isinstance(value, int) or isinstance(value, float)):
            raise NoValidMessageException(str.format("Message is no int/float: {}", value))
        value = int(value * 10 ** self._decimal)
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
        value = fast_real(value)
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
        return value.encode()

    def __str__(self):
        return "String"
