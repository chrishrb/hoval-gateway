import logging

from gateway import settings
from gateway.message import Message, Response


class Parser:
    _pending_msg = {}
    _devices = {}

    def parse(self, msg):
        if len(msg.data) < 2:
            logging.error("Message too small")
            return None

        message = Message(msg.arbitration_id, msg.data[0], msg.data[1])
        if message.operation_id in settings.OPERATIONS:
            if message.message_id == 0x1f:
                if message.message_len == 0:
                    message.put(Response(msg.data))
                    return message.parse_data()
                else:
                    self._pending_msg[message.operation_id] = {
                        message
                    }

        return None
