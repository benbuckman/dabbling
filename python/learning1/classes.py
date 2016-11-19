class Car():
    def __init__(self, make, model):
        self.make = make
        self.model = model
        self.type = ''

    def name(self):
        if len(self.type):
            return "%s %s (%s)" % (self.make, self.model, self.type)
        else:
            return "%s %s" % (self.make, self.model)

    def print_name(self):
        print(self.name())

cooper = Car('Mini', 'Cooper')
print(cooper.print_name())


class Sedan(Car):
    def __init__(self, make, model):
        super().__init__(make, model)
        self.type = 'Sedan'

corolla = Sedan('Toyota', 'Corolla')
print(corolla.print_name())


class Foo:
    pass  # filler

