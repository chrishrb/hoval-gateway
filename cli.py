import asyncio
import logging

import click
import yaml

from new_gateway import datapoint, request
from new_gateway.core import read, send_periodic
from new_gateway.source_handler import CanHandler, CandumpHandler


def parse_settings(settings_file):
    """Parse settings file"""
    settings = yaml.full_load(settings_file)

    for item, element in settings.items():
        if item == "datapoints":
            datapoint.parse_datapoints(element)
        if item == "periodic_requests":
            request.parse_periodic_requests(element)


@click.command()
@click.option('-v', '--verbose', is_flag=True, help="Debug output")
@click.option('-f', '--file', type=click.Path(resolve_path=True), help="Read can messages from file")
@click.option('-s', '--settings', required=True, type=click.File(), help="Read settings file")
def run(verbose, file, settings):
    """
    Run main application with can interface
    """
    # Settings file
    parse_settings(settings)

    # Choose right handler
    if file is None:
        can0 = CanHandler("can0")
    else:
        can0 = CandumpHandler(file)

    # Verbose output
    if verbose is True:
        logging.basicConfig(level=logging.DEBUG)
    else:
        logging.basicConfig(level=logging.INFO)

    # Start can loop
    loop = asyncio.get_event_loop()
    try:
        asyncio.ensure_future(read(can0))
        if file is None:
            asyncio.ensure_future(send_periodic(can0))
        loop.run_forever()
    except KeyboardInterrupt:
        pass
    finally:
        can0.close()


if __name__ == "__main__":
    run()
