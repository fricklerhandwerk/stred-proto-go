import pytest

from stred_proto import (Enumeration, Field, Identifier, KeyType, Map, Message,
                         OneOf, Protocol, TypedField, ValidationError,
                         ValueType)


def test_instantiation():
    test_message = Message("MyMessage")
    test_message.fields = [
        TypedField(KeyType.INT32, 1, "broogle"),
        TypedField(KeyType.UINT64, 5, "doingle"),
    ]
    test_message.fields[1].deprecated = True

    test_enum = Enumeration("MyEnum")
    test_enum.fields = [
        Field(0, "default"),
        Field(1, "some"),
        Field(2, "thing"),
        Field(2, "redundant"),
    ]

    test_message2 = Message("SomeOtherMessage")
    test_message2.fields = [
        TypedField(KeyType.STRING, 2, "foo"),
        TypedField(ValueType.BYTES, 3, "bar"),
        Map(KeyType.INT32, test_message, 4, "some_map"),
        OneOf("my_oneof"),
    ]
    test_message2.fields[1].repeated = True
    test_message2.fields[3].fields = [
        TypedField(ValueType.FLOAT, 5, "floah_dude"),
        TypedField(ValueType.DOUBLE, 6, "double_trouble"),
    ]

    test_message2.definitions = [test_enum, test_message]

    test_proto = Protocol()
    test_proto.package = "testpackage"
    test_proto.definitions = [test_message, test_enum, test_message2]

    print(test_proto)


def test_message_new_invalid_identifier():
    with pytest.raises(ValidationError):
        Message("InvalidIdent!")
    with pytest.raises(ValidationError):
        Enumeration("InvalidIdent!")
    with pytest.raises(ValidationError):
        TypedField(KeyType.BOOL, 1, "InvalidIdent!")


def test_message_set_invalid_identifier():
    m = Message("ValidIdent")
    with pytest.raises(ValidationError):
        m.label = "1InvalidIdent"
        m.label = "1InvalidIdent"
