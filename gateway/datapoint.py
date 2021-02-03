class Datapoint:
    def __init__(self, function_group, datapoint, function_name, datatype, function_number=0, **kwargs):
        self._function_group = function_group
        self._datapoint = datapoint
        self._function_name = function_name
        self._datatype = datatype
        self._function_number = function_number
        self._read = False
        self._write = False

        if "read" in kwargs:
            self._read = kwargs["read"]
        if "write" in kwargs:
            self._write = kwargs["write"]

    @property
    def write(self):
        return self._write

    @property
    def read(self):
        return self._read

    @property
    def function_name(self):
        return self._function_name

    @property
    def function_group(self):
        return self._function_group

    @property
    def function_number(self):
        return self._function_number

    @property
    def datapoint(self):
        return self._datapoint

    @property
    def datatype(self):
        return self._datatype

    def datapoint_by_bytes(self):
        return self._datapoint.to_bytes(2, byteorder='big', signed=False)


class DatapointList:
    def __init__(self, datapoint_list):
        self._datapoint_list = datapoint_list

    @property
    def datapoint_list(self):
        return self._datapoint_list

    def get_datapoint(self, function_group, function_number, point):
        for datapoint in self._datapoint_list:
            if (datapoint.function_group == function_group and datapoint.function_number == function_number
                    and datapoint.datapoint == point):
                return datapoint
        return None

    def get_datapoint_by_name(self, function_name):
        for datapoint in self._datapoint_list:
            if datapoint.function_name == function_name:
                return datapoint
        return None


class Device:
    def __init__(self, device_type, device_id):
        self._device_type = device_type
        self._device_id = device_id

    @property
    def device_type(self):
        return self._device_type

    @property
    def device_id(self):
        return self._device_id
