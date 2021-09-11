"""Tests for source_hanlder.py"""

import unittest
from os import path
from unittest import IsolatedAsyncioTestCase

from can.message import Message

from gateway.source_handler import CandumpHandler, InvalidFrame

TEST_DATA_DIR = path.abspath(path.join(path.dirname(__file__), 'data'))


class SerialHandlerTestCase(unittest.TestCase):
    # todo: add vcan test scenario
    """Test case for source_handler.SerialHandler."""

    def setUp(self):
        # todo: add test
        pass

    def tearDown(self):
        # todo: add test
        pass

    async def test_get_message(self):
        # todo: add test
        pass


class CandumpHandlerTestCase(IsolatedAsyncioTestCase):
    """Test case for source_handler.CandumpHandler."""

    maxDiff = None

    def setUp(self):
        file_path = path.join(TEST_DATA_DIR, 'test_data.log')
        self.candump_handler = CandumpHandler(file_path)

    async def test_get_message(self):
        messages = [await self.candump_handler.get_message() for _ in range(7)]

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

        with self.assertRaises(InvalidFrame):
            await self.candump_handler.get_message()

        with self.assertRaises(EOFError):
            await self.candump_handler.get_message()
