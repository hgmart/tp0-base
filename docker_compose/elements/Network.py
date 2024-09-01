from docker_compose.common.Serializable import Serializable

class Network(Serializable):
    def __init__(self, driver = None, config = None):
        self.driver = driver
        self.config = config