from docker_compose.common.Serializable import Serializable

class Container(Serializable):
    def __init__(self, name, properties = None, list = []):
        self.name=name
        self.properties=properties
        self.list=list

    def serialize(self, indentation = 0):
        
        serialization = ''

        if ( self.name != None):
            serialization += Serializable.add_indentation(indentation, f'{self.name}:\n')

        if (self.properties != None):
            serialization +=self.properties.serialize(indentation+1)

        serialization += '\n'.join(f'{element.serialize(indentation+1)}\n' for element in self.list)

        return serialization