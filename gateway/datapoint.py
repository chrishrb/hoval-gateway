from abc import abstractmethod

from pythonlangutil.overload import Overload, signature

from gateway import datatypes


class Datapoint:
    def __init__(self, function_group, datapoint, function_name, datatype, function_number=0):
        self._function_group = function_group
        self._datapoint = datapoint
        self._function_name = function_name
        self._datatype = datatype
        self._function_number = function_number

    def function_name(self):
        return self._function_name

    def function_group(self):
        return self._function_group

    def function_number(self):
        return self._function_number

    def datapoint(self):
        return self._datapoint

    def datatype(self):
        return self._datatype


class DatapointList:
    datapoint_list = {
        Datapoint(function_group=50, datapoint=0, function_name="temperature_outside_air", datatype=datatypes.Signed(1)),
        Datapoint(function_group=50, datapoint=37602, function_name="temperature_exhaust_air", datatype=datatypes.Signed(1))
    }

    @Overload
    @signature("int", "int", "int")
    def get_datapoint(self, function_group, function_number, point):
        for datapoint in self.datapoint_list:
            if (datapoint.function_group == function_group and datapoint.function_number == function_number
                    and datapoint.datapoint == point):
                return datapoint
        return None

    @Overload
    @signature("String")
    def get_datapoint(self, function_name):
        for datapoint in self.datapoint_list:
            if datapoint.function_name == function_name:
                return datapoint
        return None
