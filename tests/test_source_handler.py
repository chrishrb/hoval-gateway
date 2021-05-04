"""Tests for source_hanlder.py"""

from os import path
import unittest

import serial
from can import Message

from new_gateway.source_handler import CandumpHandler, InvalidFrame, CanHandler

TEST_DATA_DIR = path.abspath(path.join(path.dirname(__file__), 'data'))


class SerialHandlerTestCase(unittest.TestCase):
    """Test case for source_handler.SerialHandler."""

    def setUp(self):
        self.serial_handler = CanHandler("LOOP FOR TESTS")  # device_name will not be used
        self.serial_handler.can0 = serial.serial_for_url('loop://')

    def tearDown(self):
        self.serial_handler.close()

    def test_get_message(self):
        normal_frame = b"FRAME:ID=246:LEN=8:8E:62:1C:F6:1E:63:63:20\n"
        self.serial_handler.serial_device.write(normal_frame)
        frame_id, data = self.serial_handler.get_message()
        self.assertEqual(246, frame_id)
        self.assertEqual(b'\x8e\x62\x1c\xf6\x1e\x63\x63\x20', data)

        no_data_frame = b"FRAME:ID=246:LEN=0:\n"
        self.serial_handler.serial_device.write(no_data_frame)
        frame_id, data = self.serial_handler.get_message()
        self.assertEqual(246, frame_id)
        self.assertEqual(b'', data)

        wrong_length_frame = b"FRAME:ID=246:LEN=9:00:01:02:03:04:05:06:07\n"
        self.serial_handler.serial_device.write(wrong_length_frame)
        self.assertRaises(InvalidFrame, self.serial_handler.get_message)

        three_digit_data_frame = b"FRAME:ID=246:LEN=1:012\n"
        self.serial_handler.serial_device.write(three_digit_data_frame)
        self.assertRaises(InvalidFrame, self.serial_handler.get_message)

        one_digit_data_frame = b"FRAME:ID=246:LEN=1:0\n"
        self.serial_handler.serial_device.write(one_digit_data_frame)
        self.assertRaises(InvalidFrame, self.serial_handler.get_message)

        missing_id_frame = b"FRAME:LEN=1:8E\n"
        self.serial_handler.serial_device.write(missing_id_frame)
        self.assertRaises(InvalidFrame, self.serial_handler.get_message)

        missing_length_frame = b"FRAME:ID=246:8E\n"
        self.serial_handler.serial_device.write(missing_length_frame)
        self.assertRaises(InvalidFrame, self.serial_handler.get_message)


class CandumpHandlerTestCase(unittest.TestCase):
    """Test case for source_handler.CandumpHandler."""

    maxDiff = None

    def setUp(self):
        file_path = path.join(TEST_DATA_DIR, 'test_data.log')
        self.candump_handler = CandumpHandler(file_path)
        self.candump_handler.open()

    def tearDown(self):
        self.candump_handler.close()

    def test_get_message(self):
        messages = [self.candump_handler.get_message() for _ in range(7)]

        expected_messages = [
            Message(timestamp=1499184018.469421, arbitration_id=0x0, is_extended_id=True, dlc=0, data=[]),
            Message(timestamp=1499184018.471357, arbitration_id=0x0, is_extended_id=True, dlc=8,
                    data=[0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0]),
            Message(timestamp=1499184018.472199, arbitration_id=0xfff, is_extended_id=True, dlc=8,
                    data=[0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff]),
            Message(timestamp=1499184018.473352, arbitration_id=0x1, is_extended_id=True, dlc=8,
                    data=[0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1]),
            Message(timestamp=1499184018.474749, arbitration_id=0x100, is_extended_id=True, dlc=8,
                    data=[0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0]),
            Message(timestamp=1499184018.475778, arbitration_id=0xabc, is_extended_id=True, dlc=3,
                    data=[0x12, 0xaf, 0x49]),
            Message(timestamp=1499184018.478492, arbitration_id=0x743, is_extended_id=True, dlc=8,
                    data=[0x9f, 0x20, 0xa1, 0x20, 0x78, 0xbc, 0xea, 0x98])
        ]

        for i in range(0, len(messages)):
            self.assertTrue(messages[i].equals(expected_messages[i]))

        self.assertRaises(InvalidFrame, self.candump_handler.get_message)
        self.assertRaises(EOFError, self.candump_handler.get_message)
