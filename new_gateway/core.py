async def main(can0):
    can0.open()

    while True:
        msg = await can0.get_message()
        print(msg)
