from abc import abstractmethod


class Datatype:
    """Abstract Datatype"""
    @abstractmethod
    def convert(self, value):
        pass


class Unsigned(Datatype):
    """Unsigned Datatype"""
    def __init__(self, decimal=0):
        self._decimal = decimal

    def convert(self, value):
        val = int.from_bytes(value, byteorder='big', signed=False)
        return round(val * 10 ** (-self._decimal), 2)

    def __str__(self):
        return "Unsigned"


class Signed(Datatype):
    """Signed Datatype"""
    def __init__(self, decimal=0):
        self._decimal = decimal

    def convert(self, value):
        val = int.from_bytes(value, byteorder='big', signed=True)
        return round(val * 10 ** (-self._decimal), 2)

    def __str__(self):
        return "Signed"


class List(Datatype):
    """List Datatype"""
    def convert(self, value):
        val = int.from_bytes(value, byteorder='big', signed=False)
        return round(val)

    def __str__(self):
        return "List"


class String(Datatype):
    """String Datatype"""
    def convert(self, value):
        return value.decode('utf-8')

    def __str__(self):
        return "String"
