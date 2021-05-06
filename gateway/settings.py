import logging
import os

import can
from dotenv import load_dotenv

from new_gateway import datatypes
from gateway.datapoint import Datapoint, Device

load_dotenv()
if os.getenv("DEBUG"):
    logging.basicConfig(level=logging.DEBUG)
else:
    logging.basicConfig(level=logging.INFO)

MQTT = {
    "BROKER": os.getenv("MQTT_BROKER"),
    "PORT": int(os.getenv("MQTT_PORT", 1883)),
    "BROKER_USERNAME": os.getenv("MQTT_BROKER_USERNAME"),
    "BROKER_PASSWORD": os.getenv("MQTT_BROKER_PASSWORD"),
    "TOPIC": os.getenv("MQTT_TOPIC", "hoval-gw")
}

PERIODIC_TIME = os.getenv("PERIODIC_TIME", 30)

OPERATIONS = {
    "RESPONSE": 0x42,
    "GET_REQUEST": 0x40,
    "SET_REQUEST": 0x46
}

TARGET = Device(10, 8)

CAN_BUS = can.Bus(channel='can0', bustype='socketcan', receive_own_messages=False)

# todo: make optimizations
# todo: make ready for all devices??
# todo: translate to english
DATAPOINT_LIST = {
    Datapoint(function_group=50, datapoint=40650, function_name="betriebswahl_lueftung", datatype=datatypes.List(),
              read=True, write=True),
    Datapoint(function_group=50, datapoint=40651, function_name="normal_lueftungs_modulation",
              datatype=datatypes.List(), write=True, read=True),
    Datapoint(function_group=50, datapoint=40686, function_name="spar_lueftungs_modulation", datatype=datatypes.List(),
              read=True, write=True),
    Datapoint(function_group=50, datapoint=38606, function_name="lueftungs_modulation", datatype=datatypes.List(),
              read=True),
    Datapoint(function_group=50, datapoint=40687, function_name="feuchte_sollwert", datatype=datatypes.Unsigned(1),
              read=True),
    Datapoint(function_group=50, datapoint=37600, function_name="feuchtigkeit_abluft", datatype=datatypes.Unsigned(1),
              read=True),
    Datapoint(function_group=50, datapoint=37608, function_name="voc_abluft", datatype=datatypes.Unsigned(1)),
    Datapoint(function_group=50, datapoint=37611, function_name="voc_aussenluft", datatype=datatypes.Unsigned(1)),
    Datapoint(function_group=50, datapoint=39600, function_name="luftqualitaet_regulierung",
              datatype=datatypes.Unsigned(1)),
    Datapoint(function_group=50, datapoint=39652, function_name="status_lueftungsregulierung",
              datatype=datatypes.Unsigned(1)),
    Datapoint(function_group=50, datapoint=0, function_name="temperature_outside_air", datatype=datatypes.Signed(1),
              read=True),
    Datapoint(function_group=50, datapoint=37602, function_name="temperature_exhaust_air", datatype=datatypes.Signed(1),
              read=True),
    Datapoint(function_group=50, datapoint=38600, function_name="ventilator_fortluft_soll", datatype=datatypes.List()),

    # todo: Add all datapoints!
}
