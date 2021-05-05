import asyncio
import logging

import click

from new_gateway.core import main
from new_gateway.source_handler import CanHandler, CandumpHandler


@click.command()
@click.option('-v', '--verbose', is_flag=True, help="Debug output")
@click.option('-f', '--file', type=click.Path(resolve_path=True), help="Read can messages from file")
def run(verbose, file):
    """
    Run main application with can interface
    """
    loop = asyncio.get_event_loop()

    # Choose right handler
    if file is None:
        can0 = CanHandler("can0")
    else:
        can0 = CandumpHandler(file)

    # Verbose output
    if verbose is True:
        logging.basicConfig(level=logging.DEBUG)

    try:
        loop.run_until_complete(main(can0))
    except KeyboardInterrupt:
        pass
    except EOFError:
        pass
    finally:
        can0.close()
        loop.close()


if __name__ == "__main__":
    run()
