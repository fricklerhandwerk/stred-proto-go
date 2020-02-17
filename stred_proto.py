from pathlib import Path
from enum import Enum
from typing import List, Union, Mapping


class Range():
    low: int
    high: int


class Enumeration():
    label: str
    # with aliasing we may have multiple names for an enumeration item
    enumeration: Mapping[int, List[str]]
    allow_alias: bool
    reserved: List[Union[int, Range]]
    reserved_names: List[str]

    def __init__(self, label="", enumeration={}, allow_alias=False, reserved=[], reserved_names=[]):
        self.label = label
        self.enumeration = enumeration
        self.allow_alias = allow_alias
        self.reserved = reserved
        self.reserved_names = reserved_names

    def __str__(self):
        fields = indent("\n".join([f"{x} = {k};" for k, v in self.enumeration.items() for x in v]))
        return f"enum {self.label} {{{fields}}}"
        return


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


Type = Union[KeyType, ValueType, "Message", Enumeration]


class BaseField():
    label: str
    type: Type
    deprecated: bool

    def __init__(self, label="", type=None, deprecated=False):
        self.label = label
        self.type = type
        self.deprecated = deprecated


class Field(BaseField):
    repeated: bool

    def __init__(self, repeated=False, *args, **kwargs):
        self.repeated = repeated
        super().__init__(*args, **kwargs)

    def __str__(self):
        return f"{self.type.value} {self.label}"


class Oneof():
    label: str
    fields: Mapping[int, BaseField]


class MapField(BaseField):
    key: KeyType


class Message():
    label: str
    fields: Mapping[int, Union[Field, Oneof, MapField]]

    reserved: List[Union[int, Range]]
    reserved_names: List[str]

    definitions: List[Union[Enumeration, "Message"]]

    def __init__(self, label="", fields={}, reserved=[], reserved_names=[], definitions=[]):
        self.label = label
        self.fields = fields
        self.reserved = reserved
        self.reserved_names = reserved_names
        self.definitions = definitions

    def __str__(self):
        def field(k, v):
            return f"{v} = {k}{' [deprecated=true]' if v.deprecated else ''};"
        fields = indent("\n".join([field(k, v) for k, v in self.fields.items()]))
        definitions = indent("\n\n".join([str(x) for x in self.definitions]))
        return f"message {self.label} {{{fields}{definitions}}}"


def indent(x: str, level: int = 1, indent: str = "  ") -> str:
    if x == "":
        return ""
    return "\n" + indent * level + ("\n" + indent * level).join(x.split("\n")) + "\n"


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


class Proto():
    package: str
    imports: List[Import]
    services: List[Service]
    definitions: List[Union[Message, Enumeration]]

    def __init__(self, package="", imports=[], services=[], definitions=[]):
        self.package = package
        self.imports = imports
        self.services = services
        self.definitions = definitions

    def __str__(self):
        return "\n\n".join([f"package = {self.package};"] + [str(x) for x in self.definitions])


test_fields = {
    1: Field(label="broogle", type=KeyType.INT32),
    5: Field(label="doingle", type=KeyType.UINT64, deprecated=True),
}

test_fields2 = {
    2: Field(label="foo", type=KeyType.STRING),
    3: Field(label="bar", type=ValueType.BYTES),
}

test_message = Message(label="MyMessage", fields=test_fields)

test_enumeration = {
    0: ["default"],
    2: ["some"],
    3: ["thing", "redundant"],
}

test_enum = Enumeration(label="MyEnum", enumeration=test_enumeration)

test_message2 = Message(label="SomeOtherMessage", fields=test_fields2)
test_message2.definitions = [test_enum, test_message]


test_proto = Proto(package="testpackage", definitions=[test_message, test_enum, test_message2])


def main():
    print(test_proto)

