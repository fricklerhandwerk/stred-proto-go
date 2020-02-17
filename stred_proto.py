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
        fields = [f"{x} = {k};" for k, v in self.enumeration.items() for x in v]
        return "enum {self.label} {{\n  {fields}\n}}".format(self=self, fields='\n  '.join(fields))
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

    enumerations: List[Enumeration]
    messages: List["Message"]

    def __init__(self, label="", fields={}, reserved=[], reserved_names=[]):
        self.label = label
        self.fields = fields
        self.reserved = reserved
        self.reserved_names = reserved_names

    def __str__(self):
        def field(k, v):
            return f"{v} = {k}{' [deprecated=true]' if v.deprecated else ''};"
        fields = [field(k, v) for k, v in self.fields.items()]
        return "message {self.label} {{\n  {fields}\n}}".format(self=self, fields="\n  ".join(fields))


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
    imports: List[Import]
    package: str
    services: List[Service]
    definitions: List[Union[Message, Enumeration]]


test_fields = {
    1: Field(label="broogle", type=KeyType.INT32),
    5: Field(label="doingle", type=KeyType.UINT64, deprecated=True),
}

test_message = Message(label="MyMessage", fields=test_fields)

test_enumeration = {
    0: ["default"],
    2: ["some"],
    3: ["thing", "redundant"],
}

test_enum = Enumeration(label="MyEnum", enumeration=test_enumeration)


def main():
    print(test_message)
    print(test_enum)

