from new_gateway.datatypes import Unsigned, Signed, List, String
from new_gateway.exceptions import UnknownDatatypeError, NoDatapointFoundError

datapoints_by_name = {}
datapoints_by_id = {}


class Datapoint:
    """Datapoint of CAN messages"""

    def __init__(self, datapoint_from_settings):
        self.name = datapoint_from_settings["name"]
        self.function_group = datapoint_from_settings["function_group"]
        self.function_number = datapoint_from_settings["function_number"]
        self.datapoint_id = datapoint_from_settings["datapoint_id"]
        self.datatype = datapoint_from_settings["type"]
        self.decimal = _get_settings_data_safe(datapoint_from_settings, "decimal", int)

    def get_datapoint_type(self):
        """Get type of datapoint"""
        if self.datatype == "U8":
            return Unsigned(8, self.decimal)
        elif self.datatype == "U16":
            return Unsigned(16, self.decimal)
        elif self.datatype == "U32":
            return Unsigned(32, self.decimal)
        elif self.datatype == "S8":
            return Signed(8, self.decimal)
        elif self.datatype == "S16":
            return Signed(16, self.decimal)
        elif self.datatype == "S32":
            return Signed(32, self.decimal)
        elif self.datatype == "LIST":
            return List()
        elif self.datatype == "STR":
            return String()
        else:
            raise UnknownDatatypeError(str.format("Datatype {} not known", self.datatype))

    def get_datapoint_by_bytes(self):
        return self.datapoint_id.to_bytes(2, byteorder='big', signed=False)

    def __str__(self):
        return str.format("Datapoint: name: {}, function_group: {}, function_number: {}, datapoint_id: {}, "
                          "datatype: {}, decimal: {}",
                          self.name,
                          self.function_group,
                          self.function_number,
                          self.datapoint_id,
                          self.datatype,
                          self.decimal)


def _get_settings_data_safe(datapoint_from_settings, column, data_type):
    """Get settings data safe (if data not there, set default ones)"""
    if column in datapoint_from_settings:
        return datapoint_from_settings[column]
    else:
        if data_type == int or data_type == float:
            return 0
        elif data_type == str:
            return ""
        elif data_type == bool:
            return False


def parse_datapoints(element):
    """Save datapoints from settings.yaml"""
    for datapoint_item in element:
        dp = Datapoint(datapoint_item)
        datapoints_by_id[(datapoint_item["function_group"],
                          datapoint_item["function_number"],
                          datapoint_item["datapoint_id"])] = dp
        datapoints_by_name[datapoint_item["name"]] = dp


def get_datapoint_by_name(name):
    """Get datapoint by name (used for writing messages"""
    if name not in datapoints_by_name:
        raise NoDatapointFoundError(str.format("No datapoint found with name [{}]", name))

    return datapoints_by_name[name]


def get_datapoint_by_id(function_group, function_number, datapoint_id):
    """Get datapoint by ids (used for reading messages)"""
    if (function_group, function_number, datapoint_id) not in datapoints_by_id:
        raise NoDatapointFoundError(str.format("No datapoint found with [function_group: {}, function_number: {}, "
                                               "datapoint_id {}]", function_group, function_number, datapoint_id))

    return datapoints_by_id[(function_group, function_number, datapoint_id)]
