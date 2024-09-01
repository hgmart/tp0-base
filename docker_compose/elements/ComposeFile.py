from docker_compose.common.Serializable import Serializable

class ComposeFile(Serializable):
    def __init__(self, name, services = None, networks = None):
        self.name = name
        self.services = services
        self.networks = networks
    
    def serialize(self, indentation = 0):
        return super().serialize(indentation) + '\n'