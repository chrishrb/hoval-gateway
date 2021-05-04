import asyncio

from new_gateway.source_handler import CanHandler


async def main():
    can0 = CanHandler("can0")
    can0.open()
    while True:
        try:
            msg = can0.get_message()
            print(msg)
        except KeyboardInterrupt:
            can0.close()


if __name__ == '__main__':
    loop = asyncio.get_event_loop()
    loop.run_until_complete(main())
    loop.close()
