class Serializable:
    
    def create_indentation(indentation = 0):
        indentation_character = ' '
        spaces = indentation_character * indentation * 2

        return spaces

    @staticmethod
    def add_indentation(spaces = 0, value = ''):
        return f'{Serializable.create_indentation(spaces)}{value}'

    @staticmethod
    def serialize_element(name, value, indentation):
        spaces = Serializable.create_indentation(indentation)

        if (value is None):
            return ''

        # Da formato a las propiedades
        if (isinstance(value, str)):
            return f'{spaces}{name}: {value}'
        
        # Da formato especial a los elementos de una lista
        if (isinstance(value, list)):
            tab = Serializable.create_indentation(1)
            items = [f'{spaces}{tab}- {item}' for item in value]
            items.insert(0, f'{spaces}{name}:')
            return '\n'.join(items)
        
        if(isinstance(value, Serializable)):
            return value.serialize(indentation)
        
    # Serializa cada una de las propiedades ignorando aquellas que no tengan valor
    def serialize(self, indentation = 0):
        outputs=[Serializable.serialize_element(name, value, indentation) for name, value in vars(self).items() if value is not None]
        return '\n'.join(outputs)