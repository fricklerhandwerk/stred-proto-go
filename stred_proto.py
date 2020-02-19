import re
from enum import Enum
from pathlib import Path
from typing import List, Optional, Union


class ValidationError(Exception):
    pass


class Identifier(str):
    identifier = re.compile(r"^[A-Za-z][0-9A-Za-z_]*$")

    def __new__(cls, value: str):
        if not cls.identifier.match(value):
            raise ValidationError(f'Identifier must match {cls.identifier.pattern}')
        return super().__new__(cls, value)


class Declaration():
    _label: Identifier

    def __init__(self, label: Union[Identifier, str], *args, **kwargs):
        if isinstance(label, Identifier):
            self._label = label
        else:
            self._label = Identifier(label)
        super().__init__(*args, **kwargs)

    @property
    def label(self) -> Identifier:
        return self._label

    @label.setter
    def label(self, value: Union[Identifier, str]):
        if isinstance(value, Identifier):
            self._label = value
        else:
            self._label = Identifier(value)


class Definition(Declaration):
    pass


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


class Reservation():
    pass


class Range():
    start: int
    end: int


class ReservedNumbers(Reservation):
    numbers: List[Union[int, Range]]


class ReservedLabel(Reservation):
    labels: List[Identifier]


class Field(Declaration):
    number: int
    deprecated: bool

    def __init__(self, number: int, *args, **kwargs):
        self.number = number
        self.deprecated = False
        super().__init__(*args, **kwargs)

    def __str__(self):
        deprecated = " [deprecated=true]" if self.deprecated else ""
        return f"{self.label} = {self.number}{deprecated};"


class TypedField(Field):
    type: Type

    def __init__(self, field_type: Type, *args, **kwargs):
        self.type = field_type
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

    def __init__(self, key_type: KeyType, value_type: Type, *args, **kwargs):
        self.key_type = key_type
        self.value_type = value_type
        super().__init__(*args, **kwargs)

    def __str__(self):
        key_type = self.key_type.value

        def value_type():
            if isinstance(self.value_type, Definition):
                return self.value_type.label
            return self.value_type.value

        return f"map<{key_type}, {value_type()}> {super().__str__()}"


class OneOf(Declaration):
    fields: List[TypedField]

    def __str__(self):
        fields = "\n".join([str(f) for f in self.fields])
        return f"oneof {self.label} {{{indent(fields)}}}"


class Enumeration(Definition):
    fields: List[Union[Field, Reservation]]
    allow_alias: bool

    def __init__(self, *args, **kwargs):
        self.fields = []
        self.allow_alias = False
        super().__init__(*args, **kwargs)

    def __str__(self):
        fields = "\n".join([str(f) for f in self.fields])
        return f"enum {self.label} {{{indent(fields)}}}"


class Message(Definition, Container):
    fields: List[Union[RepeatableField, Map, OneOf, Reservation]]

    def __init__(self, *args, **kwargs):
        self.fields = []
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
