class NoValidMessageException(Exception):
    """If no valid message was found, e.g. message too small"""
    pass


class UnknownDatatypeError(Exception):
    """Raised if the datatype is unknown in the settings.yml"""
    pass


class NoDatapointFoundError(Exception):
    """If no datapoint was found in the settings file"""
    pass
