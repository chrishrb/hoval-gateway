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

    # todo: only int are accepted
    def check(self, value):
        return isinstance(value, int)

    def __str__(self):
        return "Unsigned"


class Signed(Datatype):
    def __init__(self, decimal=0):
        self._decimal = decimal

    def convert(self, value):
        val = int.from_bytes(value, byteorder='big', signed=True)
        return val * 10 ** (-self._decimal)

    # todo: only int are accepted
    def check(self, value):
        return isinstance(value, int)

    def __str__(self):
        return "Signed"


class List(Datatype):
    def convert(self, value):
        val = int.from_bytes(value, byteorder='big', signed=False)
        return val

    def check(self, value):
        return isinstance(value, int)

    def __str__(self):
        return "List"


class String(Datatype):
    def convert(self, value):
        return value.decode('utf-8')

    def check(self, value):
        return isinstance(value, str)

    def __str__(self):
        return "String"
