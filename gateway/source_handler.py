import asyncio
import logging
import re

from can.interface import Bus
from can.listener import AsyncBufferedReader
from can.message import Message
from can.notifier import Notifier

from gateway.exceptions import InvalidFrame


class SourceHandler:
    """Base class for classes reading CAN messages.
    This serves as a kind of interface for all classes reading CAN messages,
    whatever the source of these messages: serial port, text file etc.
    """

    def open(self):
        raise NotImplementedError

    def close(self):
        raise NotImplementedError

    def get_message(self):
        """Get CAN id and CAN data.
        Returns:
            A tuple containing the id (int) and data (bytes)
        Raises:
            InvalidFrame
        """
        raise NotImplementedError

    def send(self, message):
        """Send can message"""
        raise NotImplementedError

    def send_periodic(self, message, time):
        """Send periodic can messages"""
        raise NotImplementedError


class CanHandler(SourceHandler):
    """Handler for CAN interface"""

    def __init__(self, device_name, bus_type="socketcan"):
        self._device_name = device_name
        self._bus_type = bus_type
        self.can0, self.reader, self.notifier = self.open()

    def open(self):
        logging.debug("Open CAN Interface..")
        can0 = Bus(channel=self._device_name, bustype=self._bus_type, receive_own_messages=False)

        loop = asyncio.get_event_loop()
        reader = AsyncBufferedReader(loop)
        listeners = [reader]
        notifier = Notifier(can0, listeners, loop=loop)

        logging.debug("CAN Interface opened")
        return can0, reader, notifier

    def close(self):
        logging.debug("Close CAN Interface..")
        if self.notifier:
            self.notifier.stop()
        if self.can0:
            self.can0.stop_all_periodic_tasks()
            self.can0.shutdown()
        logging.debug("CAN Interface closed")

    def get_message(self):
        message = self.reader.get_message()
        return message

    def send(self, message):
        self.can0.send(message)

    def send_periodic(self, message, time):
        self.can0.send_periodic(message, time)


class CandumpHandler(SourceHandler):
    """Parser for text files generated by candump."""

    MSG_RE = r"\(([.0-9]+)\).* ([0-9A-F]+)\#([0-9A-F]*)"
    MSG_RGX = re.compile(MSG_RE)

    def __init__(self, file_path):
        self.file_path = file_path
        self.file_object = self.open()

    def open(self):
        logging.debug("Open fake CANDUMP Interface (file)..")
        return open(self.file_path, 'rt', encoding='utf-8')

    def close(self):
        logging.debug("Close fake CANDUMP Interface (file)..")
        if self.file_object:
            self.file_object.close()
        logging.debug("Fake CANDUMP Interface (file) closed")

    async def get_message(self):
        line = self.file_object.readline()
        if line == '':
            raise EOFError
        return self._parse_from_candump(line)

    def _parse_from_candump(self, line):
        line = line.strip('\n')

        msg_match = self.MSG_RGX.match(line)
        if msg_match is None:
            raise InvalidFrame("Wrong format: '{}'".format(line))

        can_time, hex_can_id, hex_can_data = msg_match.group(1, 2, 3)
        can_id = int(hex_can_id, 16)

        try:
            can_data = bytes.fromhex(hex_can_data)
        except ValueError as err:
            raise InvalidFrame("Can't decode message '{}': '{}'".format(line, err))

        return Message(timestamp=float(can_time), arbitration_id=can_id, data=can_data)

    def send(self, message):
        print(message)

    def send_periodic(self, message, time):
        raise NotImplementedError
