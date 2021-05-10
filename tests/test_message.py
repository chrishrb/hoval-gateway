from unittest import TestCase
from unittest.mock import Mock, MagicMock

import can

from new_gateway.message import Message


class TestMessage(TestCase):
    def test_is_valid(self):
        parsed_message = Message(123, 0, 0x42, 1)
        msg = MagicMock(can.Message)
        msg.data = bytearray(b'2\x42\x00\x92\xe2\x00\xd2')
        self.assertTrue(parsed_message._is_valid(msg))

