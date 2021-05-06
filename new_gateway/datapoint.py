from new_gateway.datatypes import UnknownDatatypeError, Unsigned, Signed, List, String

datapoints_by_name = {}
datapoints_by_id = {}


class NoDatapointFoundError(Exception):
    pass


class Datapoint:
    """Datapoint of CAN messages"""
    def __init__(self, datapoint_from_settings):
        self.name = datapoint_from_settings["name"]
        self.function_group = datapoint_from_settings["function_group"]
        self.function_number = datapoint_from_settings["function_number"]
        self.datapoint_id = datapoint_from_settings["datapoint_id"]
        self.datatype = datapoint_from_settings["type"]
        self.decimal = get_settings_data_save(datapoint_from_settings, "decimal", int)
        self.read = get_settings_data_save(datapoint_from_settings, "read", bool)
        self.write = get_settings_data_save(datapoint_from_settings, "write", bool)

    def get_datapoint_type(self):
        if self.datatype == "Unsigned":
            return Unsigned(self.decimal)
        elif self.datatype == "Signed":
            return Signed(self.decimal)
        elif self.datatype == "List":
            return List()
        elif self.datatype == "String":
            return String()
        else:
            raise UnknownDatatypeError


def get_settings_data_save(datapoint_from_settings, column, data_type):
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
    datapoints_by_id[(element["function_group"],
                      element["function_number"],
                      element["datapoint_id"])] = Datapoint(element)

    datapoints_by_name[element["name"]] = Datapoint(element)


def get_datapoint_by_name(name):
    """Get datapoint by name (used for writing messages"""
    if name not in datapoints_by_name:
        raise NoDatapointFoundError

    return datapoints_by_name[name]


def get_datapoint_by_id(function_group, function_number, datapoint_id):
    """Get datapoint by ids (used for reading messages)"""
    if (function_group, function_number, datapoint_id) not in datapoints_by_id:
        raise NoDatapointFoundError

    return datapoints_by_id[(function_group, function_number, datapoint_id)]
