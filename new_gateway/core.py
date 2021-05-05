import logging


async def read(can0):
    can0.open()

    while True:
        try:
            msg = await can0.get_message()
            logging.debug("Raw message %s", msg)
        except EOFError:
            break


async def send(can0):
    pass
