import pytest

from stred_proto import (Enumeration, Field, KeyType, Map, Message, OneOf,
                         Protocol, RepeatableField, ReservedLabels, TypedField,
                         ValidationError, ValueType)


def test_string_representation():
    message = Message()
    message.label = "MyMessage"

    message_field1 = RepeatableField()
    message_field1.label = "broogle"
    message_field1.type = KeyType.INT32
    message_field1.number = 1
    message.fields.append(message_field1)

    message_field2 = RepeatableField()
    message_field2.label = "doingle"
    message_field2.type = KeyType.UINT64
    message_field2.number = 1
    message_field2.deprecated = True
    message.fields.append(message_field2)

    enum = Enumeration()
    enum.label = "MyEnum"

    enum_field1 = Field()
    enum_field1.label = "default"
    enum_field1.number = 0
    enum.fields.append(enum_field1)

    enum_field2 = Field()
    enum_field2.label = "some"
    enum_field2.number = 1
    enum.fields.append(enum_field2)

    enum_field3 = Field()
    enum_field3.label = "thing"
    enum_field3.number = 2
    enum.fields.append(enum_field3)

    enum_field4 = Field()
    enum_field4.label = "redundant"
    enum_field4.number = 2
    enum.fields.append(enum_field3)

    message2 = Message()
    message2.label = "SomeOtherMessage"

    message2_field1 = RepeatableField()
    message2_field1.label = "foo"
    message2_field1.type = KeyType.STRING
    message2_field1.number = 2
    message2.fields.append(message2_field1)

    message2_field2 = RepeatableField()
    message2_field2.label = "bar"
    message2_field2.type = ValueType.BYTES
    message2_field2.number = 3
    message2_field2.repeated = True
    message2.fields.append(message2_field2)

    message2_map = Map()
    message2_map.label = "some_map"
    message2_map.key_type = KeyType.INT32
    message2_map.value_type = message
    message2_map.number = 4
    message2.fields.append(message2_map)

    message2_oneof = OneOf()
    message2_oneof.label = "my_oneof"

    oneof_field1 = TypedField()
    oneof_field1.label = "floah_dude"
    oneof_field1.type = ValueType.FLOAT
    oneof_field1.number = 5
    message2_oneof.fields.append(oneof_field1)

    oneof_field2 = TypedField()
    oneof_field2.label = "double_trouble"
    oneof_field2.type = ValueType.DOUBLE
    oneof_field2.number = 6
    message2_oneof.fields.append(oneof_field2)

    message2.definitions.extend([enum, message])

    test_proto = Protocol()
    test_proto.package = "testpackage"
    test_proto.definitions.extend([message, enum, message2])

    # # TODO: actually try to run this through `protoc` to check for validity
    print(test_proto)


def test_message_new_invalid_identifier():
    m = Message()
    with pytest.raises(ValidationError):
        m.label = "InvalidIdent!"
    e = Enumeration()
    with pytest.raises(ValidationError):
        e.label = "InvalidIdent!"
    f = TypedField()
    with pytest.raises(ValidationError):
        f.type = KeyType.BOOL
        f.number = 1
        f.label = "InvalidIdent!"


def test_message_set_invalid_identifier():
    m = Message()
    with pytest.raises(ValidationError):
        m.label = "1InvalidIdent"
        m.label = "1InvalidIdent"


def test_enum_fields_mutable_sequence():
    e = Enumeration()
    e.label = "ValidLabel"
    f = Field()
    f.number = 1
    f.label = "some_field"
    r = ReservedLabels()
    r.append("woah")
    e.fields.append(f)
    e.fields.append(r)
    e.fields.append(r)


def test_enum_fields_append_invalid():
    e = Enumeration()
    f = Field()
    with pytest.raises(ValidationError):
        e.fields.append(f)
