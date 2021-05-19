from unittest import TestCase
import gateway.message as message


class TestReceiveMessage(TestCase):
    # hello world in bytecode
    string = b'\x68\x65\x6c\x6c\x6f\x20\x77\x6f\x72\x6c\x64'

    def test_put(self):
        msg = message.ReceiveMessage(1, message.Operation.RESPONSE.value, 0)
        msg.put_data(self.string)

        self.assertGreater(len(self.string), 2)
        self.assertEqual(self.string, msg.data)
        self.assertEqual(0, msg.nb_remaining)
        self.assertEqual(0, msg.message_len)

    def test_put_extended_msg(self):
        msg = message.ReceiveMessage(1, message.Operation.RESPONSE.value, 0)
        msg.put_extended_data(self.string)

        self.assertGreater(len(self.string), 2)
        self.assertEqual(b'\x68\x65\x6c\x6c\x6f\x20\x77\x6f\x72', msg.data)
        self.assertEqual(0, msg.nb_remaining)
        self.assertEqual(0, msg.message_len)

    def test_parse_too_short(self):
        # todo: add test
        pass

    def test_parse(self):
        # todo: add test
        pass
