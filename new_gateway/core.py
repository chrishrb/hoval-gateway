import logging

from new_gateway.datapoint import NoDatapointFoundError
from new_gateway.message import MessageResponse, get_message_header, get_message_len, get_message_id, \
    get_operation_id, NoMessageError, Operation

_pending_msg = {}


async def read(can0):
    can0.open()

    while True:
        try:
            msg = await can0.get_message()
            logging.debug("Raw message %s", msg)
        except EOFError:
            break

        # Add message to pending
        add_to_pending_msg(msg)

        # When no more messages are expected, parse data
        if get_message_header(msg.data) in _pending_msg and \
                _pending_msg[get_message_header(msg.data)].nb_remaining == 0:
            process_message = _pending_msg[get_message_header(msg.data)]

            try:
                datapoint, operation, parsed_message = process_message.parse()
                logging.info("Message for datapoint [function_group: %s, function_number: %s, function_datapoint: "
                             "%s, operation: %s] : %s",
                             datapoint.function_group,
                             datapoint.function_number,
                             datapoint.datapoint_id,
                             operation,
                             parsed_message)

                # For mqtt
                if operation == Operation.RESPONSE:
                    pass
            except NoMessageError:
                logging.debug("No message found")
            except NoDatapointFoundError:
                logging.debug("Not datapoint in settings.yaml found")

            del _pending_msg[get_message_header(msg.data)]


async def send(can0):
    pass


def add_to_pending_msg(msg):
    if get_message_id(msg.arbitration_id) == 0x1f:
        parsed_msg = MessageResponse(msg.arbitration_id, get_message_len(msg.data), get_operation_id(msg.data),
                                     get_message_header(msg.data))
        parsed_msg.put(msg)
        _pending_msg[get_message_header(msg.data)] = parsed_msg
        logging.debug("Add new message %s", msg)
    elif get_message_header(msg.data) in _pending_msg:
        _pending_msg[get_message_header(msg.data)].put_extended_msg(msg)
        logging.debug("Add message to existent message %s", msg)
