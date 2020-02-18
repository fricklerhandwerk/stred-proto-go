import re
from enum import Enum
from pathlib import Path
from typing import List, Mapping, Optional, Union


class ValidationError(Exception):
    pass


class Range():
    low: int
    high: int


class Item():
    _label: str

    def __init__(self, label, *args, **kwargs):
        identifier = re.compile("^([A-Z]|[a-z])([0-9]|[A-Z]|[a-z]|_])*$")

        if identifier.match(label):
            self._label = label
        else:
            raise ValidationError(f'Identifier must match {identifier.pattern}')
        super().__init__(*args, **kwargs)

    @property
    def label(self):
        return self._label


class Definition(Item):
    reserved: List[Union[int, Range]]
    reserved_names: List[str]

    def __init__(self, *args, **kwargs):
        self.reserved = []
        self.reserved_names = []
        super().__init__(*args, **kwargs)


class Container():
    definitions: List[Definition]

    def __init__(self, *args, **kwargs):
        self.definitions = []
        super().__init__(*args, **kwargs)


class Enumeration(Definition):
    # with aliasing we may have multiple names for an enumeration item
    enumeration: Mapping[int, List[str]]
    allow_alias: bool

    def __init__(self, *args, **kwargs):
        self.enumeration = {}
        self.allow_alias = False
        super().__init__(*args, **kwargs)

    def __str__(self):
        fields = indent("\n".join([f"{x} = {k};" for k, v in self.enumeration.items() for x in v]))
        return f"enum {self.label} {{{fields}}}"


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


class Field(Item):
    type: Type
    deprecated: bool

    def __init__(self, field_type, *args, **kwargs):
        self.type = field_type
        self.deprecated = False
        super().__init__(*args, **kwargs)

    def __str__(self):
        return f"{self.type.value} {self.label}"


class RepeatableField(Field):
    repeated: bool

    def __init__(self, *args, **kwargs):
        self.repeated = False
        super().__init__(*args, **kwargs)

    def __str__(self):
        repeated = "repeated " if self.repeated else ""
        return f"{repeated}{super().__str__()}"


class OneOf(Item):
    fields: Mapping[int, Field]


class Map(Field):
    key: KeyType


class Message(Definition, Container):
    fields: Mapping[int, Field]

    def __init__(self, *args, **kwargs):
        self.fields = {}
        super().__init__(*args, **kwargs)

    def __str__(self):
        def field(i, f):
            deprecated = " [deprecated=true]" if f.deprecated else ""
            return f"{f} = {i}{deprecated};"
        fields = indent("\n".join([field(i, f) for i, f in self.fields.items()]))
        definitions = indent("\n\n".join([str(x) for x in self.definitions]))
        return f"message {self.label} {{{fields}{definitions}}}"

    @property
    def label(self):
        return self._label

    @label.setter
    def label(self, value):
        self._label = value


def indent(x: str, level: int = 1, prefix: str = "  ") -> str:
    if x == "":
        return ""
    return "\n" + prefix * level + ("\n" + prefix * level).join(x.split("\n")) + "\n"


class RPC():
    label: str
    requestType: Message
    streamRequest: bool
    responseType: Message
    streamResponse: bool


class Service():
    """
    https://developers.google.com/protocol-buffers/docs/reference/proto3-spec#service_definition
    """
    label: str
    rpcs: List[RPC]


class Import():
    path: Path
    public: bool


class Proto(Container):
    package: Optional[str]
    imports: List[Import]
    services: List[Service]
    definitions: List[Definition]

    def __init__(self):
        self.package = None
        self.imports = []
        self.services = []
        self.definitions = []
        super().__init__()

    def __str__(self):
        return "\n\n".join([f"package = {self.package};"] + [str(x) for x in self.definitions])


def main():
    pass
