from docker_compose.common.Serializable import Serializable

class Service(Serializable):
    def __init__(self, container_name, image, entrypoint, environment = [], networks = [], depends_on = None):
        self.container_name = container_name
        self.image = image
        self.entrypoint = entrypoint
        self.environment = environment
        self.networks = networks
        self.depends_on = depends_on
