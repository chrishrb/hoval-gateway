import setuptools

setuptools.setup(name='hoval-gateway',
                 description="Hoval Gateway - read/write to hoval CAN bus",
                 version='0.0.1',
                 author='Christoph Herb',
                 url='https://github.com/chrishrb/hoval-gateway',
                 packages=setuptools.find_packages(exclude=('tests', 'docs')))
