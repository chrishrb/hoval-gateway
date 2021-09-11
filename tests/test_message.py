from unittest import TestCase, IsolatedAsyncioTestCase

from can import Message

import gateway.message as message
from gateway.datapoint import Datapoint
from gateway.datatypes import List
from gateway.source_handler import CanHandler


class TestMessage(TestCase):
    def test_arbitration_id(self):
        arbitration_id = message.build_arbitration_id(31, 224, 71, 255)
        self.assertEqual(31, message.get_message_id(arbitration_id))
        self.assertEqual(224, message.get_message_priority(arbitration_id))
        self.assertEqual(71, message.get_message_device_type(arbitration_id))
        self.assertEqual(255, message.get_message_device_id(arbitration_id))


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


class TestSendMessage(IsolatedAsyncioTestCase):
    def setUp(self):
        arbitration_id = message.build_arbitration_id(31, 224, 10, 8)
        datapoint = Datapoint(name="example_datapoint", function_group=50, function_number=0, datapoint_id=40651,
                              datatype="U16", decimal=1)
        self.msg = message.SendMessage(arbitration_id, message.Operation.SET_REQUEST.value, datapoint)
        self.can0 = CanHandler('vcan', 'virtual')
        self.can1 = CanHandler('vcan', 'virtual')

    async def test_put_data(self):
        # given
        datatype = List()
        self.msg.put_data(datatype.convert_to_bytes("37"))
        can_msg = self.msg.to_can_message()
        self.can0.send(can_msg)

        # when
        receive_message: Message = await self.can1.get_message()
        msg_rcv = message.ReceiveMessage(receive_message.arbitration_id, message.get_operation_id(receive_message.data),
                                         message.get_message_len(receive_message.data))
        msg_rcv.put_data(receive_message.data)

        # then
        self.assertEqual(31, msg_rcv.message_id)
        self.assertEqual(224, msg_rcv.priority)
        self.assertEqual(10, msg_rcv.device_type)
        self.assertEqual(8, msg_rcv.device_id)
        self.assertEqual(message.Operation.SET_REQUEST.value, msg_rcv.operation_id)
        self.assertEqual(0, msg_rcv.nb_remaining)
        self.assertEqual(37, datatype.convert_from_bytes(msg_rcv.data[6:]))
