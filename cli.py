import asyncio
import logging

import click

from new_gateway.core import read
from new_gateway.source_handler import CanHandler, CandumpHandler


@click.command()
@click.option('-v', '--verbose', is_flag=True, help="Debug output")
@click.option('-f', '--file', type=click.Path(resolve_path=True), help="Read can messages from file")
def run(verbose, file):
    """
    Run main application with can interface
    """
    # Choose right handler
    if file is None:
        can0 = CanHandler("can0")
    else:
        can0 = CandumpHandler(file)

    # Verbose output
    if verbose is True:
        logging.basicConfig(level=logging.DEBUG)

    loop = asyncio.get_event_loop()
    try:
        asyncio.ensure_future(read(can0))
        loop.run_forever()
    except KeyboardInterrupt:
        pass
    finally:
        can0.close()


if __name__ == "__main__":
    run()
