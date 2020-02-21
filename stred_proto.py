import collections.abc as collections
import re
from abc import ABC, abstractmethod
from enum import Enum
from pathlib import Path
from typing import List, Optional, Union


class ValidationError(Exception):
    pass


class Validator(ABC):
    @abstractmethod
    def validate(self):
        pass


class Identifier(str, Validator):
    def __new__(cls, value: str):
        return super().__new__(cls, value)

    def validate(self):
        identifier = re.compile(r"^[A-Za-z][0-9A-Za-z_]*$")
        if not identifier.match(self):
            raise ValidationError(f'Identifier must match {identifier.pattern}')


class Declaration(Validator):
    _label: Identifier

    def __init__(self):
        self._label = None

    @property
    def label(self) -> Identifier:
        return self._label

    @label.setter
    def label(self, value: Union[Identifier, str]):
        if not isinstance(value, Identifier):
            value = Identifier(value)
        value.validate()
        self._label = value

    def validate(self):
        if self.label is None:
            raise ValidationError("Label is not set")
        self.label.validate()


class Fields(collections.MutableSequence):
    pass


class Definition(Declaration):
    fields: Fields


class Container():
    definitions: List[Definition]

    def __init__(self, *args, **kwargs):
        self.definitions = []
        super().__init__(*args, **kwargs)


class KeyType(Enum):
    INT32 = "int32"
    INT64 = "int64"
    UINT32 = "uint32"
    UINT64 = "uint64"
    SINT32 = "sint32"
    SINT64 = "sint64"
    FIXED32 = "fixed32"
    FIXED64 = "fixed64"
    SFIXED32 = "sfixed32"
    SFIXED64 = "sfixed64"
    BOOL = "bool"
    STRING = "string"


class ValueType(Enum):
    DOUBLE = "double"
    FLOAT = "float"
    BYTES = "bytes"


Type = Union[KeyType, ValueType, Definition]


class Reservation(Validator):
    pass


class Range():
    start: int
    end: int


class ReservedNumbers(Reservation):
    numbers: List[Union[int, Range]]


class ReservedLabels(list, Reservation):
    def __new__(cls):
        return super().__new__(cls, [])

    def __setitem__(self, key, value: Union[str, Identifier]):
        if not isinstance(value, Identifier):
            value = Identifier(value)
        value.validate()
        super().__setitem__(key, value)

    def insert(self, key, value: Union[str, Identifier]):
        if not isinstance(value, Identifier):
            value = Identifier(value)
        value.validate()
        super().insert(key, value)

    def append(self, value: Union[str, Identifier]):
        if not isinstance(value, Identifier):
            value = Identifier(value)
        value.validate()
        super().append(value)

    def validate(self):
        if not self:
            raise ValidationError("Reserved labels must contain at least one label")
        for l in self:
            l.validate()

    def __str__(self):
        labels = ", ".join(f'"{l}"' for l in self)
        return f"reserved {labels};"


class Field(Declaration):
    number: int
    deprecated: bool

    def __init__(self, *args, **kwargs):
        self.number = None
        self.deprecated = False
        super().__init__(*args, **kwargs)

    def __str__(self):
        deprecated = " [deprecated=true]" if self.deprecated else ""
        return f"{self.label} = {self.number}{deprecated};"

    def validate(self):
        super().validate()


class TypedField(Field):
    type: Type

    def __init__(self, *args, **kwargs):
        self.type = None
        super().__init__(*args, **kwargs)

    def __str__(self):
        return f"{self.type.value} {super().__str__()}"


class RepeatableField(TypedField):
    repeated: bool

    def __init__(self, *args, **kwargs):
        self.repeated = False
        super().__init__(*args, **kwargs)

    def __str__(self):
        repeated = "repeated " if self.repeated else ""
        return f"{repeated}{super().__str__()}"


class Map(Field):
    key_type: KeyType
    value_type: Type

    def __init__(self, *args, **kwargs):
        self.key_type = None
        self.value_type = None
        super().__init__(*args, **kwargs)

    def __str__(self):
        key_type = self.key_type.value

        def value_type():
            if isinstance(self.value_type, Definition):
                return self.value_type.label
            return self.value_type.value

        return f"map<{key_type}, {value_type()}> {super().__str__()}"


class TypedFields(list, Fields):
    def __setitem__(self, key, value: Union[TypedField, Reservation]):
        value.validate()
        super().__setitem__(key, value)

    def insert(self, key, value: Union[TypedField, Reservation]):
        value.validate()
        super().insert(key, value)

    def append(self, value: Union[TypedField, Reservation]):
        value.validate()
        super().append(value)


class OneOf(Declaration):
    fields: List[TypedField]

    def __init__(self):
        self.fields = TypedFields()
        super().__init__()

    def __str__(self):
        fields = "\n".join([str(f) for f in self.fields])
        return f"oneof {self.label} {{{indent(fields)}}}"


class Enumeration(Definition):
    allow_alias: bool

    class F(list, Fields):
        def __setitem__(self, key, value: Union[Field, Reservation]):
            value.validate()
            super().__setitem__(key, value)

        def insert(self, key, value: Union[Field, Reservation]):
            value.validate()
            super().insert(key, value)

        def append(self, value: Union[Field, Reservation]):
            value.validate()
            super().append(value)

    def __init__(self, *args, **kwargs):
        self.fields = Enumeration.F()
        self.allow_alias = False
        super().__init__(*args, **kwargs)

    def __str__(self):
        fields = "\n".join([str(f) for f in self.fields])
        return f"enum {self.label} {{{indent(fields)}}}"


class MessageFields(list, Fields):
    def __setitem__(self, key, value: Union[RepeatableField, Map, OneOf, Reservation]):
        value.validate()
        super().__setitem__(key, value)

    def insert(self, key, value: Union[RepeatableField, Map, OneOf, Reservation]):
        value.validate()
        super().insert(key, value)

    def append(self, value: Union[RepeatableField, Map, OneOf, Reservation]):
        value.validate()
        super().append(value)


class Message(Container, Definition):
    def __init__(self, *args, **kwargs):
        self.fields = MessageFields()
        super().__init__(*args, **kwargs)

    def __str__(self):
        fields = "\n".join([str(f) for f in self.fields])
        definitions = "\n\n".join([str(d) for d in self.definitions])
        return f"message {self.label} {{{indent(fields)}{indent(definitions)}}}"


def indent(x: str, level: int = 1, prefix: str = "  ") -> str:
    if x == "":
        return ""
    return "\n" + prefix * level + ("\n" + prefix * level).join(x.split("\n")) + "\n"


class RPC(Declaration):
    requestType: Message
    streamRequest: bool
    responseType: Message
    streamResponse: bool


class Service(Declaration):
    """
    https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#service_definition
    """
    rpcs: List[RPC]


class Import():
    path: Path
    public: bool


class Protocol(Container):
    package: Optional[Identifier]
    imports: List[Import]
    services: List[Service]

    def __init__(self):
        self.package = None
        self.imports = []
        self.services = []
        super().__init__()

    def __str__(self):
        syntax = 'syntax = "proto3";'

        def package():
            if self.package is None:
                return []
            return [f"package {self.package};"]

        return "\n\n".join([syntax] + package() + [str(x) for x in self.definitions])


def main():
    pass
