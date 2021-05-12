from unittest import TestCase
import new_gateway.message as message


class TestMessage(TestCase):
    # hello world in bytecode
    string = b'\x68\x65\x6c\x6c\x6f\x20\x77\x6f\x72\x6c\x64'

    def test_put(self):
        msg = message.Message(1, message.Operation.GET_REQUEST, 0)
        msg.put(self.string)

        self.assertGreater(len(self.string), 2)
        self.assertEqual(msg.data, self.string)
        self.assertEqual(msg.nb_remaining, 0)
        self.assertEqual(msg.message_len, 1)

    def test_put_extended_msg(self):
        msg = message.Message(1, message.Operation.GET_REQUEST, 0)
        msg.put_extended_msg(self.string)

        self.assertGreater(len(self.string), 2)
        self.assertEqual(msg.data, b'\x68\x65\x6c\x6c\x6f\x20\x77\x6f\x72')
        self.assertEqual(msg.nb_remaining, 0)
        self.assertEqual(msg.message_len, 1)

    def test_parse_too_short(self):
        # todo: add test
        pass

    def test_parse(self):
        # todo: add test
        pass
