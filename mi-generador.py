import sys
from docker_compose.common.Container import Container
from docker_compose.elements.ComposeFile import ComposeFile
from docker_compose.elements.Service import Service
from docker_compose.elements.Network import Network

# define las caracterísitcas del servidor
def define_server() -> list[Container]:
 return [Container(
    name='server',
    properties=Service(
        container_name='server',
        image='server:latest',
        entrypoint='python3 /main.py',
        environment=[
        'PYTHONUNBUFFERED=1',
        'LOGGING_LEVEL=DEBUG'
        ],
        networks=['testing_net']
    )
)]

# define las características de cada cliente
def define_clients(clients_number) -> list[Container]:
 return [Container(
        name=f'client{i+1}',
        properties=Service(
        container_name=f'client{i+1}',
        image='client:latest',
        entrypoint='/client',
        environment=[f'CLI_ID={i+1}', 'CLI_LOG_LEVEL=DEBUG'],
        networks=['testing_net'],
        depends_on=['server']
    )) for i in range(clients_number)]

# define las características de la red empleada
def define_network() -> Container:
 return Container(
            name='networks',
            properties=Container(
            name='testing_net',
            properties=Container(
                name='ipam',
                properties=Network(
                        driver= 'default',
                        config= ['subnet: 172.25.125.0/24']
                    )
                )
            )
        )

# punto de inicio de la aplicación
def main(args):
    clients_range=int(args[1])

    composeFile = ComposeFile(
        name = 'tp0',
        services = Container(
            name = 'services',
            list = define_server() + define_clients(clients_range),
        ),
        networks = define_network()
    )

    with open(args[0], 'w') as file:
        file.write(composeFile.serialize())

main(sys.argv[1:])




