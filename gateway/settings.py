import logging

from gateway import datatypes
from gateway.datapoint import Datapoint

logging.basicConfig(level=logging.DEBUG)

OPERATIONS = {
    "RESPONSE": 0x42,
    "GET_REQUEST": 0x40,
    "SET_REQUEST": 0x46
}

DATAPOINT_LIST = {
    Datapoint(function_group=50, datapoint=0, function_name="temperature_outside_air", datatype=datatypes.Signed(1)),
    Datapoint(function_group=50, datapoint=37602, function_name="temperature_exhaust_air", datatype=datatypes.Signed(1)),
    Datapoint(function_group=50, datapoint=40650, function_name="betriebswahl_lueftung", datatype=datatypes.List()),
    Datapoint(function_group=50, datapoint=37608, function_name="voc_abluft", datatype=datatypes.Unsigned(1)),
    Datapoint(function_group=50, datapoint=37611, function_name="voc_aussenluft", datatype=datatypes.Unsigned(1))
}
