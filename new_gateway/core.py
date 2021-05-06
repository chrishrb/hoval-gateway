import logging

from new_gateway.datapoint import NoDatapointFoundError
from new_gateway.datatypes import UnknownDatatypeError
from new_gateway.message import get_message_header, get_message_len, get_message_id, \
    get_operation_id, Operation, Message

_pending_msg = {}


async def read(can0):
    """Read data from CAN and export to mqtt"""
    can0.open()

    while True:
        try:
            msg = await can0.get_message()
        except EOFError:
            break

        # Add message to pending
        add_to_pending_msg(msg)

        # When no more messages are expected, parse data
        if get_message_header(msg.data) in _pending_msg and \
                _pending_msg[get_message_header(msg.data)].nb_remaining == 0:
            process_message = _pending_msg[get_message_header(msg.data)]
            logging.debug("Current message: %s", process_message)

            try:
                datapoint, operation, parsed_message = process_message.parse()
                logging.debug("Message for datapoint [function_group: %s, function_number: %s, function_datapoint: "
                              "%s, operation: %s] : %s",
                              datapoint.function_group,
                              datapoint.function_number,
                              datapoint.datapoint_id,
                              operation,
                              parsed_message)

                # For mqtt
                if operation == Operation.RESPONSE:
                    pass
            except UnknownDatatypeError as e:
                logging.error(e)
            except NoDatapointFoundError as e:
                logging.warning(e)

            del _pending_msg[get_message_header(msg.data)]


def add_to_pending_msg(msg):
    """Add message to queue"""
    if get_message_id(msg.arbitration_id) == 0x1f:
        parsed_msg = Message(msg.arbitration_id, get_message_len(msg.data), get_operation_id(msg.data),
                             get_message_header(msg.data))
        parsed_msg.put(msg)
        _pending_msg[get_message_header(msg.data)] = parsed_msg
    elif get_message_header(msg.data) in _pending_msg:
        _pending_msg[get_message_header(msg.data)].put_extended_msg(msg)


async def send(can0):
    pass
