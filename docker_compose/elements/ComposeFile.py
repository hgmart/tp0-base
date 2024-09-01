from docker_compose.common.Serializable import Serializable

class ComposeFile(Serializable):
    def __init__(self, name = None, services = None, networks = None, volumes = None):
        self.name = name
        self.services = services
        self.networks = networks
        self.volumes = volumes
    
    def serialize(self, indentation = 0):
        return super().serialize(indentation) + '\n'