import logging

from gateway import datatypes
from gateway.datapoint import Datapoint

logging.basicConfig(level=logging.DEBUG)

MQTT = {
    "BROKER": "192.168.14.2",
    "BROKER_USERNAME": "admin",
    "BROKER_PASSWORD": "admin"
}

PERIODIC_TIME = 20

OPERATIONS = {
    "RESPONSE": 0x42,
    "GET_REQUEST": 0x40,
    "SET_REQUEST": 0x46
}

DATAPOINT_LIST = {
    Datapoint(function_group=50, datapoint=40650, function_name="betriebswahl_lueftung", datatype=datatypes.List()),
    Datapoint(function_group=50, datapoint=40651, function_name="normal_lueftungs_modulation", datatype=datatypes.List()),
    Datapoint(function_group=50, datapoint=40686, function_name="spar_lueftungs_modulation", datatype=datatypes.List()),
    Datapoint(function_group=50, datapoint=38606, function_name="lueftungs_modulation", datatype=datatypes.List()),
    Datapoint(function_group=50, datapoint=40687, function_name="feuchte_sollwert", datatype=datatypes.Unsigned(1)),
    Datapoint(function_group=50, datapoint=37600, function_name="feuchtigkeit_abluft", datatype=datatypes.Unsigned(1)),
    Datapoint(function_group=50, datapoint=37608, function_name="voc_abluft", datatype=datatypes.Unsigned(1)),
    Datapoint(function_group=50, datapoint=37611, function_name="voc_aussenluft", datatype=datatypes.Unsigned(1)),
    Datapoint(function_group=50, datapoint=39600, function_name="luftqualitaet_regulierung", datatype=datatypes.Unsigned(1)),
    Datapoint(function_group=50, datapoint=39652, function_name="status_lueftungsregulierung",
              datatype=datatypes.Unsigned(1)),
    Datapoint(function_group=50, datapoint=0, function_name="temperature_outside_air", datatype=datatypes.Signed(1),
              send_periodic=True),
    Datapoint(function_group=50, datapoint=37602, function_name="temperature_exhaust_air", datatype=datatypes.Signed(1),
              send_periodic=True),
    Datapoint(function_group=50, datapoint=38600, function_name="ventilator_fortluft_soll", datatype=datatypes.List()),
    # Fehler!
}
