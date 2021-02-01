from abc import abstractmethod


class Datatype:
    @abstractmethod
    def convert(self, value):
        pass


class Unsigned(Datatype):
    def __init__(self, decimal=0):
        self._decimal = decimal

    def convert(self, value):
        val = int.from_bytes(value, byteorder='big', signed=False)
        return val * 10 ** (-self._decimal)


class Signed(Datatype):
    def __init__(self, decimal=0):
        self._decimal = decimal

    def convert(self, value):
        val = int.from_bytes(value, byteorder='big', signed=True)
        return val * 10 ** (-self._decimal)


class List(Datatype):
    def convert(self, value):
        val = int.from_bytes(value, byteorder='big', signed=False)
        return val


class String(Datatype):
    def convert(self, value):
        return value.decode('utf-8')
