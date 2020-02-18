import pytest
from stred_proto import (Enumeration, KeyType, Message, Proto, RepeatableField,
                         ValidationError, ValueType)


def test_instantiation():
    test_fields = {
        1: RepeatableField(KeyType.INT32, "broogle"),
        5: RepeatableField(KeyType.UINT64, "doingle"),
    }
    test_fields[5].deprecated = True

    test_fields2 = {
        2: RepeatableField(KeyType.STRING, "foo"),
        3: RepeatableField(ValueType.BYTES, "bar"),
    }
    test_fields2[3].repeated = True

    test_message = Message("MyMessage")
    test_message.fields = test_fields

    test_enum = Enumeration("MyEnum")
    test_enum.enumeration = {
        0: ["default"],
        2: ["some"],
        3: ["thing", "redundant"],
    }

    test_message2 = Message("SomeOtherMessage")
    test_message2.fields = test_fields2
    test_message2.definitions = [test_enum, test_message]

    test_proto = Proto()
    test_proto.package = "testpackage"
    test_proto.definitions = [test_message, test_enum, test_message2]

    print(test_proto)


def test_message_new_invalid_identifier():
    with pytest.raises(ValidationError):
        Message("InvalidIdentifier!")
    with pytest.raises(ValidationError):
        Enumeration("InvalidIdentifier!")
    with pytest.raises(ValidationError):
        RepeatableField(KeyType.BOOL, "InvalidIdentifier!")


def test_message_set_invalid_identifier():
    m = Message("ValidIdentifier")
    with pytest.raises(ValidationError):
        m.label = "1InvalidIdentifier"
