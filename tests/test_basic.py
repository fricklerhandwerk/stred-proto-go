import pytest

from stred_proto import (Enumeration, Field, Identifier, KeyType, Message,
                         Protocol, TypedField, ValidationError, ValueType)


def test_instantiation():
    test_fields = [
        TypedField(KeyType.INT32, 1, "broogle"),
        TypedField(KeyType.UINT64, 5, "doingle"),
    ]
    test_fields[1].deprecated = True

    test_fields2 = [
        TypedField(KeyType.STRING, 2, "foo"),
        TypedField(ValueType.BYTES, 3, "bar"),
    ]
    test_fields2[1].repeated = True

    test_message = Message("MyMessage")
    test_message.fields = test_fields

    test_enum = Enumeration("MyEnum")
    test_enum.fields = [
        Field(0, "default"),
        Field(1, "some"),
        Field(2, "thing"),
        Field(2, "redundant"),
    ]

    test_message2 = Message("SomeOtherMessage")
    test_message2.fields = test_fields2
    test_message2.definitions = [test_enum, test_message]

    test_proto = Protocol()
    test_proto.package = "testpackage"
    test_proto.definitions = [test_message, test_enum, test_message2]

    print(test_proto)


def test_message_new_invalid_identifier():
    with pytest.raises(ValidationError):
        Message(Identifier("InvalidIdentifier!"))
    with pytest.raises(ValidationError):
        Enumeration(Identifier("InvalidIdentifier!"))
    with pytest.raises(ValidationError):
        TypedField(KeyType.BOOL, Identifier("InvalidIdentifier!"))


def test_message_set_invalid_identifier():
    m = Message(Identifier("ValidIdentifier"))
    with pytest.raises(ValidationError):
        m.label = Identifier("1InvalidIdentifier")
