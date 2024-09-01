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
        networks=['testing_net'],
        volumes=['server_configs:/configs']
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

# valida que el primer parámetro sea un string y el segundo un número positivo
def define_volumes():
    return Container(
       list=['server_configs']
    )

# valida que el primer parámetro sea un string y el segundo un número positivo
def validate_parameters(args):
    try:
        return len(args) == 2 and isinstance(args[0], str) and isinstance(args[1], str) and int(args[1]) and int(args[1]) > 0
    
    except:
        return False   

# punto de inicio de la aplicación
def main(args):

    if (validate_parameters(args)):

        file_name = args[0]
        clients_range = int(args[1])

        composeFile = ComposeFile(
            name = 'tp0',
            services = Container(
                name = 'services',
                list = define_server() + define_clients(clients_range),
            ),
            networks = define_network(),
            volumes = define_volumes(),
        )

        with open(file_name, 'w') as file:
            file.write(composeFile.serialize())

    else:
        print('Se espera como primer parámetro un nombre de archivo y como segundo parámetro un número de clientes.', file=sys.stderr)
        sys.exit(1)

main(sys.argv[1:])
