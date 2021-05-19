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
        if not (isinstance(value, int) or isinstance(value, float)) or value < 0:
            raise NoValidMessageException(str.format("Message is no int or is smaller than 0: {}", value))
        return [value]

    def __str__(self):
        return "Unsigned"


class Signed(Datatype):
    """Signed Datatype"""

    def __init__(self, length, decimal=0):
        self._length = length
        self._decimal = decimal

    def convert_from_bytes(self, value):
        val = int.from_bytes(value, byteorder='big', signed=True)
        return round(val * 10 ** (-self._decimal), 2)

    def convert_to_bytes(self, value):
        if not (isinstance(value, int) or isinstance(value, float)) or value < 0:
            raise NoValidMessageException(str.format("Message is no int: {}", value))
        return [value]

    def __str__(self):
        return "Signed"


class List(Datatype):
    """List Datatype"""

    _length = 8

    def convert_from_bytes(self, value):
        val = int.from_bytes(value, byteorder='big', signed=False)
        return round(val)

    def convert_to_bytes(self, value):
        if not isinstance(value, int) or value < 0:
            raise NoValidMessageException(str.format("Message is no int or is smaller than 0: {}", value))
        return [value]

    def __str__(self):
        return "List"


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