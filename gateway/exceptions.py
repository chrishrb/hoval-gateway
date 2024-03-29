class InvalidFrame(Exception):
    """Invalid Frame returned by source handler"""
    pass


class NoValidMessageException(Exception):
    """If no valid message was found, e.g. message too small"""
    pass


class UnknownDatatypeError(Exception):
    """Raised if the datatype is unknown in the settings.yml"""
    pass


class NoDatapointFoundError(Exception):
    """If no datapoint was found in the settings file"""
    pass


class NoRequestFoundError(Exception):
    """If no Request was found"""
    pass


class VariableNotFoundError(Exception):
    """Variable not found in env / settings"""
    pass
